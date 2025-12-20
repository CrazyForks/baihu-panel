<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ListTodo, FileCode, Variable, Clock, Play, ScrollText } from 'lucide-vue-next'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { api, type Task, type Stats } from '@/api'

const router = useRouter()
const stats = ref<Stats>({ tasks: 0, scripts: 0, envs: 0, logs: 0, scheduled: 0, running: 0 })
const recentTasks = ref<Task[]>([])

onMounted(async () => {
  try {
    const [statsData, tasksRes] = await Promise.all([
      api.dashboard.stats(),
      api.tasks.list({ page: 1, page_size: 5 })
    ])
    stats.value = statsData
    recentTasks.value = tasksRes.data
  } catch {}
})

const statItems = [
  { key: 'tasks', label: '任务总数', icon: ListTodo, route: '/tasks' },
  { key: 'scripts', label: '脚本数量', icon: FileCode, route: '/editor' },
  { key: 'envs', label: '环境变量', icon: Variable, route: '/environments' },
  { key: 'logs', label: '日志总数', icon: ScrollText, route: '/history' },
  { key: 'scheduled', label: '调度注册', icon: Clock, route: '/tasks' },
  { key: 'running', label: '正在运行', icon: Play, route: '/tasks' },
]

function navigateTo(route?: string) {
  if (route) router.push(route)
}
</script>

<template>
  <div class="space-y-6">
    <div>
      <h2 class="text-2xl font-bold tracking-tight">数据仪表</h2>
      <p class="text-muted-foreground">查看系统运行状态和统计数据</p>
    </div>

    <div class="grid gap-4 md:grid-cols-3 lg:grid-cols-6">
      <Card v-for="item in statItems" :key="item.key" class="cursor-pointer hover:bg-accent/50 transition-colors" @click="navigateTo(item.route)">
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">{{ item.label }}</CardTitle>
          <component :is="item.icon" class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ stats[item.key as keyof Stats] }}</div>
        </CardContent>
      </Card>
    </div>

    <Card>
      <CardHeader>
        <CardTitle>最近任务</CardTitle>
        <CardDescription>最近创建的定时任务</CardDescription>
      </CardHeader>
      <CardContent>
        <div class="space-y-2">
          <div v-if="recentTasks.length === 0" class="text-sm text-muted-foreground py-8 text-center">
            暂无任务
          </div>
          <div
            v-for="task in recentTasks"
            :key="task.id"
            class="flex items-center justify-between py-2 border-b last:border-0"
          >
            <div class="flex items-center gap-3">
              <span class="text-muted-foreground text-sm">#{{ task.id }}</span>
              <span class="font-medium">{{ task.name }}</span>
            </div>
            <div class="flex items-center gap-4 text-muted-foreground">
              <code class="text-xs bg-muted px-2 py-0.5 rounded">{{ task.schedule }}</code>
              <span class="w-2 h-2 rounded-full" :class="task.enabled ? 'bg-green-500' : 'bg-gray-400'" />
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
