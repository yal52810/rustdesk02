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
        <el-table-column prop="id" label="ID" align="center" width="80"></el-table-column>
        <el-table-column prop="name" :label="T('Name')" align="center" width="150"/>
        <el-table-column prop="region" label="地区" align="center" width="100"/>
        <el-table-column prop="id_server" label="ID服务器" align="center"/>
        <el-table-column prop="relay_server" label="中继服务器" align="center"/>
        <el-table-column prop="is_default" label="默认" align="center" width="80">
          <template #default="{row}">
            <el-tag :type="row.is_default ? 'success' : 'info'" size="small">
              {{ row.is_default ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="is_active" label="状态" align="center" width="80">
          <template #default="{row}">
            <el-tag :type="row.is_active ? 'success' : 'danger'" size="small">
              {{ row.is_active ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="priority" label="优先级" align="center" width="80"/>
        <el-table-column :label="T('Actions')" align="center" width="200">
          <template #default="{row}">
            <el-button @click="toEdit(row)" size="small">{{ T('Edit') }}</el-button>
            <el-button type="danger" @click="del(row)" size="small">{{ T('Delete') }}</el-button>
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
    <el-dialog v-model="formVisible" :title="!formData.id?T('Create'):T('Update')" width="800">
      <el-form class="dialog-form" ref="form" :model="formData" label-width="120px">
        <el-form-item :label="T('Name')" prop="name" required>
          <el-input v-model="formData.name" placeholder="例如：北京1"></el-input>
        </el-form-item>
        <el-form-item label="地区" prop="region">
          <el-input v-model="formData.region" placeholder="例如：北京"></el-input>
        </el-form-item>
        <el-form-item label="ID服务器" prop="id_server" required>
          <el-input v-model="formData.id_server" placeholder="例如：example.com"></el-input>
        </el-form-item>
        <el-form-item label="中继服务器" prop="relay_server" required>
          <el-input v-model="formData.relay_server" placeholder="例如：example.com"></el-input>
        </el-form-item>
        <el-form-item label="密钥" prop="key">
          <el-input v-model="formData.key" type="textarea" :rows="2"></el-input>
        </el-form-item>
        <el-form-item label="API服务器" prop="api_server">
          <el-input v-model="formData.api_server" placeholder="例如：http://example.com/api"></el-input>
        </el-form-item>
        <el-form-item label="优先级" prop="priority">
          <el-input-number v-model="formData.priority" :min="0"></el-input-number>
        </el-form-item>
        <el-form-item label="默认服务器" prop="is_default">
          <el-switch v-model="formData.is_default"></el-switch>
        </el-form-item>
        <el-form-item label="启用" prop="is_active">
          <el-switch v-model="formData.is_active"></el-switch>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3"></el-input>
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
  import { onMounted, reactive, watch, ref, onActivated } from 'vue'
  import { list, create, update, remove } from '@/api/server'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { T } from '@/utils/i18n'

  const listRes = reactive({
    list: [], total: 0, loading: false,
  })
  const listQuery = reactive({
    page: 1,
    page_size: 10,
  })

  const getList = async () => {
    listRes.loading = true
    const res = await list(listQuery).catch(_ => false)
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
    const cf = await ElMessageBox.confirm(T('Confirm?', { param: T('Delete') }), {
      confirmButtonText: T('Confirm'),
      cancelButtonText: T('Cancel'),
      type: 'warning',
    }).catch(_ => false)
    if (!cf) {
      return false
    }

    const res = await remove({ id: row.id }).catch(_ => false)
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
  const formData = reactive({
    id: 0,
    name: '',
    region: '',
    id_server: '',
    relay_server: '',
    key: '',
    api_server: '',
    is_default: false,
    is_active: true,
    priority: 0,
    description: '',
  })

  const toEdit = (row) => {
    formVisible.value = true
    formData.id = row.id
    formData.name = row.name
    formData.region = row.region
    formData.id_server = row.id_server
    formData.relay_server = row.relay_server
    formData.key = row.key
    formData.api_server = row.api_server
    formData.is_default = row.is_default
    formData.is_active = row.is_active
    formData.priority = row.priority
    formData.description = row.description
  }
  const toAdd = () => {
    formVisible.value = true
    formData.id = 0
    formData.name = ''
    formData.region = ''
    formData.id_server = ''
    formData.relay_server = ''
    formData.key = ''
    formData.api_server = ''
    formData.is_default = false
    formData.is_active = true
    formData.priority = 0
    formData.description = ''
  }
  const submit = async () => {
    const api = formData.id ? update : create
    const res = await api(formData).catch(_ => false)
    if (res) {
      ElMessage.success(T('OperationSuccess'))
      formVisible.value = false
      getList()
    }
  }

</script>

<style scoped lang="scss">

</style>
