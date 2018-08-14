package common

import (
	"log"
	"testing"
)

func Test_Division_1(t *testing.T) {

	var r, l []byte
	l = Unpack(Packet([]byte("我啊啊啊啊啊啊")), &r)

	log.Println(string(l))
	log.Println(string(r))

}
