<template>
  <div>
    <el-card class="list-query" shadow="hover">
      <el-alert
        title="线路能力说明"
        type="info"
        :closable="false"
        show-icon
      >
        <template #default>
          标准线路使用常规 TCP 中继；公司/校园网络线路在此基础上增加 WebSocket 能力，便于客户端按网络环境自动切换。
        </template>
      </el-alert>
      <el-form inline label-width="80px" style="margin-top: 16px">
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
        <el-table-column prop="region" label="地区" align="center" width="110" />
        <el-table-column label="线路类型" align="center" width="140">
          <template #default="{ row }">
            <el-tag :type="row.support_wss ? 'success' : 'info'">
              {{ row.support_wss ? '公司/校园网络' : '标准线路' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="id_server" label="ID 服务器" min-width="180" />
        <el-table-column prop="relay_server" label="中继服务器" min-width="180" />
        <el-table-column prop="ws_host" label="WS Host" min-width="180">
          <template #default="{ row }">
            {{ row.ws_host || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="priority" label="优先级" align="center" width="90" />
        <el-table-column prop="cost_weight" label="成本权重" align="center" width="100" />
        <el-table-column prop="is_default" label="默认" align="center" width="80">
          <template #default="{ row }">
            <el-tag :type="row.is_default ? 'success' : 'info'" size="small">
              {{ row.is_default ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="is_active" label="状态" align="center" width="80">
          <template #default="{ row }">
            <el-tag :type="row.is_active ? 'success' : 'danger'" size="small">
              {{ row.is_active ? '启用' : '停用' }}
            </el-tag>
          </template>
        </el-table-column>
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
          <el-input v-model="formData.name" placeholder="例如：上海标准线" />
        </el-form-item>
        <el-form-item label="地区" prop="region">
          <el-input v-model="formData.region" placeholder="例如：华东 / 香港" />
        </el-form-item>
        <el-form-item label="ID 服务器" prop="id_server" required>
          <el-input v-model="formData.id_server" placeholder="例如：id.example.com" />
        </el-form-item>
        <el-form-item label="中继服务器" prop="relay_server" required>
          <el-input v-model="formData.relay_server" placeholder="例如：relay.example.com" />
        </el-form-item>
        <el-form-item label="API 服务器" prop="api_server">
          <el-input v-model="formData.api_server" placeholder="例如：https://api.example.com" />
        </el-form-item>
        <el-form-item label="WS Host" prop="ws_host">
          <el-input
            v-model="formData.ws_host"
            placeholder="支持公司/校园网络线路时填写，例如：wss://relay.example.com/ws/relay"
          />
        </el-form-item>
        <el-form-item label="公钥" prop="key">
          <el-input v-model="formData.key" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="支持 TCP" prop="support_tcp">
          <el-switch v-model="formData.support_tcp" />
        </el-form-item>
        <el-form-item label="支持 WebSocket" prop="support_wss">
          <el-switch v-model="formData.support_wss" />
        </el-form-item>
        <el-form-item label="优先级" prop="priority">
          <el-input-number v-model="formData.priority" :min="0" />
        </el-form-item>
        <el-form-item label="成本权重" prop="cost_weight">
          <el-input-number v-model="formData.cost_weight" :min="1" :max="10" />
        </el-form-item>
        <el-form-item label="默认线路" prop="is_default">
          <el-switch v-model="formData.is_default" />
        </el-form-item>
        <el-form-item label="启用" prop="is_active">
          <el-switch v-model="formData.is_active" />
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
import { list, create, update, remove } from '@/api/server'
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

const defaultFormData = () => ({
  id: 0,
  name: '',
  region: '',
  id_server: '',
  relay_server: '',
  key: '',
  api_server: '',
  ws_host: '',
  support_tcp: true,
  support_wss: false,
  cost_weight: 1,
  is_default: false,
  is_active: true,
  priority: 0,
  description: '',
})

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

onMounted(getList)
onActivated(getList)

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
    region: row.region,
    id_server: row.id_server,
    relay_server: row.relay_server,
    key: row.key,
    api_server: row.api_server,
    ws_host: row.ws_host,
    support_tcp: row.support_tcp,
    support_wss: row.support_wss,
    cost_weight: row.cost_weight ?? 1,
    is_default: row.is_default,
    is_active: row.is_active,
    priority: row.priority,
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
