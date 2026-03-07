import { ref, onMounted, reactive, watch } from 'vue'
import { create, detail, update, remove } from '@/api/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRouter } from 'vue-router'
import { list as groups } from '@/api/group'
import { list as packages } from '@/api/package'
import { T } from '@/utils/i18n'

export function useGetDetail (id) {
  let item = ref({})  //保留原始值
  let form = ref({
    first_login_at: null,
    valid_days: 365, // Default value
    device_limit: 10, // Default value
    package_id: null,
  })
  const groupsList = ref([])
  const packagesList = ref([])

  const getDetail = async (id) => {
    const res = await detail(id)
    item.value = { ...res.data }
    form.value = { ...res.data }
  }
  if (id > 0) {
    onMounted(_ => {getDetail(id)})
  }

  const getGroups = async () => {
    const res = await groups({ page_size: 9999 }).catch(_ => false)
    if (res) {
      groupsList.value = res.data.list
    }
  }

  const getPackages = async () => {
    const res = await packages({ page_size: 9999 }).catch(_ => false)
    if (res) {
      packagesList.value = res.data?.list || []
    }
  }

  // 监听套餐选择变化
  watch(() => form.value.package_id, (newVal) => {
    if (newVal) {
      const selectedPackage = packagesList.value.find(p => p.id === newVal)
      if (selectedPackage) {
        form.value.valid_days = selectedPackage.valid_days
        form.value.device_limit = selectedPackage.device_limit
      }
    }
  })

  onMounted(getGroups)
  onMounted(getPackages)

  return {
    form,
    item,
    getDetail,
    groupsList,
    packagesList,
  }
}

export function useSubmit (form, id) {
  const root = ref(null)
  const router = useRouter()
  const rules = reactive({
    username: [{ required: true, message: T('ParamRequired', { param: T('Username') }) }],
    // email: [{ required: true, message: T('ParamRequired', { param: T('Email') }) }],
    group_id: [{ required: true, message: T('ParamRequired', { param: T('Group') }) }],
    // nickname: [{ required: true, message: '昵称是必须的' }],
    status: [{ required: true, message: T('ParamRequired', { param: T('Status') }) }],
  })

  const validate = async () => {
    const res = await root.value.validate().catch(err => false)
    return res
  }

  const submitCreate = async () => {
    const res = await create(form.value).catch(_ => false)
    return res.code === 0
  }

  const submitUpdate = async () => {
    const res = await update(form.value).catch(_ => false)
    return res.code === 0
  }
  const submitFunc = id > 0 ? submitUpdate : submitCreate

  const submit = async () => {
    const v = await validate()
    if (!v) {
      return
    }

    const res = await submitFunc()
    if (res) {
      ElMessage.success(T('OperationSuccess'))
      router.back()
    }
  }

  const cancel = () => {
    router.back()
  }

  return {
    root,
    rules,
    validate,
    submit,
    cancel,
  }
}


