package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/extrame/xls"

	"../common"
	httpdo "github.com/546669204/golang-http-do"
	tk "github.com/546669204/taobaoke"
	"github.com/tidwall/gjson"
)

var methods = map[string]interface{}{
	"ExitHook":    ExitHook,
	"search":      search,
	"settuiguang": settuiguang,
	"islogin":     islogin,
	"login":       login,
	"checklogin":  checklogin,
	"getordernum": getordernum,
}
var IsLogin bool = false
var Client *common.Client

var RunTimeDate string
var TotalRec, TotalAlipayRec float64
var TotalAlipayNum, TotalMixClick int64
var t3 <-chan time.Time

func main() {
	os.Chdir(filepath.Dir(os.Args[0]))
	httpdo.Debug = true
	common.LogToFile()
	Client = common.InitClient()
	defer Client.Conn.Close()
	Client.InitMethods(methods)
	Client.SericeInit()
	go func() {
		Client.HeartBeat()
	}()
	go func() {
		//启动线程,定时访问alimam 保持cookies
		t1 := time.Tick(5 * 60 * time.Second)
		t2 := time.Tick(3 * 60 * time.Second)
		t3 = time.After(1 * time.Hour)
		RunTimeDate = time.Now().Format("2006-01-02")
		for {
			select {
			case <-t1:
				if IsLogin {
					tk.KeepLogin()
				}
				break
			case <-t2:
				if IsLogin {
					checkorder()
				}
				break
			case <-t3:
				if IsLogin {
					getalldata()
				}
			}

		}
	}()
	tk.ChromeUserDataDIR = path.Join(os.TempDir(), fmt.Sprintf(`%d`, time.Now().Unix()))
	copyDir(filepath.Join(filepath.Dir(os.Args[0]), "alillll"), tk.ChromeUserDataDIR)
	if tk.LoadLogin() && tk.IsLogin() {
		IsLogin = true
		tk.GetUnionPubContextInfo()
	}

	Client.ServiceHandle(Client.Conn)
}

func ExitHook() {
	//退出程序钩子
	tk.SaveLogin()
}

func login(data common.Msg) {
	if tk.LoadLogin() == false || tk.IsLogin() == false {
		/*
			var qrstr, lg string
			log.Println("即将登陆阿里妈妈请使用(淘宝)扫描二维码")
			if !tk.Login(&qrstr, &lg) {
				log.Println("阿里妈妈登陆失败 程序即将退出")
				//os.Exit(2)
				data.Data = fmt.Sprintf(`{"msg":"获取失败","status":0,"qrcode":"%s","lgtoken":"%s"}`, qrstr, lg)
				data.Method = "msgreturn"
				Client.ConnWrite(data)
				return
			}

			data.Data = fmt.Sprintf(`{"msg":"获取成功","status":1,"qrcode":"%s","lgtoken":"%s"}`, qrstr, lg)
			data.Method = "msgreturn"
			Client.ConnWrite(data)
			log.Println(data)
			return
		*/

		data.Data = fmt.Sprintf(`{"msg":"获取成功","status":1,"qrcode":"%s"}`, tk.BrowserLogin())
		data.Method = "msgreturn"
		Client.ConnWrite(data)
		return
	}
	IsLogin = true
	tk.GetUnionPubContextInfo()
	data.Data = `{"msg":"已登录!","status":2}`
	data.Method = "msgreturn"
	Client.ConnWrite(data)
}
func checklogin(data common.Msg) {
	//lgtoken := gjson.Parse(data.Data).Get("lgtoken").String()
	if status, _ := tk.BrowserCheckLogin(); status {
		tk.CookiesTotb_token()

		tk.GetUnionPubContextInfo()
		tk.SaveLogin()
		log.Println("登录成功")
		IsLogin = true
		data.Data = `{"msg":"登录成功!","status":1}`
		data.Method = "msgreturn"
		Client.ConnWrite(data)
		return
	}
	data.Data = `{"msg":"登录失败!","status":0}`
	data.Method = "msgreturn"
	Client.ConnWrite(data)

}
func search(data common.Msg) {
	p := tk.Search(gjson.Parse(data.Data).Get("keyword").String())
	if p.ID == "" {
		data.Data = `{"msg":"找不到响应返利内容!","status":0}`
		data.Method = "msgreturn"
		Client.ConnWrite(data)
		return
	}
	b, _ := json.Marshal(p)
	data.Data = fmt.Sprintf(`{"msg":"ok","status":1,"products":%s}`, string(b))
	data.Method = "msgreturn"
	Client.ConnWrite(data)
}

