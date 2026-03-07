<template>
  <div class="login-container">
    <div class="login-card">
      <img src="@/assets/logo.png" alt="logo" class="login-logo"/>

      <el-form ref="f" :model="form" label-position="top" class="login-form" :rules="rules">
        <div class="page-note">输入邮箱、验证码和新密码即可找回密码。</div>

        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email" class="login-input"></el-input>
        </el-form-item>

        <el-form-item label="验证码" prop="code">
          <el-input v-model="form.code" class="login-input code-input">
            <template #append>
              <el-button :disabled="sendingCode || countdown > 0" @click="sendCode">
                {{ countdown > 0 ? `${countdown}s` : '发送验证码' }}
              </el-button>
            </template>
          </el-input>
        </el-form-item>

        <el-form-item :label="T('NewPassword')" prop="new_password">
          <el-input
            v-model="form.new_password"
            type="password"
            show-password
            class="login-input"
            @keyup.enter.native="submit"
          ></el-input>
        </el-form-item>

        <el-form-item>
          <el-button @click="submit" class="login-button" type="primary">{{ T('ResetPassword') }}</el-button>
          <el-button @click="toLogin" class="login-button">{{ T('ToLogin') }}</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
  import { onBeforeUnmount, reactive, ref } from 'vue'
  import { ElMessage } from 'element-plus'
  import { useRouter } from 'vue-router'
  import { resetPassword, sendPasswordResetCode } from '@/api/login'
  import { T } from '@/utils/i18n'

  const router = useRouter()
  const f = ref(null)
  const sendingCode = ref(false)
  const countdown = ref(0)
  let timer = null

  const form = reactive({
    email: '',
    code: '',
    new_password: '',
  })

  const rules = {
    email: [
      { required: true, message: '邮箱不能为空', trigger: 'blur' },
    ],
    code: [
      { required: true, message: '验证码不能为空', trigger: 'blur' },
    ],
    new_password: [
      { required: true, message: '新密码不能为空', trigger: 'blur' },
      {
        validator: (rule, value, callback) => {
          if (!value || value.length < 6 || value.length > 18) {
            callback(new Error('密码长度需为 6-18 位'))
          } else {
            callback()
          }
        },
        trigger: 'blur',
      },
    ],
  }

  const startCountdown = () => {
    countdown.value = 60
    timer = setInterval(() => {
      countdown.value -= 1
      if (countdown.value <= 0) {
        clearInterval(timer)
        timer = null
      }
    }, 1000)
  }

  const sendCode = async () => {
    if (!form.email.trim()) {
      ElMessage.warning('请先输入邮箱')
      return
    }
    if (countdown.value > 0 || sendingCode.value) {
      return
    }
    sendingCode.value = true
    const res = await sendPasswordResetCode({
      email: form.email.trim(),
    }).catch(_ => false)
    sendingCode.value = false
    if (!res) {
      return
    }
    ElMessage.success('验证码已发送，请检查邮箱。')
    startCountdown()
  }

  const submit = async () => {
    const valid = await f.value.validate().catch(_ => false)
    if (!valid) {
      return
    }
    const res = await resetPassword({
      email: form.email.trim(),
      code: form.code.trim(),
      new_password: form.new_password,
    }).catch(_ => false)
    if (!res) {
      return
    }
    ElMessage.success('密码已重置，请重新登录。')
    router.push('/login')
  }

  const toLogin = () => {
    router.push('/login')
  }

  onBeforeUnmount(() => {
    if (timer) {
      clearInterval(timer)
      timer = null
    }
  })
</script>

<style scoped lang="scss">
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #2d3a4b;
  padding: 20px;
  box-sizing: border-box;
}

.login-card {
  width: 360px;
  background-color: #283342;
  padding: 40px;
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  text-align: center;
}

.login-form {
  margin-bottom: 20px;
}

.login-input {
  width: 100%;
}

.code-input {
  :deep(.el-input-group__append) {
    padding: 0;
  }
}

.login-button {
  width: 100%;
  height: 40px;
  margin-bottom: 20px;
  margin-top: 20px;
  margin-left: 0;
}

.login-logo {
  width: 80px;
  height: 80px;
  margin: 0 auto 20px;
  display: block;
}

.page-note {
  margin-bottom: 20px;
  color: rgba(255, 255, 255, 0.72);
  font-size: 13px;
  line-height: 1.6;
  text-align: left;
}

.el-form-item {
  ::v-deep(.el-form-item__label) {
    color: #fff;
  }

  .el-input {
    ::v-deep(.el-input__wrapper) {
      border: 1px solid rgba(255, 255, 255, 0.1);
      background: transparent;
    }

    ::v-deep(input) {
      color: #fff;
    }
  }
}
</style>
