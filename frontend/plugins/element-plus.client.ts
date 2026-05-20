export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.vueApp.config.globalProperties.$ELEMENT = { zIndex: 3000 }
})
