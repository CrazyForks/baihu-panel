<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { api, type SiteSettings } from '@/api'
import { toast } from 'vue-sonner'
import { useSiteSettings } from '@/composables/useSiteSettings'

const { refreshSettings } = useSiteSettings()

const form = ref<SiteSettings>({
  title: '',
  subtitle: '',
  icon: '',
  page_size: '10',
  cookie_days: '7'
})
const loading = ref(false)

const iconPreview = computed(() => {
  if (!form.value.icon) return ''
  // 简单验证是否是 SVG
  if (form.value.icon.trim().startsWith('<svg')) {
    return form.value.icon
  }
  return ''
})

async function loadSettings() {
  try {
    const res = await api.settings.getSite()
    form.value = res
  } catch {}
}

async function saveSettings() {
  loading.value = true
  try {
    await api.settings.updateSite(form.value)
    await refreshSettings()
    toast.success('保存成功')
  } catch {
    toast.error('保存失败')
  } finally {
    loading.value = false
  }
}

onMounted(loadSettings)
</script>

<template>
  <div class="space-y-4">
    <div class="grid grid-cols-4 items-center gap-4">
      <Label class="text-right">站点标题</Label>
      <Input v-model="form.title" placeholder="白虎面板" class="col-span-3" />
    </div>
    <div class="grid grid-cols-4 items-center gap-4">
      <Label class="text-right">站点标语</Label>
      <Input v-model="form.subtitle" placeholder="轻量级定时任务管理系统" class="col-span-3" />
    </div>
    <div class="grid grid-cols-4 items-center gap-4">
      <Label class="text-right">站点图标</Label>
      <div class="col-span-3 flex items-center gap-2">
        <Input v-model="form.icon" placeholder="<svg>...</svg>" class="flex-1 font-mono text-xs" />
        <div v-if="iconPreview" class="p-1.5 border rounded bg-muted/50 w-8 h-8 flex items-center justify-center shrink-0 [&>svg]:w-5 [&>svg]:h-5" v-html="iconPreview" />
      </div>
    </div>
    <div class="grid grid-cols-4 items-center gap-4">
      <Label class="text-right">分页/Cookie</Label>
      <div class="col-span-3 flex items-center gap-4">
        <div class="flex items-center gap-2">
          <Input v-model="form.page_size" type="number" class="w-20" />
          <span class="text-sm text-muted-foreground">条/页</span>
        </div>
        <div class="flex items-center gap-2">
          <Input v-model="form.cookie_days" type="number" class="w-20" />
          <span class="text-sm text-muted-foreground">天过期</span>
        </div>
      </div>
    </div>
    <div class="flex justify-end pt-2">
      <Button @click="saveSettings" :disabled="loading">
        {{ loading ? '保存中...' : '保存设置' }}
      </Button>
    </div>
  </div>
</template>
