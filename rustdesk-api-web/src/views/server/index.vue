<template>
  <div>
    <!-- ======================================== -->
    <!-- 状态卡片行 -->
    <!-- ======================================== -->
    <el-row :gutter="16" style="margin-bottom: 16px">
      <el-col :span="5">
        <el-card shadow="hover" class="status-card">
          <div class="status-row">
            <span class="status-dot" :class="idOnline ? 'online' : 'offline'"></span>
            <div>
              <div class="status-title">ID 信令服务器</div>
              <div class="status-addr">{{ idAddr || '未配置' }}</div>
            </div>
            <el-tag :type="idOnline ? 'success' : 'danger'" size="small" effect="dark">
              {{ idOnline ? '在线' : '离线' }}
            </el-tag>
          </div>
        </el-card>
      </el-col>
      <el-col :span="5">
        <el-card shadow="hover" class="status-card">
          <div class="status-row">
            <div>
              <div class="status-title">API 管理后台</div>
              <div class="status-addr">{{ apiAddr }}</div>
            </div>
            <el-tag type="success" size="small" effect="dark">运行中</el-tag>
          </div>
        </el-card>
      </el-col>
      <el-col :span="5">
        <el-card shadow="hover" class="status-card">
          <div class="status-row">
            <div>
              <div class="status-title">中继节点</div>
              <div class="status-addr">{{ relayNodes.length }} 个在线 / {{ listRes.total }} 个总计</div>
            </div>
            <el-tag :type="relayNodes.length > 0 ? 'success' : 'warning'" size="small" effect="dark">
              {{ relayNodes.length > 0 ? '正常' : '无节点' }}
            </el-tag>
          </div>
        </el-card>
      </el-col>
      <el-col :span="5">
        <el-card shadow="hover" class="status-card key-card">
          <div class="status-row key-row">
            <div style="flex:1;min-width:0">
              <div class="status-title">全局配置</div>
              <div class="status-addr mono-key">{{ globalServerConfig.key || '密钥未配置' }}</div>
            </div>
            <el-button text size="small" @click="showKeyDialog" type="primary">编辑</el-button>
          </div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card shadow="hover" class="status-card">
          <div style="display:flex;align-items:center;justify-content:space-between">
            <div class="status-title">发卡网设置</div>
            <el-button text size="small" type="primary" @click="showCardShop">设置</el-button>
          </div>
          <div style="font-size:12px;color:#909399;margin-top:4px">
            {{ globalServerConfig.card_shop_url ? '已配置' : '未配置' }}
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- ======================================== -->
    <!-- 中继节点列表 -->
    <!-- ======================================== -->
    <el-card shadow="hover">
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center">
          <span><b>中继节点管理</b></span>
          <div style="display:flex;gap:8px">
            <el-button type="primary" size="small" @click="toAdd">+ 添加节点</el-button>
          </div>
        </div>
      </template>

      <el-table :data="listRes.list" v-loading="listRes.loading" border stripe>
        <el-table-column prop="name" label="节点名称" min-width="140" />
        <el-table-column prop="region" label="地区" width="100" align="center" />
        <el-table-column label="线路类型" width="140" align="center">
          <template #default="{ row }">
            <div style="display:flex;gap:4px;justify-content:center;flex-wrap:wrap">
              <el-tag type="" size="small">TCP 标准</el-tag>
              <el-tag v-if="row.support_wss" type="success" size="small">WSS 专业</el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="relay_server" label="中继地址" min-width="200" />
        <el-table-column label="在线状态" width="120" align="center">
          <template #default="{ row }">
            <el-switch
              :model-value="row.is_online"
              :active-value="true"
              :inactive-value="false"
              active-text="在线"
              inactive-text="离线"
              @change="(val) => handleToggleOnline(row, val)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="last_check_at" label="最近检测" width="160" align="center">
          <template #default="{ row }">
            {{ row.last_check_at || '尚未检测' }}
          </template>
        </el-table-column>
        <el-table-column label="启用" width="70" align="center">
          <template #default="{ row }">
            <el-tag :type="row.is_active ? 'success' : 'info'" size="small">
              {{ row.is_active ? '开' : '关' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160" align="center" fixed="right">
          <template #default="{ row }">
            <el-button text type="primary" size="small" @click="toEdit(row)">编辑</el-button>
            <el-button text type="danger" size="small" @click="del(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div style="margin-top: 16px; display: flex; justify-content: flex-end">
        <el-pagination
          background
          layout="prev, pager, next, sizes, jumper"
          :page-sizes="[10, 20, 50, 100]"
          v-model:page-size="listQuery.page_size"
          v-model:current-page="listQuery.page"
          :total="listRes.total"
          small
        />
      </div>
    </el-card>

    <!-- ======================================== -->
    <!-- 添加 / 编辑中继节点对话框 -->
    <!-- ======================================== -->
    <el-dialog
      v-model="formVisible"
      :title="formData.id ? '编辑中继节点' : '添加中继节点'"
      width="560px"
      :close-on-click-modal="false"
    >
      <el-form :model="formData" label-width="100px">
        <el-form-item label="节点名称" required>
          <el-input v-model="formData.name" placeholder="例如：北京-联通" maxlength="50" />
        </el-form-item>
        <el-form-item label="地区代码" required>
          <el-select v-model="formData.region" placeholder="选择地区" style="width: 100%" filterable allow-create>
            <el-option label="🇨🇳 华东 (CN-East)" value="CN-East" />
            <el-option label="🇨🇳 华北 (CN-North)" value="CN-North" />
            <el-option label="🇨🇳 华南 (CN-South)" value="CN-South" />
            <el-option label="🇨🇳 西南 (CN-West)" value="CN-West" />
            <el-option label="🇭🇰 香港 (HK)" value="HK" />
            <el-option label="🇯🇵 日本 (JP)" value="JP" />
            <el-option label="🇺🇸 美国 (US)" value="US" />
            <el-option label="🇪🇺 欧洲 (EU)" value="EU" />
            <el-option label="🇸🇬 新加坡 (SG)" value="SG" />
          </el-select>
        </el-form-item>
        <el-form-item label="ID 服务器" required>
          <el-input v-model="formData.id_server" placeholder="例如：id.example.com:21116" />
        </el-form-item>
        <el-form-item label="中继地址" required>
          <el-input v-model="formData.relay_server" placeholder="例如：relay-bj.example.com:21117" />
        </el-form-item>
        <el-form-item label="专业线路地址" :required="formData.support_wss">
          <el-input v-model="formData.ws_host" placeholder="例如：relay-bj.example.com:21119" />
          <div style="font-size:12px;color:#909399;margin-top:2px">开启下方 WSS 专业线路后必须填写，用于 WebSocket 加密传输</div>
        </el-form-item>
        <el-divider content-position="left" style="margin: 8px 0">高级选项</el-divider>
        <el-form-item label="线路能力">
          <div style="display:flex;align-items:center;gap:12px">
            <el-tag type="" size="default">TCP 标准线路</el-tag>
            <span style="color:#909399;font-size:12px">（默认开启）</span>
          </div>
        </el-form-item>
        <el-form-item label="WSS 专业线路">
          <el-switch v-model="formData.support_wss" active-text="启用" inactive-text="未启用" />
          <div style="font-size:12px;color:#909399;margin-top:4px">在 TCP 基础上额外支持 WebSocket 加密线路，需填写上方「专业线路地址」</div>
        </el-form-item>
        <el-form-item label="启用节点">
          <el-switch v-model="formData.is_active" active-text="启用" inactive-text="停用" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="formVisible = false">取消</el-button>
        <el-button type="primary" @click="submit" :loading="submitting">
          {{ formData.id ? '保存' : '添加' }}
        </el-button>
      </template>
    </el-dialog>

    <!-- ======================================== -->
    <!-- 发卡网设置对话框 -->
    <!-- ======================================== -->
    <el-dialog v-model="cardShopVisible" title="发卡网设置" width="460px" :close-on-click-modal="false">
      <el-form label-width="80px">
        <el-form-item label="发卡网址">
          <el-input v-model="cardShopForm.url" placeholder="例如：https://shop.example.com" />
        </el-form-item>
        <div class="key-hint">设置后在用户中心"购买套餐"处显示"前往发卡网购买"按钮。</div>
      </el-form>
      <template #footer>
        <el-button @click="cardShopVisible = false">取消</el-button>
        <el-button type="primary" @click="saveCardShop" :loading="cardShopSaving">保存</el-button>
      </template>
    </el-dialog>

    <!-- ======================================== -->
    <!-- 服务器全局配置对话框 -->
    <!-- ======================================== -->
    <el-dialog v-model="keyDialogVisible" title="服务器全局配置" width="560px" :close-on-click-modal="false">
      <el-form :model="keyForm" label-width="120px">
        <el-form-item label="ID 服务器">
          <el-input v-model="keyForm.id_server" placeholder="例如：id.example.com:21116" />
        </el-form-item>
        <el-form-item label="中继服务器">
          <el-input v-model="keyForm.relay_server" placeholder="例如：relay.example.com:21117" />
        </el-form-item>
        <el-form-item label="API 服务器">
          <el-input v-model="keyForm.api_server" placeholder="例如：https://api.example.com" />
        </el-form-item>
        <el-form-item label="统一密钥">
          <el-input v-model="keyForm.key" placeholder="输入密钥，留空则自动生成" />
        </el-form-item>
        <el-form-item label="专业线路地址">
          <el-input v-model="keyForm.ws_host" placeholder="例如：wss://api.example.com" />
        </el-form-item>
        <div class="key-hint" style="margin-left: 120px">修改配置后需重启 hbbs/hbbr 服务生效。<br/>客户端配置生成将使用此处填写的服务器信息。</div>
      </el-form>
      <template #footer>
        <el-button @click="keyDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveKey" :loading="keySaving">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { list, create, update, remove, checkServers, toggleOnline } from '@/api/server'
import { getServerConfig, updateServerKey } from '@/api/config'
import { ElMessage, ElMessageBox } from 'element-plus'

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
  api_server: window.location.origin,
  ws_host: '',
  support_wss: false,
  cost_weight: 1,
  is_default: false,
  is_active: true,
  priority: 0,
  description: '',
})

