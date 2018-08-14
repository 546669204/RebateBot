<template>
    <el-row type="flex" class="row-bg" justify="center">
        <el-col :span="24">
            <el-tabs v-model="activeName" @tab-click="handleTabClick">
                <el-tab-pane label="Home" name="home"><Home :services.sync="services" ref="home"/></el-tab-pane>
                <el-tab-pane label="Login" name="login"><Login :services="services" @getList="servicesGetList" /></el-tab-pane>
                <el-tab-pane label="User" name="user"><User/></el-tab-pane>
                <el-tab-pane label="Order" name="order"><Order/></el-tab-pane>
                <el-tab-pane label="Withdraw" name="withdraw"><Withdraw/></el-tab-pane>
                <el-tab-pane label="Template" name="template"><Template/></el-tab-pane>

                <!-- <el-tab-pane label="Log" name="log"><Log/></el-tab-pane> -->
            </el-tabs>
        </el-col>
        <el-col :span="6"><div class="grid-content bg-purple-light"></div></el-col>
        <el-col :span="6"><div class="grid-content bg-purple"></div></el-col>
    </el-row>
</template>
<style>
.row-bg{
    flex-wrap: wrap;
    max-width: 1000px;
    margin: 0 auto;
}
</style>

<script>

import api from "@/api";


import Log  from './Log'
import User  from './User'
import Order  from './Order'
import Login  from './login'
import Home  from './Home'
import Withdraw  from './withdraw'
import Template  from './template'

export default {
  name: "IndexIndex",
  data() {
    return {
        activeName:"home",
        services:{}
    };
  },
  components:{Log,Home,Login,Order,User,Withdraw,Template},
  methods: {
    handleTabClick(tab, event) {
      //console.log(tab, event);
    },
    servicesGetList(){
        api.getService().then((res) => {
            this.services = res.data.services;
        });
    },

  },
  created(){
      this.servicesGetList()
  }
};
</script>