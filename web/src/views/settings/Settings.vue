<script setup lang="ts">
import { ref } from 'vue'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import PasswordSettings from './PasswordSettings.vue'
import LogsSettings from './LogsSettings.vue'
import SiteSettings from './SiteSettings.vue'
import AboutSettings from './AboutSettings.vue'

const activeTab = ref('password')
</script>

<template>
  <div class="space-y-6">
    <div>
      <h2 class="text-2xl font-bold tracking-tight">系统设置</h2>
      <p class="text-muted-foreground">管理系统配置和账户安全</p>
    </div>

    <Tabs v-model="activeTab" class="w-full">
      <TabsList>
        <TabsTrigger value="password">密码修改</TabsTrigger>
        <TabsTrigger value="logs">日志清理</TabsTrigger>
        <TabsTrigger value="site">站点设置</TabsTrigger>
        <TabsTrigger value="about">关于</TabsTrigger>
      </TabsList>

      <TabsContent value="password" class="mt-6">
        <Card class="max-w-lg">
          <CardHeader>
            <CardTitle>修改密码</CardTitle>
            <CardDescription>更新您的账户密码</CardDescription>
          </CardHeader>
          <CardContent>
            <PasswordSettings />
          </CardContent>
        </Card>
      </TabsContent>

      <TabsContent value="logs" class="mt-6">
        <Card class="max-w-lg">
          <CardHeader>
            <CardTitle>日志清理</CardTitle>
            <CardDescription>清理历史执行日志以释放存储空间</CardDescription>
          </CardHeader>
          <CardContent>
            <LogsSettings />
          </CardContent>
        </Card>
      </TabsContent>

      <TabsContent value="site" class="mt-6">
        <Card class="max-w-xl">
          <CardHeader>
            <CardTitle>站点设置</CardTitle>
            <CardDescription>配置站点标题、图标和系统参数</CardDescription>
          </CardHeader>
          <CardContent>
            <SiteSettings />
          </CardContent>
        </Card>
      </TabsContent>

      <TabsContent value="about" class="mt-6">
        <Card class="max-w-xl">
          <CardContent class="pt-6">
            <AboutSettings />
          </CardContent>
        </Card>
      </TabsContent>
    </Tabs>
  </div>
</template>
