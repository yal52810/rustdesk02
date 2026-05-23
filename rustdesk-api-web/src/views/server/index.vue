<template>
  <div>
    <el-card class="list-query" shadow="hover">
      <el-alert
        title="线路互通说明"
        type="info"
        :closable="false"
        show-icon
      >
        <template #default>
          所有线路建议共用同一套 ID 服务、API 和 Key，仅区分不同中继入口与受限网络接入方式。这样控制端切美国、日本线路时，仍可与未手动切线的被控端互通。
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
        <el-table-column prop="topology_group" label="互通组" align="center" width="130">
          <template #default="{ row }">
            {{ row.topology_group || 'default' }}
          </template>
        </el-table-column>
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
        <el-table-column label="启用状态" align="center" width="90">
          <template #default="{ row }">
            <el-tag :type="row.is_active ? 'success' : 'danger'" size="small">
              {{ row.is_active ? '启用' : '停用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="在线状态" align="center" width="90">
          <template #default="{ row }">
            <el-tag :type="row.is_online ? 'success' : 'danger'" size="small">
              {{ row.is_online ? '在线' : '离线' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="最近检测" align="center" min-width="160">
          <template #default="{ row }">
            {{ row.last_check_at || '-' }}
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
          <el-input v-model="formData.name" placeholder="例如：美国标准线路" />
        </el-form-item>
        <el-form-item label="地区" prop="region">
          <el-input v-model="formData.region" placeholder="例如：US / JP / HK" />
        </el-form-item>
        <el-form-item label="互通组" prop="topology_group">
          <el-input v-model="formData.topology_group" placeholder="例如：global-main，同组共享 hbbs/api/key" />
        </el-form-item>
        <el-form-item label="ID 服务器" prop="id_server" required>
          <el-input v-model="formData.id_server" placeholder="例如：id.example.com:21116" />
        </el-form-item>
        <el-form-item label="中继服务器" prop="relay_server" required>
          <el-input v-model="formData.relay_server" placeholder="例如：us-relay.example.com:21117" />
        </el-form-item>
        <el-form-item label="WS Host" prop="ws_host">
          <el-input
            v-model="formData.ws_host"
            placeholder="公司/校园线路填写，例如：wss://relay.example.com/ws/relay"
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
  topology_group: 'default',
  id_server: '',
  relay_server: '',
  key: '',
  api_server: window.location.origin,
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
    topology_group: row.topology_group || 'default',
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