// ============================================
// 全局服务器配置（从 /config/server 加载）
// ============================================
const globalServerConfig = ref({
  id_server: '',
  relay_server: '',
  api_server: '',
  key: '',
  ws_host: '',
  card_shop_url: '',
})

const idAddr = computed(() => {
  // 优先从全局配置读取，其次从节点列表
  return globalServerConfig.value.id_server || (listRes.list.find(s => s.id_server)?.id_server) || ''
})

const apiAddr = computed(() => {
  return window.location.host
})

const idOnline = computed(() => {
  if (listRes.list.length === 0) return null
  const idServers = listRes.list.filter(s => s.id_server)
  if (idServers.length === 0) return null
  return idServers.every(s => s.is_online)
})

const relayNodes = computed(() => {
  return listRes.list.filter(s => s.is_online)
})

// ============================================
// 数据获取
// ============================================
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
  const confirmed = await ElMessageBox.confirm(`确定删除中继节点「${row.name}」？`, '确认删除', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning',
  }).catch(() => false)
  if (!confirmed) return

  const res = await remove({ id: row.id }).catch(() => false)
  if (res) {
    ElMessage.success('已删除')
    getList()
  }
}

// ============================================
// 统一密钥管理
// ============================================
const serverKey = ref('')
const keyDialogVisible = ref(false)
const keySaving = ref(false)
const rechecking = ref(false)
const cardShopVisible = ref(false)
const cardShopSaving = ref(false)
const cardShopForm = reactive({ url: '' })
const keyForm = reactive({
  id_server: '',
  relay_server: '',
  api_server: '',
  key: '',
  ws_host: '',
  card_shop_url: '',
})

