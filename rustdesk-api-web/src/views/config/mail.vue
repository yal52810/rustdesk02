<template>
  <div>
    <el-card shadow="hover">
      <el-alert
        title="邮件配置"
        type="info"
        :closable="false"
        show-icon
      >
        <template #default>
          用于注册成功通知、邮箱验证码找回密码等功能。密码留空时保持当前已保存的授权码不变。
        </template>
      </el-alert>

      <el-form :model="form" label-width="130px" style="margin-top: 20px; max-width: 760px">
        <el-form-item label="SMTP 主机" required>
          <el-input v-model="form.host" placeholder="例如：smtp.163.com" />
        </el-form-item>
        <el-form-item label="SMTP 端口" required>
          <el-input-number v-model="form.port" :min="1" :max="65535" />
        </el-form-item>
        <el-form-item label="登录账号">
          <el-input v-model="form.username" placeholder="例如：xxx@163.com" />
        </el-form-item>
        <el-form-item label="授权码/密码">
          <el-input
            v-model="form.password"
            type="password"
            show-password
            placeholder="留空保持当前密码"
          />
          <div v-if="form.password_set" class="form-tip">当前已保存授权码；不修改可保持为空。</div>
        </el-form-item>
        <el-form-item label="发信邮箱" required>
          <el-input v-model="form.from" placeholder="例如：xxx@163.com" />
        </el-form-item>
        <el-form-item label="发信人名称">
          <el-input v-model="form.from_name" placeholder="例如：RustDesk" />
        </el-form-item>
        <el-form-item label="SSL 直连">
          <el-switch v-model="form.use_ssl" />
        </el-form-item>
        <el-form-item label="跳过证书校验">
          <el-switch v-model="form.skip_verify" />
        </el-form-item>
        <el-form-item label="测试收件邮箱">
          <el-input v-model="form.test_to" placeholder="留空默认发到发信邮箱" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSave" :loading="saving">保存配置</el-button>
          <el-button @click="handleTest" :loading="testing">发送测试邮件</el-button>
          <el-button @click="loadConfig">重新加载</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { reactive, onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { mail, updateMail, testMail } from '@/api/config'

const defaultForm = () => ({
  host: '',
  port: 465,
  username: '',
  password: '',
  from: '',
  from_name: 'RustDesk',
  use_ssl: true,
  skip_verify: false,
  password_set: false,
  test_to: '',
})

const form = reactive(defaultForm())
const saving = ref(false)
const testing = ref(false)

const assignForm = (data = {}) => {
  Object.assign(form, defaultForm(), data, { password: '', test_to: '' })
}

const loadConfig = async () => {
  const res = await mail().catch(() => false)
  if (res) {
    assignForm(res.data || {})
  }
}

const handleSave = async () => {
  saving.value = true
  const payload = {
    host: form.host.trim(),
    port: form.port,
    username: form.username.trim(),
    password: form.password,
    from: form.from.trim(),
    from_name: form.from_name.trim(),
    use_ssl: form.use_ssl,
    skip_verify: form.skip_verify,
  }
  const res = await updateMail(payload).catch(() => false)
  saving.value = false
  if (res) {
    ElMessage.success('邮件配置已保存')
    assignForm(res.data || {})
  }
}

const handleTest = async () => {
  testing.value = true
  const payload = {
    host: form.host.trim(),
    port: form.port,
    username: form.username.trim(),
    password: form.password,
    from: form.from.trim(),
    from_name: form.from_name.trim(),
    use_ssl: form.use_ssl,
    skip_verify: form.skip_verify,
    to: form.test_to.trim(),
  }
  const res = await testMail(payload).catch(() => false)
  testing.value = false
  if (res) {
    ElMessage.success(res.data?.message || 'Test email sent')
  }
}

onMounted(loadConfig)
</script>

<style scoped lang="scss">
.form-tip {
  color: #909399;
  font-size: 12px;
  line-height: 1.6;
  margin-top: 6px;
}
</style>
