package common

import (
	"testing"
	"time"
)

func Test_Conn(t *testing.T) {

	for {
		time.Sleep(5 * time.Second)
		Client := InitClient()
		data := `{"method":"init","data":"` + "wechat" + `"}`
		Client.WriteChan <- Packet([]byte(data))
		go Client.ServiceHandle(Client.Conn)
	}
}
