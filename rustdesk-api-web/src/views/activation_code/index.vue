<template>
  <div>
    <el-card shadow="hover">
      <el-alert
        title="激活码说明"
        type="info"
        :closable="false"
        show-icon
      >
        <template #default>
          选择套餐后，激活码会自动继承套餐对应的线路权限。客户端兑换后，账号有效期、设备限制和默认线路会一起生效。
        </template>
      </el-alert>

      <el-row :gutter="20" style="margin-top: 16px">
        <el-col :span="12">
          <el-button type="primary" @click="showCreateDialog">创建</el-button>
          <el-button type="success" @click="showBatchCreateDialog">批量创建</el-button>
          <el-button :disabled="!lastBatchCodes.length" @click="exportBatchCodes()">导出最近一批</el-button>
        </el-col>
      </el-row>

      <el-table :data="list" style="width: 100%; margin-top: 20px" border>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="code" label="激活码" min-width="280" />
        <el-table-column label="套餐" width="150">
          <template #default="{ row }">
            {{ row.package?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="继承线路" min-width="220">
          <template #default="{ row }">
            <el-space wrap>
              <el-tag
                v-for="server in row.package?.servers || []"
                :key="server.id"
                :type="server.support_wss ? 'success' : 'info'"
                size="small"
              >
                {{ server.name }}
              </el-tag>
              <span v-if="!(row.package?.servers || []).length">-</span>
            </el-space>
          </template>
        </el-table-column>
        <el-table-column prop="valid_days" label="有效天数" width="100" />
        <el-table-column prop="device_limit" label="设备限制" width="100" />
        <el-table-column label="过期时间" width="180">
          <template #default="{ row }">
            {{ row.expires_at || '永久' }}
          </template>
        </el-table-column>
        <el-table-column label="使用状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.used_by > 0 ? 'success' : 'info'">
              {{ row.used_by > 0 ? '已使用' : '未使用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="160" />
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button type="danger" size="small" @click="handleDelete(row)" :disabled="row.used_by > 0">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        @current-change="getList"
        @size-change="getList"
        layout="total, sizes, prev, pager, next"
        style="margin-top: 20px"
      />
    </el-card>

    <el-dialog v-model="createDialogVisible" title="创建激活码" width="520px">
      <el-form :model="createForm" label-width="100px">
        <el-form-item label="套餐">
          <el-select v-model="createForm.package_id" clearable placeholder="选择套餐">
            <el-option
              v-for="item in packagesList"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="有效天数" required>
          <el-input-number v-model="createForm.valid_days" :min="1" />
        </el-form-item>
        <el-form-item label="设备限制" required>
          <el-input-number v-model="createForm.device_limit" :min="1" :max="1000" />
        </el-form-item>
        <el-form-item label="过期时间">
          <el-date-picker
            v-model="createForm.expires_at"
            type="datetime"
            placeholder="选择过期时间"
            format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="createForm.remark" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleCreate">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="batchCreateDialogVisible" title="批量创建激活码" width="520px">
      <el-form :model="batchCreateForm" label-width="100px">
        <el-form-item label="套餐">
          <el-select v-model="batchCreateForm.package_id" clearable placeholder="选择套餐">
            <el-option
              v-for="item in packagesList"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="数量" required>
          <el-input-number v-model="batchCreateForm.count" :min="1" :max="1000" />
        </el-form-item>
        <el-form-item label="有效天数" required>
          <el-input-number v-model="batchCreateForm.valid_days" :min="1" />
        </el-form-item>
        <el-form-item label="设备限制" required>
          <el-input-number v-model="batchCreateForm.device_limit" :min="1" :max="1000" />
        </el-form-item>
        <el-form-item label="过期时间">
          <el-date-picker
            v-model="batchCreateForm.expires_at"
            type="datetime"
            placeholder="选择过期时间"
            format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="batchCreateForm.remark" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="batchCreateDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleBatchCreate">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { list as getActivationCodeList, create, batchCreate, remove } from '@/api/activationCode'
import { list as packages } from '@/api/package'
import { ElMessage, ElMessageBox } from 'element-plus'

const list = ref([])
const lastBatchCodes = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const createDialogVisible = ref(false)
const batchCreateDialogVisible = ref(false)
const packagesList = ref([])

const defaultCreateForm = () => ({
  package_id: null,
  valid_days: 365,
  device_limit: 10,
  expires_at: null,
  remark: '',
})

const defaultBatchCreateForm = () => ({
  package_id: null,
  count: 10,
  valid_days: 365,
  device_limit: 10,
  expires_at: null,
  remark: '',
})

const createForm = ref(defaultCreateForm())
const batchCreateForm = ref(defaultBatchCreateForm())

const getList = async () => {
  const res = await getActivationCodeList({ page: page.value, page_size: pageSize.value }).catch(() => false)
  if (res && res.code === 0) {
    list.value = res.data.list || []
    total.value = res.data.total || 0
  }
}

const getPackages = async () => {
  const res = await packages({ page_size: 9999 }).catch(() => false)
  if (res) {
    packagesList.value = res.data?.list || []
  }
}

const applyPackageDefaults = (formRef, packageId) => {
  if (!packageId) {
    return
  }
  const selectedPackage = packagesList.value.find((item) => item.id === packageId)
  if (!selectedPackage) {
    return
  }
  formRef.value.valid_days = selectedPackage.valid_days
  formRef.value.device_limit = selectedPackage.device_limit
}

watch(() => createForm.value.package_id, (newVal) => applyPackageDefaults(createForm, newVal))
watch(() => batchCreateForm.value.package_id, (newVal) => applyPackageDefaults(batchCreateForm, newVal))

const showCreateDialog = () => {
  createForm.value = defaultCreateForm()
  createDialogVisible.value = true
}

const showBatchCreateDialog = () => {
  batchCreateForm.value = defaultBatchCreateForm()
  batchCreateDialogVisible.value = true
}

const formatDate = (value) => {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }
  const yyyy = date.getFullYear()
  const mm = String(date.getMonth() + 1).padStart(2, '0')
  const dd = String(date.getDate()).padStart(2, '0')
  const hh = String(date.getHours()).padStart(2, '0')
  const mi = String(date.getMinutes()).padStart(2, '0')
  const ss = String(date.getSeconds()).padStart(2, '0')
  return `${yyyy}-${mm}-${dd} ${hh}:${mi}:${ss}`
}

const csvEscape = (value) => {
  const text = (value ?? '').toString().replace(/"/g, '""')
  return `"${text}"`
}

const exportBatchCodes = (rows = lastBatchCodes.value) => {
  if (!rows.length) {
    ElMessage.warning('暂无可导出的批量激活码')
    return
  }
  const header = ['code', 'package_name', 'valid_days', 'device_limit', 'expires_at', 'remark']
  const lines = [
    header.join(','),
    ...rows.map((row) => [
      csvEscape(row.code),
      csvEscape(row.package?.name || ''),
      csvEscape(row.valid_days),
      csvEscape(row.device_limit),
      csvEscape(row.expires_at ? formatDate(row.expires_at) : ''),
      csvEscape(row.remark || ''),
    ].join(',')),
  ]
  const csvContent = `\ufeff${lines.join('\n')}`
  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  const stamp = new Date()
  const fileName = `activation-codes-${stamp.getFullYear()}${String(stamp.getMonth() + 1).padStart(2, '0')}${String(stamp.getDate()).padStart(2, '0')}-${String(stamp.getHours()).padStart(2, '0')}${String(stamp.getMinutes()).padStart(2, '0')}${String(stamp.getSeconds()).padStart(2, '0')}.csv`
  link.href = url
  link.download = fileName
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}

const handleCreate = async () => {
  const res = await create(createForm.value).catch(() => false)
  if (res && res.code === 0) {
    ElMessage.success('创建成功')
    createDialogVisible.value = false
    getList()
  }
}

const handleBatchCreate = async () => {
  const res = await batchCreate(batchCreateForm.value).catch(() => false)
  if (res && res.code === 0) {
    lastBatchCodes.value = Array.isArray(res.data) ? res.data : []
    ElMessage.success('批量创建成功，已导出 CSV')
    batchCreateDialogVisible.value = false
    exportBatchCodes(lastBatchCodes.value)
    getList()
  }
}

const handleDelete = async (row) => {
  const confirmed = await ElMessageBox.confirm('确定删除这个激活码吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  }).catch(() => false)
  if (!confirmed) {
    return
  }

  const res = await remove(row.id).catch(() => false)
  if (res && res.code === 0) {
    ElMessage.success('删除成功')
    getList()
  }
}

onMounted(() => {
  getList()
  getPackages()
})
</script>
