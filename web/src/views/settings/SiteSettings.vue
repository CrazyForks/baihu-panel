<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { api } from '@/api'

const siteName = ref('')
const sitePort = ref(0)

async function loadSiteSettings() {
  try {
    const res = await api.settings.getSite()
    siteName.value = res.site_name || '白虎面板'
    sitePort.value = res.port
  } catch {}
}

onMounted(loadSiteSettings)
</script>

<template>
  <div class="space-y-4">
    <div class="space-y-2">
      <Label>站点名称</Label>
      <Input v-model="siteName" disabled class="max-w-sm" />
    </div>
    <div class="space-y-2">
      <Label>端口</Label>
      <Input v-model.number="sitePort" type="number" class="w-32" disabled />
    </div>
    <p class="text-sm text-muted-foreground">站点设置需要修改配置文件后重启服务</p>
  </div>
</template>
