// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  ssr: false,

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
      ],
    },
  },
});
