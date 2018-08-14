<template>
    <div>
        <el-row>
            <el-col :span="24">
                <div v-for="(item,key,index) in services" :key="index" style="text-align:left">
                    <h3>{{ key }}</h3>
                    <ul class="ui list ">
                        <li v-for="(itemli,keyli,indexli) in item" class="item" :key="indexli">
                            服务状态：{{itemli.isrun? "服务已运行":"服务未运行"}} 
                            <a v-on:click="reboot(key + '|' + keyli+1)" class="reboot">重启服务</a><br>
                            登录状态：{{itemli.islogin? "服务已登录":"服务未登录"}} <br>
                            登录名称：{{itemli.name}} |{{itemli.runid}}<br>
                            运行ID：{{keyli+1}} 
                        </li>
                    </ul>
                    <div class="ui  divider"></div>
                </div>
            </el-col>
        </el-row>
    </div>
</template>
<style>
.reboot:hover{
    color:dodgerblue;
    cursor: pointer;
}
</style>

<script>
import api from "@/api";
export default {
  name: "IndexHome",
  props:[
      "services"
  ],
  data() {
    return {
      //services: {}
    };
  },
  methods: {
    reboot(id) {
      const loading = this.$loading({
        lock: true,
        text: "Loading",
        spinner: "el-icon-loading",
        background: "rgba(0, 0, 0, 0.7)"
      });
      api.reboot({ data: id }).then(() => {
          loading.close();
      });
    }
  },
  created() {

  },
};
</script>