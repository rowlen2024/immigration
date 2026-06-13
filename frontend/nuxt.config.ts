// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  css: ['~/assets/css/variables.css', '~/assets/css/global.css', '~/assets/css/admin.css'],

  modules: ['@pinia/nuxt', '@element-plus/nuxt'],

  experimental: {
    appManifest: false,
  },

  devServer: {
    port: 3000,
  },

  nitro: {
    devProxy: {
      '/api': 'http://localhost:8080',
      '/uploads': 'http://localhost:8080',
    },
    routeRules: {
      // Proxy API/uploads to Go backend — required for SSR data fetching (site-config etc.)
      '/api/**': { proxy: 'http://backend:8080' },
      '/uploads/**': { proxy: 'http://backend:8080' },

      // ISR — SSR on first request / after 5min TTL, cached static HTML otherwise
      '/': { swr: 300 },
      '/projects': { swr: 300 },
      '/projects/*': { swr: 300 },
      '/pages/*': { swr: 300 },
      // '/projects/cies': { swr: 300 },
      // '/projects/panama': { swr: 300 },
      '/cases': { swr: 300 },
      '/faq': { swr: 300 },
      '/contact': { swr: 300 },
      '/compare': { swr: 300 },

      // Admin — never SSR
      '/admin/**': { ssr: false },

      // All other routes — SPA fallback (no runtime SSR)
      '/**': { ssr: false },

      // Legacy redirects
      '/usa/eb5': { redirect: '/projects/eb5' },
      '/hongkong/cies': { redirect: '/projects/cies' },
      '/panama/property': { redirect: '/projects/panama' },
    },
  },

  vite: {
    server: {
      proxy: {
        '/api': {
          target: 'http://localhost:8080',
          changeOrigin: true,
        },
        '/uploads': {
          target: 'http://localhost:8080',
          changeOrigin: true,
        },
      },
    },
    css: {
      preprocessorOptions: {
        scss: {
          additionalData: '',
        },
      },
    },
  },

  app: {
    head: {
      title: '北极星移民 | 专业投资移民服务',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        {
          name: 'description',
          content:
            '北极星移民提供美国EB-5、香港投资移民、巴拿马购房移民等专业投资移民服务',
        },
      ],
      link: [
        {
          rel: 'icon',
          type: 'image/x-icon',
          href: '/favicon.ico',
        },
        {
          rel: 'preconnect',
          href: 'https://fonts.googleapis.com',
        },
        {
          rel: 'preconnect',
          href: 'https://fonts.gstatic.com',
          crossorigin: 'anonymous',
        },
        {
          rel: 'stylesheet',
          href: 'https://fonts.googleapis.com/css2?family=Noto+Sans+SC:wght@300;400;500;700&family=Noto+Serif+SC:wght@400;600;700&display=swap',
        },
      ],
    },
  },
});
