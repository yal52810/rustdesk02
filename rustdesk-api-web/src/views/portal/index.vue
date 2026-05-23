<template>
  <div class="portal-shell">
    <div class="portal-bg"></div>
    <div class="portal-page">
      <section class="hero-card">
        <div class="hero-copy">
          <div class="hero-brand">
            <img src="@/assets/logo.png" alt="logo" class="hero-logo" />
            <div>
              <div class="hero-title">RustDesk 用户服务中心</div>
              <div class="hero-subtitle">注册、查询、激活和手动服务器配置统一在这里完成。</div>
            </div>
          </div>
          <div class="hero-note">
            蓝白风格与客户端保持一致。无法使用定制客户端时，可在这里查看可用线路和连接参数。
          </div>
        </div>
        <div class="hero-status" v-if="isLoggedIn && userInfo">
          <div class="status-label">当前账号</div>
          <div class="status-value">{{ userInfo.username || userInfo.name || userStore.username }}</div>
          <div class="status-meta">{{ packageLabel }}</div>
          <el-button text class="logout-btn" @click="logout">退出登录</el-button>
        </div>
      </section>

      <div class="portal-grid" v-if="!isLoggedIn">
        <section class="auth-card">
          <div class="card-title">账号入口</div>
          <div class="card-desc">首次使用先注册。已有账号直接登录，忘记密码可通过邮箱验证码重置。</div>
          <div class="auth-tabs">
            <button
              v-for="tab in filteredAuthTabs"
              :key="tab.key"
              class="auth-tab"
              :class="{ active: activeAuthTab === tab.key }"
              @click="switchAuthTab(tab.key)"
            >
              {{ tab.label }}
            </button>
          </div>

          <el-form v-if="activeAuthTab === 'login'" class="portal-form" label-position="top">
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
            <el-button type="primary" class="action-btn" @click="submitLogin">登录</el-button>
          </el-form>

          <el-form v-else-if="activeAuthTab === 'register'" class="portal-form" label-position="top">
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
              <el-input v-model="registerForm.confirmPassword" type="password" show-password placeholder="再次输入密码" @keyup.enter="submitRegister" />
            </el-form-item>
            <el-button type="primary" class="action-btn" @click="submitRegister">注册</el-button>
          </el-form>

          <el-form v-else class="portal-form" label-position="top">
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
            <el-button type="primary" class="action-btn" @click="submitReset">重置密码</el-button>
          </el-form>
        </section>

        <section class="preview-card">
          <div class="card-title">使用说明</div>
          <div class="info-list">
            <div class="info-item">
              <span class="info-dot"></span>
              <span>普通版、专业版由套餐决定，线路只是接入方式，不是第三种套餐。</span>
            </div>
            <div class="info-item">
              <span class="info-dot"></span>
              <span>普通网络优先使用标准线路，公司或校园网络受限时再切换对应线路。</span>
            </div>
            <div class="info-item">
              <span class="info-dot"></span>
              <span>登录后可以查看当前套餐、到期时间、激活入口和服务器配置。</span>
            </div>
          </div>

            <div class="preview-panel">
              <div class="preview-label">公开线路预览</div>
              <div class="interop-note">
                所有线路属于同一套互通架构，共享账号体系、ID 服务和 Key，只区分中继入口与受限网络接入方式。
              </div>
              <el-empty v-if="publicServers.length === 0" description="暂无可展示线路" />
              <div v-else class="preview-server-list">
                <button
                v-for="server in publicServers"
                :key="server.id"
                class="server-pill"
                :class="{ active: selectedPreviewServerId === server.id }"
                @click="selectedPreviewServerId = server.id"
              >
                {{ server.name }}
              </button>
            </div>
            <div v-if="previewServer" class="config-grid compact">
                <div class="config-item">
                  <div class="config-key">ID</div>
                  <div class="config-val">{{ previewServer.id_server || '-' }}</div>
                </div>
              <div class="config-item">
                <div class="config-key">互通组</div>
                <div class="config-val">{{ previewServer.topology_group || 'default' }}</div>
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
                  <div class="config-key">中继</div>
                  <div class="config-val">{{ previewServer.relay_server || '-' }}</div>
              </div>
              <div class="config-item">
                <div class="config-key">API</div>
                <div class="config-val">{{ previewServer.api_server || currentOrigin }}</div>
              </div>
              <div class="config-item">
                <div class="config-key">Key</div>
                <div class="config-val mono">{{ previewServer.key || '-' }}</div>
              </div>
            </div>
          </div>
        </section>
      </div>

      <div v-else class="portal-grid logged-grid">
        <section class="summary-card">
          <div class="card-title">账号信息</div>
          <div class="summary-grid">
            <div class="summary-item">
              <div class="summary-label">账号</div>
              <div class="summary-value">{{ userInfo.username || userInfo.name || '-' }}</div>
            </div>
            <div class="summary-item">
              <div class="summary-label">邮箱</div>
              <div class="summary-value">{{ userInfo.email || '未绑定' }}</div>
            </div>
            <div class="summary-item">
              <div class="summary-label">当前套餐</div>
              <div class="summary-value">{{ packageLabel }}</div>
            </div>
            <div class="summary-item">
              <div class="summary-label">到期时间</div>
              <div class="summary-value">{{ userInfo.expired_at || (userInfo.valid_days > 0 ? `${userInfo.valid_days} 天` : '未开通') }}</div>
            </div>
            <div class="summary-item">
              <div class="summary-label">设备限制</div>
              <div class="summary-value">{{ userInfo.device_limit || 0 }}</div>
            </div>
            <div class="summary-item">
              <div class="summary-label">传输上限</div>
              <div class="summary-value">{{ transferLimitText }}</div>
            </div>
          </div>
        </section>

        <section class="ops-card">
          <div class="card-title">账号操作</div>
          <div class="ops-grid">
            <div class="mini-card">
              <div class="mini-title">激活套餐</div>
              <div class="mini-desc">输入激活码后立即生效，刷新当前账号套餐和线路。</div>
              <el-input v-model="redeemForm.code" placeholder="请输入激活码" />
              <el-button type="primary" class="action-btn" @click="submitRedeem">立即激活</el-button>
            </div>
            <div class="mini-card">
              <div class="mini-title">修改密码</div>
              <div class="mini-desc">已登录账号直接修改密码，修改完成后需要重新登录。</div>
              <el-form class="mini-form" label-position="top">
                <el-form-item label="旧密码">
                  <el-input v-model="changePwdForm.oldPassword" type="password" show-password placeholder="请输入旧密码" />
                </el-form-item>
                <el-form-item label="新密码">
                  <el-input v-model="changePwdForm.newPassword" type="password" show-password placeholder="请输入新密码" />
                </el-form-item>
                <el-form-item label="确认密码">
                  <el-input v-model="changePwdForm.confirmPassword" type="password" show-password placeholder="再次输入新密码" />
                </el-form-item>
                <el-button type="primary" class="action-btn" @click="submitChangePassword">保存新密码</el-button>
              </el-form>
            </div>
          </div>
        </section>

        <section class="server-card">
          <div class="card-head">
            <div>
              <div class="card-title">服务器配置</div>
              <div class="card-desc">用于普通客户端手动填写。这里只展示当前账号可用线路。</div>
            </div>
            <el-button text @click="refreshPortal">刷新</el-button>
          </div>

          <div v-if="availableServers.length > 0" class="server-selector">
            <button
              v-for="server in availableServers"
              :key="server.id || server.name"
              class="server-option"
              :class="{ active: selectedServerId === (server.id || server.name) }"
              @click="selectedServerId = server.id || server.name"
            >
              <span>{{ server.name }}</span>
              <small>{{ server.region || '默认线路' }}</small>
            </button>
          </div>
          <el-empty v-else description="当前账号暂无可用线路" />

          <div v-if="selectedServer" class="config-wrap">
            <div class="config-grid">
              <div class="config-item">
                <div class="config-key">互通组</div>
                <div class="config-val">{{ selectedServer.topology_group || 'default' }}</div>
              </div>
              <div class="config-item">
                <div class="config-key">线路状态</div>
                <div class="config-val">
                  <span class="line-status" :class="selectedServer.is_online ? 'online' : 'offline'">
                    {{ selectedServer.is_online ? '在线' : '离线' }}
                  </span>
                </div>
              </div>
              <div class="config-item">
                <div class="config-key">ID 服务器</div>
                <div class="config-val">{{ selectedServer.id_server || '-' }}</div>
                <el-button text @click="copyText(selectedServer.id_server)">复制</el-button>
              </div>
              <div class="config-item">
                <div class="config-key">中继服务器</div>
                <div class="config-val">{{ selectedServer.relay_server || '-' }}</div>
                <el-button text @click="copyText(selectedServer.relay_server)">复制</el-button>
              </div>
              <div class="config-item">
                <div class="config-key">API 服务器</div>
                <div class="config-val">{{ selectedServer.api_server || currentOrigin }}</div>
                <el-button text @click="copyText(selectedServer.api_server || currentOrigin)">复制</el-button>
              </div>
              <div class="config-item">
                <div class="config-key">Key</div>
                <div class="config-val mono">{{ selectedServer.key || '-' }}</div>
                <el-button text @click="copyText(selectedServer.key)">复制</el-button>
              </div>
            </div>
            <div class="server-hint">
              <span>线路说明：</span>
              <span>{{ selectedServer.description || '标准线路适合普通网络，公司/校园受限网络可切换支持 WebSocket 的线路。所有线路属于同一套互通架构，可直接远控同账号体系下的设备。' }}</span>
            </div>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { computed, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
  import { useRouter, useRoute } from 'vue-router'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { useUserStore } from '@/store/user'
  import { register as registerUser, current, changeCurPwd } from '@/api/user'
  import { sendPasswordResetCode, resetPassword, loginOptions, captcha } from '@/api/login'
  import { listVipServers, redeemActivationCode } from '@/api/portal'

  const router = useRouter()
  const route = useRoute()
  const userStore = useUserStore()
  const currentOrigin = window.location.origin

  const authTabs = [
    { key: 'login', label: '登录' },
    { key: 'register', label: '注册' },
    { key: 'forgot', label: '找回密码' },
  ]

  const loginForm = reactive({
    username: '',
    password: '',
    platform: 'web',
    captcha: '',
    captcha_id: '',
  })
  const registerForm = reactive({
    username: '',
    email: '',
    password: '',
    confirmPassword: '',
  })
  const resetForm = reactive({
    email: '',
    code: '',
    newPassword: '',
  })
  const redeemForm = reactive({
    code: '',
  })
  const changePwdForm = reactive({
    oldPassword: '',
    newPassword: '',
    confirmPassword: '',
  })

  const activeAuthTab = ref('login')
  const userInfo = ref(null)
  const publicServers = ref([])
  const selectedPreviewServerId = ref(null)
  const selectedServerId = ref(null)
  const codeCountdown = ref(0)
  const resetLoading = ref(false)
  const captchaCode = ref(null)
  const allowRegister = ref(true)
  let codeTimer = null

  const syncTabFromRoute = () => {
    if (route.path === '/register') {
      activeAuthTab.value = 'register'
      return
    }
    if (route.path === '/forgot-password') {
      activeAuthTab.value = 'forgot'
      return
    }
    activeAuthTab.value = route.query.tab || 'login'
  }

  const filteredAuthTabs = computed(() => {
    return allowRegister.value ? authTabs : authTabs.filter(tab => tab.key !== 'register')
  })
  const isLoggedIn = computed(() => !!userStore.token)
  const packageLabel = computed(() => {
    const pkgName = userInfo.value?.package?.name || userInfo.value?.info?.package_name
    return pkgName || '未开通'
  })
  const availableServers = computed(() => {
    const packageServers = userInfo.value?.package?.servers || []
    return packageServers.length > 0 ? packageServers : publicServers.value
  })
  const selectedServer = computed(() => {
    return availableServers.value.find(server => (server.id || server.name) === selectedServerId.value) || availableServers.value[0] || null
  })
  const previewServer = computed(() => {
    return publicServers.value.find(server => server.id === selectedPreviewServerId.value) || publicServers.value[0] || null
  })
  const transferLimitText = computed(() => {
    const mb = userInfo.value?.info?.file_transfer_limit_mb || 100
    return `${mb} MB`
  })

  watch(() => route.fullPath, syncTabFromRoute)
  watch(availableServers, (list) => {
    if (!list.length) {
      selectedServerId.value = null
      return
    }
    const exists = list.some(server => (server.id || server.name) === selectedServerId.value)
    if (!exists) {
      selectedServerId.value = list[0].id || list[0].name
    }
  }, { immediate: true })
  watch(publicServers, (list) => {
    if (!list.length) {
      selectedPreviewServerId.value = null
      return
    }
    const exists = list.some(server => server.id === selectedPreviewServerId.value)
    if (!exists) {
      selectedPreviewServerId.value = list[0].id
    }
  }, { immediate: true })

  const switchAuthTab = (tab) => {
    activeAuthTab.value = tab
    const path = tab === 'register' ? '/register' : tab === 'forgot' ? '/forgot-password' : '/login'
    router.replace({ path })
  }

  const refreshPortal = async () => {
    await loadPublicServers()
    if (isLoggedIn.value) {
      await loadUserInfo()
    }
  }

  const loadPublicServers = async () => {
    const res = await listVipServers().catch(() => false)
    if (res?.list) {
      publicServers.value = res.list
    }
  }

  const loadLoginOptions = async () => {
    const res = await loginOptions().catch(() => false)
    if (!res?.data) {
      return
    }
    allowRegister.value = !!res.data.register
    if (!allowRegister.value && activeAuthTab.value === 'register') {
      switchAuthTab('login')
    }
    if (res.data.need_captcha) {
      await loadCaptcha()
    } else {
      captchaCode.value = null
      loginForm.captcha = ''
      loginForm.captcha_id = ''
    }
  }

  const loadCaptcha = async () => {
    const res = await captcha().catch(() => false)
    if (!res?.data?.captcha) {
      return
    }
    captchaCode.value = res.data.captcha
    loginForm.captcha_id = res.data.captcha.id
  }

  const loadUserInfo = async () => {
    const res = await current().catch(() => false)
    if (res?.data) {
      userInfo.value = res.data
    }
  }

  const submitLogin = async () => {
    if (!loginForm.username || !loginForm.password) {
      ElMessage.error('请输入账号和密码')
      return
    }
    const res = await userStore.login(loginForm).catch(err => err)
    if (res?.code) {
      if (res.code === 110) {
        await loadCaptcha()
      }
      return
    }
    if (!res) {
      return
    }
    ElMessage.success('登录成功')
    await loadUserInfo()
    router.replace(route.query.redirect || '/portal')
  }

  const submitRegister = async () => {
    if (!allowRegister.value) {
      ElMessage.error('当前站点暂未开放注册')
      return
    }
    if (!registerForm.username || !registerForm.password) {
      ElMessage.error('请填写账号和密码')
      return
    }
    if (registerForm.password !== registerForm.confirmPassword) {
      ElMessage.error('两次输入的密码不一致')
      return
    }
    const res = await registerUser({
      username: registerForm.username,
      email: registerForm.email,
      password: registerForm.password,
      confirm_password: registerForm.confirmPassword,
    }).catch(() => false)
    if (!res) {
      return
    }
    userStore.saveUserData(res.data)
    ElMessage.success('注册成功')
    await loadUserInfo()
    router.replace('/portal')
  }

  const startCodeCountdown = () => {
    codeCountdown.value = 60
    codeTimer = setInterval(() => {
      codeCountdown.value -= 1
      if (codeCountdown.value <= 0) {
        clearInterval(codeTimer)
        codeTimer = null
      }
    }, 1000)
  }

  const sendResetCode = async () => {
    if (!resetForm.email) {
      ElMessage.error('请输入邮箱')
      return
    }
    resetLoading.value = true
    const res = await sendPasswordResetCode({ email: resetForm.email }).catch(() => false)
    resetLoading.value = false
    if (!res) {
      return
    }
    ElMessage.success('验证码已发送')
    if (!codeTimer) {
      startCodeCountdown()
    }
  }

  const submitReset = async () => {
    if (!resetForm.email || !resetForm.code || !resetForm.newPassword) {
      ElMessage.error('请填写完整找回信息')
      return
    }
    const res = await resetPassword({
      email: resetForm.email,
      code: resetForm.code,
      new_password: resetForm.newPassword,
    }).catch(() => false)
    if (!res) {
      return
    }
    ElMessage.success('密码已重置，请重新登录')
    resetForm.code = ''
    resetForm.newPassword = ''
    switchAuthTab('login')
  }

  const submitRedeem = async () => {
    if (!redeemForm.code) {
      ElMessage.error('请输入激活码')
      return
    }
    const res = await redeemActivationCode({ code: redeemForm.code }).catch(() => false)
    if (!res) {
      return
    }
    ElMessage.success('激活成功')
    redeemForm.code = ''
    await loadUserInfo()
  }

  const submitChangePassword = async () => {
    if (!changePwdForm.oldPassword || !changePwdForm.newPassword || !changePwdForm.confirmPassword) {
      ElMessage.error('请填写完整密码信息')
      return
    }
    if (changePwdForm.newPassword !== changePwdForm.confirmPassword) {
      ElMessage.error('两次输入的新密码不一致')
      return
    }
    const confirmed = await ElMessageBox.confirm('修改密码后需要重新登录，是否继续？', '修改密码', {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning',
    }).catch(() => false)
    if (!confirmed) {
      return
    }
    const res = await changeCurPwd({
      old_password: changePwdForm.oldPassword,
      new_password: changePwdForm.newPassword,
      confirm_password: changePwdForm.confirmPassword,
    }).catch(() => false)
    if (!res) {
      return
    }
    ElMessage.success('密码修改成功，请重新登录')
    logout()
  }

  const copyText = async (value) => {
    if (!value) {
      return
    }
    let copied = false
    try {
      await navigator.clipboard.writeText(value)
      copied = true
    } catch (_) {}
    if (!copied) {
      const textarea = document.createElement('textarea')
      textarea.value = value
      textarea.style.position = 'fixed'
      textarea.style.opacity = '0'
      document.body.appendChild(textarea)
      textarea.focus()
      textarea.select()
      document.execCommand('copy')
      document.body.removeChild(textarea)
    }
    ElMessage.success('已复制')
  }

  const logout = () => {
    userStore.logout()
    userInfo.value = null
    router.replace('/login')
  }

  onMounted(async () => {
    syncTabFromRoute()
    await loadLoginOptions()
    await loadPublicServers()
    if (isLoggedIn.value) {
      await loadUserInfo()
      router.replace('/portal')
    }
  })

  onBeforeUnmount(() => {
    if (codeTimer) {
      clearInterval(codeTimer)
    }
  })
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

.hero-card,
.auth-card,
.preview-card,
.summary-card,
.ops-card,
.server-card {
  background: rgba(255, 255, 255, 0.92);
  border: 1px solid rgba(58, 121, 255, 0.12);
  box-shadow: 0 18px 40px rgba(27, 76, 160, 0.08);
  backdrop-filter: blur(8px);
}

.hero-card {
  display: flex;
  justify-content: space-between;
  gap: 24px;
  padding: 28px 32px;
  border-radius: 28px;
  margin-bottom: 24px;
}

.hero-brand {
  display: flex;
  align-items: center;
  gap: 18px;
}

.hero-logo {
  width: 58px;
  height: 58px;
}

.hero-title {
  font-size: 30px;
  font-weight: 700;
  color: #1d3d8f;
  letter-spacing: 0.02em;
}

.hero-subtitle,
.hero-note,
.card-desc,
.mini-desc {
  color: #5f6f91;
  line-height: 1.7;
}

.hero-note {
  margin-top: 14px;
  max-width: 640px;
}

.hero-status {
  min-width: 220px;
  padding: 22px 24px;
  border-radius: 22px;
  background: linear-gradient(180deg, #edf4ff, #f8fbff);
}

.status-label,
.summary-label,
.config-key,
.preview-label {
  font-size: 13px;
  color: #6f84ab;
}

.status-value,
.summary-value,
.mini-title,
.card-title {
  color: #193b88;
  font-weight: 700;
}

.status-value {
  font-size: 24px;
  margin-top: 8px;
}

.status-meta {
  margin-top: 6px;
  color: #376df3;
}

.logout-btn {
  margin-top: 12px;
}

.portal-grid {
  display: grid;
  grid-template-columns: 1.05fr 0.95fr;
  gap: 24px;
}

.logged-grid {
  grid-template-columns: 1fr;
}

.auth-card,
.preview-card,
.summary-card,
.ops-card,
.server-card {
  border-radius: 26px;
  padding: 28px;
}

.card-title {
  font-size: 24px;
  margin-bottom: 8px;
}

.auth-tabs {
  display: flex;
  gap: 10px;
  margin: 22px 0 20px;
}

.auth-tab,
.server-option,
.server-pill {
  border: 0;
  cursor: pointer;
  transition: all 0.2s ease;
}

.auth-tab {
  padding: 11px 18px;
  border-radius: 999px;
  background: #e8f1ff;
  color: #2a58b8;
  font-weight: 600;
}

.auth-tab.active {
  background: linear-gradient(135deg, #1d78ff, #1457eb);
  color: #fff;
  box-shadow: 0 10px 20px rgba(29, 120, 255, 0.2);
}

.portal-form {
  margin-top: 8px;
}

.action-btn {
  width: 100%;
  min-height: 44px;
  border-radius: 14px;
  font-weight: 600;
}

.inline-field {
  display: grid;
  grid-template-columns: 1fr 128px;
  gap: 12px;
  width: 100%;
}

.captcha-field {
  display: grid;
  grid-template-columns: 1fr 120px;
  gap: 12px;
  width: 100%;
}

.captcha-image {
  width: 120px;
  height: 44px;
  border-radius: 14px;
  object-fit: cover;
  cursor: pointer;
  border: 1px solid rgba(58, 121, 255, 0.14);
}

.info-list {
  display: flex;
  flex-direction: column;
  gap: 14px;
  margin-top: 16px;
}

.info-item {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  color: #40527b;
  line-height: 1.7;
}

.info-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #2879ff;
  margin-top: 10px;
  flex-shrink: 0;
}

.preview-panel {
  margin-top: 24px;
  padding: 20px;
  border-radius: 20px;
  background: linear-gradient(180deg, #f7fbff, #ffffff);
  border: 1px solid rgba(58, 121, 255, 0.1);
}

.interop-note {
  margin-top: 10px;
  color: #58729f;
  line-height: 1.7;
}

.preview-server-list,
.server-selector {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 16px;
  margin-bottom: 18px;
}

.server-pill,
.server-option {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 12px 16px;
  border-radius: 16px;
  background: #eef5ff;
  color: #244c99;
  font-weight: 600;
}

.server-option {
  align-items: flex-start;
  flex-direction: column;
  min-width: 160px;
  text-align: left;
}

.server-option small {
  color: #7690bd;
}

.server-pill.active,
.server-option.active {
  background: linear-gradient(135deg, #1c7bff, #1557ea);
  color: #fff;
  box-shadow: 0 10px 18px rgba(21, 87, 234, 0.22);
}

.server-option.active small {
  color: rgba(255, 255, 255, 0.84);
}

.summary-grid,
.ops-grid,
.config-grid {
  display: grid;
  gap: 16px;
}

.summary-grid {
  grid-template-columns: repeat(3, minmax(0, 1fr));
  margin-top: 18px;
}

.summary-item,
.mini-card,
.config-item {
  padding: 18px 20px;
  border-radius: 18px;
  background: linear-gradient(180deg, #f8fbff, #ffffff);
  border: 1px solid rgba(58, 121, 255, 0.1);
}

.summary-value {
  margin-top: 8px;
  font-size: 20px;
}

.ops-grid {
  grid-template-columns: repeat(2, minmax(0, 1fr));
  margin-top: 18px;
}

.mini-form {
  margin-top: 14px;
}

.card-head {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  align-items: flex-start;
}

.config-wrap {
  margin-top: 8px;
}

.config-grid {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.config-grid.compact {
  margin-top: 16px;
}

.config-val {
  margin-top: 8px;
  color: #224587;
  font-weight: 600;
  line-height: 1.6;
  word-break: break-all;
}

.line-status {
  display: inline-flex;
  align-items: center;
  padding: 6px 12px;
  border-radius: 999px;
  font-size: 13px;
  font-weight: 700;
}

.line-status.online {
  color: #11783c;
  background: rgba(24, 181, 94, 0.12);
}

.line-status.offline {
  color: #c33d42;
  background: rgba(240, 72, 72, 0.12);
}

.mono {
  font-family: Consolas, Monaco, monospace;
}

.server-hint {
  margin-top: 16px;
  color: #5f7195;
  line-height: 1.7;
}

:deep(.el-input__wrapper) {
  min-height: 44px;
  border-radius: 14px;
  box-shadow: 0 0 0 1px rgba(58, 121, 255, 0.12) inset;
}

:deep(.el-button--primary) {
  background: linear-gradient(135deg, #1c7bff, #1557ea);
  border-color: #1c7bff;
}

@media (max-width: 960px) {
  .portal-grid,
  .summary-grid,
  .ops-grid,
  .config-grid {
    grid-template-columns: 1fr;
  }

  .hero-card {
    flex-direction: column;
  }

  .portal-page {
    padding: 20px 16px 36px;
  }
}
</style>
