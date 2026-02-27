<script setup lang="ts">
import { ref, onMounted, computed, onUnmounted } from 'vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter, DialogDescription } from '@/components/ui/dialog'
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from '@/components/ui/alert-dialog'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import {
  RefreshCw, Trash2, Edit, Copy, Server, Search, Download, RotateCw,
  Plus, Ticket, ListTodo, Eye, Wifi as WifiIcon, WifiOff as WifiOffIcon,
  Zap as ZapIcon, ZapOff as ZapOffIcon, Check, X
} from 'lucide-vue-next'
import { api, type Agent, type AgentToken } from '@/api'
import { toast } from 'vue-sonner'
import { useRouter } from 'vue-router'
import { AGENT_STATUS } from '@/constants'

const router = useRouter()

const agents = ref<Agent[]>([])
const tokens = ref<AgentToken[]>([])
const loading = ref(false)
const searchQuery = ref('')
const activeTab = ref('agents')
const agentVersion = ref('')
const platforms = ref<{ os: string; arch: string; filename: string }[]>([])
const showEditDialog = ref(false)
const showDeleteDialog = ref(false)
const showDownloadDialog = ref(false)
const showTokenDialog = ref(false)
const showDetailDialog = ref(false)
const formData = ref({ name: '', description: '' })
const tokenForm = ref({ remark: '', max_uses: 0, expires_at: '' })
const editingAgent = ref<Agent | null>(null)
const deletingAgent = ref<Agent | null>(null)
const viewingAgent = ref<Agent | null>(null)
let refreshTimer: ReturnType<typeof setInterval> | null = null

const filteredAgents = computed(() => {
  if (!searchQuery.value) return agents.value
  const q = searchQuery.value.toLowerCase()
  return agents.value.filter(a =>
    a.name.toLowerCase().includes(q) ||
    a.hostname?.toLowerCase().includes(q) ||
    a.ip?.toLowerCase().includes(q)
  )
})

function isOnline(agent: Agent): boolean {
  return agent.status === AGENT_STATUS.ONLINE
}

async function loadAgents() {
  loading.value = true
  try {
    const [agentList, versionInfo, tokenList] = await Promise.all([
      api.agents.list(),
      api.agents.getVersion(),
      api.agents.listTokens()
    ])
    agents.value = agentList
    agentVersion.value = versionInfo.version || ''
    platforms.value = versionInfo.platforms || []
    tokens.value = tokenList
  } catch {
    toast.error('加载失败')
  } finally {
    loading.value = false
  }
}

function viewDetail(agent: Agent) {
  viewingAgent.value = agent
  showDetailDialog.value = true
}

function openEditDialog(agent: Agent) {
  editingAgent.value = agent
  formData.value = { name: agent.name, description: agent.description }
  showEditDialog.value = true
}

async function updateAgent() {
  if (!editingAgent.value || !formData.value.name.trim()) return
  try {
    await api.agents.update(editingAgent.value.id, { ...formData.value, enabled: editingAgent.value.enabled })
    showEditDialog.value = false
    await loadAgents()
    toast.success('更新成功')
  } catch (e: unknown) {
    toast.error((e as Error).message || '更新失败')
  }
}

async function toggleEnabled(agent: Agent) {
  try {
    const newEnabled = !agent.enabled
    await api.agents.update(agent.id, { name: agent.name, description: agent.description, enabled: newEnabled })
    await loadAgents()
    toast.success(`${agent.name} 已${newEnabled ? '启用' : '禁用'}`)
  } catch (e: unknown) {
    toast.error((e as Error).message || '操作失败')
  }
}

function confirmDelete(agent: Agent) {
  deletingAgent.value = agent
  showDeleteDialog.value = true
}

async function deleteAgent() {
  if (!deletingAgent.value) return
  try {
    await api.agents.delete(deletingAgent.value.id)
    showDeleteDialog.value = false
    await loadAgents()
    toast.success('删除成功')
  } catch (e: unknown) {
    toast.error((e as Error).message || '删除失败')
  }
}

async function forceUpdate(agent: Agent) {
  try {
    await api.agents.forceUpdate(agent.id)
    toast.success('已标记强制更新')
  } catch (e: unknown) {
    toast.error((e as Error).message || '操作失败')
  }
}

