<template>
  <div class="app-container">
    <div class="filter-container">
      <el-button class="filter-item" style="margin-left: 10px;" type="primary" icon="el-icon-edit" @click="handleCreate">
        新增
      </el-button>
    </div>

    <el-table
      :key="tableKey"
      v-loading="listLoading"
      :data="list"
      border
      fit
      highlight-current-row
      style="width: 100%;"
    >
      <el-table-column label="ID" align="center" width="80">
        <template slot-scope="{row}">
          <span>{{ row.id }}</span>
        </template>
      </el-table-column>
      <el-table-column label="网站域名" min-width="150px">
        <template slot-scope="{row}">
          <span class="link-type" @click="handleUpdate(row)">{{ row.domain }}</span>
        </template>
      </el-table-column>
      <el-table-column label="网站路径" min-width="150px"> align="center">
        <template slot-scope="{row}">
          <span>{{ row.domain }}</span>
        </template>
      </el-table-column>
      <el-table-column label="状态" class-name="status-col" width="100">
        <template slot-scope="{row}">
          <el-tag>
            {{ row.status | statusFilter }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" width="230" class-name="small-padding fixed-width">
        <template slot-scope="{row,$index}">
          <el-button type="primary" size="mini" @click="handleUpdate(row)">
            编辑
          </el-button>
          <el-button v-if="row.status!='deleted'" size="mini" type="danger" @click="handleDelete(row,$index)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.limit" @pagination="getList" />

    <el-dialog :title="textMap[dialogStatus]" :visible.sync="dialogFormVisible">
      <el-form ref="dataForm" :rules="rules" :model="temp" label-position="left" label-width="70px" style="width: 400px; margin-left:50px;">
        <el-form-item label="PHP" prop="php_version">
          <el-select v-model="temp.php_version" class="filter-item" placeholder="请选择">
            <el-option v-for="item in phpVersionOptions" :key="item.key" :label="item.display_name" :value="item.key" />
          </el-select>
        </el-form-item>
        <el-form-item label="域名" prop="domain">
          <el-input v-model="temp.domain" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="协议" prop="is_ssl">
          <el-radio-group v-model="temp.is_ssl">
            <el-radio :label="1">https</el-radio>
            <el-radio :label="0">http</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="temp.is_ssl == 1" label="邮箱" prop="email">
          <el-input v-model="temp.email" placeholder="请输入" />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">
          取消
        </el-button>
        <el-button type="primary" @click="dialogStatus==='create'?createData():updateData()">
          保存
        </el-button>
      </div>
    </el-dialog>

    <el-dialog :title="'删除 ['+temp.domain+']'" :visible.sync="dialogDelVisible" width="30%">
      <el-form ref="delForm" :model="temp" label-position="left" label-width="70px">
        <template>
          <el-checkbox v-model="delDatabase" lable="1">删除对应数据库</el-checkbox>
          <el-checkbox v-model="delCode" lable="1">删除对应程序代码</el-checkbox>
        </template>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogDelVisible = false">
          取消
        </el-button>
        <el-button type="primary" @click="delData()">
          确定删除
        </el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { fetchList, createSite, updateSite, deleteSite } from '@/api/site'
import Pagination from '@/components/Pagination' // secondary package based on el-pagination
import { validDomain, validEmail } from '@/utils/validate'

const phpVersionOptions = [
  { key: '5.6', display_name: 'php-5.6' },
  { key: '5.6-sec', display_name: 'php-5.6 安全版本' },
  { key: '7.0', display_name: 'php-7.0' },
  { key: '7.0-sec', display_name: 'php-7.0 安全版本' },
  { key: '7.1', display_name: 'php-7.1' },
  { key: '7.1-sec', display_name: 'php-7.1 安全版本' },
  { key: '7.2', display_name: 'php-7.2' },
  { key: '7.2-sec', display_name: 'php-7.2 安全版本' },
  { key: '7.3', display_name: 'php-7.3' },
  { key: '7.3-sec', display_name: 'php-7.3 安全版本' }
]

export default {
  name: 'SiteList',
  components: { Pagination },
  filters: {
    statusFilter(status) {
      const statusMap = {
        1: '运行中',
        2: '已暂停',
        0: '启动中'
      }
      return statusMap[status]
    }
  },
  data() {
    const validateDomain = (rule, value, callback) => {
      if (!validDomain(value)) {
        callback(new Error('域名格式错误'))
      } else {
        callback()
      }
    }
    const validateEmail = (rule, value, callback) => {
      if (!validEmail(value)) {
        callback(new Error('邮箱格式错误'))
      } else {
        callback()
      }
    }
    return {
      tableKey: 0,
      delDatabase: 0,
      delCode: 0,
      list: null,
      total: 0,
      listLoading: true,
      listQuery: {
        page: 1,
        limit: 20,
        php_version: undefined,
        email: undefined,
        is_ssl: undefined,
        status: undefined,
        domain: undefined
      },
      phpVersionOptions,
      temp: {
        id: undefined,
        php_version: '5.6',
        is_ssl: 1,
        domain: '',
        status: 0
      },
      dialogFormVisible: false,
      dialogDelVisible: false,
      dialogStatus: '',
      textMap: {
        update: '编辑网站',
        create: '新建网站'
      },
      dialogPvVisible: false,
      pvData: [],
      rules: {
        php_version: [
          { required: true, message: 'php版本 必选', trigger: 'change' }
        ],
        domain: [
          { required: true, message: '域名必填', trigger: 'blur' },
          { trigger: 'blur', validator: validateDomain }
        ],
        is_ssl: [
          { required: true, message: '请选择网站协议', trigger: 'change' }
        ],
        email: [
          { required: true, trigger: 'blur', validator: validateEmail }
        ]
      },
      downloadLoading: false
    }
  },
  created() {
    this.getList()
  },
  methods: {
    getList() {
      this.listLoading = true
      fetchList(this.listQuery).then(response => {
        this.list = response.data.list
        this.total = response.data.total
        this.listLoading = false
      })
    },
    handleFilter() {
      this.listQuery.page = 1
      this.getList()
    },
    handleModifyStatus(row, status) {
      this.$message({
        message: '操作成功',
        type: 'success'
      })
      row.status = status
    },
    resetTemp() {
      this.temp = {
        id: undefined,
        php_version: '5.6',
        is_ssl: 0,
        status: 0,
        domain: ''
      }
    },
    handleCreate() {
      this.resetTemp()
      this.dialogStatus = 'create'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    createData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          this.temp.id = 0
          const _that = this
          createSite(this.temp).then(() => {
            this.dialogFormVisible = false
            this.$notify({
              title: 'Success',
              message: '创建成功',
              type: 'success',
              duration: 2000,
              onClose: function() {
                // 从新渲染
                _that.getList()
              }
            })
          })
        }
      })
    },
    handleUpdate(row) {
      this.temp = Object.assign({}, row) // copy obj
      this.dialogStatus = 'update'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    updateData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          const tempData = Object.assign({}, this.temp)
          updateSite(tempData).then(() => {
            const index = this.list.findIndex(v => v.id === this.temp.id)
            this.list.splice(index, 1, this.temp)
            this.dialogFormVisible = false
            this.$notify({
              title: 'Success',
              message: '保存成功',
              type: 'success',
              duration: 2000
            })
          })
        }
      })
    },
    handleDelete(row, index) {
      this.temp = Object.assign({}, row) // copy obj
      this.dialogDelVisible = true
      this.delDatabase = 0
      this.delCode = 0
    },
    delData() {
      this.temp.delDatabase = this.delDatabase
      this.temp.delCode = this.delCode

      deleteSite(this.temp).then(() => {
        this.dialogDelVisible = false
        const _that = this
        this.$notify({
          title: '提示',
          message: '删除成功',
          type: 'success',
          duration: 2000,
          onClose: function() {
            // 从新渲染
            _that.getList()
          }
        })
      })
    }
  }
}
</script>
