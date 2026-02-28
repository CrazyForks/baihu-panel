<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { api, type SiteSettings } from '@/api'
import { toast } from 'vue-sonner'
import { useSiteSettings } from '@/composables/useSiteSettings'
import { Badge } from '@/components/ui/badge'
import { RefreshCw, Copy } from 'lucide-vue-next'

const { refreshSettings } = useSiteSettings()

const form = ref<SiteSettings>({
  title: '',
  subtitle: '',
  icon: '',
  page_size: '10',
  cookie_days: '7',
  api_token: '',
  api_token_expire: ''
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
    await api.settings.updateSite({
      ...form.value,
      page_size: String(form.value.page_size),
      cookie_days: String(form.value.cookie_days)
    })
    await refreshSettings()
    toast.success('保存成功')
  } catch {
    toast.error('保存失败')
  } finally {
    loading.value = false
  }
}

async function generateToken() {
  try {
    const res = await api.settings.generateApiToken()
    form.value.api_token = res.token
    
    // 如果没有设置过期时间，默认给一年后
    if (!form.value.api_token_expire) {
      const d = new Date()
      d.setFullYear(d.getFullYear() + 1)
      form.value.api_token_expire = d.toISOString().split('T')[0]
    }
  } catch {
    toast.error('生成 Token 失败')
  }
}

async function copyToken() {
  if (!form.value.api_token) return
  try {
    await navigator.clipboard.writeText(form.value.api_token)
    toast.success('Token 已复制到剪贴板')
  } catch {
    toast.error('复制失败，请手动复制')
  }
}

onMounted(loadSettings)
</script>

<template>
  <div class="space-y-4">
    <div class="grid grid-cols-1 sm:grid-cols-4 items-center gap-2 sm:gap-4">
      <Label class="sm:text-right">站点标题</Label>
      <Input v-model="form.title" placeholder="白虎面板" class="sm:col-span-3" />
    </div>
    <div class="grid grid-cols-1 sm:grid-cols-4 items-center gap-2 sm:gap-4">
      <Label class="sm:text-right">站点标语</Label>
      <Input v-model="form.subtitle" placeholder="轻量级定时任务管理系统" class="sm:col-span-3" />
    </div>
    <div class="grid grid-cols-1 sm:grid-cols-4 items-center gap-2 sm:gap-4">
      <Label class="sm:text-right">站点图标</Label>
      <div class="sm:col-span-3 flex items-center gap-2">
        <Input v-model="form.icon" placeholder="<svg>...</svg>" class="flex-1 font-mono text-xs" />
        <div v-if="iconPreview" class="p-1.5 border rounded bg-white dark:bg-white w-8 h-8 flex items-center justify-center shrink-0 [&>svg]:w-5 [&>svg]:h-5" v-html="iconPreview" />
      </div>
    </div>
    <div class="grid grid-cols-1 sm:grid-cols-4 items-center gap-2 sm:gap-4">
      <Label class="sm:text-right">分页/Cookie</Label>
      <div class="sm:col-span-3 flex flex-wrap items-center gap-4">
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
    
    <div class="pt-6 border-t mt-6">
      <div class="flex items-center gap-2 mb-4">
        <h3 class="text-lg font-medium text-foreground">API Token</h3>
        <Badge variant="secondary" class="font-normal text-xs bg-amber-500/10 text-amber-600 dark:text-amber-400 border-amber-500/20">实验特性，可能变更</Badge>
      </div>
      <p class="text-sm text-muted-foreground mb-4">开启全局 API 直接访问能力，配置后可通过请求头 <code class="bg-muted px-1.5 py-0.5 rounded text-xs select-all font-sans">X-API-Token: &lt;在此生成的Token&gt;</code> 无需登录直接调用系统的所有接口，请妥善保管并设置合理的有效期。</p>
      
      <div class="grid grid-cols-1 sm:grid-cols-4 items-center gap-2 sm:gap-4 mb-4">
        <Label class="sm:text-right text-muted-foreground">Token 密钥</Label>
        <div class="sm:col-span-3 flex w-full max-w-sm items-center space-x-2">
          <Input v-model="form.api_token" placeholder="点击右侧按钮生成 32 位随机 Token" class="text-sm" />
          <Button type="button" variant="outline" size="icon" @click="generateToken" title="随机生成">
            <RefreshCw class="w-4 h-4" />
          </Button>
          <Button type="button" variant="outline" size="icon" @click="copyToken" title="复制" :disabled="!form.api_token">
            <Copy class="w-4 h-4" />
          </Button>
        </div>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-4 items-center gap-2 sm:gap-4">
        <Label class="sm:text-right text-muted-foreground">截止有效期</Label>
        <div class="sm:col-span-3">
          <Input v-model="form.api_token_expire" type="date" class="w-full max-w-xs dark:[color-scheme:dark]" />
          <div class="text-xs text-muted-foreground mt-1.5 ml-1">超过此日期后该 Token 将失效，置空代表该特性完全关闭。</div>
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
