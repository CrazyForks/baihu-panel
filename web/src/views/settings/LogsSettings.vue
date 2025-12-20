<script setup lang="ts">
import { ref } from 'vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { api } from '@/api'
import { toast } from 'vue-sonner'

const cleanDays = ref(30)
const cleanResult = ref<number | null>(null)

async function cleanLogs() {
  if (cleanDays.value < 1) {
    toast.error('天数必须大于0')
    return
  }
  try {
    const res = await api.settings.cleanLogs(cleanDays.value)
    cleanResult.value = res.deleted
    toast.success(`已清理 ${res.deleted} 条日志`)
  } catch {
    toast.error('清理失败')
  }
}
</script>

<template>
  <div class="space-y-4">
    <div class="space-y-2">
      <Label>清理多少天前的日志</Label>
      <Input v-model.number="cleanDays" type="number" class="w-32" min="1" />
    </div>
    <Button variant="destructive" @click="cleanLogs">清理日志</Button>
    <p v-if="cleanResult !== null" class="text-sm text-muted-foreground">
      上次清理了 {{ cleanResult }} 条日志
    </p>
  </div>
</template>
