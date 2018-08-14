package main

import (
	crand "crypto/rand"
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	db "upper.io/db.v3"
	"upper.io/db.v3/mysql"
)

func TestReverse(t *testing.T) {
	data := "1236571236821"
	s := ""
	sbyte := []byte(data)
	for index := 1; index <= len(sbyte); index++ {
		s = s + string(sbyte[len(sbyte)-index:len(sbyte)-index+1])
	}
	log.Println(s)
}

func TestGetServiceName(t *testing.T) {
	dir := filepath.Dir(os.Args[0])
	log.Println(filepath.Base(dir))
}

func TestJson(t *testing.T) {
	var s MemberItemModel
	sss := `{
		"UserName": "@44aadbe25a93b21041526809cce23a3f",
		"NickName": "明",
		"Sex": 1,
		"HeadImgUpdateFlag": 1,
		"ContactType": 0,
		"Alias": "",
		"ChatRoomOwner": "",
		"HeadImgUrl": "/cgi-bin/mmwebwx-bin/webwxgeticon?seq=691720364&username=@44aadbe25a93b21041526809cce23a3f&skey=",
		"ContactFlag": 3,
		"MemberCount": 0,
		"MemberList": [],
		"HideInputBarFlag": 0,
		"Signature": "",
		"VerifyFlag": 0,
		"RemarkName": "weid5",
		"Statues": 0,
		"AttrStatus": 50529183,
		"Province": "其他",
		"City": "德国",
		"SnsFlag": 1,
		"KeyWord": "zzx"
		}`
	err := json.Unmarshal([]byte(sss), &s)
	log.Println(s, err)
}

func TestRand(t *testing.T) {
	var seed int64
	binary.Read(crand.Reader, binary.LittleEndian, &seed)
	rand.Seed(seed)
	m := rand.Float64()
	m = math.Trunc(m*1e2) * 1e-2
	log.Println(m)
}

func TestUnicodeEmojiCode(t *testing.T) {
	s := "1F30D"
	p, _ := strconv.ParseInt(s, 16, 32)
	log.Println(string(rune(p)))
}

func TestTemplate(t *testing.T) {
	var err error
	sess, err = mysql.Open(mysqlsettings)
	if err != nil {
		log.Fatal("Conn ==> ", err)
	}
	defer sess.Close()
	sqltxt, err := ioutil.ReadFile(filepath.Join(`C:\Users\Administrator\Desktop\git\fanlibot\msgprocess`, "sql", "1.sql"))
	if err != nil {
		log.Panicln(err.Error())
	}

	sqltxts := strings.Split(string(sqltxt), ";")
	log.Println(sqltxts)
	tx, _ := sess.NewTx(sess.Context())
	for index := 0; index < len(sqltxts); index++ {
		log.Println("=" + sqltxts[index] + "*")
		_, errexec := tx.QueryRow(sqltxts[index])
		if errexec != nil {
			tx.Rollback()
			log.Panicln(err.Error())
		}
	}

	tx.Commit()
}

func TestMysql(t *testing.T) {
	var settings = mysql.ConnectionURL{
		Database: `webot`,
		Host:     `127.0.0.1:3306`,
		User:     `root`,
		Password: `a866552241`,
	}
	sess, err := mysql.Open(settings)
	if err != nil {
		log.Fatal("CONN ", err)
	}
	defer sess.Close()

	sess.SetLogging(false)
	var j []MysqlUser
	log.Println(sess.Collection("xm_user").Find().All(&j))
	js, _ := json.Marshal(j)
	log.Println(string(js))
	return
	var userconn = sess.Collection("user")
	var where db.Cond
	where = make(db.Cond)
	//全部查询
	var all []MysqlUser
	err = userconn.Find(where).All(&all)
	log.Println(all, err)
	//单条查询
	var one MysqlUser
	where["id"] = 8
	err = userconn.Find(where).One(&one)
	log.Println(one, err)
	//插入数据
	var inster MysqlUser
	inster.Money = 50
	inster.Weid = "weid98"
	inster.Wename = "test"
	id, err := (userconn.Insert(inster))
	if err != nil {
		log.Println(err)
	} else {

		log.Println(id)
	}
	//修改数据
	var update MysqlUser
	where = make(db.Cond)
	where["id"] = id
	update.Money = 55
	userconn.Find(db.Cond{"id": id}).Update(update)
	//删除数据
	where = make(db.Cond)
	where["id"] = 21
	userconn.Find(where).Delete()
}

func TestFunc(t *testing.T) {
	log.Println("123123")
	reg := regexp.MustCompile(`(\%| )`)
	log.Println(strconv.ParseFloat(reg.ReplaceAllString("2.50 %", ""), 10))
	log.Println(`asd`)
}
