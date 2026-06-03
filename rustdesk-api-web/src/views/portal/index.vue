<template>
  <div class="portal-shell">
    <div class="portal-bg"></div>
    <div class="portal-page">

      <!-- ======================================== -->
      <!-- Hero Header -->
      <!-- ======================================== -->
      <section class="hero-card">
        <div class="hero-copy">
          <div class="hero-brand">
            <img src="@/assets/logo.png" alt="logo" class="hero-logo" />
            <div>
              <div class="hero-title">RustDesk 用户中心</div>
            </div>
          </div>
        </div>
        <div class="hero-status" v-if="isLoggedIn && userInfo">
          <div class="status-label">当前账号</div>
          <div class="status-value">{{ userInfo.username || userStore.username }}</div>
          <div class="status-meta">
            <el-tag :type="packageTypeTag" size="small" effect="dark">{{ packageLabel }}</el-tag>
          </div>
          <div class="status-meta expire-info" v-if="userInfo.expired_at">
            到期：<strong>{{ userInfo.expired_at }}</strong>
            <el-tag :type="expireStatusTag" size="small" effect="plain" style="margin-left: 8px">{{ expireStatusText }}</el-tag>
          </div>
          <div class="status-meta" v-else>未开通套餐</div>
          <div class="hero-actions">
            <el-button v-if="isAdmin" type="primary" size="small" @click="goAdmin">管理后台</el-button>
            <el-button text class="logout-btn" @click="logout">退出登录</el-button>
          </div>
        </div>
      </section>

      <!-- ======================================== -->
      <!-- Logged-Out: Auth + Packages Preview -->
      <!-- ======================================== -->
      <div class="portal-grid" v-if="!isLoggedIn">
        <!-- Auth Card -->
        <section class="auth-card">
          <div class="card-title">账号中心</div>
          <div class="card-desc">登录或注册账号，管理套餐与客户端配置</div>
          <div class="auth-tabs">
            <button v-for="tab in filteredAuthTabs" :key="tab.key" class="auth-tab"
              :class="{ active: activeAuthTab === tab.key }" @click="switchAuthTab(tab.key)">
              {{ tab.label }}
            </button>
          </div>

          <!-- Login Form -->
          <el-form v-if="activeAuthTab === 'login'" class="portal-form" label-position="top" @submit.prevent="submitLogin">
            <el-form-item label="用户名">
              <el-input v-model="loginForm.username" placeholder="请输入账号" />
            </el-form-item>
            <el-form-item label="密码">
              <el-input v-model="loginForm.password" type="password" show-password placeholder="请输入密码" @keyup.enter="submitLogin" />
            </el-form-item>
            <el-form-item v-if="captchaCode" label="验证码">
              <div class="captcha-field">
                <el-input v-model="loginForm.captcha" placeholder="请输入验证码" />
                <img :src="captchaCode.b64" alt="captcha" class="captcha-image" @click="loadCaptcha" />
              </div>
            </el-form-item>
            <el-button type="primary" class="action-btn" @click="submitLogin" :loading="authLoading">登录</el-button>
          </el-form>

          <!-- Register Form -->
          <el-form v-else-if="activeAuthTab === 'register'" class="portal-form" label-position="top" @submit.prevent="submitRegister">
            <el-form-item label="用户名">
              <el-input v-model="registerForm.username" placeholder="2-18 位账号" />
            </el-form-item>
            <el-form-item label="邮箱">
              <el-input v-model="registerForm.email" placeholder="可选，用于找回密码和接收通知" />
            </el-form-item>
            <el-form-item label="密码">
              <el-input v-model="registerForm.password" type="password" show-password placeholder="6-18 位密码" />
            </el-form-item>
            <el-form-item label="确认密码">
              <el-input v-model="registerForm.confirmPassword" type="password" show-password placeholder="再次输入密码" />
            </el-form-item>
            <el-button type="primary" class="action-btn" @click="submitRegister" :loading="authLoading">注册</el-button>
          </el-form>

          <!-- Reset Password Form -->
          <el-form v-else class="portal-form" label-position="top" @submit.prevent="submitReset">
            <el-form-item label="邮箱">
              <el-input v-model="resetForm.email" placeholder="请输入注册邮箱" />
            </el-form-item>
            <el-form-item label="验证码">
              <div class="inline-field">
                <el-input v-model="resetForm.code" maxlength="6" placeholder="6 位验证码" />
                <el-button :disabled="codeCountdown > 0 || resetLoading" @click="sendResetCode">
                  {{ codeCountdown > 0 ? `${codeCountdown}s` : '发送验证码' }}
                </el-button>
              </div>
            </el-form-item>
            <el-form-item label="新密码">
              <el-input v-model="resetForm.newPassword" type="password" show-password placeholder="输入新密码" @keyup.enter="submitReset" />
            </el-form-item>
            <el-button type="primary" class="action-btn" @click="submitReset" :loading="resetLoading">重置密码</el-button>
          </el-form>
        </section>

        <!-- Public Packages + Server Preview -->
        <section class="preview-card">
          <div class="card-title">套餐预览</div>
          <div class="card-desc">以下为可购买套餐。购买后获得激活码，注册时填写或登录后激活。</div>

          <div v-if="packages.length > 0" class="pkg-grid">
            <div v-for="pkg in packages" :key="pkg.id" class="pkg-item">
              <div class="pkg-name">{{ pkg.name }}</div>
              <div class="pkg-meta">
                <span>{{ pkg.valid_days }} 天</span>
                <span>{{ pkg.device_limit }} 设备</span>
              </div>
              <div class="pkg-price" v-if="pkg.price > 0">¥{{ pkg.price }}</div>
              <div class="pkg-price free" v-else>免费</div>
              <div class="pkg-desc">{{ pkg.description || '标准套餐' }}</div>
            </div>
          </div>
          <el-empty v-else description="暂无可购买套餐" />

          <div class="card-title" style="margin-top: 28px">公开线路</div>
          <div class="card-desc">所有线路互通，共享账号体系和 ID 服务。</div>
          <el-empty v-if="publicServers.length === 0" description="暂无可用线路" />
          <div v-else class="preview-server-list">
            <button v-for="server in publicServers" :key="server.id"
              class="server-pill" :class="{ active: selectedPreviewServerId === server.id }"
              @click="selectedPreviewServerId = server.id">
              {{ server.name }}
            </button>
          </div>
          <div v-if="previewServer" class="config-grid compact">
            <div class="config-item">
              <div class="config-key">ID 服务器</div>
              <div class="config-val">{{ previewServer.id_server || '-' }}</div>
            </div>
            <div class="config-item">
              <div class="config-key">中继服务器</div>
              <div class="config-val">{{ previewServer.relay_server || '-' }}</div>
            </div>
            <div class="config-item">
              <div class="config-key">状态</div>
              <div class="config-val">
                <span class="line-status" :class="previewServer.is_online ? 'online' : 'offline'">
                  {{ previewServer.is_online ? '在线' : '离线' }}
                </span>
              </div>
            </div>
            <div class="config-item">
              <div class="config-key">API</div>
              <div class="config-val">{{ previewServer.api_server || currentOrigin }}</div>
            </div>
          </div>
        </section>
      </div>

      <!-- ======================================== -->
      <!-- Logged-In: Dashboard -->
      <!-- ======================================== -->
      <div v-else class="portal-grid logged-grid">
        <!-- Row 1: Account Summary + Plan Details -->
        <div class="logged-row">
          <section class="summary-card">
            <div class="card-title">账号概览</div>
            <div class="summary-grid">
              <div class="summary-item">
                <div class="summary-label">账号</div>
                <div class="summary-value">{{ userInfo.username || '-' }}</div>
              </div>
              <div class="summary-item">
                <div class="summary-label">邮箱</div>
                <div class="summary-value">{{ userInfo.email || '未绑定' }}</div>
              </div>
              <div class="summary-item">
                <div class="summary-label">设备 / 传输</div>
                <div class="summary-value">{{ userInfo.device_limit || 0 }} 台 / {{ transferLimitText }}</div>
              </div>
            </div>
          </section>

          <section class="plan-card">
            <div class="card-title">当前套餐</div>
            <div v-if="!userInfo.package && !userInfo.package_id" class="plan-empty">
              <span>尚未开通套餐</span>
              <span class="plan-hint">注册时可填写激活码，或登录后在下方激活码处激活</span>
            </div>
            <div v-else class="plan-detail">
              <div class="plan-name">{{ packageLabel }}</div>
              <div class="plan-days">
                <div class="days-badge">
                  <span class="days-num">{{ daysRemaining }}</span>
                  <span class="days-unit">天</span>
                </div>
                <div class="days-bar">
                  <div class="days-bar-fill" :style="{ width: daysPercent + '%' }"></div>
                </div>
              </div>
              <div class="plan-meta">
                <span>到期时间：{{ userInfo.expired_at || '永久' }}</span>
                <span>设备限制：{{ userInfo.device_limit || 0 }} 台</span>
                <span v-if="userInfo.package">传输上限：{{ userInfo.package.file_transfer_limit_mb || 100 }} MB</span>
              </div>
              <el-tag :type="daysRemaining > 30 ? 'success' : daysRemaining > 7 ? 'warning' : 'danger'" size="small" effect="dark">
                {{ daysRemaining > 30 ? '正常' : daysRemaining > 7 ? '即将到期' : daysRemaining > 0 ? '紧急续费' : '已过期' }}
              </el-tag>
            </div>
          </section>
        </div>

        <!-- Row 2: Purchase/Renew + Activation -->
        <div class="logged-row">
          <section class="purchase-card">
            <div class="card-title">购买 / 续费套餐</div>
            <div class="card-desc">选择套餐后联系管理员购买激活码，在右侧激活码框激活生效</div>
            <div v-if="packages.length > 0" class="pkg-grid">
              <div v-for="pkg in packages" :key="pkg.id"
                class="pkg-item" :class="{ current: userInfo.package_id === pkg.id }">
                <div class="pkg-name">
                  {{ pkg.name }}
                  <el-tag v-if="userInfo.package_id === pkg.id" size="small" type="success" effect="dark">当前</el-tag>
                </div>
                <div class="pkg-meta">
                  <span>{{ pkg.valid_days }} 天</span>
                  <span>{{ pkg.device_limit }} 设备</span>
                  <span>{{ pkg.file_transfer_limit_mb || 100 }} MB 传输</span>
                </div>
                <div class="pkg-price" v-if="pkg.price > 0">¥{{ pkg.price }}</div>
                <div class="pkg-price free" v-else>免费</div>
                <div class="pkg-desc">{{ pkg.description || '标准套餐，包含基础功能和线路接入' }}</div>
              </div>
            </div>
            <el-empty v-else description="暂无可购买套餐" />
          </section>

          <section class="redeem-card">
            <div class="card-title">激活码</div>
            <div class="card-desc">输入激活码立即开通或续费套餐</div>
            <el-input v-model="redeemForm.code" placeholder="请输入激活码" size="large" @keyup.enter="submitRedeem" />
            <el-button type="primary" class="action-btn" @click="submitRedeem" :loading="redeemLoading">立即激活</el-button>
            <el-divider style="margin: 16px 0" />
            <div class="card-title" style="font-size: 15px">修改密码</div>
            <el-form class="mini-form" label-position="top" @submit.prevent="submitChangePassword">
              <el-form-item label="旧密码">
                <el-input v-model="changePwdForm.oldPassword" type="password" show-password placeholder="输入旧密码" />
              </el-form-item>
              <el-form-item label="新密码">
                <el-input v-model="changePwdForm.newPassword" type="password" show-password placeholder="输入新密码" />
              </el-form-item>
              <el-form-item label="确认密码">
                <el-input v-model="changePwdForm.confirmPassword" type="password" show-password placeholder="再次输入" />
              </el-form-item>
              <el-button type="primary" class="action-btn" @click="submitChangePassword" :loading="changingPwd">保存新密码</el-button>
            </el-form>
          </section>
        </div>

        <!-- Row 3: Server Config + Client Config Generation -->
        <div class="logged-row">
          <section class="server-card">
            <div class="card-head">
              <div>
                <div class="card-title">服务器配置</div>
                <div class="card-desc">选择线路后复制配置或生成客户端配置文件</div>
              </div>
              <el-button text @click="refreshPortal">刷新</el-button>
            </div>

            <div v-if="availableServers.length > 0" class="server-selector">
              <button v-for="server in availableServers" :key="server.id || server.name"
                class="server-option" :class="{ active: selectedServerId === (server.id || server.name) }"
                @click="selectedServerId = server.id || server.name">
                <span>{{ server.name }}</span>
                <small>{{ server.region || '默认' }} · {{ server.is_online ? '在线' : '离线' }}</small>
              </button>
            </div>
            <el-empty v-else description="暂无可用线路" />

            <div v-if="selectedServer" class="config-wrap">
              <div class="config-grid">
                <div class="config-item">
                  <div class="config-key">ID 服务器</div>
                  <div class="config-val">{{ selectedServer.id_server || '-' }}</div>
                  <el-button text size="small" @click="copyText(selectedServer.id_server)">复制</el-button>
                </div>
                <div class="config-item">
                  <div class="config-key">中继服务器</div>
                  <div class="config-val">{{ selectedServer.relay_server || '-' }}</div>
                  <el-button text size="small" @click="copyText(selectedServer.relay_server)">复制</el-button>
                </div>
                <div class="config-item">
                  <div class="config-key">API 服务器</div>
                  <div class="config-val">{{ selectedServer.api_server || currentOrigin }}</div>
                  <el-button text size="small" @click="copyText(selectedServer.api_server || currentOrigin)">复制</el-button>
                </div>
                <div class="config-item">
                  <div class="config-key">Key</div>
                  <div class="config-val mono">{{ selectedServer.key || '-' }}</div>
                  <el-button text size="small" @click="copyText(selectedServer.key)">复制</el-button>
                </div>
                <div class="config-item" v-if="selectedServer.ws_host">
                  <div class="config-key">WebSocket 地址</div>
                  <div class="config-val mono">{{ selectedServer.ws_host }}</div>
                  <el-button text size="small" @click="copyText(selectedServer.ws_host)">复制</el-button>
                </div>
                <div class="config-item">
                  <div class="config-key">连接类型</div>
                  <div class="config-val">{{ selectedServer.support_wss ? '专业线路' : 'TCP 标准线路' }}</div>
                </div>
              </div>

              <el-divider>客户端配置导入</el-divider>
              <div class="client-config-section">
                <div class="config-hint">将以下配置字符串导入客户端，或复制 ID/中继信息到客户端手动填写。</div>
                <div class="config-textarea-wrap">
                  <el-input v-model="clientConfigStr" type="textarea" :rows="4" readonly resize="none"
                    class="config-textarea" />
                </div>
                <el-button type="primary" size="small" @click="generateClientConfig" :loading="configLoading">
                  生成配置
                </el-button>
                <el-button size="small" @click="copyText(clientConfigStr)" :disabled="!clientConfigStr">
                  复制配置
                </el-button>
                <el-button size="small" @click="downloadConfig" :disabled="!clientConfigStr">
                  下载配置
                </el-button>
              </div>
            </div>
          </section>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '@/store/user'
