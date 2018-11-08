<template>
  <div>
    <el-row>
      <el-col :span="12">
        <el-card class="box-card" v-for="(todo, index) in services.wechat" :key="index" v-if="false">
          <div slot="header" class="clearfix">
            <span>微信</span>
            <el-button style="float: right; padding: 3px 0" type="text" icon="el-icon-tickets" v-if="!todo.islogin" @click="login('wechat|' + todo.runid)">扫码登录</el-button>
          </div>
          <div class="image">
            <img v-bind:src="todo.islogin ? todo.pic : '/images/notlogin.png' ">
          </div>
          <div class="content">
            <a class="header">{{todo.islogin ? todo.name: '未登录'}}</a>
            <div class="meta">
              <span class="date">最后登录时间{{todo.lastlogin}}</span>
            </div>
          </div>
          <div class="extra content">
            <a>
              <i class="user icon"></i> {{todo.friend}} 个 好友 </a>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card class="box-card" v-for="(todo, index) in services.alimama" :key="index">
          <div slot="header" class="clearfix">
            <span>阿里巴巴</span>
            <el-button style="float: right; padding: 3px 0" type="text" icon="el-icon-tickets" v-if="!todo.islogin" @click="login('alibaba|' + todo.runid)">扫码登录</el-button>
          </div>
          <div class="image">
            <img v-bind:src="todo.pic" v-if="todo.islogin">
            <img src='../../assets/notlogin.png' v-else>
          </div>
          <div class="content">
            <a class="header">{{todo.islogin ? todo.name: '未登录'}}</a>
            <div class="meta">
              <span class="date">最后登录时间{{todo.lastlogin}}</span>
            </div>
          </div>
          <div class="extra content">
            <a>
              <i class="el-icon-view"></i> {{todo.friend}} 个 好友 </a>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-dialog :title="dialogLogin.title" :visible.sync="dialogLogin.visible">
      <div>
        <div v-loading="dialogLogin.loading" class="qrcode">
          <img :src="dialogLogin.src" >
        </div>
        <p>{{dialogLogin.description}}</p>
      </div>
    </el-dialog>
  </div>
</template>
<style>
  .box-card {
    text-align: left;
  }

  img {
    max-width: 100%;
  }

  .qrcode {
    width: 300px;
    height: 300px;
    border-radius: 15px;
  }

</style>

<script>
  import api from "@/api";
  export default {
    name: "IndexLogin",
    props: ["services"],
    data() {
      return {
        dialogLogin: {
          visible: false,
          title: "",
          description: "",
          loading: false,
          src:null
        },
        wetimeout: null,
        tbtimeout: null,
        tbcheck: false,
        wecheck: false
      };
    },
    methods: {
      login(id) {
        var sname = id.split("|")[0];
        var runid = id.split("|")[1];
        if (sname == "wechat") {
          this.dialogLogin.description = "请使用微信客户端扫码";
          this.dialogLogin.title = "微信登录";
        }
        if (sname == "alimama") {
          this.dialogLogin.description = "请使用淘宝客户端扫码";
          this.dialogLogin.title = "淘宝登录";
        }
        this.dialogLogin.visible = true;
        this.dialogLogin.loading = true;
        var method = sname == "wechat" ? "welogin" : "tblogin";
        api[method]({data: runid}).then(res => {
            if (res.data.status == 0) {
              alert("获取失败");
              this.dialogLogin.visible = false;
              return;
            }
            if (res.data.status == 2) {
              alert("登录成功");
              this.$emit("getList");
              return;
            }

            this.dialogLogin.src = JSON.parse(res.data.qrcode).src;
            this.dialogLogin.loading = false;
            if (sname == "wechat") {
              this.wechecklogin(res.uuid);
            } else {
              this.tbchecklogin(runid);
            }
          });
      },
      wechecklogin(uuid) {
        let _this = this;
        clearTimeout(_this.wetimeout);
        _this.wetimeout = setInterval(function () {
          if (_this.wecheck) {
            return;
          }
          _this.wecheck = true;
          api.v1({
            method: "wechecklogin",
            data: `{"uuid":"` + uuid + `"}`
          }).then(res => {
            _this.wecheck = false;
            if (!res.status) {
              return;
            }
            clearInterval(_this.wetimeout);
            _this.dialogLogin.visible = false;
            _this.$emit("getList");
          });
        }, 1000);
      },
      tbchecklogin(runid) {
        let _this = this;
        clearTimeout(_this.tbtimeout);
        _this.tbtimeout = setInterval(function () {
          if (_this.tbcheck) {
            return;
          }
          _this.tbcheck = true;
          api.tbchecklogin({data: runid}).then(res => {
            _this.tbcheck = false;
            if (!res.data.status) {
              return;
            }
            clearInterval(_this.tbtimeout);
            this.dialogLogin.visible = false;
            this.$emit("getList");
          });
        }, 1000);
      }
    }
  };

</script>
