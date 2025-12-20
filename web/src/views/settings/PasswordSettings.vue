<script setup lang="ts">
import { ref } from 'vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { api } from '@/api'
import { toast } from 'vue-sonner'

const oldPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')

async function changePassword() {
  if (!oldPassword.value || !newPassword.value) {
    toast.error('请填写完整')
    return
  }
  if (newPassword.value.length < 6) {
    toast.error('新密码至少6位')
    return
  }
  if (newPassword.value !== confirmPassword.value) {
    toast.error('两次密码不一致')
    return
  }
  try {
    await api.settings.changePassword({ old_password: oldPassword.value, new_password: newPassword.value })
    toast.success('密码修改成功')
    oldPassword.value = ''
    newPassword.value = ''
    confirmPassword.value = ''
  } catch (e: any) {
    toast.error(e.message || '修改失败')
  }
}
</script>

<template>
  <div class="space-y-4">
    <div class="space-y-2">
      <Label>原密码</Label>
      <Input v-model="oldPassword" type="password" placeholder="请输入原密码" class="max-w-sm" />
    </div>
    <div class="space-y-2">
      <Label>新密码</Label>
      <Input v-model="newPassword" type="password" placeholder="请输入新密码（至少6位）" class="max-w-sm" />
    </div>
    <div class="space-y-2">
      <Label>确认密码</Label>
      <Input v-model="confirmPassword" type="password" placeholder="请再次输入新密码" class="max-w-sm" />
    </div>
    <Button @click="changePassword">修改密码</Button>
  </div>
</template>
