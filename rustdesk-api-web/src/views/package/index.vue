<template>
  <div>
    <el-card class="list-query" shadow="hover">
      <el-form inline label-width="80px">
        <el-form-item>
          <el-button type="primary" @click="handlerQuery">{{ T('Filter') }}</el-button>
          <el-button type="danger" @click="toAdd">{{ T('Add') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="list-body" shadow="hover">
      <el-table :data="listRes.list" v-loading="listRes.loading" border>
        <el-table-column prop="id" label="ID" align="center" width="80" />
        <el-table-column prop="name" :label="T('Name')" align="center" min-width="140" />
        <el-table-column prop="valid_days" label="有效天数" align="center" width="100" />
        <el-table-column prop="device_limit" label="设备限制" align="center" width="100" />
        <el-table-column prop="file_transfer_limit_mb" label="传输上限" align="center" width="110">
          <template #default="{ row }">
            {{ row.file_transfer_limit_mb || 100 }} MB
          </template>
        </el-table-column>
        <el-table-column prop="price" label="价格" align="center" width="100">
          <template #default="{ row }">
            ￥{{ row.price }}
          </template>
        </el-table-column>
        <el-table-column label="可用线路" min-width="220">
          <template #default="{ row }">
            <el-space wrap>
              <el-tag
                v-for="server in row.servers || []"
                :key="server.id"
                :type="server.support_wss ? 'success' : 'info'"
                size="small"
              >
                {{ server.name }}
              </el-tag>
              <span v-if="!(row.servers || []).length">-</span>
            </el-space>
          </template>
        </el-table-column>
        <el-table-column prop="is_active" label="状态" align="center" width="100">
          <template #default="{ row }">
            <el-tag :type="row.is_active ? 'success' : 'danger'">
              {{ row.is_active ? '启用' : '停用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="priority" label="优先级" align="center" width="100" />
        <el-table-column :label="T('Actions')" align="center" width="200" fixed="right">
          <template #default="{ row }">
            <el-button @click="toEdit(row)" size="small">{{ T('Edit') }}</el-button>
            <el-button type="danger" @click="del(row)" size="small">{{ T('Delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-card class="list-page" shadow="hover">
      <el-pagination
        background
        layout="prev, pager, next, sizes, jumper"
        :page-sizes="[10, 20, 50, 100]"
        v-model:page-size="listQuery.page_size"
        v-model:current-page="listQuery.page"
        :total="listRes.total"
      />
    </el-card>

    <el-dialog v-model="formVisible" :title="!formData.id ? T('Create') : T('Update')" width="820">
      <el-form class="dialog-form" :model="formData" label-width="120px">
        <el-form-item :label="T('Name')" prop="name" required>
          <el-input v-model="formData.name" placeholder="例如：普通版 / 专业版" />
        </el-form-item>
        <el-form-item label="有效天数" prop="valid_days" required>
          <el-input-number v-model="formData.valid_days" :min="1" :max="3650" />
        </el-form-item>
        <el-form-item label="设备限制" prop="device_limit" required>
          <el-input-number v-model="formData.device_limit" :min="1" :max="1000" />
        </el-form-item>
        <el-form-item label="传输上限(MB)" prop="file_transfer_limit_mb" required>
          <el-input-number v-model="formData.file_transfer_limit_mb" :min="1" :max="102400" />
        </el-form-item>
        <el-form-item label="价格" prop="price">
          <el-input-number v-model="formData.price" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="优先级" prop="priority">
          <el-input-number v-model="formData.priority" :min="0" />
        </el-form-item>
        <el-form-item label="启用" prop="is_active">
          <el-switch v-model="formData.is_active" />
        </el-form-item>
        <el-form-item label="可用线路" prop="server_ids">
          <el-select v-model="formData.server_ids" multiple placeholder="选择该套餐可使用的线路" style="width: 100%">
            <el-option
              v-for="server in servers"
              :key="server.id"
              :label="`${server.name}${server.support_wss ? '（专业线路）' : ''}`"
              :value="server.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item>
          <el-button @click="formVisible = false">{{ T('Cancel') }}</el-button>
          <el-button @click="submit" type="primary">{{ T('Submit') }}</el-button>
        </el-form-item>
      </el-form>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref, watch, onActivated } from 'vue'
import { list, create, update, remove } from '@/api/package'
import { list as serverList } from '@/api/server'
import { ElMessage, ElMessageBox } from 'element-plus'
import { T } from '@/utils/i18n'

const listRes = reactive({
  list: [],
  total: 0,
  loading: false,
})

const listQuery = reactive({
  page: 1,
  page_size: 10,
})

const servers = ref([])

const defaultFormData = () => ({
  id: 0,
  name: '',
  valid_days: 30,
  device_limit: 10,
  file_transfer_limit_mb: 100,
  price: 0,
  priority: 0,
  is_active: true,
  server_ids: [],
  description: '',
})

const getServerList = async () => {
  const res = await serverList({ page: 1, page_size: 1000 }).catch(() => false)
  if (res) {
    servers.value = res.data?.list || []
  }
}

const getList = async () => {
  listRes.loading = true
  const res = await list(listQuery).catch(() => false)
  listRes.loading = false
  if (res) {
    listRes.list = res.data?.list || []
    listRes.total = res.data?.total || 0
  }
}

const handlerQuery = () => {
  if (listQuery.page === 1) {
    getList()
  } else {
    listQuery.page = 1
  }
}

const del = async (row) => {
  const confirmed = await ElMessageBox.confirm(T('Confirm?', { param: T('Delete') }), {
    confirmButtonText: T('Confirm'),
    cancelButtonText: T('Cancel'),
    type: 'warning',
  }).catch(() => false)
  if (!confirmed) {
    return
  }

  const res = await remove({ id: row.id }).catch(() => false)
  if (res) {
    ElMessage.success(T('OperationSuccess'))
    getList()
  }
}

onMounted(() => {
  getServerList()
  getList()
})

onActivated(() => {
  getServerList()
  getList()
})

watch(() => listQuery.page, getList)
watch(() => listQuery.page_size, handlerQuery)

const formVisible = ref(false)
const formData = reactive(defaultFormData())

const assignFormData = (row = defaultFormData()) => {
  Object.assign(formData, defaultFormData(), row)
}

const toEdit = (row) => {
  formVisible.value = true
  assignFormData({
    id: row.id,
    name: row.name,
    valid_days: row.valid_days,
    device_limit: row.device_limit,
    file_transfer_limit_mb: row.file_transfer_limit_mb || 100,
    price: row.price,
    priority: row.priority,
    is_active: row.is_active,
    server_ids: row.servers?.map((server) => server.id) || [],
    description: row.description,
  })
}

const toAdd = () => {
  formVisible.value = true
  assignFormData()
}

const submit = async () => {
  const api = formData.id ? update : create
  const res = await api({ ...formData }).catch(() => false)
  if (res) {
    ElMessage.success(T('OperationSuccess'))
    formVisible.value = false
    getList()
  }
}
</script>

<style scoped lang="scss">
</style>
