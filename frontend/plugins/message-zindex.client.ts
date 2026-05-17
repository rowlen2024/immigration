import { ElMessage } from 'element-plus'

const HIGH_ZINDEX = 3000

export default defineNuxtPlugin(() => {
  const methods = ['error', 'success', 'warning', 'info'] as const

  for (const method of methods) {
    const original = (ElMessage as any)[method] as Function
    if (!original) continue
    ;(ElMessage as any)[method] = (msg: string) => {
      return original({ message: msg, zIndex: HIGH_ZINDEX })
    }
  }
})
