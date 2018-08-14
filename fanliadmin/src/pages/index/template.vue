<template>
    <div>
        <el-select v-model="parm.type" placeholder="请选择" @change="handleSelectChange">
            <el-option
            v-for="item in option"
            :key="item.value"
            :label="item.label"
            :value="item.value">
            </el-option>
        </el-select>
        <el-button icon="el-icon-search" circle @click="getList"></el-button>
        <el-button type="primary" size="mini" @click="handleEdit(0)">新增模版</el-button>
        <el-table :data="tableData" style="width: 100%">
            <!-- <el-table-column prop="type" label="模版类型" width="180">
            </el-table-column> -->
  
            <el-table-column prop="preview" label="模版预览" width="300">
                <template slot-scope="scope">
                    <div v-html="scope.row.preview" style="text-align:left"></div>
                    
                </template>
            </el-table-column>
      

            <el-table-column label="操作">
            <template slot-scope="scope">
                <el-button type="primary" size="mini" @click="handleEdit(scope.$index+1)">修改</el-button>
                <el-button type="primary" size="mini" @click="handleDelete(scope.$index+1)">删除</el-button>
            </template>
            </el-table-column>


        </el-table>

<el-dialog
  title="提示"
  :visible.sync="centerDialogVisible"
  width="30%"
  center>

<el-form  label-width="80px">
<el-form-item label="选择类型">
           <el-select v-model="temp.type" placeholder="请选择" @change="handleSelectChange">
            <el-option
            v-for="item in option"
            :key="item.value"
            :label="item.label"
            :value="item.value">
            </el-option>
        </el-select>
  </el-form-item>
  <el-form-item label="可用标签" v-if="rawData[temp.type]">
      <el-tag v-for="(item,index) in rawData[temp.type].props" :key="index" @click.native="handleInsert(index)">{{item}}</el-tag>
 </el-form-item>
  <el-form-item label="原代码">
    <el-input type="textarea" autosize v-model="temp.text"></el-input>
  </el-form-item>
  <el-form-item label="预览" v-if="rawData[temp.type]">
    <el-input type="textarea" autosize :value="preview(temp.text,rawData[temp.type].props)"></el-input>
  </el-form-item>
</el-form>


  <span slot="footer" class="dialog-footer">
    <el-button @click="centerDialogVisible = false">取 消</el-button>
    <el-button type="primary" @click="handleDialogOk">确 定</el-button>
  </span>
</el-dialog>

    </div>
</template>
<style>
</style>

<script>
import api from "@/api";
import _formatDate from "@/common/date";
export default {
  name: "Withdraw",
  data() {
    return {
      rawData:{},
      tableData: [],
      pageCount: 0,
      pageSize: 5,
      pageCurrent: 0,
      parm:{
          type:"check"
      },
      temp:{
          type:"check",
          text:"",
          preview:"",
          id:0
      },
      centerDialogVisible:false,
      option:[
        {label:"签到",value:"check"},
        {label:"帮助",value:"help"},
        {label:"返利信息",value:"fanli"},
        {label:"支付成功",value:"payok"},
        {label:"绑定订单提醒",value:"bindorder"},
        {label:"绑定订单成功",value:"bindorderok"},
        {label:"自动绑定订单成功",value:"autobindorder"},
        {label:"收货后提醒",value:"shouhuo"},
        {label:"查询余额",value:"money"},
        {label:"提现失败",value:"withdrawoff"},
        {label:"提现成功",value:"withdrawok"},
      ]
    };
  },
  methods: {
    getList() {
      api.getTemplateData().then(res => {
            this.rawData = res.data.data;
            this.handleSelectChange();
        });
    },
    saveData(){
        api.setTemplateData({data:JSON.stringify(this.rawData)}).then(res => {
            if (res.data.code == 0){
                this.$message({
                    message: '保存成功',
                    type: 'success'
                });
            }else{
                this.$message({
                    message: res.data.msg,
                    type: 'error'
                });
            }
   
        });
    },
    handleSelectChange(e) {
          var temp = this.rawData[this.parm.type];
          var list = Object.assign([]);
          temp.text.forEach((v)=>{
              list.push({
                  type:this.parm.type,
                  text:v,
                  preview:this.preview(v,temp.props,true)
              })
          })
          this.tableData = list;
    },
    preview(t,p,no){
        Object.keys(p).forEach((i)=>{
            t = t.replace(new RegExp(`\{\{\.${i}\}\}`,"gi"),"{"+p[i]+"}");
        })
        no && (t = t.replace(/\n/gi,"<br/>"))
        return t
    },
    handleEdit(id){
        if (id === 0) {
            this.temp = Object.assign({},{
                type:this.parm.type,
                text:"",
                preview:"",
                id:0
            })
        }else{
            this.temp = Object.assign({},{
                type:this.parm.type,
                text:this.rawData[this.parm.type].text[id-1],
                preview:"",
                id:id
            })
        }
        this.centerDialogVisible = true;
    },
    handleInsert(item){
        this.temp.text += `{{.${item}}}`;
    },
    handleDelete(id){
        this.rawData[this.parm.type].text.splice(id-1,1);
        this.handleSelectChange();
        this.saveData();
    },
    handleDialogOk(){
        this.centerDialogVisible = false;
        if(this.temp.id == 0){
            this.rawData[this.parm.type].text.push(this.temp.text);
        }else{
            this.rawData[this.parm.type].text[this.temp.id-1] = (this.temp.text);
        }
        this.handleSelectChange();
        this.saveData();
    }
  },
  filters: {
    formatDate(d) {
      return _formatDate(d, "yyyy-MM-dd h:m:s");
    }
  },
  created() {
    this.getList();
  }
};
</script>