<template>
    <div>
      <el-button icon="el-icon-search" circle @click="getList"></el-button>
        <el-table :data="tableData" style="width: 100%">
            <el-table-column prop="id" label="ID" width="180">
            </el-table-column>
            <el-table-column prop="weid" label="微信ID" width="180">
            </el-table-column>
            <el-table-column prop="wename" label="微信昵称" width="180">
            </el-table-column>
            <el-table-column prop="money" label="余额" width="180">
            </el-table-column>
            <el-table-column prop="update_time" label="活跃时间" width="180">
              <template slot-scope="scope">
                <span>{{ scope.row.update_time*1000 | formatDate }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="create_time" label="创建时间" width="180">
              <template slot-scope="scope">
                <span>{{ scope.row.create_time*1000 | formatDate }}</span>
              </template>
            </el-table-column>
        </el-table>
        <el-pagination background layout="total,sizes, prev, pager, next" :page-sizes="[5,15,30,50]" :total="pageCount" @current-change="pageCurrentChange" @size-change="pageSizeChange">
        </el-pagination>
    </div>
</template>
<style>
</style>

<script>
import api from "@/api";
import _formatDate from '@/common/date'
export default {
  name: "IndexUser",
  data() {
    return {
      tableData: [],
      pageCount: 0,
      pageSize: 5,
      pageCurrent: 0
    };
  },
  methods: {
    getList() {
      api.getUserData({"page":this.pageCurrent,"pageSize":this.pageSize}).then(res => {
        this.pageCount = res.data.count;
        this.tableData = res.data.data;
      });
    },
    pageCurrentChange(e) {
      this.pageCurrent = e;
      this.getList();
    },
    pageSizeChange(e) {
      this.pageSize = e;
      this.getList();
    },

  },
  filters:{
    formatDate(d){
      return _formatDate(d,"yyyy-MM-dd h:m:s")
    }
  },
  mounted() {
    this.getList();
  },
  created() {
    //this.getList();
  }
};
</script>