import { current, changeCurPwd } from '@/api/user'
import { sendPasswordResetCode, resetPassword, loginOptions, captcha } from '@/api/login'
import { listVipServers, listVipPackages, getClientConfig, redeemActivationCode, vipRegister } from '@/api/portal'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const currentOrigin = window.location.origin

const authTabs = [
  { key: 'login', label: '登录' },
  { key: 'register', label: '注册' },
  { key: 'forgot', label: '找回密码' },
]

const loginForm = reactive({ username: '', password: '', platform: 'web', captcha: '', captcha_id: '' })
const registerForm = reactive({ username: '', email: '', password: '', confirmPassword: '', activation_code: '' })
const resetForm = reactive({ email: '', code: '', newPassword: '' })
const redeemForm = reactive({ code: '' })
const changePwdForm = reactive({ oldPassword: '', newPassword: '', confirmPassword: '' })

const activeAuthTab = ref('login')
const userInfo = ref(null)
const publicServers = ref([])
const packages = ref([])
const selectedPreviewServerId = ref(null)
const selectedServerId = ref(null)
const codeCountdown = ref(0)
const resetLoading = ref(false)
const captchaCode = ref(null)
const allowRegister = ref(false)
const authLoading = ref(false)
const redeemLoading = ref(false)
const changingPwd = ref(false)
const configLoading = ref(false)
const clientConfigStr = ref('')
let codeTimer = null

