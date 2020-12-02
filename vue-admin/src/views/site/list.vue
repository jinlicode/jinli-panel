<template>
  <div class="app-container">
    <div class="filter-container">
      <el-button
        class="filter-item"
        style="margin-left: 10px;"
        type="primary"
        icon="el-icon-edit"
        @click="handleCreate"
      >
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
      <el-table-column
        label="ID"
        align="center"
        width="80"
      >
        <template slot-scope="{row}">
          <span>{{ row.id }}</span>
        </template>
      </el-table-column>
      <el-table-column
        label="网站域名"
        min-width="150"
      >
        <template slot-scope="{row}">
          {{ row.domain }}
        </template>
      </el-table-column>
      <el-table-column
        label="网站路径"
        min-width="200"
        align="left"
      >
        <template slot-scope="{row}">
          <span>/var/jinli_panel/code/{{ row.domain }}</span>
        </template>
      </el-table-column>
      <el-table-column
        label="状态"
        class-name="status-col"
        width="100"
      >
        <template slot-scope="{row}">
          <el-tag @click="handleUpdate(row, 'status')">
            <i
              v-if="row.status === 0"
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
        <template slot-scope="scope">
          <el-button
            size="small"
            type="primary"
            @click="handleUpdate(scope.row, 'conf')"
          >配置文件</el-button>
          <el-button
            size="small"
            type="primary"
            @click="handleUpdate(scope.row, 'domain')"
          >域名</el-button>
          <el-button
            size="small"
            type="primary"
            @click="handleUpdate(scope.row, 'rewrite')"
          >伪静态</el-button>
          <el-button
            size="small"
            type="primary"
            @click="handleUpdate(scope.row, 'php')"
          >PHP</el-button>
          <el-button
            size="small"
            type="primary"
            @click="handleUpdate(scope.row, 'basepath')"
          >根目录</el-button>
          <el-button
            size="small"
            type="danger"
            @click="delData(scope.row)"
          >删除</el-button>

        </template>
      </el-table-column>
    </el-table>
    <pagination
      v-show="total>0"
      :total="total"
      :page.sync="listQuery.page"
      :limit.sync="listQuery.limit"
      @pagination="getList"
    />

    <el-dialog
      :title="textMap[dialogStatus]"
      :visible.sync="dialogFormVisible"
    >
      <el-form
        ref="dataForm"
        :rules="rules"
        :model="temp"
        label-position="left"
        label-width="70px"
        style="width: 400px; margin-left:50px;"
      >
        <el-form-item
          label="PHP"
          prop="php_version"
        >
          <el-select
            v-model="temp.php_version"
            class="filter-item"
            placeholder="请选择"
          >
            <el-option
              v-for="item in phpVersionOptions"
              :key="item.name"
              :label="item.desc"
              :value="item.name"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          label="域名"
          prop="domain"
        >
          <el-input
            v-model="temp.domain"
            placeholder="请输入"
          />
        </el-form-item>
        <el-form-item
          label="协议"
          prop="is_ssl"
        >
          <el-radio-group v-model="temp.is_ssl">
            <el-radio :label="1">https</el-radio>
            <el-radio :label="0">http</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item
          v-if="temp.is_ssl == 1"
          label="邮箱"
          prop="email"
        >
          <el-input
            v-model="temp.email"
            placeholder="请输入"
          />
        </el-form-item>
      </el-form>
      <div
        slot="footer"
        class="dialog-footer"
      >
        <el-button @click="dialogFormVisible = false">
          取消
        </el-button>
        <el-button
          type="primary"
          @click="createData()"
        >
          保存
        </el-button>
      </div>
    </el-dialog>
    <el-dialog
      :title="textMap[dialogStatus]"
      :visible.sync="dialogOptVisible"
    >
      <el-input
        v-model="dataText"
        type="textarea"
        placeholder="请输入内容"
        :rows="10"
      />
      <div
        slot="footer"
        class="dialog-footer"
      >
        <el-button @click="dialogOptVisible = false">
          取消
        </el-button>
        <el-button
          type="primary"
          @click="updateData(handleId, dialogStatus)"
        >
          保存
        </el-button>
      </div>
    </el-dialog>
    <el-dialog
      :title="textMap[dialogStatus]"
      :visible.sync="dialogBasepathVisible"
      width="30%"
    >
      请选择运行根目录：
      <el-select
        v-model="basepath"
        placeholder="请选择"
      >
        <el-option
          v-for="item in basepathData"
          :key="item"
          :value="item"
        />
      </el-select>

      <div
        slot="footer"
        class="dialog-footer"
      >
        <el-button @click="dialogBasepathVisible = false">
          取消
        </el-button>
        <el-button
          type="primary"
          @click="updateData(handleId, dialogStatus)"
        >
          保存
        </el-button>
      </div>
    </el-dialog>

    <el-dialog
      :title="textMap[dialogStatus]"
      :visible.sync="dialogPhpVisible"
      width="30%"
    >
      请选择运行php版本：
      <el-select
        v-model="phpcur"
        placeholder="请选择"
      >
        <el-option
          v-for="item in phpVersionOptions"
          :key="item.name"
          :label="item.desc"
          :value="item.name"
        />
      </el-select>
      <div
        slot="footer"
        class="dialog-footer"
      >
        <el-button @click="dialogPhpVisible = false">
          取消
        </el-button>
        <el-button
          type="primary"
          @click="updateData(handleId, dialogStatus)"
        >
          保存
        </el-button>
      </div>
    </el-dialog>
    <el-dialog
      :title="textMap[dialogStatus]"
      :visible.sync="dialogDomainVisible"
      width="30%"
    >
      <el-input
        v-model="dataText"
        type="textarea"
        placeholder="请输入新域名，一行一个"
        :rows="10"
      />
      <div>
        <el-table
          :key="tableKey"
          v-model="domainData"
          :data="domainData"
          style="width: 100%;"
        >
          <el-table-column label="域名">
            <template slot-scope="sc">
              {{ sc.row.name }}
            </template>
          </el-table-column>
          <el-table-column
            align="right"
            label="操作"
          >
            <template slot-scope="sc">
              <a
                href="javascript:;"
                @click="delSiteDomainHandle(sc.row, sc.$index)"
              >删除</a>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <div
        slot="footer"
        class="dialog-footer"
      >
        <el-button @click="dialogDomainVisible = false">
          取消
        </el-button>
        <el-button
          type="primary"
          @click="updateData(handleId, dialogStatus)"
        >
          保存
        </el-button>
      </div>
    </el-dialog>

  </div>
