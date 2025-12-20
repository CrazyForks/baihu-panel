<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Badge } from '@/components/ui/badge'
import { ExternalLink } from 'lucide-vue-next'
import { api, type AboutInfo } from '@/api'

const aboutInfo = ref<AboutInfo | null>(null)

const techStack = ['Golang', 'Vue 3', 'TypeScript', 'Vite', 'Tailwind CSS', 'Shadcn/ui']
const features = ['脚本管理', '定时任务', '在线终端', '执行日志', '环境变量', 'Docker部署']

async function loadAbout() {
  try {
    aboutInfo.value = await api.settings.getAbout()
  } catch {}
}

onMounted(loadAbout)
</script>

<template>
  <div>
    <!-- 站点关于 -->
    <div class="mb-8">
      <h3 class="text-xl font-semibold mb-2">白虎面板</h3>
      <p class="text-muted-foreground">一个轻量级的定时任务管理系统，支持脚本管理、定时执行和日志追踪。</p>
    </div>

    <div class="grid md:grid-cols-2 gap-x-16 gap-y-8">
      <!-- 左侧：技术栈和功能特性 -->
      <div class="space-y-8">
        <div>
          <h4 class="text-sm font-medium mb-4">技术栈</h4>
          <div class="flex flex-wrap gap-2">
            <Badge v-for="tech in techStack" :key="tech" variant="secondary">{{ tech }}</Badge>
          </div>
        </div>

        <div>
          <h4 class="text-sm font-medium mb-4">功能特性</h4>
          <div class="flex flex-wrap gap-2">
            <Badge v-for="feature in features" :key="feature" variant="outline">{{ feature }}</Badge>
          </div>
        </div>
      </div>

      <!-- 右侧：系统信息 -->
      <div>
        <h4 class="text-sm font-medium mb-4">系统信息</h4>
        <div class="space-y-3">
          <div class="flex justify-between items-center">
            <span class="text-muted-foreground text-sm">系统版本:</span>
            <Badge variant="outline" class="font-mono">{{ aboutInfo?.version || 'dev' }}</Badge>
          </div>
          <div class="flex justify-between items-center">
            <span class="text-muted-foreground text-sm">构建时间:</span>
            <span class="text-sm">{{ aboutInfo?.build_time || '-' }}</span>
          </div>
          <div class="flex justify-between items-center">
            <span class="text-muted-foreground text-sm">内存使用:</span>
            <span class="text-sm">{{ aboutInfo?.mem_usage || '-' }}</span>
          </div>
          <div class="flex justify-between items-center">
            <span class="text-muted-foreground text-sm">运行时间:</span>
            <span class="text-sm">{{ aboutInfo?.uptime || '-' }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部：版权和链接 -->
    <div class="mt-10 pt-6 border-t flex items-center justify-center gap-2 text-sm text-muted-foreground">
      <span>© 2025 保留所有权利。</span>
      <a href="https://github.com" target="_blank" class="inline-flex items-center gap-1 text-primary hover:underline">
        <ExternalLink class="h-3.5 w-3.5" />
        GitHub 仓库
      </a>
    </div>
  </div>
</template>