const syncTabFromRoute = () => {
  if (route.path === '/register') { activeAuthTab.value = 'register'; return }
  if (route.path === '/forgot-password') { activeAuthTab.value = 'forgot'; return }
  activeAuthTab.value = route.query.tab || 'login'
}

const filteredAuthTabs = computed(() => {
  return allowRegister.value ? authTabs : authTabs.filter(tab => tab.key !== 'register')
})
const isLoggedIn = computed(() => !!userStore.token)
const isAdmin = computed(() => userStore.route_names.includes('*') || userStore.route_names.includes('UserList'))
const packageLabel = computed(() => {
  return userInfo.value?.package?.name || userInfo.value?.info?.package_name || '未开通'
})
const daysRemaining = computed(() => {
  if (!userInfo.value?.expired_at) return userInfo.value?.valid_days || 0
  const now = new Date()
  const expire = new Date(userInfo.value.expired_at)
  const diff = Math.ceil((expire - now) / (1000 * 60 * 60 * 24))
  return Math.max(0, diff)
})
const daysPercent = computed(() => {
  const total = userInfo.value?.package?.valid_days || 365
  return Math.min(100, Math.max(0, (daysRemaining.value / total) * 100))
})
const availableServers = computed(() => {
  const pkgServers = userInfo.value?.package?.servers || []
  return pkgServers.length > 0 ? pkgServers : publicServers.value
})
const selectedServer = computed(() => {
  return availableServers.value.find(s => (s.id || s.name) === selectedServerId.value) || availableServers.value[0] || null
})
const previewServer = computed(() => {
  return publicServers.value.find(s => s.id === selectedPreviewServerId.value) || publicServers.value[0] || null
})
const transferLimitText = computed(() => {
  const mb = userInfo.value?.package?.file_transfer_limit_mb || userInfo.value?.info?.file_transfer_limit_mb || 100
  return `${mb} MB`
})
const packageTypeTag = computed(() => {
  const name = (userInfo.value?.package?.name || userInfo.value?.info?.package_name || '').toLowerCase()
  if (name.includes('pro') || name.includes('专业') || name.includes('premium')) return 'success'
  if (name.includes('free') || name.includes('免费')) return 'info'
  return ''
})
const expireStatusTag = computed(() => {
  if (daysRemaining.value > 30) return 'success'
  if (daysRemaining.value > 7) return 'warning'
  return 'danger'
})
const expireStatusText = computed(() => {
  if (daysRemaining.value > 30) return '正常'
  if (daysRemaining.value > 7) return '即将到期'
  if (daysRemaining.value > 0) return '即将过期'
  return '已过期'
})