func settuiguang(data common.Msg) {
	id := gjson.Parse(data.Data).Get("id").String()
	pstr := gjson.Parse(data.Data).Get("pid").String()

	sa := tk.NewSelfAdzone2(id)
	if len(sa) == 0 {
		data.Data = fmt.Sprintf(`{"msg":"ok","status":0}`)
		data.Method = "msgreturn"
		Client.ConnWrite(data)
		return
	}
	a := ""
	b := ""
	if len(pstr) != 0 {
		for _, value := range sa {
			for _, value2 := range value.Adzoneid {
				tmp := value.Siteid + "-" + value2
				if strings.Index(pstr, tmp) <= -1 {
					a = value.Siteid
					b = value2
					break
				}
			}
		}
		log.Println("重复使用pid")
		a = sa[0].Siteid
		b = sa[0].Adzoneid[0]
	} else {
		a = sa[0].Siteid
		b = sa[0].Adzoneid[0]
	}

	tk.SelfAdzoneCreate(a, b)
	l := tk.GetAuctionCode(id, a, b)
	lstr, _ := json.Marshal(l)
	data.Data = fmt.Sprintf(`{"msg":"ok","status":1,"link":%s,"pid":"%s"}`, string(lstr), a+"-"+b)
	data.Method = "msgreturn"
	Client.ConnWrite(data)
}

func islogin(data common.Msg) {

	if !IsLogin {
		data.Data = fmt.Sprintf(`{"isrun":1,"islogin":0,"name":"","lastlogin":"","runid":%s}`, os.Args[1])
	} else {
		var d map[string]string
		d = make(map[string]string)
		d["isrun"] = "1"
		d["islogin"] = "1"
		d["name"] = tk.UserInfo.Name
		d["lastlogin"] = tk.UserInfo.LastLogin
		d["runid"] = os.Args[1]
		d["pic"] = tk.UserInfo.Pic
		f, _ := json.Marshal(d)
		data.Data = string(f)
	}
	data.Method = "msgreturn"
	Client.ConnWrite(data)
}

func checkorder() {
	EndTime := time.Now().Format("2006-01-02")
	RunTimeDate = time.Now().AddDate(0, 0, -90).Format("2006-01-02")
	b := tk.GetTbkPaymentDetails(RunTimeDate, EndTime)
	if TotalAlipayNum != gjson.ParseBytes(b).Get("data").Get("paginator").Get("items").Int() {
		TotalAlipayNum = gjson.ParseBytes(b).Get("data").Get("paginator").Get("items").Int()
		getalldata()
	}
}

func getXls(b []byte) string {
	if xlFile, err := xls.OpenReader(bytes.NewReader(b), "utf-8"); err == nil {
		if sheet1 := xlFile.GetSheet(0); sheet1 != nil {
			fmt.Print("Total Lines ", sheet1.MaxRow, sheet1.Name, "\n")
			row0 := sheet1.Row(0)
			var head []string
			for i2 := row0.FirstCol(); i2 <= row0.LastCol(); i2++ {
				head = append(head, row0.Col(i2))
			}
			var ss []map[string]string
			for i := 1; i <= (int(sheet1.MaxRow)); i++ {
				var sss map[string]string
				sss = make(map[string]string)
				row1 := sheet1.Row(i)
				for i2 := row1.FirstCol(); i2 <= row1.LastCol(); i2++ {
					sss[head[i2]] = row1.Col(i2)
				}
				ss = append(ss, sss)
			}
			jsonbyte, _ := json.Marshal(ss)
			return string(jsonbyte)
		}
	}
	return ""
}

func getalldata() {
	RunTimeDate = time.Now().AddDate(0, 0, -90).Format("2006-01-02")
	EndTime := time.Now().Format("2006-01-02")
	b := tk.Download(RunTimeDate, EndTime)
	j := getXls(b)
	var rdata common.Msg
	rdata.Data = j
	rdata.Method = "parsexls"
	Client.ConnWrite(rdata)
	t3 = time.After(time.Hour)
}

func getordernum(data common.Msg) {
	var tdata map[string]interface{}
	tdata = make(map[string]interface{})
	tdata["TotalAlipayNum"] = TotalAlipayNum
	EndTime := time.Now().Format("2006-01-02")
	RunTimeDate = time.Now().AddDate(0, 0, -90).Format("2006-01-02")
	b := tk.GetTbkPaymentDetails(RunTimeDate, EndTime)

	tdata["remote"] = string(b)
	data.Method = "msgreturn"
	jj, _ := json.Marshal(tdata)
	data.Data = string(string(jj))
	Client.ConnWrite(data)
}

func copyFile(src, dst string) error {
	srcFile, err := os.OpenFile(src, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, srcFile)
	return err
}
func copyDir(src, dst string) error {
	log.Println(dst)
	filepath.Walk(src, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		if err != nil { //忽略错误
			return err
		}

		if fi.IsDir() { // 忽略目录
			return nil
		}
		filename, _ = filepath.Rel(src, filename)
		filename = strings.Replace(filename, "\\", "/", -1)

		os.MkdirAll(filepath.Join(dst, filepath.Dir(filename)), 0666)
		copyFile(filepath.Join(src, filename), filepath.Join(dst, filename))
		return nil
	})
	return nil
}
