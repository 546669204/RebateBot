<template>
    <div>
        <el-button icon="el-icon-search" circle @click="getList"></el-button>
        <el-table :data="tableData" style="width: 100%">
            <el-table-column prop="id" label="订单ID" width="180">
            </el-table-column>
            <el-table-column prop="orderid" label="订单号" width="180">
            </el-table-column>
            <el-table-column prop="create_time" label="创建时间" width="180">
              <template slot-scope="scope">
                <span>{{ scope.row.create_time*1000 | formatDate }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="weid" label="微信昵称" width="180">
            </el-table-column>
            <el-table-column prop="goodsname" label="商品标题" width="180">
            </el-table-column>
            <el-table-column prop="goodsid" label="商品ID" width="180">
            </el-table-column>
            <el-table-column prop="paymoney" label="付款金额" width="180">
            </el-table-column>
            <el-table-column prop="pubmoney" label="联盟佣金" width="180">
            </el-table-column>
            <el-table-column prop="pubbili" label="联盟佣金比例" width="180">
            </el-table-column>
            <el-table-column prop="buymoney" label="买家佣金" width="180">
            </el-table-column>
            <el-table-column prop="buybili" label="买家佣金比例" width="180">
            </el-table-column>
            <el-table-column prop="income" label="预计收入" width="180">
            </el-table-column>
            <el-table-column prop="status" label="状态" width="180">
            </el-table-column>
            <el-table-column prop="pid" label="PID" width="180">
            </el-table-column>
        </el-table>
        <el-pagination background layout="total,sizes, prev, pager, next" :page-sizes="[5,15,30,50]" :total="pageCount" @current-change="pageCurrentChange" @size-change="pageSizeChange">
        </el-pagination>
    </div>
</template>
<style>
</style>

<script>
import api from '@/api'
import _formatDate from '@/common/date'
export default {
  name: "IndexOrder",
  data() {
    return {
      tableData: [],
      pageCount: 0,
      pageSize: 5,
      pageCurrent:0
    };
  },
  methods: {
    getList(){
      api.getOrderData({"page":this.pageCurrent,"pageSize":this.pageSize}).then(res => {
       this.pageCount = res.data.count;
        this.tableData = res.data.data;
      })
    },
    pageCurrentChange(e) {
      this.pageCurrent = e;
      this.getList();
    },
    pageSizeChange(e) {
      this.pageSize = e;
      this.getList();
    }
  },
  filters:{
    formatDate(d){
      return _formatDate(d,"yyyy-MM-dd h:m:s")
    }
  },
  created() {
    this.getList()
  },
};
</script>