watch(() => route.fullPath, syncTabFromRoute)
watch(availableServers, (list) => {
  if (!list.length) { selectedServerId.value = null; return }
  const exists = list.some(s => (s.id || s.name) === selectedServerId.value)
  if (!exists) selectedServerId.value = list[0].id || list[0].name
}, { immediate: true })
watch(publicServers, (list) => {
  if (!list.length) { selectedPreviewServerId.value = null; return }
  const exists = list.some(s => s.id === selectedPreviewServerId.value)
  if (!exists) selectedPreviewServerId.value = list[0].id
}, { immediate: true })

const switchAuthTab = (tab) => {
  activeAuthTab.value = tab
  const path = tab === 'register' ? '/register' : tab === 'forgot' ? '/forgot-password' : '/login'
  router.replace({ path })
}

const refreshPortal = async () => {
  await Promise.all([loadPublicServers(), loadPackages()])
  if (isLoggedIn.value) await loadUserInfo()
}

const loadPublicServers = async () => {
  const res = await listVipServers().catch(() => false)
  if (res?.list) publicServers.value = res.list
}
const loadPackages = async () => {
  const res = await listVipPackages().catch(() => false)
  if (res?.list) packages.value = res.list
}
const loadLoginOptions = async () => {
  const res = await loginOptions().catch(() => false)
  if (!res?.data) return
  allowRegister.value = !!res.data.register
  if (!allowRegister.value && activeAuthTab.value === 'register') switchAuthTab('login')
  if (res.data.need_captcha) { await loadCaptcha() } else {
    captchaCode.value = null; loginForm.captcha = ''; loginForm.captcha_id = ''
  }
}
const loadCaptcha = async () => {
  const res = await captcha().catch(() => false)
  if (!res?.data?.captcha) return
  captchaCode.value = res.data.captcha
  loginForm.captcha_id = res.data.captcha.id
}
const loadUserInfo = async () => {
  const res = await current().catch(() => false)
  if (res?.data) userInfo.value = res.data
}

