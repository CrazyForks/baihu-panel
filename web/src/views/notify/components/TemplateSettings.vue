<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Badge } from '@/components/ui/badge'
import { api } from '@/api'
import { toast } from 'vue-sonner'
import { Save, RefreshCcw, Info, ChevronDown, ChevronUp } from 'lucide-vue-next'

const props = defineProps<{
  activeTab?: string,
}>()

const emit = defineEmits(['update:activeTab'])

const loading = ref(false)
const saving = ref(false)

const prefix = ref('')
const templates = ref<Record<string, string>>({})
const expandedEvents = ref<Record<string, boolean>>({})

const eventGroups = [
  {
    title: '系统与安全事件',
    events: [
      {
        id: 'user_login',
        name: '用户登录 (成功/失败)',
        keys: { title: 'notify_template_user_login_title', text: 'notify_template_user_login_text' },
        variables: ['username', 'ip', 'status_label', 'message']
      },
      {
        id: 'brute_force_login',
        name: '密码尝试破解',
        keys: { title: 'notify_template_brute_force_login_title', text: 'notify_template_brute_force_login_text' },
        variables: ['ip', 'username']
      },
      {
        id: 'password_changed',
        name: '密码修改',
        keys: { title: 'notify_template_password_changed_title', text: 'notify_template_password_changed_text' },
        variables: ['username']
      }
    ]
  },
  {
    title: '任务执行事件',
    events: [
      {
        id: 'task_success',
        name: '任务成功',
        keys: { title: 'notify_template_task_success_title', text: 'notify_template_task_success_text' },
        variables: ['task_id', 'task_name', 'start_time', 'duration', 'output']
      },
      {
        id: 'task_failed',
        name: '任务失败',
        keys: { title: 'notify_template_task_failed_title', text: 'notify_template_task_failed_text' },
        variables: ['task_id', 'task_name', 'start_time', 'duration', 'error', 'output']
      },
      {
        id: 'task_timeout',
        name: '任务超时',
        keys: { title: 'notify_template_task_timeout_title', text: 'notify_template_task_timeout_text' },
        variables: ['task_id', 'task_name', 'start_time', 'duration', 'output']
      }
    ]
  }
]

async function loadSettings() {
  loading.value = true
  try {
    const res = await api.settings.getSection('notify')
    prefix.value = res.notify_prefix || '[白虎面板]'
    templates.value = res
  } catch (e: any) {
    toast.error('加载配置失败: ' + e.message)
  } finally {
    loading.value = false
  }
}

async function saveSettings() {
  saving.value = true
  try {
    const data: Record<string, string> = {
      notify_prefix: prefix.value,
      ...templates.value
    }
    await api.settings.setSection('notify', data)
    toast.success('模板配置已保存')
  } catch (e: any) {
    toast.error('保存失败: ' + e.message)
  } finally {
    saving.value = false
  }
}

function insertVariable(eventKey: string, variable: string) {
  const current = templates.value[eventKey] || ''
  templates.value[eventKey] = current + ` {{${variable}}}`
}

function toggleExpand(id: string) {
  expandedEvents.value[id] = !expandedEvents.value[id]
}

onMounted(() => {
  loadSettings()
  // 默认展开第一个
  expandedEvents.value['user_login'] = true
})
</script>

