package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/546669204/RebateBot/common"

	"github.com/gin-gonic/gin"
)

func initWebApi() {
	router := gin.New()
	api := router.Group("api")
	api.GET("/getService", getService)
	api.GET("/getUserData", getUserData)
	api.GET("/getOrderData", getOrderData)
	api.GET("/getWithdrawData", getWithdrawData)
	api.POST("/withdrawPay", withdrawPay)
	api.GET("/getTemplateData", getTemplateData)
	api.POST("/setTemplateData", setTemplateData)
	api.POST("/tblogin", tblogin)
	api.POST("/tbchecklogin", tbchecklogin)
	api.POST("/reboot", reboot)
	router.GET("/", func(c *gin.Context) {
		c.File("./fanliadmin/dist/index.html")
	})
	router.Static("/static", "./fanliadmin/dist/static")
	router.Run(":1778")
}
func getService(c *gin.Context) {
	var m common.Msg
	m.Data = ""
	m.Method = "getservice"
	resp := Client.ConnWriteReturn(m)
	c.Data(200, "application/json", []byte(resp.Data))
}
func tblogin(c *gin.Context) {
	type parmModel struct {
		Data string `form:"data"`
	}
	var parm parmModel
	c.ShouldBind(&parm)
	var m common.Msg
	m.Data = ""
	m.To = parm.Data
	m.Method = "tblogin"
	resp := Client.ConnWriteReturn(m)
	c.Data(200, "application/json", []byte(resp.Data))
}
func tbchecklogin(c *gin.Context) {
	type parmModel struct {
		Data string `form:"data"`
	}
	var parm parmModel
	c.ShouldBind(&parm)
	var m common.Msg
	m.Data = ""
	m.To = parm.Data
	m.Method = "tbchecklogin"
	resp := Client.ConnWriteReturn(m)
	c.Data(200, "application/json", []byte(resp.Data))
}
func reboot(c *gin.Context) {
	type parmModel struct {
		Data string `form:"data"`
	}
	var parm parmModel
	c.ShouldBind(&parm)
	var m common.Msg
	m.Data = parm.Data
	m.Method = "reboot"
	resp := Client.ConnWriteReturn(m)
	c.Data(200, "application/json", []byte(resp.Data))
}
func getUserData(c *gin.Context) {
	type parmModel struct {
		Page     uint `form:"page"`
		PageSize uint `form:"pageSize"`
	}
	var parm parmModel
	c.ShouldBind(&parm)
	var j []MysqlUser
	err := sess.Collection("xm_user").Find().Page(parm.Page).Paginate(parm.PageSize).All(&j)
	count, _ := sess.Collection("xm_user").Find().Count()
	if err != nil {
		log.Println("getuserdata", err)
		c.JSON(200, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 0, "msg": "", "data": j, "count": count})
}

func getOrderData(c *gin.Context) {
	type parmModel struct {
		Page     uint `form:"page"`
		PageSize uint `form:"pageSize"`
	}
	var parm parmModel
	c.ShouldBind(&parm)
	var j []MysqlOrders
	err := sess.Collection("xm_orders").Find().Page(parm.Page).Paginate(parm.PageSize).OrderBy(`create_time desc`).All(&j)
	count, _ := sess.Collection("xm_orders").Find().Count()
	if err != nil {
		c.JSON(200, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 0, "msg": "", "data": j, "count": count})
}

func getWithdrawData(c *gin.Context) {
	type parmModel struct {
		Page     uint `form:"page"`
		PageSize uint `form:"pageSize"`
	}
	var parm parmModel
	c.ShouldBind(&parm)
	var j []MysqlOrders
	err := sess.Collection("xm_withdraw").Find().Where(`status = 0`).Page(parm.Page).Paginate(parm.PageSize).OrderBy(`create_time desc`).All(&j)
	count, _ := sess.Collection("xm_withdraw").Find().Where(`status = 0`).Count()
	if err != nil {
		c.JSON(200, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 0, "msg": "", "data": j, "count": count})
}

func withdrawPay(c *gin.Context) {
	type parmModel struct {
		ID uint `form:"id"`
	}
	var parm parmModel
	c.ShouldBind(&parm)
	_, err := sess.Update("xm_withdraw").Where(`id = ?`, parm.ID).Set(`status = 1`).Exec()
	if err != nil {
		c.JSON(200, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 0, "msg": "", "data": ""})
}

func getTemplateData(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "", "data": TemplateList})
}

func setTemplateData(c *gin.Context) {
	type parmModel struct {
		Data string `form:"data"`
	}
	var parm parmModel
	c.ShouldBind(&parm)
	err := ioutil.WriteFile(filepath.Join(filepath.Dir(os.Args[0]), "template.json"), []byte(parm.Data), 0666)
	if err != nil {
		c.JSON(200, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 0, "msg": ""})
}
