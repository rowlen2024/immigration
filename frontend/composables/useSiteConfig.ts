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
  contact_email: string;
  contact_address: string;
  contact_wechat: string;
  ga_tracking_id: string;
  baidu_tongji_id: string;
  custom_head_code: string;
  custom_body_code: string;
  copyright_text: string;
  icp_number: string;
  footer_tagline: string;
}

const siteConfig = ref<SiteConfig | null>(null);
let fetchPromise: Promise<void> | null = null;

export const useSiteConfig = () => {
  const fetch = async () => {
    if (siteConfig.value) return;
    if (fetchPromise) {
      await fetchPromise;
      return;
    }

    fetchPromise = (async () => {
      try {
        const api = useApi();
        siteConfig.value = await api<SiteConfig>('/site-config');
      } catch {
        siteConfig.value = null;
      }
    })();

    await fetchPromise;
    fetchPromise = null;
  };

  return {
    siteConfig: readonly(siteConfig),
    fetch,
  };
};