function viewTasks(agent: Agent) {
  router.push({ path: '/tasks', query: { agent_id: String(agent.id) } })
}

function copyToken(token: string) {
  navigator.clipboard.writeText(token)
  toast.success('已复制')
}

async function createToken() {
  try {
    await api.agents.createToken({
      remark: tokenForm.value.remark,
      max_uses: tokenForm.value.max_uses,
      expires_at: tokenForm.value.expires_at || undefined
    })
    showTokenDialog.value = false
    tokenForm.value = { remark: '', max_uses: 0, expires_at: '' }
    await loadAgents()
    toast.success('创建成功')
  } catch (e: unknown) {
    toast.error((e as Error).message || '创建失败')
  }
}

async function deleteToken(id: number) {
  try {
    await api.agents.deleteToken(id)
    await loadAgents()
    toast.success('删除成功')
  } catch (e: unknown) {
    toast.error((e as Error).message || '删除失败')
  }
}

function isTokenExpired(token: AgentToken) {
  if (!token.expires_at) return false
  return new Date(token.expires_at) < new Date()
}

function isTokenExhausted(token: AgentToken) {
  return token.max_uses > 0 && token.used_count >= token.max_uses
}

function downloadAgent(os: string, arch: string) {
  window.open(api.agents.downloadUrl(os, arch), '_blank')
}

function getPlatformLabel(os: string, arch: string) {
  const osLabels: Record<string, string> = { linux: 'Linux', windows: 'Windows', darwin: 'macOS' }
  const archLabels: Record<string, string> = { amd64: 'x64', arm64: 'ARM64', '386': 'x86' }
  return `${osLabels[os] || os} ${archLabels[arch] || arch}`
}

onMounted(() => {
  loadAgents()
  refreshTimer = setInterval(loadAgents, 10000)
})

onUnmounted(() => {
  if (refreshTimer) clearInterval(refreshTimer)
})
</script>


