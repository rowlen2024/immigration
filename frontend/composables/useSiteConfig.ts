interface SiteConfig {
  site_name: string;
  site_logo: string;
  site_favicon: string;
  seo_title: string;
  seo_description: string;
  seo_keywords: string;
  og_image: string;
  canonical_base: string;
  organization_name: string;
  organization_description: string;
  organization_logo: string;
  organization_url: string;
  same_as: string[];
  contact_phone: string;
  contact_phone_2: string;
  contact_email: string;
  contact_address: string;
  contact_wechat: string;
  contact_wechat_mp: string;
  contact_wechat_video: string;
  ga_tracking_id: string;
  baidu_tongji_id: string;
  custom_head_code: string;
  custom_body_code: string;
  copyright_text: string;
  icp_number: string;
  footer_tagline: string;
}

export const useMygoSiteConfig = () => {
  const { data } = useFetch('/api/v1/site-config', {
    key: 'public:site-config',
    transform: (response: any) => response?.data ?? response,
  })

  const siteConfig = computed<SiteConfig | null>(() => (data.value as SiteConfig) ?? null)
  usePublicDataFreshness(['public:site-config'])

  // 客户端强制刷新（绕过 Nuxt payload 缓存）
  const refreshSiteConfig = () => {
    refreshNuxtData('public:site-config')
  }

  return { siteConfig, refreshSiteConfig, data }
}
