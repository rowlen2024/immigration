// Plugin: 为 v-html 富文本内容中的图片自动添加 lazy loading
// 处理场景：案例详情 (case/[slug].vue)、CMS 页面 (pages/[...slug].vue)、预览页面 (preview/page/[...slug].vue)

export default defineNuxtPlugin(() => {
  if (typeof window === 'undefined') return

  function processImage(img: HTMLImageElement) {
    if (img.hasAttribute('data-lazy-processed')) return
    if (img.closest('[data-no-lazy]')) return

    img.setAttribute('data-lazy-processed', '')
    img.setAttribute('loading', 'lazy')
    img.setAttribute('decoding', 'async')
  }

  function scanContainers() {
    const selectors = [
      '.case-content',
      '.page-content',
      '.detail-content',
    ]
    for (const selector of selectors) {
      document.querySelectorAll(selector).forEach((el) => {
        el.querySelectorAll('img').forEach(processImage)
      })
    }
  }

  function setupMutationObserver() {
    const observer = new MutationObserver((mutations) => {
      for (const mutation of mutations) {
        for (const node of mutation.addedNodes) {
          if (!(node instanceof HTMLElement)) continue
          if (node.tagName === 'IMG') {
            processImage(node as HTMLImageElement)
          }
          if (node.querySelectorAll) {
            node.querySelectorAll('img').forEach(processImage)
          }
        }
      }
    })

    observer.observe(document.body, { childList: true, subtree: true })
    return observer
  }

  // Initial scan
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => {
      scanContainers()
    })
  } else {
    scanContainers()
  }

  const observer = setupMutationObserver()

  // Re-scan after SPA navigation
  const router = useRouter()
  router.afterEach(() => {
    setTimeout(scanContainers, 100)
  })

  // Cleanup (belt-and-suspenders)
  onUnmounted(() => {
    observer.disconnect()
  })
})
