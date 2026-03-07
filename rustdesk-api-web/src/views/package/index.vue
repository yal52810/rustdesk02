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
        <el-table-column prop="name" :label="T('Name')" align="center"/>
        <el-table-column prop="valid_days" label="有效天数" align="center" width="100"/>
        <el-table-column prop="device_limit" label="设备限制" align="center" width="100"/>
        <el-table-column prop="price" label="价格" align="center" width="100">
          <template #default="{row}">
            ¥{{ row.price }}
          </template>
        </el-table-column>
        <el-table-column prop="is_active" label="状态" align="center" width="100">
          <template #default="{row}">
            <el-tag :type="row.is_active ? 'success' : 'danger'">
              {{ row.is_active ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="priority" label="优先级" align="center" width="100"/>
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
          <el-input v-model="formData.name"></el-input>
        </el-form-item>
        <el-form-item label="有效天数" prop="valid_days" required>
          <el-input-number v-model="formData.valid_days" :min="1" :max="3650"></el-input-number>
        </el-form-item>
        <el-form-item label="设备限制" prop="device_limit" required>
          <el-input-number v-model="formData.device_limit" :min="1" :max="1000"></el-input-number>
        </el-form-item>
        <el-form-item label="价格" prop="price">
          <el-input-number v-model="formData.price" :min="0" :precision="2"></el-input-number>
        </el-form-item>
        <el-form-item label="优先级" prop="priority">
          <el-input-number v-model="formData.priority" :min="0"></el-input-number>
        </el-form-item>
        <el-form-item label="状态" prop="is_active">
          <el-switch v-model="formData.is_active"></el-switch>
        </el-form-item>
        <el-form-item label="服务器" prop="server_ids">
          <el-select v-model="formData.server_ids" multiple placeholder="选择服务器" style="width: 100%">
            <el-option v-for="server in servers" :key="server.id" :label="server.name" :value="server.id"/>
          </el-select>
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
  import { list, create, update, remove } from '@/api/package'
  import { list as serverList } from '@/api/server'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { T } from '@/utils/i18n'

  const listRes = reactive({
    list: [], total: 0, loading: false,
  })
  const listQuery = reactive({
    page: 1,
    page_size: 10,
  })

  const servers = ref([])

  const getServerList = async () => {
    const res = await serverList({ page: 1, page_size: 1000 }).catch(_ => false)
    if (res) {
      servers.value = res.data?.list || []
    }
  }

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
  const formData = reactive({
    id: 0,
    name: '',
    valid_days: 30,
    device_limit: 10,
    price: 0,
    priority: 0,
    is_active: true,
    server_ids: [],
    description: '',
  })

  const toEdit = (row) => {
    formVisible.value = true
    formData.id = row.id
    formData.name = row.name
    formData.valid_days = row.valid_days
    formData.device_limit = row.device_limit
    formData.price = row.price
    formData.priority = row.priority
    formData.is_active = row.is_active
    formData.server_ids = row.servers?.map(s => s.id) || []
    formData.description = row.description
  }
  const toAdd = () => {
    formVisible.value = true
    formData.id = 0
    formData.name = ''
    formData.valid_days = 30
    formData.device_limit = 10
    formData.price = 0
    formData.priority = 0
    formData.is_active = true
    formData.server_ids = []
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
