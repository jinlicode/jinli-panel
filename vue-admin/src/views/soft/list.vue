<template>
  <div class="app-container">
    <el-table
      :key="tableKey"
      v-loading="listLoading"
      :data="list"
      border
      fit
      highlight-current-row
      style="width: 100%;"
    >
      <el-table-column
        label="容器名"
        min-width="200"
      >
        <template slot-scope="{row}">
          {{ row.name }}
        </template>
      </el-table-column>
      <el-table-column
        label="描述"
        min-width="480"
      >
        <template slot-scope="{row}">
          {{ row.desc }}
        </template>
      </el-table-column>
      <el-table-column
        label="状态"
        class-name="status-col"
        width="100"
      >
        <template slot-scope="{row}">
          <el-tag>
            <i
              v-if="row.status === 2"
              class="el-icon-loading"
            />
            {{ row.status | statusFilter }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column
        fixed="right"
        label="操作"
        min-width="460"
      >
        <template slot-scope="{row}">
          <el-button
            v-if="row.status === 0"
            size="small"
            type="primary"
            @click="handleInstall(row)"
          >安装</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
import { fetchList, InstallData } from '@/api/soft'

export default {
  name: 'SoftList',
  filters: {
    statusFilter(status) {
      const statusMap = {
        1: '已安装',
        2: '安装中',
        0: '未安装'
      }
      return statusMap[status]
    }
  },
  data() {
    return {
      tableKey: 0,
      list: null,
      listLoading: true,
      loading: null
    }
  },
  created() {
    this.getList()
  },
  mounted() {
    this.timer = window.setInterval(() => {
      setTimeout(() => {
        this.getList()
      }, 0)
    }, 3000)
  },
  destroyed() {
    window.clearInterval(this.timer)
  },
  methods: {
    getList() {
      this.listLoading = true
      fetchList().then(response => {
        this.list = response.data.list
        this.listLoading = false
      })
    },
    handleInstall(row) {
      const _that = this
      this.$confirm('确定安装此镜像吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        _that.openFullLoader()
        InstallData({
          name: row.name
        }).then(response => {
          _that.$notify({
            title: '提示',
            message: '已加入安装队列',
            type: 'success',
            duration: 2000,
            onClose() {
              _that.getList()
              _that.loading.close()
            }
          })
        })
      })
    },
    // 全屏遮罩层
    openFullLoader() {
      this.loading = this.$loading({
        lock: true,
        text: 'Loading',
        spinner: 'el-icon-loading',
        background: 'rgba(0, 0, 0, 0.7)'
      })
    }
  }
}
</script>
