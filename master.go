package main

import (
	crand "crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"./common"
	"github.com/fsnotify/fsnotify"
	"github.com/tidwall/gjson"
)

var methods = map[string]interface{}{
	"init":       ServiceInit,
	"reverse":    reverse,
	"getservice": getservice,
	"msgreturn":  MsgReturn,
	"reboot":     reboot,
}

// WorkDir 工作目录
var WorkDir string

// RunSuffix 运行目录
var RunSuffix = filepath.Ext(strings.Join(os.Args, ""))

// ConnMap 连接池
var ConnMap map[string]net.Conn

// ConnMapMutex 互斥锁
var ConnMapMutex sync.Mutex

// ConnToName 连接转名称
var ConnToName map[string]string

// ConnToNameMutex 互斥锁
var ConnToNameMutex sync.Mutex

// ServicePid PidMap
var ServicePid map[string]*os.Process

// Config 配置引用
var Config common.ConfigModel

// Msgreturn 返回Map
var Msgreturn map[string]chan common.Msg

// MsgreturnMutex 互斥锁
var MsgreturnMutex sync.Mutex

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

func init() {

	file, _ := os.OpenFile("log", os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0666)
	log.SetOutput(file)

	WorkDir = common.GetRunDir()
	ConnMap = make(map[string]net.Conn)
	ServicePid = make(map[string]*os.Process)
	ConnToName = make(map[string]string)
	Msgreturn = make(map[string]chan common.Msg)
	InitConfig()

}

