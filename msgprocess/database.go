package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"

	"upper.io/db.v3"

	"upper.io/db.v3/mysql"
)

var mysqlsettings = mysql.ConnectionURL{
	Database: `webot`,          //数据库名字
	Host:     `127.0.0.1:3306`, //数据库Host
	User:     `root`,           //数据库用户名
	Password: ``,               //数据库密码
}

func initMysqlSetting() {
	jsstr, err := ioutil.ReadFile(filepath.Join(filepath.Dir(os.Args[0]), "database.json"))
	if err != nil {
		log.Println("读取配置项出错.程序自动退出")
		os.Exit(1)
	}
	mysqlsettings.Database = gjson.ParseBytes(jsstr).Get("database").String()
	mysqlsettings.Host = gjson.ParseBytes(jsstr).Get("host").String()
	mysqlsettings.User = gjson.ParseBytes(jsstr).Get("user").String()
	mysqlsettings.Password = gjson.ParseBytes(jsstr).Get("password").String()
}
func databaseVersioCheck() {
	var err error
	if !sess.Collection("xm_config").Exists() {

		sqltxt, err := ioutil.ReadFile(filepath.Join(filepath.Dir(os.Args[0]), "sql", "1.sql"))
		if err != nil {
			log.Panicln(err.Error())
		}
		err = multiLineExec(string(sqltxt))
		if err != nil {
			log.Panicln(err.Error())
		}
		_, err = sess.InsertInto("xm_config").Columns("key", "value").Values("version", "1").Exec()
		if err != nil {
			log.Panicln(err.Error())
		}

	} else {
		files, _ := ioutil.ReadDir(filepath.Join(filepath.Dir(os.Args[0]), "sql"))

		var versions []int
		for _, file := range files {
			num, _ := strconv.Atoi(strings.Replace(file.Name(), ".sql", "", -1))
			versions = append(versions, num)
		}
		sort.Ints(versions)

		var aa MysqlConfig
		var where db.Cond
		where = make(db.Cond)
		where["key"] = "version"
		sess.Collection("xm_config").Find().Where(where).One(&aa)
		num, _ := strconv.Atoi(aa.Value)

		for _, i := range versions {
			if num < i {
				sql, err := ioutil.ReadFile(filepath.Join(filepath.Dir(os.Args[0]), "sql", fmt.Sprintf("%d.sql", i)))
				if err != nil {
					log.Panicln(err.Error())
				}
				err = multiLineExec(string(sql))
				if err != nil {
					log.Panicln(err.Error())
				}
				num++
			}
			continue
		}
		if aa.Value != fmt.Sprintf(`%d`, num) {
			_, err = sess.Update("xm_config").Where(where).Set(`value = ?`, num).Exec()
			if err != nil {
				log.Panicln(err.Error())
			}
		}

	}
}

func multiLineExec(sqltxt string) error {
	sqltxts := strings.Split(sqltxt, ";")
	tx, _ := sess.NewTx(sess.Context())
	for index := 0; index < len(sqltxts); index++ {
		log.Println("=" + sqltxts[index] + "*")
		_, err := tx.QueryRow(sqltxts[index])
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}