const submitLogin = async () => {
  if (!loginForm.username || !loginForm.password) { ElMessage.error('请输入账号和密码'); return }
  authLoading.value = true
  const res = await userStore.login(loginForm).catch(err => err)
  authLoading.value = false
  if (res?.code) { if (res.code === 110) await loadCaptcha(); return }
  if (!res) return
  ElMessage.success('登录成功')
  await loadUserInfo()
  if (isAdmin.value) { router.replace('/user/index'); return }
  router.replace(route.query.redirect || '/portal')
}

const submitRegister = async () => {
  if (!allowRegister.value) { ElMessage.error('暂未开放注册'); return }
  if (!registerForm.username || !registerForm.password) { ElMessage.error('请填写账号和密码'); return }
  if (registerForm.password !== registerForm.confirmPassword) { ElMessage.error('两次密码不一致'); return }
  authLoading.value = true
  const res = await vipRegister({
    username: registerForm.username,
    password: registerForm.password,
    confirm_password: registerForm.confirmPassword,
    activation_code: registerForm.activation_code,
  }).catch(() => null)
  authLoading.value = false
  if (!res) { ElMessage.error('注册失败，请检查网络连接或联系管理员'); return }
  if (res.error) { ElMessage.error(res.error); return }
  if (res.data) {
    userStore.saveUserData(res.data)
  }
  ElMessage.success('注册成功')
  await loadUserInfo()
  router.replace('/portal')
}

