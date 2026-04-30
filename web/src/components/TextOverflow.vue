<script setup lang="ts">
import { ref } from 'vue'
import BaihuDialog from '@/components/ui/BaihuDialog.vue'

const props = withDefaults(
  defineProps<{
    text: string
    title?: string
  }>(),
  {
    title: '详情'
  }
)

const showDialog = ref(false)

function handleClick() {
  if (props.text && props.text !== '-') {
    showDialog.value = true
  }
}
</script>

<template>
  <span v-bind="$attrs" class="truncate block cursor-pointer hover:text-primary transition-colors" :title="text || '-'"
    @click.stop="handleClick">
    {{ text || '-' }}
  </span>

  <BaihuDialog v-model:open="showDialog" :title="title">
    <div class="max-h-[60vh] overflow-y-auto custom-scrollbar">
      <div class="p-4 bg-muted/30 rounded-xl border border-border/50">
        <p class="text-[13.5px] leading-relaxed text-foreground/90 break-all whitespace-pre-wrap">
          {{ text }}
        </p>
      </div>
    </div>
  </BaihuDialog>
</template>
