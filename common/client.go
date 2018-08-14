package common

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/tidwall/gjson"
)

var ServiceName = GetServiceName(os.Args[0])

type Client struct {
	Conn      net.Conn
	Msgreturn map[string]chan Msg
	Methods   map[string]interface{}
	WriteChan chan []byte
}

type ConnInsID struct {
	ID   int64
	lock sync.Mutex
}

func (self *ConnInsID) Get() int64 {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.ID++
	return self.ID
}

var CID ConnInsID

func InitClient() *Client {
	conn, err := net.Dial("tcp", ":188")
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	c := new(Client)
	c.Conn = conn
	c.Msgreturn = make(map[string]chan Msg)
	c.WriteChan = make(chan []byte, 16)
	return c
}

func (this *Client) InitMethods(m map[string]interface{}) {
	this.Methods = m
	this.Methods["msgreturn"] = this.MsgReturn
}

func (this *Client) MsgReturn(data Msg) {
	this.Msgreturn[data.ID] <- data
}
func (this *Client) ConnWrite(data Msg) {
	b, _ := json.Marshal(data)
	this.WriteChan <- Packet(b)
}
func (this *Client) ConnWriteReturn(data Msg) (msg Msg) {
	id := fmt.Sprintf(`%d`, CID.Get()) //fmt.Sprintf(`%s%d`, ServiceName, time.Now().UnixNano())
	data.ID = id
	this.Msgreturn[id] = make(chan Msg, 2)
	this.ConnWrite(data)
	t := time.After(time.Second * 60)
	msg.Method = "msgreturn"
	for {
		select {
		case resp, ok := <-this.Msgreturn[id]:
			if ok {
				close(this.Msgreturn[id])
				delete(this.Msgreturn, id)
				//log.Println("return", resp)
				msg = resp
				goto ForEnd
			}
		case <-t:
			goto ForEnd
		}
	}
ForEnd:
	return
}

func (this *Client) SericeInit() {
	data := `{"method":"init","data":"` + ServiceName + `"}`
	if len(os.Args) > 1 && os.Args[1] != "1" {
		data = `{"method":"init","data":"` + ServiceName + "|" + os.Args[1] + `"}`
	}
	this.WriteChan <- Packet([]byte(data))
}
func (this *Client) ServiceProcess(str []byte) {
	if !gjson.Valid(string(str)) {
		return
	}
	var resp Msg
	err := json.Unmarshal(str, &resp)
	if err != nil {
		log.Println(err)
	}
	//log.Println(this.Methods, resp.Method, resp)
	_, err = Call(this.Methods, resp.Method, resp)
	if err != nil {
		log.Println(err)
	}
}

func (this *Client) ServiceHandle(conn net.Conn) {
	go this.ServiceWriteHandle(conn)
	tmpBuffer := make([]byte, 0)
	buffer := make([]byte, 1024)
	var reader []byte
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}
		reader = nil
		tmpBuffer = Unpack(append(tmpBuffer, buffer[:n]...), &reader)
		if reader != nil {
			log.Printf("==> %s \n", reader)
			go this.ServiceProcess(reader)
		}
	}
}

func (this *Client) ServiceWriteHandle(conn net.Conn) {
	defer func() {
		fmt.Println("断线退出")
		conn.Close()
		Call(this.Methods, "ExitHook")
		os.Exit(1)
	}()
	var err error
	var data []byte
	var ok bool
	for {
		select {
		case data, ok = <-this.WriteChan:
			if ok {
				_, err = conn.Write(data)
				log.Printf("<== %s \n", string(data))
				if err != nil {
					return
				}
			}
			break
		}
	}
}

func (this *Client) HeartBeat() {
	heartbyte := Packet([]byte("HeartBoom"))
	t := time.Tick(time.Second * 10)
	for {
		select {
		case <-t:
			this.WriteChan <- heartbyte
		}
	}
}
