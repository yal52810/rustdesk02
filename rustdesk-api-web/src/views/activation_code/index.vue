<template>
  <div>
    <el-card shadow="hover">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-button type="primary" @click="showCreateDialog">创建激活码</el-button>
          <el-button type="success" @click="showBatchCreateDialog">批量创建</el-button>
        </el-col>
      </el-row>

      <el-table :data="list" style="width: 100%; margin-top: 20px" border>
        <el-table-column prop="id" label="ID" width="80"></el-table-column>
        <el-table-column prop="code" label="激活码" width="300"></el-table-column>
        <el-table-column label="套餐" width="150">
          <template #default="scope">
            {{ scope.row.package?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="valid_days" label="有效天数" width="100"></el-table-column>
        <el-table-column prop="device_limit" label="设备限制" width="100"></el-table-column>
        <el-table-column label="过期时间" width="180">
          <template #default="scope">
            {{ scope.row.expires_at || '永久' }}
          </template>
        </el-table-column>
        <el-table-column label="使用状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.used_by > 0 ? 'success' : 'info'">
              {{ scope.row.used_by > 0 ? '已使用' : '未使用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注"></el-table-column>
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="scope">
            <el-button type="danger" size="small" @click="handleDelete(scope.row)" :disabled="scope.row.used_by > 0">删除</el-button>
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

    <!-- 创建激活码对话框 -->
    <el-dialog v-model="createDialogVisible" title="创建激活码" width="500px">
      <el-form :model="createForm" label-width="100px">
        <el-form-item label="套餐">
          <el-select v-model="createForm.package_id" clearable placeholder="选择套餐">
            <el-option
              v-for="item in packagesList"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            ></el-option>
          </el-select>
          <el-text class="mx-1" type="info">选择套餐后自动填充</el-text>
        </el-form-item>
        <el-form-item label="有效天数" required>
          <el-input-number v-model="createForm.valid_days" :min="1"></el-input-number>
        </el-form-item>
        <el-form-item label="设备限制" required>
          <el-input-number v-model="createForm.device_limit" :min="1" :max="1000"></el-input-number>
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
          <el-input v-model="createForm.remark" type="textarea"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleCreate">确定</el-button>
      </template>
    </el-dialog>

    <!-- 批量创建对话框 -->
    <el-dialog v-model="batchCreateDialogVisible" title="批量创建激活码" width="500px">
      <el-form :model="batchCreateForm" label-width="100px">
        <el-form-item label="套餐">
          <el-select v-model="batchCreateForm.package_id" clearable placeholder="选择套餐">
            <el-option
              v-for="item in packagesList"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            ></el-option>
          </el-select>
          <el-text class="mx-1" type="info">选择套餐后自动填充</el-text>
        </el-form-item>
        <el-form-item label="数量" required>
          <el-input-number v-model="batchCreateForm.count" :min="1" :max="1000"></el-input-number>
        </el-form-item>
        <el-form-item label="有效天数" required>
          <el-input-number v-model="batchCreateForm.valid_days" :min="1"></el-input-number>
        </el-form-item>
        <el-form-item label="设备限制" required>
          <el-input-number v-model="batchCreateForm.device_limit" :min="1" :max="1000"></el-input-number>
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
          <el-input v-model="batchCreateForm.remark" type="textarea"></el-input>
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
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const createDialogVisible = ref(false)
const batchCreateDialogVisible = ref(false)
const packagesList = ref([])

const createForm = ref({
  package_id: null,
  valid_days: 365,
  device_limit: 10,
  expires_at: null,
  remark: ''
})

const batchCreateForm = ref({
  package_id: null,
  count: 10,
  valid_days: 365,
  device_limit: 10,
  expires_at: null,
  remark: ''
})

const getList = async () => {
  const res = await getActivationCodeList({ page: page.value, page_size: pageSize.value })
  if (res.code === 0) {
    list.value = res.data.list || []
    total.value = res.data.total
  }
}

const getPackages = async () => {
  const res = await packages({ page_size: 9999 }).catch(_ => false)
  if (res) {
    packagesList.value = res.data?.list || []
  }
}

// 监听套餐选择变化 - 创建表单
watch(() => createForm.value.package_id, (newVal) => {
  if (newVal) {
    const selectedPackage = packagesList.value.find(p => p.id === newVal)
    if (selectedPackage) {
      createForm.value.valid_days = selectedPackage.valid_days
      createForm.value.device_limit = selectedPackage.device_limit
    }
  }
})

// 监听套餐选择变化 - 批量创建表单
watch(() => batchCreateForm.value.package_id, (newVal) => {
  if (newVal) {
    const selectedPackage = packagesList.value.find(p => p.id === newVal)
    if (selectedPackage) {
      batchCreateForm.value.valid_days = selectedPackage.valid_days
      batchCreateForm.value.device_limit = selectedPackage.device_limit
    }
  }
})

const showCreateDialog = () => {
  createForm.value = {
    package_id: null,
    valid_days: 365,
    device_limit: 10,
    expires_at: null,
    remark: ''
  }
  createDialogVisible.value = true
}

const showBatchCreateDialog = () => {
  batchCreateForm.value = {
    package_id: null,
    count: 10,
    valid_days: 365,
    device_limit: 10,
    expires_at: null,
    remark: ''
  }
  batchCreateDialogVisible.value = true
}

const handleCreate = async () => {
  const res = await create(createForm.value)
  if (res.code === 0) {
    ElMessage.success('创建成功')
    createDialogVisible.value = false
    getList()
  }
}

const handleBatchCreate = async () => {
  const res = await batchCreate(batchCreateForm.value)
  if (res.code === 0) {
    ElMessage.success('批量创建成功')
    batchCreateDialogVisible.value = false
    getList()
  }
}

const handleDelete = async (row) => {
  await ElMessageBox.confirm('确定删除该激活码吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
  const res = await remove(row.id)
  if (res.code === 0) {
    ElMessage.success('删除成功')
    getList()
  }
}

onMounted(() => {
  getList()
  getPackages()
})
</script>
