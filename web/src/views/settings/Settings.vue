<script setup lang="ts">
import { ref } from 'vue'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import PasswordSettings from './PasswordSettings.vue'
import SiteSettings from './SiteSettings.vue'
import SchedulerSettings from './SchedulerSettings.vue'
import BackupSettings from './BackupSettings.vue'
import AboutSettings from './AboutSettings.vue'

const activeTab = ref('password')
</script>

<template>
  <div class="space-y-6">
    <div>
      <h2 class="text-xl sm:text-2xl font-bold tracking-tight">系统设置</h2>
      <p class="text-muted-foreground text-sm">管理系统配置和账户安全</p>
    </div>

    <Tabs v-model="activeTab" class="max-w-2xl">
      <TabsList class="w-full sm:w-auto flex sm:inline-flex h-10 gap-1 p-1 bg-muted/50 rounded-lg overflow-x-auto no-scrollbar">
        <TabsTrigger value="password" class="flex-1">密码修改</TabsTrigger>
        <TabsTrigger value="site" class="flex-1">站点设置</TabsTrigger>
        <TabsTrigger value="scheduler" class="flex-1">调度设置</TabsTrigger>
        <TabsTrigger value="backup" class="flex-1">备份恢复</TabsTrigger>
        <TabsTrigger value="about" class="flex-1">关于</TabsTrigger>
      </TabsList>

      <TabsContent value="password" class="mt-6">
        <Card>
          <CardHeader>
            <CardTitle>修改密码</CardTitle>
            <CardDescription>更新您的账户密码</CardDescription>
          </CardHeader>
          <CardContent>
            <PasswordSettings />
          </CardContent>
        </Card>
      </TabsContent>

      <TabsContent value="site" class="mt-6">
        <Card>
          <CardHeader>
            <CardTitle>站点设置</CardTitle>
            <CardDescription>配置站点标题、图标和系统参数</CardDescription>
          </CardHeader>
          <CardContent>
            <SiteSettings />
          </CardContent>
        </Card>
      </TabsContent>

      <TabsContent value="scheduler" class="mt-6">
        <Card>
          <CardHeader>
            <CardTitle>调度设置</CardTitle>
            <CardDescription>配置任务调度器参数</CardDescription>
          </CardHeader>
          <CardContent>
            <SchedulerSettings />
          </CardContent>
        </Card>
      </TabsContent>

      <TabsContent value="backup" class="mt-6">
        <Card>
          <CardHeader>
            <CardTitle>备份恢复</CardTitle>
            <CardDescription>备份和恢复系统数据</CardDescription>
          </CardHeader>
          <CardContent>
            <BackupSettings />
          </CardContent>
        </Card>
      </TabsContent>

      <TabsContent value="about" class="mt-6">
        <Card>
          <CardContent class="pt-6">
            <AboutSettings />
          </CardContent>
        </Card>
      </TabsContent>
    </Tabs>
  </div>
</template>
