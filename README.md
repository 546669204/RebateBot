# RebateBot 
 
## project description
Keywords: rebate WeChat Alimama robot cross platform

RebateBot, based on WeChat to establish robot channels and users to quickly generate rebate links through chat


---

You only need the minimal computer performance and an idle WeChat can to turn on a 24-hour unmanned shift RebateBot
Shopping only needs to send a link to the robot, and the robot can immediately give you a discount price and link.

---

## Function realization

### WeChat robot
This module can see the latest code [WeChatBot](https://github.com/546669204/wechatbot-xposed)
- [x] message callback
- [x] auto reply message
- [x] Add friend callback
- [x] Default to add new friends
- [x] Auto Reply (text, emoticon)


### Alimama
- [x] link identify
- [x] word command identify
- [x] Convert links to discount links
- [x] automatically assigns different pids
- [x] Scan QRcode login
- [x] Order scheduled download
- [x] automatically calculate the rebate point

### Message Processing
- [x] Sign in
- [x] Help
- [x] Automatic Binding Order
- [x] Manually bind the order
- [x] withdrawal
- [x] Receiving reminder
- [x] Payment reminder


## Run

- Download [latest version of server](https://github.com/546669204/RebateBot/releases) and [latest version of WeChat robot](https://github.com/546669204/wechatbot-xposed/releases)
- Open the primary server
- Install WeChatBot in mobile phones
- Visit http://IP:1778 for alimama login
- Test communication


## Development and build

## Environmental requirements

- golang
- mysql

## Pull code
```
Git clone https://github.com/546669204/RebateBot.git
Cd RebateBot
```

## Installation dependencies
```
Go get
Cd msgprocess
Go get
Cd ../alimama
Go get
Cd ../
```
## Configuring mysql

```
Cd msgprocess
Vim database.json

{
    "database": "webot", //database name
    "host": "127.0.0.1:3306", //Database Host
    "user": "root", //database username
    "password": "" / / database password
}
```

## Commissioning
```
Go run master.go
```

## Update log

[CHANGELOG](https://github.com/546669204/RebateBot/CHANGELOG)