</template>

<script>
import { fetchList, createSite, deleteSite, getSiteConf, updateSiteConf, getSiteRewrite, updateSiteRewrite, getSitePhp, updateSitePhp, getSiteDomain, updateSiteDomain, delSiteDomain, getSiteBasepath, updateSiteBasepath, updateSiteStatus } from '@/api/site'
import { getPHPList } from '@/api/soft'
import Pagination from '@/components/Pagination' // secondary package based on el-pagination
import { validDomain, validEmail } from '@/utils/validate'

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
      list: null,
      total: 0,
      listLoading: true,
      loading: null,
      listQuery: {
        page: 1,
        limit: 20,
        php_version: undefined,
        email: undefined,
        is_ssl: undefined,
        status: undefined,
        domain: undefined
      },
      phpVersionOptions: [],
      dataText: '',
      temp: {
        id: undefined,
        php_version: '',
        is_ssl: 1,
        domain: '',
        status: 0
      },
      phpcur: '',
      basepath: '/',
      handleId: null,
      basepathData: null,
      domainData: null,
      dialogFormVisible: false,
      dialogOptVisible: false,
      dialogBasepathVisible: false,
      dialogPhpVisible: false,
      dialogDomainVisible: false,
      dialogStatus: '',
      textMap: {
        update: '编辑网站',
        create: '新建网站',
        conf: '配置文件',
        domain: '域名编辑',
        rewrite: '伪静态',
        php: 'php版本',
        basepath: '根目录绑定'
      },
      dialogPvVisible: false,
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
      }
    }
  },
  created() {
    this.getPhpVersion()
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
        php_version: '',
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
            this.openFullLoader()
            this.dialogFormVisible = false
            this.$notify({
              title: 'Success',
              message: '创建成功',
              type: 'success',
              duration: 2000,
              onClose() {
                // 从新渲染
                _that.getList()
                _that.loading.close()
              }
            })
          })
        }
      })
    },
    handleUpdate(row, status) {
      // 加载动画
      const _that = this
      if (status !== 'status') {
        this.listLoading = true
      }
      this.handleId = row.id
      if (status === 'conf') {
        getSiteConf(row.id).then(response => {
          this.dialogOptVisible = true
          this.dataText = response.data.text
          this.dialogStatus = status
          this.listLoading = false
        })
      } else if (status === 'rewrite') {
        getSiteRewrite(row.id).then(response => {
          this.dialogOptVisible = true
          this.dataText = response.data.text
          this.dialogStatus = status
          this.listLoading = false
        })
      } else if (status === 'basepath') {
        getSiteBasepath(row.id).then(response => {
          this.dialogBasepathVisible = true
          this.basepath = response.data.basepath
          this.basepathData = response.data.list
          this.dialogStatus = status
          this.listLoading = false
        })
      } else if (status === 'php') {
        getSitePhp(row.id).then(response => {
          this.dialogPhpVisible = true
          this.phpcur = response.data.text
          this.dialogStatus = status
          this.listLoading = false
        })
      } else if (status === 'domain') {
        getSiteDomain(row.id).then(response => {
          this.dialogDomainVisible = true
          this.dataText = ''
          this.domainData = response.data.list
          this.dialogStatus = status
          this.listLoading = false
        })
      } else if (status === 'status') {
        if (row.status === 1) {
          this.$confirm('确定暂停网站吗?', '提示', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }).then(() => {
            _that.openFullLoader()
            updateSiteStatus({
              id: row.id
            }).then((result) => {
              _that.$notify({
                title: '提示',
                message: '暂停成功',
                type: 'success',
                duration: 2000,
                onClose() {
                  // 从新渲染
                  _that.getList()
                  _that.loading.close()
                }
              })
            })
          })
        } else if (row.status === 2) {
          this.$confirm('确定启用网站吗?', '提示', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }).then(() => {
            _that.openFullLoader()
            updateSiteStatus({
              id: row.id
            }).then((result) => {
              _that.$notify({
                title: '提示',
                message: '启用成功',
                type: 'success',
                duration: 2000,
                onClose() {
                  // 从新渲染
                  _that.getList()
                  _that.loading.close()
                }
              })
            })
          })
        }
      }
    },
    // 更新conf
    updateData(id, dialogStatus) {
      const _that = this
      _that.openFullLoader()
      if (dialogStatus === 'conf') {
        updateSiteConf({
          id: id,
          host_conf: _that.dataText
        }).then((result) => {
          _that.$notify({
            title: 'Success',
            message: '保存成功',
            type: 'success',
            duration: 2000,
            onClose() {
              _that.loading.close()
            }
          })
        })
      } else if (dialogStatus === 'rewrite') {
        updateSiteRewrite({
          id: id,
          rewrite_conf: _that.dataText
        }).then((result) => {
          _that.$notify({
            title: 'Success',
            message: '保存成功',
            type: 'success',
            duration: 2000,
            onClose() {
              _that.loading.close()
            }
          })
        })
      } else if (dialogStatus === 'php') {
        updateSitePhp({
          id: id,
          php: _that.phpcur
        }).then((result) => {
          _that.$notify({
            title: 'Success',
            message: '保存成功',
            type: 'success',
            duration: 2000,
            onClose() {
              _that.loading.close()
            }
          })
        })
      } else if (dialogStatus === 'basepath') {
        updateSiteBasepath({
          id: id,
          basepath: _that.basepath
        }).then((result) => {
          _that.$notify({
            title: 'Success',
            message: '保存成功',
            type: 'success',
            duration: 2000,
            onClose() {
              _that.loading.close()
            }
          })
        })
      } else if (dialogStatus === 'domain') {
        updateSiteDomain({
          id: id,
          text: _that.dataText
        }).then((result) => {
          _that.$notify({
            title: 'Success',
            message: '保存成功',
            type: 'success',
            duration: 2000,
            onClose() {
              _that.dataText = ''
              getSiteDomain(id).then(response => {
                _that.dataText = ''
                _that.domainData = response.data.list
              })
              _that.loading.close()
            }
          })
        })
      }
    },
    // 删除站
    delData(row) {
      const _that = this
      this.$confirm('此操作将永久删除此网站，包含数据库，程序代码', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        _that.openFullLoader()
        deleteSite({
          id: row.id
        }).then((result) => {
          _that.$notify({
            title: '提示',
            message: '删除成功',
            type: 'success',
            duration: 2000,
            onClose() {
              // 从新渲染
              _that.getList()
              _that.loading.close()
            }
          })
        })
      })
    },
    // 删除域名
    delSiteDomainHandle(row, index) {
      const _that = this
      this.$confirm('确定删除绑定域名吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        _that.openFullLoader()
        delSiteDomain({
          id: row.id
        }).then((result) => {
          _that.$notify({
            title: '提示',
            message: '删除成功',
            type: 'success',
            duration: 2000,
            onClose() {
              _that.domainData.splice(index, 1)
              _that.loading.close()
            }
          })
        })
      })
    },
    // 获取php版本
    getPhpVersion() {
      getPHPList().then(response => {
        this.phpVersionOptions = response.data.list
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
