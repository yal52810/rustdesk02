<template>
  <div>
    <el-alert
      v-if="groups.length === 0"
      title="提示：未找到用户组"
      type="warning"
      :closable="false"
      style="margin-bottom: 20px;"
    >
      <p>批量创建用户需要先创建用户组。请前往"用户组管理"创建一个用户组后再使用此功能。</p>
      <p>系统首次初始化时会自动创建"默认组"，如果没有，请手动创建。</p>
    </el-alert>

    <el-card shadow="hover">
      <template #header>
        <span>批量创建用户</span>
      </template>

      <el-form :model="form" label-width="120px" :rules="rules" ref="formRef">
        <el-form-item label="用户名前缀" prop="username_prefix">
          <el-input v-model="form.username_prefix" placeholder="例如: user" />
          <div class="form-tip">生成的用户名格式: user01, user02, ...</div>
        </el-form-item>

        <el-form-item label="创建数量" prop="count">
          <el-input-number v-model="form.count" :min="1" :max="100" />
          <div class="form-tip">最多一次创建100个用户</div>
        </el-form-item>

        <el-form-item label="密码位数" prop="password_length">
          <el-input-number v-model="form.password_length" :min="6" :max="32" />
          <div class="form-tip">生成随机数字密码的位数</div>
        </el-form-item>

        <el-form-item label="用户组" prop="group_id">
          <el-select v-model="form.group_id" placeholder="请选择用户组">
            <el-option
              v-for="group in groups"
              :key="group.id"
              :label="group.name"
              :value="group.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="套餐选择">
          <el-select v-model="form.package_id" placeholder="请选择套餐（可选）" clearable>
            <el-option
              v-for="pkg in packages"
              :key="pkg.id"
              :label="pkg.name"
              :value="pkg.id"
            />
          </el-select>
          <div class="form-tip">选择套餐后将自动填充有效天数和设备限制，也可手动修改</div>
        </el-form-item>

        <el-form-item label="账号状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="2">禁用</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="有效天数" prop="valid_days">
          <el-input-number v-model="form.valid_days" :min="-1" />
          <div class="form-tip">-1 表示永久有效</div>
        </el-form-item>

        <el-form-item label="设备限制" prop="device_limit">
          <el-input-number v-model="form.device_limit" :min="1" :max="100" />
          <div class="form-tip">允许同时在线的设备数量</div>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleCreate" :loading="loading">
            创建用户
          </el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card shadow="hover" v-if="result" style="margin-top: 20px;">
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <span>创建结果</span>
          <div>
            <el-button type="success" @click="handleCopyText" size="small">
              复制文本
            </el-button>
            <el-button type="primary" @click="handleDownloadCSV" size="small">
              下载CSV
            </el-button>
            <el-button type="warning" @click="handleDownloadTXT" size="small">
              下载TXT
            </el-button>
          </div>
        </div>
      </template>

      <el-alert
        :title="`成功创建 ${result.success_count} 个用户${result.error_count > 0 ? '，失败 ' + result.error_count + ' 个' : ''}`"
        :type="result.error_count > 0 ? 'warning' : 'success'"
        :closable="false"
        style="margin-bottom: 20px;"
      />

      <el-table :data="result.users" border max-height="400">
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

      <div v-if="result.errors && result.errors.length > 0" style="margin-top: 20px;">
        <el-alert
          title="错误信息"
          type="error"
          :closable="false"
        >
          <div v-for="(error, index) in result.errors" :key="index">
            {{ error }}
          </div>
        </el-alert>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { quickBatchCreate } from '@/api/user'
import { list as getGroupList } from '@/api/group'
import { list as getPackageList } from '@/api/package'

const formRef = ref(null)
const loading = ref(false)
const groups = ref([])
const packages = ref([])
const result = ref(null)

const form = reactive({
  username_prefix: 'user',
  count: 10,
  password_length: 8,
  group_id: null,
  package_id: null,
  status: 1,
  valid_days: 365,
  device_limit: 10
})

const rules = {
  username_prefix: [
    { required: true, message: '请输入用户名前缀', trigger: 'blur' },
    { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  count: [
    { required: true, message: '请输入创建数量', trigger: 'blur' }
  ],
  password_length: [
    { required: true, message: '请输入密码位数', trigger: 'blur' }
  ],
  group_id: [
    { required: true, message: '请选择用户组', trigger: 'change' }
  ]
}

onMounted(async () => {
  const res = await getGroupList({ page: 1, page_size: 100 })
  groups.value = res.data?.list || []
  if (groups.value.length > 0) {
    form.group_id = groups.value[0].id
  } else {
    ElMessage.warning('未找到用户组，请先在"用户组管理"中创建用户组')
  }

  const pkgRes = await getPackageList({ page: 1, page_size: 100 })
  packages.value = pkgRes.data?.list || []
})

watch(() => form.package_id, (newPackageId) => {
  if (newPackageId) {
    const selectedPackage = packages.value.find(pkg => pkg.id === newPackageId)
    if (selectedPackage) {
      form.valid_days = selectedPackage.valid_days
      form.device_limit = selectedPackage.device_limit
    }
  }
})

const handleCreate = async () => {
  await formRef.value.validate()

  loading.value = true
  try {
    const res = await quickBatchCreate(form)
    result.value = res.data
    ElMessage.success(res.data.message || '创建成功')
  } catch (error) {
    ElMessage.error(error.message || '创建失败')
  } finally {
    loading.value = false
  }
}

const handleReset = () => {
  formRef.value.resetFields()
  result.value = null
}

const copyText = (text) => {
  if (!text) {
    ElMessage.warning('内容为空，无法复制')
    return
  }

  // 尝试使用现代 Clipboard API
  if (navigator.clipboard && navigator.clipboard.writeText) {
    navigator.clipboard.writeText(text).then(() => {
      ElMessage.success('已复制到剪贴板')
    }).catch((err) => {
      console.error('复制失败:', err)
      // 降级到传统方法
      fallbackCopyText(text)
    })
  } else {
    // 降级到传统方法
    fallbackCopyText(text)
  }
}

// 降级复制方法（兼容旧浏览器）
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

const handleCopyText = () => {
  if (!result.value || !result.value.export_text) return
  copyText(result.value.export_text)
}

const handleDownloadCSV = () => {
  if (!result.value || !result.value.export_csv) return
  downloadFile(result.value.export_csv, 'users.csv', 'text/csv')
}

const handleDownloadTXT = () => {
  if (!result.value || !result.value.export_text) return
  downloadFile(result.value.export_text, 'users.txt', 'text/plain')
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
