<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { Terminal } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import '@xterm/xterm/css/xterm.css'
import { Button } from '@/components/ui/button'
import { RefreshCw } from 'lucide-vue-next'

const terminalRef = ref<HTMLDivElement | null>(null)
let terminal: Terminal | null = null
let fitAddon: FitAddon | null = null
let ws: WebSocket | null = null
let isPtyMode = false // 是否是 PTY 模式（Unix）
let inputBuffer = ''
let commandHistory: string[] = []
let historyIndex = -1

function initTerminal() {
  if (!terminalRef.value || terminal) return

  terminal = new Terminal({
    cursorBlink: true,
    fontSize: 13,
    fontFamily: 'Consolas, Monaco, monospace',
    theme: {
      background: '#1e1e1e',
      foreground: '#d4d4d4',
      cursor: '#d4d4d4',
    }
  })

  fitAddon = new FitAddon()
  terminal.loadAddon(fitAddon)
  terminal.open(terminalRef.value)
  fitAddon.fit()
  terminal.focus()

  connectWebSocket()

  // 清除当前输入行（Windows 模式用）
  function clearLine() {
    for (let i = 0; i < inputBuffer.length; i++) {
      terminal?.write('\b \b')
    }
  }

  // 处理用户输入
  terminal.onData((data) => {
    if (!ws || ws.readyState !== WebSocket.OPEN) return

    // PTY 模式：直接透传所有输入
    if (isPtyMode) {
      ws.send(data)
      return
    }

    // Windows 模式：本地处理输入和历史记录
    // 回车键
    if (data === '\r') {
      terminal?.write('\r\n')
      if (inputBuffer.trim()) {
        commandHistory.push(inputBuffer)
        historyIndex = commandHistory.length
        ws.send(inputBuffer + '\r\n')
      }
      inputBuffer = ''
    }
    // 上箭头 - 上一条历史
    else if (data === '\x1b[A') {
      if (commandHistory.length > 0 && historyIndex > 0) {
        clearLine()
        historyIndex--
        inputBuffer = commandHistory[historyIndex] ?? ''
        terminal?.write(inputBuffer)
      }
    }
    // 下箭头 - 下一条历史
    else if (data === '\x1b[B') {
      clearLine()
      if (historyIndex < commandHistory.length - 1) {
        historyIndex++
        inputBuffer = commandHistory[historyIndex] ?? ''
        terminal?.write(inputBuffer)
      } else {
        historyIndex = commandHistory.length
        inputBuffer = ''
      }
    }
    // 退格键
    else if (data === '\x7f' || data === '\b') {
      if (inputBuffer.length > 0) {
        inputBuffer = inputBuffer.slice(0, -1)
        terminal?.write('\b \b')
      }
    }
    // Ctrl+C
    else if (data === '\x03') {
      ws.send('\x03')
      inputBuffer = ''
      historyIndex = commandHistory.length
      terminal?.write('^C\r\n')
    }
    // 普通字符
    else if (data >= ' ' || data === '\t') {
      inputBuffer += data
      terminal?.write(data)
    }
  })
}

function connectWebSocket() {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const wsUrl = `${protocol}//${window.location.host}/api/terminal/ws`
  
  ws = new WebSocket(wsUrl)
  
  ws.onopen = () => {
    terminal?.writeln('\x1b[32m已连接到终端\x1b[0m')
    terminal?.writeln('')
    terminal?.focus()
  }
  
  ws.onmessage = (event) => {
    // 检查是否是 PTY 模式标识
    if (event.data === '__PTY_MODE__') {
      isPtyMode = true
      return
    }
    if (event.data === '__PIPE_MODE__') {
      isPtyMode = false
      return
    }
    terminal?.write(event.data)
  }
  
  ws.onclose = () => {
    terminal?.writeln('')
    terminal?.writeln('\x1b[31m连接已断开\x1b[0m')
  }
  
  ws.onerror = () => {
    terminal?.writeln('\x1b[31m连接错误\x1b[0m')
  }
}

function reconnect() {
  if (ws) {
    ws.close()
  }
  inputBuffer = ''
  isPtyMode = false
  terminal?.clear()
  connectWebSocket()
}

function handleResize() {
  fitAddon?.fit()
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
  setTimeout(initTerminal, 100)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  ws?.close()
  terminal?.dispose()
  terminal = null
  ws = null
})
</script>

<template>
  <div class="flex flex-col h-[calc(100vh-100px)]">
    <div class="flex items-center justify-between p-2 border rounded-t-md bg-[#252526]">
      <span class="text-xs font-medium text-gray-300">终端</span>
      <Button variant="ghost" size="icon" class="h-6 w-6 text-gray-400 hover:text-white" @click="reconnect" title="重新连接">
        <RefreshCw class="h-3 w-3" />
      </Button>
    </div>
    <div ref="terminalRef" class="terminal-container flex-1 border border-t-0 rounded-b-md bg-[#1e1e1e] p-1" />
  </div>
</template>

<style scoped>
.terminal-container :deep(.xterm-viewport) {
  scrollbar-width: thin;
  scrollbar-color: #4a4a4a #1e1e1e;
}

.terminal-container :deep(.xterm-viewport::-webkit-scrollbar) {
  width: 8px;
}

.terminal-container :deep(.xterm-viewport::-webkit-scrollbar-track) {
  background: #1e1e1e;
}

.terminal-container :deep(.xterm-viewport::-webkit-scrollbar-thumb) {
  background: #4a4a4a;
  border-radius: 4px;
}

.terminal-container :deep(.xterm-viewport::-webkit-scrollbar-thumb:hover) {
  background: #5a5a5a;
}
</style>
