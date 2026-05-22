/**
 * SSR 期间 API 代理中间件
 * Nuxt 服务端渲染时 useFetch('/api/...') 是进程内请求，
 * 不经过 nginx/Vite 反向代理，需要此中间件转发到 Go 后端。
 */
export default defineEventHandler(async (event) => {
  const path = event.path
  if (!path.startsWith('/api/')) return

  const backend = process.env.BACKEND_URL || 'http://localhost:8080'
  const start = Date.now()

  try {
    const data = await $fetch(event.path, {
      baseURL: backend,
      method: event.method,
      query: getQuery(event),
      headers: getHeaders(event),
      body: event.method !== 'GET' && event.method !== 'HEAD'
        ? await readBody(event).catch(() => undefined)
        : undefined,
    })
    const ms = Date.now() - start
    console.log(`[SSR API] [${start}] ${event.method} ${event.path} → ${backend} [${ms}ms]`)
    return data
  } catch (err: any) {
    const ms = Date.now() - start
    console.error(`[SSR API] [${start}] ${event.method} ${event.path} → ${backend} 失败 [${ms}ms]:`, err.message)
    throw err
  }
})
