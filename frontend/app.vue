<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>

<script setup lang="ts">
const { siteConfig, fetch: fetchSiteConfig } = useSiteConfig();
const route = useRoute();

onMounted(() => { fetchSiteConfig(); });

const isPublicPage = computed(() => !route.path.startsWith('/admin'));

// Head scripts (GA, Baidu, custom_head_code) & favicon
useHead(() => {
  const links: Record<string, unknown>[] = [];
  const scripts: Record<string, unknown>[] = [];

  if (siteConfig.value?.site_favicon) {
    links.push({ rel: 'icon', type: 'image/x-icon', href: siteConfig.value.site_favicon });
  }

  if (siteConfig.value?.ga_tracking_id) {
    scripts.push({
      async: true,
      src: `https://www.googletagmanager.com/gtag/js?id=${siteConfig.value.ga_tracking_id}`,
    });
    scripts.push({
      innerHTML: `window.dataLayer=window.dataLayer||[];function gtag(){dataLayer.push(arguments);}gtag('js',new Date());gtag('config','${siteConfig.value.ga_tracking_id}');`,
    });
  }

  if (siteConfig.value?.baidu_tongji_id) {
    scripts.push({
      innerHTML: `var _hmt=_hmt||[];(function(){var hm=document.createElement("script");hm.src="https://hm.baidu.com/hm.js?${siteConfig.value.baidu_tongji_id}";var s=document.getElementsByTagName("script")[0];s.parentNode.insertBefore(hm,s);})();`,
    });
  }

  if (siteConfig.value?.custom_head_code) {
    scripts.push({
      innerHTML: siteConfig.value.custom_head_code,
    });
  }

  return { link: links, script: scripts };
});

// Body-close scripts (custom_body_code)
useHead(() => {
  if (!siteConfig.value?.custom_body_code) return {};

  return {
    script: [
      {
        innerHTML: siteConfig.value.custom_body_code,
        tagPosition: 'bodyClose',
      },
    ],
  };
});

// Organization schema (public pages only)
useHead(() => {
  if (!isPublicPage.value) return {};
  const org = siteConfig.value;
  if (!org?.organization_name) return {};

  const jsonLd: Record<string, unknown> = {
    '@context': 'https://schema.org',
    '@type': 'Organization',
    name: org.organization_name,
    url: org.organization_url || '',
  };

  if (org.organization_description) {
    jsonLd.description = org.organization_description;
  }
  if (org.organization_logo) {
    jsonLd.logo = org.organization_logo;
  }
  if (org.same_as && org.same_as.length > 0) {
    jsonLd.sameAs = org.same_as.filter((s: string) => s.trim() !== '');
  }
  if (org.contact_phone || org.contact_email) {
    jsonLd.contactPoint = {
      '@type': 'ContactPoint',
      telephone: org.contact_phone || '',
      email: org.contact_email || '',
      contactType: 'customer service',
    };
  }
  if (org.contact_address) {
    jsonLd.address = {
      '@type': 'PostalAddress',
      streetAddress: org.contact_address,
    };
  }

  return {
    script: [
      {
        type: 'application/ld+json',
        innerHTML: JSON.stringify(jsonLd),
      },
    ],
  };
});

// WebSite schema with SearchAction (public pages only)
useHead(() => {
  if (!isPublicPage.value) return {};
  const config = siteConfig.value;
  if (!config?.site_name) return {};

  const siteUrl = config.canonical_base || config.organization_url || '';
  const webSite: Record<string, unknown> = {
    '@context': 'https://schema.org',
    '@type': 'WebSite',
    name: config.site_name,
    url: siteUrl,
  };

  if (siteUrl) {
    webSite.potentialAction = {
      '@type': 'SearchAction',
      target: {
        '@type': 'EntryPoint',
        urlTemplate: `${siteUrl}/search?q={search_term_string}`,
      },
      'query-input': 'required name=search_term_string',
    };
  }

  return {
    script: [
      {
        type: 'application/ld+json',
        innerHTML: JSON.stringify(webSite),
      },
    ],
  };
});
</script>