//InitConfig 初始化配置
func InitConfig() {
	//初始化配置
	jsstr, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Println("读取配置项出错.程序自动退出")
		os.Exit(1)
	}
	err = json.Unmarshal(jsstr, &Config)
	if err != nil {
		log.Println("解析配置项出错.程序自动退出")
		os.Exit(1)
	}
}
func main() {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	go func() {
		for {
			select {
			case <-watcher.Events:
				//log.Println("修改文件：" + event.Name)
				//log.Println("修改类型：" + event.Op.String())
				log.Println("配置更新")
				InitConfig()
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()
	err = watcher.Add("config.json")
	if err != nil {
		log.Fatal(err)
	}
	//time.Sleep(time.Second * 1e4)

	// 监听
	l, err := net.Listen("tcp", ":188")
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer l.Close()
	RunAllService()
	for {
		// 接收一个client
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err)
			os.Exit(1)
		}

		//fmt.Printf("Received message %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())
		go handleRequest(conn)
		// 执行

	}
}
func handleRequest(conn net.Conn) {
	ipStr := conn.RemoteAddr().String()
	defer func() {
		ConnMapMutex.Lock()
		delete(ConnMap, ConnToName[conn.RemoteAddr().String()])
		delete(ConnToName, conn.RemoteAddr().String())
		ConnMapMutex.Unlock()
		fmt.Println("Disconnected :" + ipStr)
		ServiceExit(conn)
		conn.Close()
	}()
	tmpBuffer := make([]byte, 0)
	buffer := make([]byte, 1024)
	var reader []byte
	overtime := time.NewTimer(20 * time.Second)
	endflag := make(chan int)
	go func() {
		for {
			n, err := conn.Read(buffer)
			if err != nil {
				fmt.Println(time.Now().Format("2006-01-02 15:04:05"), ConnToName[conn.RemoteAddr().String()], " connection error: ", err)
				endflag <- 1
				return
			}
			if n > 0 {
				reader = nil
				tmpBuffer = common.Unpack(append(tmpBuffer, buffer[:n]...), &reader)
				if reader != nil {
					if string(reader) == "HeartBoom" {
						overtime.Reset(20 * time.Second)
					} else {
						log.Println(fmt.Sprintf(`%s ==> %s %s`, ConnToName[conn.RemoteAddr().String()], "Master", string(reader)))
						go ServiceProcess(reader, conn)
					}
				}
			}
		}
	}()
	for {
		select {
		case <-overtime.C:
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"), ConnToName[conn.RemoteAddr().String()], " Over Time Close: ")
			return
		case <-endflag:
			return
		}
	}

}

//ServiceProcess 服务请求处理
func ServiceProcess(str []byte, conn net.Conn) {
	if !gjson.Valid(string(str)) {
		//log.Println(`非json数据：`, string(str))
		return
	}
	var resp common.Msg
	err := json.Unmarshal(str, &resp)
	if err != nil {
		log.Println(err)
		return
	}
	if _, ok := Config.Routing[resp.Method]; !ok {
		log.Println("请求方法不存在请检查路由表", resp.Method)
		return
	}
	r := Config.Routing[resp.Method]

	if b := strings.Index(r, "."); b != -1 {
		s := r[:b]
		f := r[b+1:]
		address := conn.RemoteAddr().String()
		t1 := time.After(30 * time.Second)
		ctnok := false
		for {
			select {
			case <-t1:
				goto ForEnd
				break
			default:
				_, ctnok = ConnToName[address]
				if ctnok {
					goto ForEnd
				}
				break
			}
		}
	ForEnd:

		runid, _ := strconv.Atoi(resp.To)
		if runid > 1 {
			s += "|" + resp.To
		}
		if s == "wechat" && resp.To != "" {
			s += "|" + resp.To
		}

		if ctnok && (strings.Split(ConnToName[address], "|")[0] == "alimama" || strings.Split(ConnToName[address], "|")[0] == "wechat") {
			if len(strings.Split(ConnToName[address], "|")) > 1 {
				resp.To = strings.Split(ConnToName[address], "|")[1]
			} else {
				resp.To = "1"
			}
		}

		resp.Method = f
		connn, ok := ConnMap[s]
		if ok {
			if resp.ID != "" {
				tempID := resp.ID
				resp.ID = fmt.Sprintf(`%d`, CID.Get())
				resp2 := ConnWriteReturn(resp, connn)
				resp2.ID = tempID
				resp.Method = "msgreturn"
				ConnWrite(resp2, conn)
			} else {
				ConnWrite(resp, connn)
			}

		} else {
			resp.Data = `{"isrun":0}`
			resp.Method = "msgreturn"
			ConnWrite(resp, conn)
		}

		return
	}
	_, err = common.Call(methods, resp.Method, resp, conn)
	if err != nil {
		log.Println(err)
	}
}

//ServiceInit 服务注册调用
func ServiceInit(data common.Msg, conn net.Conn) {
	connname := data.Data
	if strings.Index(connname, "|") == -1 {
		for {
			_, ok := ConnMap[connname]
			if ok {
				var seed int64
				binary.Read(crand.Reader, binary.LittleEndian, &seed)
				rand.Seed(seed)
				connname = data.Data + "|" + strconv.Itoa(rand.Intn(100))
			} else {
				break
			}
		}
	}

	if _, ok := ConnMap[connname]; ok {
		conn := ConnMap[connname]
		ServiceExit(conn)
		conn.Close()
		delete(ConnMap, ConnToName[conn.RemoteAddr().String()])
		delete(ConnToName, conn.RemoteAddr().String())
	}

	ConnMapMutex.Lock()
	defer ConnMapMutex.Unlock()
	ConnMap[connname] = conn

	ConnToNameMutex.Lock()
	defer ConnToNameMutex.Unlock()
	ConnToName[conn.RemoteAddr().String()] = connname

	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "服务已上线 ", connname)
}

// CompileAndRun 编译并运行
func CompileAndRun(ServiceName string, id int) (bool, *os.Process) {
	if common.FileIsExist(filepath.Join(WorkDir, ServiceName, "main.go")) {
		compile(filepath.Join(WorkDir, ServiceName))
	}
	if common.FileIsExist(filepath.Join(WorkDir, ServiceName, ServiceName+RunSuffix)) {
		ok, pid := runservice(filepath.Join(WorkDir, ServiceName, ServiceName+RunSuffix), id)
		return ok, pid
	}
	return false, &os.Process{}
}

func compile(ServiceDir string) bool {
	cmd := exec.Command("go", "build")
	cmd.Dir = ServiceDir

	b, err := cmd.CombinedOutput()
	if len(b) > 0 {
		fmt.Println(string(b))
	}
	if err != nil {
		return false
	}
	return true
}
func runservice(ServiceFile string, id int) (bool, *os.Process) {
	cmd := exec.Command(ServiceFile, fmt.Sprintf(`%d`, id))
	go func() {
		b, _ := cmd.CombinedOutput()
		log.Println(string(b))
	}()
	/*if err != nil {
		return false, &os.Process{}
	}*/
	return true, cmd.Process
}

