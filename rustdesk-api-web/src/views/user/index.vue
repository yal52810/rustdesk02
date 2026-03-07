<template>
  <div>
    <el-card class="list-query" shadow="hover">
      <el-form inline label-width="80px">
        <el-form-item :label="T('Username')">
          <el-input v-model="listQuery.username"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handlerQuery">{{ T('Filter') }}</el-button>
          <el-button type="danger" @click="toAdd">{{ T('Add') }}</el-button>
          <el-button type="warning" @click="openBatchDialog">批量创建</el-button>
          <el-button type="success" @click="toExport">{{ T('Export') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>
    <el-card class="list-body" shadow="hover">
      <el-table :data="listRes.list" v-loading="listRes.loading" border>
        <el-table-column prop="id" label="ID" align="center"></el-table-column>
        <el-table-column prop="username" :label="T('Username')" align="center"/>
        <el-table-column prop="email" :label="T('Email')" align="center"/>
        <el-table-column prop="nickname" :label="T('Nickname')" align="center"/>
        <el-table-column :label="T('Group')" align="center">
          <template #default="{row}">
            <span v-if="row.group_id"> <el-tag>{{ listRes.groups?.find(g => g.id === row.group_id)?.name }} </el-tag> </span>
            <span v-else> - </span>
          </template>
        </el-table-column>
        <el-table-column :label="T('ActivationDate')" align="center">
          <template #default="{row}">
            <span v-if="row.first_login_at">{{ row.first_login_at }}</span>
            <el-tag v-else type="info">{{ T('NotActivated') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="T('ExpirationInfo')" align="center">
          <template #default="{row}">
            <el-tag v-if="row.valid_days === -1" type="success">{{ T('Lifetime') }}</el-tag>
            <span v-else-if="!row.first_login_at">{{ T('NotActivated') }}</span>
            <el-tag v-else-if="isExpired(row)" type="danger">{{ T('Expired') }}</el-tag>
            <span v-else>{{ daysRemaining(row) }} {{ T('DaysRemaining') }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="remark" :label="T('Remark')" align="center"/>
        <el-table-column prop="created_at" :label="T('CreatedAt')" align="center"/>
        <el-table-column prop="updated_at" :label="T('UpdatedAt')" align="center"/>
        <el-table-column :label="T('Actions')" align="center" width="850">
          <template #default="{row}">
            <el-button @click="toTag(row)">{{ T('UserTags') }}</el-button>
            <el-button @click="toAddressBook(row)">{{ T('UserAddressBook') }}</el-button>
            <el-button @click="toEdit(row)">{{ T('Edit') }}</el-button>
            <el-button type="primary" @click="quickEditValidDays(row)">{{ T('EditValidDays') }}</el-button>
            <el-button
              :type="row.status === ENABLE_STATUS ? 'warning' : 'success'"
              @click="toggleBan(row)">
              {{ row.status === ENABLE_STATUS ? T('Ban') : T('Unban') }}
            </el-button>
            <el-button type="warning" @click="changePass(row)">{{ T('ResetPassword') }}</el-button>
            <el-button type="danger" @click="remove(row)">{{ T('Delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    <el-card class="list-page" shadow="hover">
      <el-pagination background
                     layout="prev, pager, next, sizes, jumper"
                     :page-sizes="[10,20,50,100]"
                     v-model:page-size="listQuery.page_size"
                     v-model:current-page="listQuery.page"
                     :total="listRes.total">
      </el-pagination>
    </el-card>

    <!-- Quick Edit Valid Days Dialog -->
    <el-dialog v-model="validDaysDialog.visible" :title="T('EditValidDays')" width="500px">
      <el-form label-width="120px">
        <el-form-item :label="T('Username')">
          <el-input v-model="validDaysDialog.username" disabled></el-input>
        </el-form-item>
        <el-form-item :label="T('CurrentValidDays')">
          <el-tag v-if="validDaysDialog.currentValidDays === -1" type="success">{{ T('Lifetime') }}</el-tag>
          <el-tag v-else>{{ validDaysDialog.currentValidDays }} {{ T('Days') }}</el-tag>
        </el-form-item>
        <el-form-item :label="T('NewValidDays')">
          <el-input-number
            v-model="validDaysDialog.newValidDays"
            :min="-1"
            :placeholder="T('ValidDaysHint')"
          ></el-input-number>
          <el-text class="mx-1" type="info">{{ T('ValidDaysDesc') }}</el-text>
        </el-form-item>
        <el-form-item :label="T('DeviceLimit')">
          <el-input-number
            v-model="validDaysDialog.deviceLimit"
            :min="0"
          ></el-input-number>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="validDaysDialog.visible = false">{{ T('Cancel') }}</el-button>
        <el-button type="primary" @click="saveValidDays">{{ T('Confirm') }}</el-button>
      </template>
    </el-dialog>

    <!-- Batch Create Dialog -->
    <el-dialog v-model="batchDialog.visible" title="批量创建用户" width="900px">
      <el-form :model="batchDialog.form" label-width="120px" :rules="batchDialog.rules" ref="batchFormRef">
        <el-form-item label="用户名前缀" prop="username_prefix">
          <el-input v-model="batchDialog.form.username_prefix" placeholder="例如: user" />
          <div class="form-tip">生成的用户名格式: user01, user02, ...</div>
        </el-form-item>
        <el-form-item label="创建数量" prop="count">
          <el-input-number v-model="batchDialog.form.count" :min="1" :max="100" />
          <div class="form-tip">最多一次创建100个用户</div>
        </el-form-item>
        <el-form-item label="密码位数" prop="password_length">
          <el-input-number v-model="batchDialog.form.password_length" :min="6" :max="32" />
          <div class="form-tip">生成随机数字密码的位数</div>
        </el-form-item>
        <el-form-item label="用户组" prop="group_id">
          <el-select v-model="batchDialog.form.group_id" placeholder="请选择用户组">
            <el-option
              v-for="group in listRes.groups"
              :key="group.id"
              :label="group.name"
              :value="group.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="账号状态" prop="status">
          <el-radio-group v-model="batchDialog.form.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="2">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="有效天数" prop="valid_days">
          <el-input-number v-model="batchDialog.form.valid_days" :min="-1" />
          <div class="form-tip">-1 表示永久有效</div>
        </el-form-item>
        <el-form-item label="设备限制" prop="device_limit">
          <el-input-number v-model="batchDialog.form.device_limit" :min="1" :max="100" />
          <div class="form-tip">允许同时在线的设备数量</div>
        </el-form-item>
      </el-form>

      <div v-if="batchDialog.result" style="margin-top: 20px;">
        <el-alert
          :title="`成功创建 ${batchDialog.result.success_count} 个用户${batchDialog.result.error_count > 0 ? '，失败 ' + batchDialog.result.error_count + ' 个' : ''}`"
          :type="batchDialog.result.error_count > 0 ? 'warning' : 'success'"
          :closable="false"
          style="margin-bottom: 20px;" />

        <div style="margin-bottom: 10px;">
          <el-button type="success" @click="copyBatchText" size="small">复制文本</el-button>
          <el-button type="primary" @click="downloadBatchCSV" size="small">下载CSV</el-button>
          <el-button type="warning" @click="downloadBatchTXT" size="small">下载TXT</el-button>
        </div>

        <el-table :data="batchDialog.result.users" border max-height="300">
          <el-table-column type="index" label="序号" width="60" align="center" />
          <el-table-column prop="username" label="账号" align="center">
            <template #default="{row}">
              <el-input v-model="row.username" readonly>
                <template #append>
                  <el-button @click="copyText(row.username)">复制</el-button>
                </template>
              </el-input>
            </template>
          </el-table-column>
          <el-table-column prop="password" label="密码" align="center">
            <template #default="{row}">
              <el-input v-model="row.password" readonly>
                <template #append>
                  <el-button @click="copyText(row.password)">复制</el-button>
                </template>
              </el-input>
            </template>
          </el-table-column>
        </el-table>

        <div v-if="batchDialog.result.errors && batchDialog.result.errors.length > 0" style="margin-top: 20px;">
          <el-alert title="错误信息" type="error" :closable="false">
            <div v-for="(error, index) in batchDialog.result.errors" :key="index">
              {{ error }}
            </div>
          </el-alert>
        </div>
      </div>

      <template #footer>
        <el-button @click="closeBatchDialog">关闭</el-button>
        <el-button type="primary" @click="handleBatchCreate" :loading="batchDialog.loading">创建用户</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
  import { useRepositories, useDel, useToEditOrAdd, useChangePwd } from '@/views/user/composables'
  import { T } from '@/utils/i18n'
  import { DISABLE_STATUS, ENABLE_STATUS } from '@/utils/common_options'
  import { update, quickBatchCreate } from '@/api/user'
  import { ElMessageBox, ElMessage } from 'element-plus'
  import { onMounted, watch, reactive, ref } from 'vue'
  //列表
  const {
    listRes,
    listQuery,
    handlerQuery,
    getList,
    getGroups,
    toExport,
  } = useRepositories()

  onMounted(getGroups)

  onMounted(getList)

  watch(() => listQuery.page, getList)
  watch(() => listQuery.page_size, handlerQuery)

  const { toEdit, toAdd, toBatchCreate, toAddressBook, toTag } = useToEditOrAdd()

  const { changePass } = useChangePwd()

  //删除
  const { del } = useDel()
  const remove = async (row) => {
    const res = await del(row.id)
    if (res) {
      getList(listQuery)
    }
  }

  // Quick Edit Valid Days Dialog
  const validDaysDialog = reactive({
    visible: false,
    userId: null,
    username: '',
    currentValidDays: 0,
    newValidDays: 0,
    deviceLimit: 10,
    userData: null
  })

  const quickEditValidDays = (row) => {
    validDaysDialog.visible = true
    validDaysDialog.userId = row.id
    validDaysDialog.username = row.username
    validDaysDialog.currentValidDays = row.valid_days
    validDaysDialog.newValidDays = row.valid_days
    validDaysDialog.deviceLimit = row.device_limit || 10
    validDaysDialog.userData = { ...row }
  }

  const saveValidDays = async () => {
    const res = await update({
      ...validDaysDialog.userData,
      valid_days: validDaysDialog.newValidDays,
      device_limit: validDaysDialog.deviceLimit
    }).catch(_ => false)

    if (res) {
      ElMessage.success(T('OperationSuccess'))
      validDaysDialog.visible = false
      getList(listQuery)
    }
  }

  // Utility functions for expiration
  const isExpired = (row) => {
    if (row.valid_days === -1) return false // Lifetime
    if (!row.first_login_at) return false // Not activated yet
    const activationDate = new Date(row.first_login_at)
    const expirationDate = new Date(activationDate)
    expirationDate.setDate(expirationDate.getDate() + row.valid_days)
    return new Date() > expirationDate
  }

  const daysRemaining = (row) => {
    if (!row.first_login_at) return 'N/A'
    const activationDate = new Date(row.first_login_at)
    const expirationDate = new Date(activationDate)
    expirationDate.setDate(expirationDate.getDate() + row.valid_days)
    const today = new Date()
    const diffTime = expirationDate - today
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
    return diffDays > 0 ? diffDays : 0
  }

  // Ban/Unban toggle function
  const toggleBan = async (row) => {
    const newStatus = row.status === ENABLE_STATUS ? DISABLE_STATUS : ENABLE_STATUS
    const action = newStatus === ENABLE_STATUS ? 'Unban' : 'Ban'

    const confirm = await ElMessageBox.confirm(
      T('Confirm?', { param: T(action) }),
      {
        confirmButtonText: T('Confirm'),
        cancelButtonText: T('Cancel'),
      }
    ).catch(_ => false)

    if (!confirm) return false

    const res = await update({ ...row, status: newStatus }).catch(_ => false)
    if (res) {
      ElMessage.success(T('OperationSuccess'))
      getList(listQuery)
    }
  }

  // Batch Create Dialog
  const batchFormRef = ref(null)
  const batchDialog = reactive({
    visible: false,
    loading: false,
    result: null,
    form: {
      username_prefix: 'user',
      count: 10,
      password_length: 8,
      group_id: null,
      status: 1,
      valid_days: 365,
      device_limit: 10
    },
    rules: {
      username_prefix: [
        { required: true, message: '请输入用户名前缀', trigger: 'blur' },
        { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
      ],
      count: [{ required: true, message: '请输入创建数量', trigger: 'blur' }],
      password_length: [{ required: true, message: '请输入密码位数', trigger: 'blur' }],
      group_id: [{ required: true, message: '请选择用户组', trigger: 'change' }]
    }
  })

  const openBatchDialog = () => {
    batchDialog.visible = true
    batchDialog.result = null
    if (listRes.groups && listRes.groups.length > 0) {
      batchDialog.form.group_id = listRes.groups[0].id
    }
  }

  const closeBatchDialog = () => {
    batchDialog.visible = false
    batchDialog.result = null
    getList(listQuery)
  }

  const handleBatchCreate = async () => {
    await batchFormRef.value.validate()
    batchDialog.loading = true
    try {
      const res = await quickBatchCreate(batchDialog.form)
      batchDialog.result = res.data
      ElMessage.success(res.data.message || '创建成功')
    } catch (error) {
      ElMessage.error(error.message || '创建失败')
    } finally {
      batchDialog.loading = false
    }
  }

  const copyText = (text) => {
    if (!text) {
      ElMessage.warning('内容为空，无法复制')
      return
    }

    if (navigator.clipboard && navigator.clipboard.writeText) {
      navigator.clipboard.writeText(text).then(() => {
        ElMessage.success('已复制到剪贴板')
      }).catch((err) => {
        console.error('复制失败:', err)
        fallbackCopyText(text)
      })
    } else {
      fallbackCopyText(text)
    }
  }

  const fallbackCopyText = (text) => {
    const textArea = document.createElement('textarea')
    textArea.value = text
    textArea.style.position = 'fixed'
    textArea.style.left = '-999999px'
    textArea.style.top = '-999999px'
    document.body.appendChild(textArea)
    textArea.focus()
    textArea.select()

    try {
      const successful = document.execCommand('copy')
      if (successful) {
        ElMessage.success('已复制到剪贴板')
      } else {
        ElMessage.error('复制失败，请手动复制')
      }
    } catch (err) {
      console.error('复制失败:', err)
      ElMessage.error('复制失败，请手动复制')
    }

    document.body.removeChild(textArea)
  }

  const copyBatchText = () => {
    if (!batchDialog.result || !batchDialog.result.export_text) return
    copyText(batchDialog.result.export_text)
  }

  const downloadBatchCSV = () => {
    if (!batchDialog.result || !batchDialog.result.export_csv) return
    downloadFile(batchDialog.result.export_csv, 'users.csv', 'text/csv')
  }

  const downloadBatchTXT = () => {
    if (!batchDialog.result || !batchDialog.result.export_text) return
    downloadFile(batchDialog.result.export_text, 'users.txt', 'text/plain')
  }

  const downloadFile = (content, filename, type) => {
    const blob = new Blob([content], { type: `${type};charset=utf-8;` })
    const link = document.createElement('a')
    link.href = URL.createObjectURL(blob)
    link.download = filename
    link.click()
    URL.revokeObjectURL(link.href)
    ElMessage.success('下载成功')
  }

</script>

<style scoped>
.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
}
</style>
