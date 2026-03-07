<template>
  <div class="form-card">
    <el-form ref="root" label-width="120px" :model="form" :rules="rules">
      <el-form-item :label="T('Username')" prop="username">
        <el-input v-model="form.username"></el-input>
      </el-form-item>
      <el-form-item :label="T('Email')" prop="email">
        <el-input v-model="form.email"></el-input>
      </el-form-item>
      <el-form-item :label="T('Nickname')" prop="nickname">
        <el-input v-model="form.nickname"></el-input>
      </el-form-item>
      <el-form-item :label="T('Group')" prop="group_id">
        <el-select v-model="form.group_id">
          <el-option
              v-for="item in groupsList"
              :key="item.id"
              :label="item.name"
              :value="item.id"
          ></el-option>
        </el-select>
      </el-form-item>
      <el-form-item :label="T('IsAdmin')" prop="is_admin">
        <el-switch v-model="form.is_admin"
                   :active-value="true"
                   :inactive-value="false"
        ></el-switch>
      </el-form-item>
      <el-form-item :label="T('Status')" prop="status">
        <el-switch v-model="form.status"
                   :active-value="ENABLE_STATUS"
                   :inactive-value="DISABLE_STATUS"
        ></el-switch>
      </el-form-item>
      <el-form-item label="套餐" prop="package_id">
        <el-select v-model="form.package_id" clearable placeholder="选择套餐">
          <el-option
              v-for="item in packagesList"
              :key="item.id"
              :label="item.name"
              :value="item.id"
          ></el-option>
        </el-select>
        <el-text class="mx-1" type="info">选择套餐后自动填充有效天数和设备限制</el-text>
      </el-form-item>
      <el-form-item :label="T('ValidDays')" prop="valid_days">
        <el-input-number
          v-model="form.valid_days"
          :min="-1"
          :placeholder="T('ValidDaysHint')"
        ></el-input-number>
        <el-text class="mx-1" type="info">{{ T('ValidDaysDesc') }}</el-text>
      </el-form-item>
      <el-form-item label="设备限制" prop="device_limit">
        <el-input-number
          v-model="form.device_limit"
          :min="1"
          :max="1000"
        ></el-input-number>
        <el-text class="mx-1" type="info">允许同时登录的设备数量</el-text>
      </el-form-item>
      <el-form-item :label="T('ActivationDate')" v-if="form.id">
        <el-input v-model="form.first_login_at" disabled></el-input>
      </el-form-item>
      <el-form-item :label="T('Remark')" prop="remark">
          <el-input v-model="form.remark"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button @click="cancel">{{ T('Cancel') }}</el-button>
        <el-button @click="submit" type="primary">{{ T('Submit') }}</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
  import { useRoute } from 'vue-router'
  import { useGetDetail, useSubmit } from '@/views/user/composables/edit'
  import { ENABLE_STATUS, DISABLE_STATUS } from '@/utils/common_options'
  import { T } from '@/utils/i18n'

  const route = useRoute()
  const { form, item, getDetail, groupsList, packagesList } = useGetDetail(route.params.id)

  const { root, rules, validate, submit, cancel } = useSubmit(form, route.params.id)

</script>

<style lang="scss" scoped>
.form-card {
}
</style>
