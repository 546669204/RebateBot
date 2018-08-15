# RebateBot 返利机器人
 
## 项目描述
关键词: 返利 微信 阿里妈妈 机器人 跨平台

返利机器人，基于微信建立机器人通道与用户通过聊天快速生成返利链接


---

利用闲置微信和极小的电脑性能开启24小时无人轮值返利机器人  
购物只需要发送链接给机器人，机器人能马上给你回复优惠价格及链接 

<img src="https://github.com/546669204/RebateBot/blob/master/screenshots/demo.png" style="width:500px" />  

---

## 功能实现

### 微信机器人
这个模块在这里可以看到最新的代码[微信机器人](https://github.com/546669204/wechatbot-xposed)
- [x] 消息回调
- [x] 自动回复消息
- [x] 新增好友回调
- [x] 默认同意新增好友
- [x] 自动回复(文字,表情)


### 阿里妈妈
- [x] 链接识别
- [x] 淘口令识别
- [x] 链接转换返利链接
- [x] 自动分配不同pid
- [x] 扫码登录
- [x] 订单定时下载
- [x] 自动计算返利反点

### 消息处理
- [x] 签到
- [x] 帮助
- [x] 自动绑定订单
- [x] 手动绑定订单
- [x] 提现
- [x] 收货提醒
- [x] 支付提醒


## 运行

- 下载[最新版服务器](https://github.com/546669204/RebateBot/releases) 和 [最新版微信机器人](https://github.com/546669204/wechatbot-xposed/releases)
- 启动主服务器
- 为手机安装微信机器人
- 访问http://IP:1778为alimama登录
- 测试通讯


## 开发和构建

## 环境要求

- golang
- mysql 

## 拉取代码
```
git clone https://github.com/546669204/RebateBot.git
cd RebateBot
```

## 安装依赖
```
go get 
cd msgprocess
go get 
cd ../alimama
go get 
cd ../
```
## 配置mysql

```
cd msgprocess
vim database.json

{
    "database": "webot",		//数据库名字
    "host": "127.0.0.1:3306",	//数据库Host
    "user": "root",				//数据库用户名
    "password": "" 				//数据库密码
}
```

## 调试运行
```
go run master.go
```

## 更新日志 

[CHANGELOG](https://github.com/546669204/RebateBot/CHANGELOG)