const startCodeCountdown = () => {
  codeCountdown.value = 60
  codeTimer = setInterval(() => { codeCountdown.value -= 1; if (codeCountdown.value <= 0) { clearInterval(codeTimer); codeTimer = null } }, 1000)
}
const sendResetCode = async () => {
  if (!resetForm.email) { ElMessage.error('请输入邮箱'); return }
  resetLoading.value = true
  const res = await sendPasswordResetCode({ email: resetForm.email }).catch(() => false)
  resetLoading.value = false
  if (!res) return
  ElMessage.success('验证码已发送')
  if (!codeTimer) startCodeCountdown()
}
const submitReset = async () => {
  if (!resetForm.email || !resetForm.code || !resetForm.newPassword) { ElMessage.error('请填写完整信息'); return }
  resetLoading.value = true
  const res = await resetPassword({ email: resetForm.email, code: resetForm.code, new_password: resetForm.newPassword }).catch(() => false)
  resetLoading.value = false
  if (!res) return
  ElMessage.success('密码已重置，请重新登录')
  resetForm.code = ''; resetForm.newPassword = ''
  switchAuthTab('login')
}

const submitRedeem = async () => {
  if (!redeemForm.code) { ElMessage.error('请输入激活码'); return }
  redeemLoading.value = true
  const res = await redeemActivationCode({ code: redeemForm.code }).catch(() => false)
  redeemLoading.value = false
  if (!res) return
  ElMessage.success(`激活成功，有效期增加 ${res.valid_days || 0} 天`)
  redeemForm.code = ''
  await loadUserInfo()
}

const submitChangePassword = async () => {
  if (!changePwdForm.oldPassword || !changePwdForm.newPassword || !changePwdForm.confirmPassword) { ElMessage.error('请填写完整信息'); return }
  if (changePwdForm.newPassword !== changePwdForm.confirmPassword) { ElMessage.error('两次新密码不一致'); return }
  const confirmed = await ElMessageBox.confirm('修改密码后需重新登录，继续？', '修改密码', {
    confirmButtonText: '确认', cancelButtonText: '取消', type: 'warning',
  }).catch(() => false)
  if (!confirmed) return
  changingPwd.value = true
  const res = await changeCurPwd({
    old_password: changePwdForm.oldPassword, new_password: changePwdForm.newPassword, confirm_password: changePwdForm.confirmPassword,
  }).catch(() => false)
  changingPwd.value = false
  if (!res) return
  ElMessage.success('密码修改成功，请重新登录')
  logout()
}

const generateClientConfig = async () => {
  configLoading.value = true
  const res = await getClientConfig().catch(() => false)
  configLoading.value = false
  if (res?.config?.config_str) {
    clientConfigStr.value = res.config.config_str
    ElMessage.success('配置已生成')
  } else {
    ElMessage.error('生成配置失败')
  }
}