const recheckAll = async () => {
  rechecking.value = true
  const res = await checkServers().catch(() => false)
  rechecking.value = false
  if (res) {
    ElMessage.success('检测已触发，稍后刷新查看结果')
    setTimeout(() => getList(), 3000)
  }
}

const handleToggleOnline = async (row, val) => {
  const res = await toggleOnline({ id: row.id, is_online: val }).catch(() => false)
  if (res) {
    row.is_online = val
    ElMessage.success(`节点「${row.name}」已设为${val ? '在线' : '离线'}`)
  }
}

const showCardShop = () => {
  cardShopForm.url = globalServerConfig.value.card_shop_url || ''
  cardShopVisible.value = true
}

const saveCardShop = async () => {
  cardShopSaving.value = true
  const res = await updateServerKey({
    id_server: globalServerConfig.value.id_server,
    relay_server: globalServerConfig.value.relay_server,
    api_server: globalServerConfig.value.api_server,
    key: globalServerConfig.value.key,
    ws_host: globalServerConfig.value.ws_host,
    card_shop_url: cardShopForm.url,
  }).catch(() => false)
  cardShopSaving.value = false
  if (res) {
    globalServerConfig.value.card_shop_url = cardShopForm.url
    ElMessage.success('发卡网地址已保存')
    cardShopVisible.value = false
  }
}