<template>
  <div class="space-y-6">
    <Card class="border-none shadow-none bg-transparent">
      <CardHeader class="px-0 pt-0 pb-4">
        <div class="space-y-4">
          <div class="flex items-center justify-between gap-4">
            <CardTitle class="text-xl sm:text-2xl font-bold tracking-tight">通知模板管理</CardTitle>
            <div class="flex items-center gap-2 shrink-0">
              <Button variant="outline" size="sm" @click="loadSettings" :disabled="loading" class="h-9 w-9 sm:w-auto p-0 sm:px-3">
                <RefreshCcw class="w-4 h-4 sm:mr-2" :class="{ 'animate-spin': loading }" />
                <span class="hidden sm:inline">刷新</span>
              </Button>
              <Button size="sm" @click="saveSettings" :disabled="saving" class="h-9 px-4">
                <Save class="w-4 h-4 sm:mr-2" />
                <span class="hidden sm:inline">{{ saving ? '保存中...' : '提交修改' }}</span>
                <span class="sm:hidden">保存</span>
              </Button>
            </div>
          </div>
          <CardDescription class="text-xs sm:text-sm leading-relaxed max-w-2xl">
            定制通知的消息格式，支持局前缀与动态变量内置变量
          </CardDescription>
        </div>
      </CardHeader>

      <CardContent class="px-0 space-y-6">
        <!-- 全局前缀 -->
        <div class="p-5 rounded-2xl bg-accent/20 border border-accent/30 space-y-4">
          <div class="flex items-center gap-2 text-sm font-bold text-foreground">
             <div class="w-1.5 h-4 bg-primary rounded-full" />
             全局消息前缀
          </div>
          <div class="flex flex-col lg:flex-row gap-4 items-start lg:items-end">
            <div class="flex-1 w-full space-y-2">
              <Label class="text-[11px] text-muted-foreground ml-1 font-medium tracking-wide uppercase">该前缀会添加在所有通知标题的最前面</Label>
              <Input v-model="prefix" placeholder="例如: [生产环境]" class="bg-background/60 h-10 border-accent/20 focus:border-primary/50" />
            </div>
            <div class="flex items-center gap-3 px-4 h-10 rounded-lg bg-background/40 border border-dashed border-muted-foreground/20 text-[11px] text-muted-foreground shrink-0 w-full lg:w-auto">
              预览效果: <span class="font-mono text-primary font-bold tracking-tight">{{ prefix }} 用户登录成功</span>
            </div>
          </div>
        </div>

        <!-- 模板详情 -->
        <div class="w-full space-y-8 mt-4">
          <div v-for="group in eventGroups" :key="group.title" class="space-y-4">
            <div class="flex items-center gap-3 ml-1">
              <h3 class="text-[11px] font-black text-muted-foreground uppercase tracking-[0.2em]">{{ group.title }}</h3>
              <div class="h-[1px] flex-1 bg-gradient-to-r from-border/60 to-transparent" />
            </div>
            
            <!-- 任务执行事件专属提示 -->
            <div v-if="group.title === '任务执行事件'" class="p-3 rounded-lg bg-yellow-500/5 border border-yellow-500/10 flex gap-3 text-xs text-yellow-600/80 leading-relaxed mb-6">
              <Info class="w-4 h-4 shrink-0 mt-0.5 text-yellow-500/50" />
              <div>
                <span class="font-bold">配置提示：</span>
                部分变量（如 <code v-pre class="bg-yellow-500/10 px-1 rounded text-yellow-600">{{ output }}</code>）依赖于的具体绑定配置。请确保在 
                <span class="font-bold underline decoration-dotted cursor-pointer hover:text-yellow-600" @click="emit('update:activeTab', 'events')">事件绑定</span> 
                的高级设置中开启了 <span class="text-foreground/80">“发送任务日志”</span>，否则通知消息中该变量将为空。
              </div>
            </div>
            
            <div v-for="event in group.events" :key="event.id" 
              class="border rounded-xl bg-background/50 overflow-hidden transition-all duration-200"
              :class="{ 'ring-1 ring-primary/20 bg-accent/5': expandedEvents[event.id] }">
              
              <div class="flex items-center justify-between p-4 cursor-pointer select-none" @click="toggleExpand(event.id)">
                <div class="flex items-center gap-3">
                  <div class="w-8 h-8 rounded-lg bg-primary/5 flex items-center justify-center text-primary">
                    <Info class="w-4 h-4" />
                  </div>
                  <div class="text-left">
                    <div class="text-sm font-semibold">{{ event.name }}</div>
                    <div class="text-[10px] text-muted-foreground font-normal uppercase tracking-tight">ID: {{ event.id }}</div>
                  </div>
                </div>
                <component :is="expandedEvents[event.id] ? ChevronUp : ChevronDown" class="w-4 h-4 text-muted-foreground" />
              </div>

              <div v-if="expandedEvents[event.id]" class="px-4 pb-6 space-y-4 pt-2 border-t border-dashed">
                <!-- 变量提示 -->
                <div class="flex flex-wrap items-center gap-2 py-2">
                  <span class="text-[10px] font-bold text-muted-foreground mr-1 uppercase">可用参数:</span>
                  <Badge v-for="v in event.variables" :key="v" variant="secondary" 
                    class="cursor-pointer hover:bg-primary/10 hover:text-primary transition-colors py-0.5 px-2 font-medium text-[11px] border-none"
                    @click="insertVariable(event.keys.text, v)"
                    v-text="'{{' + v + '}}'">
                  </Badge>
                </div>

                <div class="flex flex-col gap-5">
                  <!-- 标题模板 -->
                  <div class="space-y-2.5">
                    <Label class="text-[11px] font-bold flex items-center gap-1.5 text-muted-foreground uppercase tracking-wide">
                      推送标题模板
                    </Label>
                    <Input v-model="templates[event.keys.title]" placeholder="通知标题" class="bg-background/80 h-10 border-accent/20 focus:border-primary/50" />
                  </div>
                  <!-- 内容模板 -->
                  <div class="space-y-2.5">
                    <Label class="text-[11px] font-bold flex items-center gap-1.5 text-muted-foreground uppercase tracking-wide">
                      推送正文模板
                    </Label>
                    <Textarea v-model="templates[event.keys.text]" 
                      :rows="4"
                      placeholder="通知详细内容..." 
                      class="resize-none font-sans leading-relaxed bg-background/80 border-accent/20 focus:border-primary/50" />
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</template>

<style scoped>
/* 移除不必要的移动端过重样式，保持清爽 */
</style>