//RunAllService 运行全部服务
func RunAllService() {
	for _, k := range Config.Services {
		for i := 1; i <= k.Run; i++ {
			ok, pid := CompileAndRun(k.Name, i)
			if ok {
				log.Println(`服务运行成功 `, k, `PID `, pid)
				if i > 1 {
					ServicePid[fmt.Sprintf(`%s|%d`, k.Name, i)] = pid
				} else {
					ServicePid[k.Name] = pid
				}
			}
		}
	}
}

//ConnWrite Socket写函数
func ConnWrite(data common.Msg, conn net.Conn) {
	b, _ := json.Marshal(data)
	log.Println(fmt.Sprintf(`%s ==> %s %s`, "Master", ConnToName[conn.RemoteAddr().String()], string(b)))
	conn.Write(common.Packet(b))
}

//ServiceExit Socket推出回调
func ServiceExit(conn net.Conn) {
	//服务被关闭或意外断开连接 删除变量 强制Kill
	name := ConnToName[conn.RemoteAddr().String()]
	ConnMapMutex.Lock()
	defer ConnMapMutex.Unlock()
	delete(ConnMap, name)
	delete(ServicePid, name)
}

//ConnWriteReturn 带回调的发送
func ConnWriteReturn(data common.Msg, conn net.Conn) (msg common.Msg) {
	id := fmt.Sprintf(`%d`, CID.Get()) //fmt.Sprintf(`%s%d%s`, "Master", time.Now().UnixNano(), conn.RemoteAddr().String())
	data.ID = id
	MsgreturnMutex.Lock()
	Msgreturn[id] = make(chan common.Msg, 2)
	MsgreturnMutex.Unlock()

	ConnWrite(data, conn)
	t := time.After(time.Second * 60)
	msg = common.Msg{}
	msg.Method = "msgreturn"

	for {
		select {
		case resp, ok := <-Msgreturn[id]:
			if ok {
				MsgreturnMutex.Lock()
				close(Msgreturn[id])
				delete(Msgreturn, id)
				MsgreturnMutex.Unlock()
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

//MsgReturn 消息返回
func MsgReturn(data common.Msg, conn net.Conn) {
	Msgreturn[data.ID] <- data
}

func reverse(data common.Msg, conn net.Conn) {
	s := ""
	sbyte := []byte(data.Data)
	for index := 1; index <= len(sbyte); index++ {
		s = s + string(sbyte[len(sbyte)-index:len(sbyte)-index+1])
	}
	data.Data = s
	data.Method = "msgreturn"
	ConnWrite(data, conn)
}

func getservice(data common.Msg, conn net.Conn) {
	var ser []string

	for _, k := range Config.Services {
		var str []string

		for ck := range ConnMap {
			if strings.Index(ck, k.Name) == 0 {
				var msg common.Msg
				msg.Method = "islogin"
				msg.Data = ""
				conn, ok := ConnMap[ck]
				if !ok {
					str = append(str, `{"isrun":0,"islogin":0,"name":"","lastlogin":""}`)
				} else {
					resp := ConnWriteReturn(msg, conn)
					str = append(str, resp.Data)
				}
			}
		}
		ser = append(ser, fmt.Sprintf(`"%s":[%s]`, k.Name, strings.Join(str, ",")))
	}
	data.Data = fmt.Sprintf(`{"status":1,"services":{%s}}`, strings.Join(ser, ","))
	data.Method = "msgreturn"
	ConnWrite(data, conn)
}

func reboot(data common.Msg, conn net.Conn) {
	index := strings.Index(data.Data, "|")
	sname := data.Data[:index]
	runid, _ := strconv.Atoi(data.Data[index+1:])
	ssname := ""
	if runid > 1 {
		ssname = fmt.Sprintf(`%s|%d`, sname, runid)
	} else {
		ssname = sname
	}
	if _, ok := ConnMap[ssname]; ok {
		ConnMapMutex.Lock()
		defer ConnMapMutex.Unlock()
		ConnMap[ssname].Close()
	}
	time.Sleep(time.Second)
	ok, pid := CompileAndRun(sname, runid)
	if ok {
		log.Println(`服务运行成功 `, sname, `PID `, pid)
		ServicePid[ssname] = pid
	}
	data.Data = fmt.Sprintf(`{"status":1}`)
	data.Method = "msgreturn"
	ConnWrite(data, conn)
}