const loadServerKey = async () => {
  const res = await getServerConfig().catch(() => false)
  if (res && res.data) {
    globalServerConfig.value = {
      id_server: res.data.id_server || '',
      relay_server: res.data.relay_server || '',
      api_server: res.data.api_server || '',
      key: res.data.key || '',
      ws_host: res.data.ws_host || '',
      card_shop_url: res.data.card_shop_url || '',
    }
    serverKey.value = res.data.key || ''
    keyForm.id_server = res.data.id_server || ''
    keyForm.relay_server = res.data.relay_server || ''
    keyForm.api_server = res.data.api_server || ''
    keyForm.key = res.data.key || ''
    keyForm.ws_host = res.data.ws_host || ''
    keyForm.card_shop_url = res.data.card_shop_url || ''
  }
}

const showKeyDialog = () => {
  keyDialogVisible.value = true
}

const saveKey = async () => {
  keySaving.value = true
  const res = await updateServerKey({
    id_server: keyForm.id_server,
    relay_server: keyForm.relay_server,
    api_server: keyForm.api_server,
    key: keyForm.key,
    ws_host: keyForm.ws_host,
    card_shop_url: keyForm.card_shop_url,
  }).catch(() => false)
  keySaving.value = false
  if (res) {
    serverKey.value = keyForm.key
    ElMessage.success('服务器配置已更新，请重启 hbbs/hbbr 服务生效')
    keyDialogVisible.value = false
  }
}

onMounted(loadServerKey)
onMounted(getList)
watch(() => listQuery.page, getList)
watch(() => listQuery.page_size, handlerQuery)

// ============================================
// 表单逻辑
// ============================================
const formVisible = ref(false)
const submitting = ref(false)
const formData = reactive(defaultFormData())

const assignFormData = (row = defaultFormData()) => {
  Object.assign(formData, defaultFormData(), {
    id: row.id || 0,
    name: row.name || '',
    region: row.region || '',
    relay_server: row.relay_server || '',
    ws_host: row.ws_host || '',
    support_wss: row.support_wss ?? (!!row.ws_host),
    is_active: row.is_active ?? true,
    id_server: row.id_server || idAddr.value,
    api_server: row.api_server || window.location.origin,
    cost_weight: row.cost_weight ?? 1,
    priority: row.priority ?? 0,
  })
}

const toEdit = (row) => {
  formVisible.value = true
  assignFormData(row)
}

const toAdd = () => {
  formVisible.value = true
  assignFormData()
}

const submit = async () => {
  if (!formData.name.trim()) return ElMessage.warning('请输入节点名称')
  if (!formData.region.trim()) return ElMessage.warning('请选择地区')
  if (!formData.relay_server.trim()) return ElMessage.warning('请输入中继地址')
  if (formData.support_wss && !formData.ws_host.trim()) return ElMessage.warning('启用 WSS 专业线路时必须填写专业线路地址')

  submitting.value = true
  const api = formData.id ? update : create
  const res = await api({ ...formData }).catch(() => false)
  submitting.value = false

  if (res) {
    ElMessage.success(formData.id ? '已保存' : '已添加')
    formVisible.value = false
    getList()
  }
}
</script>

<style scoped lang="scss">
.status-card {
  .status-row {
    display: flex;
    align-items: center;
    gap: 12px;
  }
  .status-title {
    font-size: 13px;
    color: #909399;
    margin-bottom: 2px;
  }
  .status-addr {
    font-size: 14px;
    font-weight: 500;
    word-break: break-all;
  }
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
  &.online {
    background: #67c23a;
    box-shadow: 0 0 6px rgba(103, 194, 58, 0.6);
  }
  &.offline {
    background: #f56c6c;
    box-shadow: 0 0 6px rgba(245, 108, 108, 0.6);
  }
}

.key-card .key-row {
  justify-content: space-between;
}
.mono-key {
  font-family: Consolas, Monaco, monospace;
  font-size: 11px;
  word-break: break-all;
}
.key-hint {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
  line-height: 1.5;
}
</style>