<template>
  <div class="space-y-4">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h2 class="text-xl sm:text-2xl font-bold tracking-tight">Agent 管理</h2>
        <p class="text-muted-foreground text-sm">管理远程执行代理</p>
      </div>
      <div class="flex items-center gap-2">
        <div class="relative flex-1 sm:flex-none">
          <Search class="absolute left-2.5 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input v-model="searchQuery" placeholder="搜索..." class="h-9 pl-8 w-full sm:w-48 text-sm" />
        </div>
        <Button variant="outline" size="sm" class="h-9" @click="showDownloadDialog = true">
          <Download class="h-4 w-4 mr-1.5" />下载
        </Button>
        <Button variant="outline" size="icon" class="h-9 w-9 shrink-0" @click="loadAgents" :disabled="loading">
          <RefreshCw class="h-4 w-4" :class="{ 'animate-spin': loading }" />
        </Button>
      </div>
    </div>

    <Tabs v-model="activeTab">
      <TabsList>
        <TabsTrigger value="agents">Agent 列表</TabsTrigger>
        <TabsTrigger value="regcodes">
          <Ticket class="h-4 w-4 mr-1" />令牌
        </TabsTrigger>
      </TabsList>

      <TabsContent value="agents" class="mt-4">
        <div class="rounded-lg border bg-card overflow-x-auto">
          <div class="min-w-full w-max">
            <!-- 大屏表头 -->
            <div
              class="hidden sm:flex items-center gap-2 sm:gap-4 px-3 sm:px-4 py-2 border-b bg-muted/50 text-xs sm:text-sm text-muted-foreground font-medium">
            <span class="w-10 sm:w-12 shrink-0">ID</span>
            <span class="w-6 shrink-0"></span>
            <span class="w-24 sm:w-32 shrink-0">名称</span>
            <span class="w-24 sm:w-28 shrink-0">IP</span>
            <span class="w-20 sm:w-32 shrink-0 hidden md:block">主机名</span>
            <span class="w-20 sm:w-36 shrink-0 hidden lg:block">版本</span>
            <span class="w-40 shrink-0 hidden xl:block">心跳时间</span>
            <span class="w-40 shrink-0 hidden xl:block">创建时间</span>
            <span class="flex-1 min-w-[180px] text-center sm:text-right sm:pr-2">操作</span>
          </div>
          <div class="divide-y">
            <div v-if="filteredAgents.length === 0" class="text-center py-8 text-muted-foreground">
              <Server class="h-8 w-8 mx-auto mb-2 opacity-50" />
              {{ searchQuery ? '无匹配结果' : '暂无 Agent' }}
            </div>
            <!-- 小屏布局 -->
            <div v-for="agent in filteredAgents" :key="agent.id"
              class="sm:hidden p-3 hover:bg-muted/50 transition-colors">
              <div class="flex items-start justify-between mb-2">
                <div class="flex items-center gap-2 flex-1 min-w-0">
                  <span class="text-xs text-muted-foreground shrink-0">#{{ agent.id }}</span>
                  <span class="flex items-center shrink-0" :title="isOnline(agent) ? '在线' : '离线'">
                    <div v-if="isOnline(agent)"
                      class="h-5 w-5 rounded-full bg-green-500/10 flex items-center justify-center">
                      <WifiIcon class="h-3 w-3 text-green-500" />
                    </div>
                    <div v-else class="h-5 w-5 rounded-full bg-muted flex items-center justify-center">
                      <WifiOffIcon class="h-3 w-3 text-muted-foreground" />
                    </div>
                  </span>
                  <span class="font-medium text-sm truncate cursor-pointer hover:text-primary"
                    @click="viewDetail(agent)" :title="agent.name">{{ agent.name }}</span>
                </div>
                <div class="flex items-center gap-2 shrink-0 ml-2">
                  <span class="cursor-pointer group" @click="toggleEnabled(agent)"
                    :title="agent.enabled ? '点击禁用' : '点击启用'">
                    <div v-if="agent.enabled"
                      class="h-6 w-6 rounded-md bg-green-500/10 flex items-center justify-center group-hover:bg-green-500/20 transition-colors">
                      <ZapIcon class="h-3.5 w-3.5 text-green-500 fill-green-500" />
                    </div>
                    <div v-else
                      class="h-6 w-6 rounded-md bg-muted flex items-center justify-center group-hover:bg-muted/80 transition-colors">
                      <ZapOffIcon class="h-3.5 w-3.5 text-muted-foreground" />
                    </div>
                  </span>
                  <Button variant="ghost" size="icon" class="h-7 w-7" @click="viewDetail(agent)" title="详情">
                    <Eye class="h-3.5 w-3.5" />
                  </Button>
                  <Button variant="ghost" size="icon" class="h-7 w-7" @click="viewTasks(agent)" title="查看任务">
                    <ListTodo class="h-3.5 w-3.5" />
                  </Button>
                </div>
              </div>
              <div class="space-y-1 text-xs text-muted-foreground">
                <div class="flex items-center gap-2">
                  <span class="w-12 shrink-0">IP:</span>
                  <span class="truncate">{{ agent.ip || '-' }}</span>
                </div>
                <div class="flex items-center gap-2">
                  <span class="w-12 shrink-0">主机:</span>
                  <span class="truncate">{{ agent.hostname || '-' }}</span>
                </div>
                <div class="flex items-center gap-2">
                  <span class="w-12 shrink-0">版本:</span>
                  <span class="truncate">{{ agent.version || '-' }}</span>
                </div>
              </div>
              <div class="flex items-center justify-end gap-1 mt-2 pt-2 border-t">
                <Button variant="ghost" size="sm" class="h-7 text-xs" @click="forceUpdate(agent)">
                  <RotateCw class="h-3 w-3 mr-1" />更新
                </Button>
                <Button variant="ghost" size="sm" class="h-7 text-xs" @click="openEditDialog(agent)">
                  <Edit class="h-3 w-3 mr-1" />编辑
                </Button>
                <Button variant="ghost" size="sm" class="h-7 text-xs text-destructive" @click="confirmDelete(agent)">
                  <Trash2 class="h-3 w-3 mr-1" />删除
                </Button>
              </div>
            </div>
            <!-- 大屏布局 -->
            <div v-for="agent in filteredAgents" :key="`desktop-${agent.id}`"
              class="hidden sm:flex items-center gap-2 sm:gap-4 px-3 sm:px-4 py-2 hover:bg-muted/50 transition-colors">
              <span class="w-10 sm:w-12 shrink-0 text-muted-foreground text-xs sm:text-sm">#{{ agent.id }}</span>
              <span class="w-6 shrink-0 flex justify-center">
                <span class="flex justify-center shrink-0" :title="isOnline(agent) ? '在线' : '离线'">
                  <div v-if="isOnline(agent)"
                    class="h-6 w-6 rounded-full bg-green-500/10 flex items-center justify-center">
                    <WifiIcon class="h-3.5 w-3.5 text-green-500" />
                  </div>
                  <div v-else class="h-6 w-6 rounded-full bg-muted flex items-center justify-center">
                    <WifiOffIcon class="h-3.5 w-3.5 text-muted-foreground" />
                  </div>
                </span>
              </span>
              <span
                class="w-24 sm:w-32 shrink-0 font-medium text-xs sm:text-sm truncate cursor-pointer hover:text-primary"
                @click="viewDetail(agent)" :title="agent.name">{{ agent.name }}</span>
              <span class="w-24 sm:w-28 shrink-0 text-xs sm:text-sm text-muted-foreground truncate">{{ agent.ip || '-'
              }}</span>
              <span class="w-20 sm:w-32 shrink-0 text-xs sm:text-sm text-muted-foreground truncate hidden md:block">{{
                agent.hostname || '-' }}</span>
              <span class="w-20 sm:w-36 shrink-0 text-xs sm:text-sm text-muted-foreground truncate hidden lg:block">{{
                agent.version || '-' }}</span>
              <span class="w-40 shrink-0 text-xs sm:text-sm text-muted-foreground hidden xl:block">{{ agent.last_seen ||
                '-' }}</span>
              <span class="w-40 shrink-0 text-xs sm:text-sm text-muted-foreground hidden xl:block">{{ agent.created_at
                || '-' }}</span>
              <span class="flex-1 min-w-[180px] flex justify-end gap-1 sm:gap-2 items-center">
                <span class="cursor-pointer group shrink-0" @click="toggleEnabled(agent)"
                  :title="agent.enabled ? '点击禁用' : '点击启用'">
                  <div v-if="agent.enabled"
                    class="h-6 w-6 rounded-md bg-green-500/10 flex items-center justify-center group-hover:bg-green-500/20 transition-colors">
                    <ZapIcon class="h-3.5 w-3.5 text-green-500 fill-green-500" />
                  </div>
                  <div v-else
                    class="h-6 w-6 rounded-md bg-muted flex items-center justify-center group-hover:bg-muted/80 transition-colors">
                    <ZapOffIcon class="h-3.5 w-3.5 text-muted-foreground" />
                  </div>
                </span>
                <Button variant="ghost" size="icon" class="h-7 w-7 shrink-0" @click="viewDetail(agent)" title="详情">
                  <Eye class="h-3.5 w-3.5" />
                </Button>
                <Button variant="ghost" size="icon" class="h-7 w-7 shrink-0" @click="viewTasks(agent)" title="查看任务">
                  <ListTodo class="h-3.5 w-3.5" />
                </Button>
                <Button variant="ghost" size="icon" class="h-7 w-7 shrink-0" @click="forceUpdate(agent)" title="强制更新">
                  <RotateCw class="h-3.5 w-3.5" />
                </Button>
                <Button variant="ghost" size="icon" class="h-7 w-7 shrink-0" @click="openEditDialog(agent)" title="编辑">
                  <Edit class="h-3.5 w-3.5" />
                </Button>
                <Button variant="ghost" size="icon" class="h-7 w-7 text-destructive shrink-0" @click="confirmDelete(agent)"
                  title="删除">
                  <Trash2 class="h-3.5 w-3.5" />
                </Button>
              </span>
            </div>
          </div>
          </div>
        </div>
      </TabsContent>

      <TabsContent value="regcodes" class="mt-4">
        <div class="rounded-lg border bg-card overflow-x-auto hide-scrollbar">
          <div class="min-w-full w-max">
            <div
              class="flex items-center gap-2 sm:gap-4 px-3 sm:px-4 py-2 border-b bg-muted/50 text-xs sm:text-sm text-muted-foreground font-medium">
            <span class="w-6 shrink-0"></span>
            <span class="flex-1 min-w-[200px]">令牌</span>
            <span class="w-24 sm:w-32 shrink-0">备注</span>
            <span class="w-16 sm:w-20 shrink-0 text-center">使用次数</span>
            <span class="w-28 sm:w-36 shrink-0 hidden sm:block">过期时间</span>
            <span class="w-20 shrink-0 flex justify-center">
              <Button size="sm" class="h-7" @click="showTokenDialog = true">
                <Plus class="h-3.5 w-3.5 mr-1" />生成
              </Button>
            </span>
          </div>
          <div class="divide-y">
            <div v-if="tokens.length === 0" class="text-center py-8 text-muted-foreground">
              <Ticket class="h-8 w-8 mx-auto mb-2 opacity-50" />暂无令牌
            </div>
            <div v-for="token in tokens" :key="token.id"
              class="flex items-center gap-2 sm:gap-4 px-3 sm:px-4 py-2 hover:bg-muted/50 transition-colors">
              <span class="w-6 shrink-0 flex justify-center">
                <div v-if="!isTokenExpired(token) && !isTokenExhausted(token)"
                  class="h-5 w-5 rounded-full bg-green-500/10 flex items-center justify-center">
                  <Check class="h-3 w-3 text-green-500 stroke-[3]" />
                </div>
                <div v-else class="h-5 w-5 rounded-full bg-red-500/10 flex items-center justify-center">
                  <X class="h-3 w-3 text-red-500 stroke-[3]" />
                </div>
              </span>
              <code
                class="flex-1 min-w-[200px] font-mono text-xs bg-muted px-2 py-0.5 rounded truncate">{{ token.token }}</code>
              <span class="w-24 sm:w-32 shrink-0 text-xs sm:text-sm text-muted-foreground truncate">{{ token.remark ||
                '-' }}</span>
              <span class="w-16 sm:w-20 shrink-0 text-xs sm:text-sm text-muted-foreground text-center">
                {{ token.used_count }}/{{ token.max_uses === 0 ? '∞' : token.max_uses }}
              </span>
              <span class="w-28 sm:w-36 shrink-0 text-xs sm:text-sm text-muted-foreground hidden sm:block truncate">
                {{ token.expires_at || '永不过期' }}
              </span>
              <span class="w-20 shrink-0 flex justify-center gap-1">
                <Button variant="ghost" size="icon" class="h-7 w-7" @click="copyToken(token.token)" title="复制">
                  <Copy class="h-3.5 w-3.5" />
                </Button>
                <Button variant="ghost" size="icon" class="h-7 w-7 text-destructive" @click="deleteToken(token.id)"
                  title="删除">
                  <Trash2 class="h-3.5 w-3.5" />
                </Button>
              </span>
            </div>
          </div>
          </div>
        </div>
      </TabsContent>
    </Tabs>

    <!-- 详情对话框 -->
    <Dialog v-model:open="showDetailDialog">
      <DialogContent class="sm:max-w-md md:max-w-lg" @openAutoFocus.prevent>
        <DialogHeader>
          <DialogTitle>Agent 详情</DialogTitle>
          <DialogDescription class="sr-only">显示 Agent 的详细配置和状态信息</DialogDescription>
        </DialogHeader>
        <div v-if="viewingAgent" class="space-y-3">
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
            <div class="flex items-center justify-between sm:block">
              <Label class="text-muted-foreground text-xs">ID</Label>
              <div class="text-sm font-medium">#{{ viewingAgent.id }}</div>
            </div>
            <div class="flex items-center justify-between sm:block">
              <Label class="text-muted-foreground text-xs">名称</Label>
              <div class="text-sm font-medium">{{ viewingAgent.name }}</div>
            </div>
            <div class="flex items-center justify-between sm:block">
              <Label class="text-muted-foreground text-xs">IP 地址</Label>
              <div class="text-sm">{{ viewingAgent.ip || '-' }}</div>
            </div>
            <div class="flex items-center justify-between sm:block">
              <Label class="text-muted-foreground text-xs">主机名</Label>
              <div class="text-sm">{{ viewingAgent.hostname || '-' }}</div>
            </div>
            <div class="flex items-center justify-between sm:block">
              <Label class="text-muted-foreground text-xs">操作系统</Label>
              <div class="text-sm">{{ viewingAgent.os || '-' }}</div>
            </div>
            <div class="flex items-center justify-between sm:block">
              <Label class="text-muted-foreground text-xs">架构</Label>
              <div class="text-sm">{{ viewingAgent.arch || '-' }}</div>
            </div>
            <div class="flex items-center justify-between sm:block">
              <Label class="text-muted-foreground text-xs">版本</Label>
              <div class="text-sm">{{ viewingAgent.version || '-' }}</div>
            </div>
            <div class="flex items-center justify-between sm:block">
              <Label class="text-muted-foreground text-xs">构建时间</Label>
              <div class="text-sm">{{ viewingAgent.build_time || '-' }}</div>
            </div>
            <div class="flex items-center justify-between sm:block">
              <Label class="text-muted-foreground text-xs">在线状态</Label>
              <div class="flex items-center gap-2">
                <WifiIcon v-if="isOnline(viewingAgent)" class="h-4 w-4 text-green-500" />
                <WifiOffIcon v-else class="h-4 w-4 text-muted-foreground" />
                <span class="text-sm">{{ isOnline(viewingAgent) ? '在线' : '离线' }}</span>
              </div>
            </div>
            <div class="flex items-center justify-between sm:block">
              <Label class="text-muted-foreground text-xs">启用状态</Label>
              <div class="text-sm">{{ viewingAgent.enabled ? '已启用' : '已禁用' }}</div>
            </div>
            <div class="flex items-center justify-between sm:block">
              <Label class="text-muted-foreground text-xs">最后心跳</Label>
              <div class="text-sm">{{ viewingAgent.last_seen || '-' }}</div>
            </div>
            <div class="flex items-center justify-between sm:block">
              <Label class="text-muted-foreground text-xs">注册时间</Label>
              <div class="text-sm">{{ viewingAgent.created_at || '-' }}</div>
            </div>
          </div>
          <div v-if="viewingAgent.description" class="pt-2 border-t">
            <Label class="text-muted-foreground text-xs">描述</Label>
            <div class="text-sm mt-1">{{ viewingAgent.description }}</div>
          </div>
        </div>
      </DialogContent>
    </Dialog>

    <!-- 编辑对话框 -->
    <Dialog v-model:open="showEditDialog">
      <DialogContent @openAutoFocus.prevent>
        <DialogHeader>
          <DialogTitle>编辑 Agent</DialogTitle>
          <DialogDescription class="sr-only">修改 Agent 的名称和描述信息</DialogDescription>
        </DialogHeader>
        <div class="space-y-4">
          <div>
            <Label>名称</Label>
            <Input v-model="formData.name" placeholder="Agent 名称" />
          </div>
          <div>
            <Label>描述</Label>
            <Input v-model="formData.description" placeholder="描述信息（可选）" />
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="showEditDialog = false">取消</Button>
          <Button @click="updateAgent">保存</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- 删除确认对话框 -->
    <AlertDialog v-model:open="showDeleteDialog">
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>确认删除</AlertDialogTitle>
          <AlertDialogDescription>
            确定要删除 Agent "{{ deletingAgent?.name }}" 吗？此操作无法撤销。
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>取消</AlertDialogCancel>
          <AlertDialogAction class="bg-destructive text-white hover:bg-destructive/90" @click="deleteAgent">删除
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>

    <!-- 下载对话框 -->
    <Dialog v-model:open="showDownloadDialog">
      <DialogContent class="sm:max-w-lg" @openAutoFocus.prevent>
        <DialogHeader>
          <DialogTitle>下载 Agent</DialogTitle>
          <DialogDescription>当前版本: {{ agentVersion }}</DialogDescription>
        </DialogHeader>
        <div class="space-y-4">
          <div class="bg-blue-500/10 text-blue-600 dark:text-blue-400 p-3 rounded-md text-sm border border-blue-500/20">
            <p class="font-medium mb-1">💡 下载说明：</p>
            <ul class="list-disc list-inside space-y-1 text-xs opacity-90">
              <li>若主程序为 <strong class="font-semibold">Docker 部署</strong>：支持直接在此处下载包含配置的自动打包程序。</li>
              <li>若主程序为 <strong class="font-semibold">单文件二进制部署</strong>：面板无法直接提供完整打包下载，请前往 <a href="https://github.com/engigu/baihu-panel/releases" target="_blank" class="underline font-medium hover:text-blue-500 transition-colors">GitHub Releases</a> 手动下载对应的 Agent。</li>
            </ul>
          </div>
          <div class="space-y-2">
            <div v-for="platform in platforms" :key="`${platform.os}-${platform.arch}`"
              class="flex items-center justify-between p-3 border rounded-lg hover:bg-muted/50 transition-colors">
              <span class="font-medium">{{ getPlatformLabel(platform.os, platform.arch) }}</span>
              <Button size="sm" @click="downloadAgent(platform.os, platform.arch)">
                <Download class="h-4 w-4 mr-1.5" />下载
              </Button>
            </div>
          </div>
          <div class="border-t pt-4">
            <h4 class="font-medium mb-2">使用说明</h4>
            <ol class="text-sm text-muted-foreground space-y-1.5 list-decimal list-inside">
              <li>下载对应平台的 Agent 压缩包并解压</li>
              <li>复制 <code class="bg-muted px-1.5 py-0.5 rounded text-foreground">config.example.ini</code> 为 <code
                  class="bg-muted px-1.5 py-0.5 rounded text-foreground">config.ini</code></li>
              <li>编辑 <code class="bg-muted px-1.5 py-0.5 rounded text-foreground">config.ini</code>，填写服务器地址和注册令牌</li>
              <li>运行 <code class="bg-muted px-1.5 py-0.5 rounded text-foreground">./baihu-agent start</code> 启动（后台运行）
              </li>
            </ol>
            <div class="mt-3 text-sm text-muted-foreground">
              <p class="font-medium text-foreground mb-1.5">常用命令：</p>
              <div class="space-y-1">
                <div><code class="bg-muted px-1.5 py-0.5 rounded text-foreground text-xs">baihu-agent start</code> <span
                    class="text-xs">- 后台启动</span></div>
                <div><code class="bg-muted px-1.5 py-0.5 rounded text-foreground text-xs">baihu-agent stop</code> <span
                    class="text-xs">- 停止运行</span></div>
                <div><code class="bg-muted px-1.5 py-0.5 rounded text-foreground text-xs">baihu-agent status</code>
                  <span class="text-xs">- 查看状态</span>
                </div>
                <div><code class="bg-muted px-1.5 py-0.5 rounded text-foreground text-xs">baihu-agent logs</code> <span
                    class="text-xs">- 查看日志</span></div>
                <div><code class="bg-muted px-1.5 py-0.5 rounded text-foreground text-xs">baihu-agent run</code> <span
                    class="text-xs">- 前台运行</span></div>
              </div>
            </div>
          </div>
        </div>
      </DialogContent>
    </Dialog>

    <!-- 创建令牌对话框 -->
    <Dialog v-model:open="showTokenDialog">
      <DialogContent @openAutoFocus.prevent>
        <DialogHeader>
          <DialogTitle>生成令牌</DialogTitle>
          <DialogDescription class="sr-only">创建一个新的注册令牌，用于 Agent 认证</DialogDescription>
        </DialogHeader>
        <div class="space-y-4">
          <div>
            <Label>备注</Label>
            <Input v-model="tokenForm.remark" placeholder="备注信息（可选）" />
          </div>
          <div>
            <Label>最大使用次数</Label>
            <Input v-model.number="tokenForm.max_uses" type="number" placeholder="0 表示无限制" />
          </div>
          <div>
            <Label>过期时间</Label>
            <Input v-model="tokenForm.expires_at" type="datetime-local" />
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="showTokenDialog = false">取消</Button>
          <Button @click="createToken">生成</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