const downloadConfig = () => {
  if (!clientConfigStr.value) { ElMessage.error('请先生成配置'); return }
  const blob = new Blob([clientConfigStr.value], { type: 'text/plain;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url; a.download = 'rustdesk-config.txt'
  document.body.appendChild(a); a.click()
  document.body.removeChild(a); URL.revokeObjectURL(url)
  ElMessage.success('配置已下载')
}

const copyText = async (value) => {
  if (!value) return
  let copied = false
  try { await navigator.clipboard.writeText(value); copied = true } catch (_) { }
  if (!copied) {
    const textarea = document.createElement('textarea')
    textarea.value = value; textarea.style.position = 'fixed'; textarea.style.opacity = '0'
    document.body.appendChild(textarea); textarea.focus(); textarea.select()
    document.execCommand('copy'); document.body.removeChild(textarea)
  }
  ElMessage.success('已复制')
}

const goAdmin = () => { router.replace('/user/index') }
const logout = () => {
  userStore.logout()
  userInfo.value = null
  router.replace('/login')
}

onMounted(async () => {
  syncTabFromRoute()
  await Promise.all([loadLoginOptions(), loadPublicServers(), loadPackages()])
  if (isLoggedIn.value) { await loadUserInfo(); router.replace('/portal') }
})
onBeforeUnmount(() => { if (codeTimer) { clearInterval(codeTimer); codeTimer = null } })
</script>

<style scoped lang="scss">
.portal-shell {
  position: relative;
  min-height: 100vh;
  overflow: hidden;
  background: linear-gradient(180deg, #eff6ff 0%, #f8fbff 42%, #ffffff 100%);
}

.portal-bg {
  position: absolute;
  inset: 0;
  background:
    radial-gradient(circle at top left, rgba(44, 120, 255, 0.15), transparent 30%),
    radial-gradient(circle at right 20%, rgba(114, 185, 255, 0.18), transparent 24%),
    linear-gradient(135deg, rgba(255, 255, 255, 0.9), rgba(239, 246, 255, 0.72));
}

.portal-page {
  position: relative;
  z-index: 1;
  max-width: 1200px;
  margin: 0 auto;
  padding: 36px 24px 56px;
}

// Cards
.hero-card, .auth-card, .preview-card, .summary-card, .plan-card, .purchase-card, .redeem-card, .server-card {
  background: rgba(255, 255, 255, 0.92);
  border: 1px solid rgba(58, 121, 255, 0.12);
  box-shadow: 0 18px 40px rgba(27, 76, 160, 0.08);
  backdrop-filter: blur(8px);
}

// Hero
.hero-card {
  display: flex;
  justify-content: space-between;
  gap: 24px;
  padding: 28px 32px;
  border-radius: 28px;
  margin-bottom: 24px;
}
.hero-brand { display: flex; align-items: center; gap: 18px; }
.hero-logo { width: 58px; height: 58px; }
.hero-title { font-size: 30px; font-weight: 700; color: #1d3d8f; letter-spacing: 0.02em; }
.hero-subtitle, .card-desc { color: #5f6f91; line-height: 1.7; }
.hero-status {
  min-width: 220px;
  padding: 22px 24px;
  border-radius: 22px;
  background: linear-gradient(180deg, #edf4ff, #f8fbff);
}
.status-label, .summary-label, .config-key { font-size: 13px; color: #6f84ab; }
.status-value { font-size: 24px; margin-top: 8px; color: #193b88; font-weight: 700; }
.status-meta { margin-top: 6px; color: #376df3; }
.expire-info strong { color: #e63946; }
.logout-btn { margin-top: 12px; }
.hero-actions { display: flex; gap: 10px; align-items: center; justify-content: flex-end; margin-top: 8px; }

// Grid
.portal-grid { display: grid; grid-template-columns: 1.05fr 0.95fr; gap: 24px; }
.logged-grid { display: flex; flex-direction: column; gap: 24px; }
.logged-row { display: grid; grid-template-columns: 1fr 1fr; gap: 24px; }
@media (max-width: 960px) {
  .portal-grid, .logged-row { grid-template-columns: 1fr; }
}

// Section cards
.auth-card, .preview-card, .summary-card, .plan-card, .purchase-card, .redeem-card, .server-card {
  border-radius: 26px;
  padding: 28px;
}
.card-title { font-size: 22px; font-weight: 700; color: #193b88; margin-bottom: 8px; }
.card-desc { margin-bottom: 16px; }

// Auth tabs
.auth-tabs { display: flex; gap: 10px; margin: 22px 0 20px; }
.auth-tab {
  padding: 11px 18px;
  border-radius: 999px;
  border: 0;
  cursor: pointer;
  transition: all 0.2s ease;
  background: #e8f1ff;
  color: #2a58b8;
  font-weight: 600;
}
.auth-tab.active {
  background: linear-gradient(135deg, #1d78ff, #1457eb);
  color: #fff;
  box-shadow: 0 10px 20px rgba(29, 120, 255, 0.2);
}

// Forms
.portal-form { margin-top: 8px; }
.action-btn { width: 100%; min-height: 44px; border-radius: 14px; font-weight: 600; }
.inline-field { display: grid; grid-template-columns: 1fr 128px; gap: 12px; width: 100%; }
.captcha-field { display: grid; grid-template-columns: 1fr 120px; gap: 12px; width: 100%; }
.captcha-image { width: 120px; height: 44px; border-radius: 14px; object-fit: cover; cursor: pointer; border: 1px solid rgba(58, 121, 255, 0.14); }

// Summary
.summary-grid { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 16px; margin-top: 8px; }
.summary-item {
  padding: 18px 20px;
  border-radius: 18px;
  background: linear-gradient(180deg, #f8fbff, #ffffff);
  border: 1px solid rgba(58, 121, 255, 0.1);
}
.summary-value { margin-top: 8px; font-size: 18px; color: #224587; font-weight: 600; }

// Plan detail
.plan-empty {
  display: flex; flex-direction: column; gap: 8px; padding: 24px 0;
  color: #8e99b3; text-align: center;
}
.plan-hint { font-size: 13px; color: #b0bcd0; }
.plan-name { font-size: 22px; font-weight: 700; color: #193b88; }
.plan-days { display: flex; align-items: center; gap: 16px; margin: 16px 0; }
.days-badge {
  padding: 12px 20px;
  border-radius: 18px;
  background: linear-gradient(135deg, #e8f1ff, #d4e4ff);
  min-width: 90px;
  text-align: center;
}
.days-num { font-size: 36px; font-weight: 800; color: #1d3d8f; }
.days-unit { font-size: 14px; color: #5f78ad; margin-left: 4px; }
.days-bar {
  flex: 1;
  height: 12px;
  border-radius: 999px;
  background: #e8f1ff;
  overflow: hidden;
}
.days-bar-fill {
  height: 100%;
  border-radius: 999px;
  background: linear-gradient(90deg, #1d78ff, #3d93ff);
  transition: width 0.5s ease;
}
.plan-meta {
  display: flex; flex-wrap: wrap; gap: 16px; margin-bottom: 12px;
  color: #5f7195; font-size: 14px;
}

// Packages
.pkg-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); gap: 12px; }
.pkg-item {
  padding: 18px;
  border-radius: 18px;
  background: linear-gradient(180deg, #f8fbff, #ffffff);
  border: 1px solid rgba(58, 121, 255, 0.1);
  transition: all 0.2s ease;
}
.pkg-item:hover { border-color: #1d78ff; box-shadow: 0 8px 24px rgba(29, 120, 255, 0.1); }
.pkg-item.current { border-color: #1d78ff; background: linear-gradient(180deg, #edf4ff, #f8fbff); }
.pkg-name { font-size: 17px; font-weight: 700; color: #193b88; display: flex; align-items: center; justify-content: space-between; gap: 8px; }
.pkg-meta { display: flex; gap: 12px; margin-top: 8px; color: #5f7195; font-size: 13px; }
.pkg-price { margin-top: 10px; font-size: 22px; font-weight: 800; color: #1d78ff; }
.pkg-price.free { color: #67c23a; }
.pkg-desc { margin-top: 6px; font-size: 13px; color: #8e99b3; line-height: 1.6; }

// Server selector
.preview-server-list, .server-selector {
  display: flex; flex-wrap: wrap; gap: 12px; margin: 16px 0;
}
.server-pill, .server-option {
  display: inline-flex; align-items: center; justify-content: center;
  gap: 6px; padding: 12px 16px; border-radius: 16px;
  background: #eef5ff; color: #244c99; font-weight: 600;
  border: 0; cursor: pointer; transition: all 0.2s ease;
}
.server-option { align-items: flex-start; flex-direction: column; min-width: 160px; text-align: left; }
.server-option small { color: #7690bd; }
.server-pill.active, .server-option.active {
  background: linear-gradient(135deg, #1c7bff, #1557ea);
  color: #fff; box-shadow: 0 10px 18px rgba(21, 87, 234, 0.22);
}
.server-option.active small { color: rgba(255, 255, 255, 0.84); }

// Config
.card-head { display: flex; justify-content: space-between; gap: 16px; align-items: flex-start; }
.config-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 12px; }
.config-grid.compact { margin-top: 16px; }
.config-item {
  padding: 14px 16px; border-radius: 14px;
  background: linear-gradient(180deg, #f8fbff, #ffffff);
  border: 1px solid rgba(58, 121, 255, 0.08);
}
.config-val { margin-top: 6px; color: #224587; font-weight: 600; word-break: break-all; font-size: 14px; }
.config-val.mono { font-family: Consolas, Monaco, monospace; font-size: 13px; }

.line-status {
  display: inline-flex; align-items: center;
  padding: 4px 10px; border-radius: 999px; font-size: 12px; font-weight: 700;
}
.line-status.online { color: #11783c; background: rgba(24, 181, 94, 0.12); }
.line-status.offline { color: #c33d42; background: rgba(240, 72, 72, 0.12); }

// Client config
.client-config-section { display: flex; flex-direction: column; gap: 12px; }
.config-hint { color: #5f7195; font-size: 13px; line-height: 1.6; }
.config-textarea-wrap { width: 100%; }
.config-textarea :deep(textarea) { font-family: Consolas, Monaco, monospace; font-size: 13px; border-radius: 14px; }

// Redeem card mini-form
.mini-form { margin-top: 4px; }

// Element overrides
:deep(.el-input__wrapper) { min-height: 44px; border-radius: 14px; box-shadow: 0 0 0 1px rgba(58, 121, 255, 0.12) inset; }
:deep(.el-button--primary) { background: linear-gradient(135deg, #1c7bff, #1557ea); border-color: #1c7bff; }

@media (max-width: 960px) {
  .hero-card { flex-direction: column; }
  .portal-page { padding: 20px 16px 36px; }
  .summary-grid { grid-template-columns: 1fr; }
}
</style>
