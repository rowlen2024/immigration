interface SeoOptions {
  title?: string;
  description?: string;
  jsonLd?: Record<string, unknown>;
  breadcrumbLabel?: string;
  robots?: string;
}

export const useSeo = (options: SeoOptions) => {
  const { siteConfig } = useSiteConfig();
  const { getBreadcrumb } = useNavigation();
  const route = useRoute();
  const siteName = computed(() => siteConfig.value?.site_name || '北极星移民');

  const fullTitle = computed(() => {
    if (options.title) {
      return `${options.title} | ${siteName.value}`;
    }
    const seoTitle = siteConfig.value?.seo_title || '{site_name} | 专业投资移民服务';
    return seoTitle.replace('{site_name}', siteName.value);
  });

  const head = computed(() => {
    const h: Record<string, unknown> = {
      title: fullTitle.value,
      meta: [] as Record<string, unknown>[],
      link: [] as Record<string, unknown>[],
      script: [] as Record<string, unknown>[],
    };

    const metas: Record<string, unknown>[] = [];
    const links: Record<string, unknown>[] = [];
    const scripts: Record<string, unknown>[] = [];

    const desc = options.description || siteConfig.value?.seo_description || '';

    if (desc) {
      metas.push(
        { name: 'description', content: desc },
        { property: 'og:title', content: fullTitle.value },
        { property: 'og:description', content: desc },
        { name: 'twitter:title', content: fullTitle.value },
        { name: 'twitter:description', content: desc },
      );
    }

    // Keywords
    if (siteConfig.value?.seo_keywords) {
      metas.push({ name: 'keywords', content: siteConfig.value.seo_keywords });
    }

    // Open Graph
    metas.push({ property: 'og:type', content: 'website' });
    metas.push({ property: 'og:site_name', content: siteName.value });
    if (siteConfig.value?.canonical_base) {
      metas.push({ property: 'og:url', content: siteConfig.value.canonical_base + route.path });
    }
    if (siteConfig.value?.og_image) {
      metas.push({ property: 'og:image', content: siteConfig.value.og_image });
    }

    // Twitter Card
    metas.push({ name: 'twitter:card', content: 'summary_large_image' });
    if (siteConfig.value?.og_image) {
      metas.push({ name: 'twitter:image', content: siteConfig.value.og_image });
    }

    // Canonical (link, not meta)
    if (siteConfig.value?.canonical_base) {
      links.push({ rel: 'canonical', href: siteConfig.value.canonical_base + route.path });
    }

    // Robots
    metas.push({ name: 'robots', content: options.robots || 'index, follow' });

    (h.meta as Record<string, unknown>[]).push(...metas);
    (h.link as Record<string, unknown>[]).push(...links);

    // JSON-LD: custom
    if (options.jsonLd) {
      scripts.push({
        type: 'application/ld+json',
        innerHTML: JSON.stringify(options.jsonLd),
      });
    }

    // JSON-LD: BreadcrumbList (auto-derived from navigation)
    const crumbs = computed(() => getBreadcrumb(route.path, options.breadcrumbLabel));
    if (crumbs.value.length > 0) {
      const base = siteConfig.value?.canonical_base || '';
      // Fix breadcrumb labels that are raw path segments (no Chinese characters)
      // This happens when navigation API data hasn't loaded yet in SSR context
      const hasCJK = (s: string) => /[一-鿿]/.test(s);
      const breadcrumbItems = crumbs.value.map((item, index) => {
        const isLast = index === crumbs.value.length - 1;
        let name = item.label;
        // If label looks like a raw slug/path segment, try to use breadcrumbLabel for last item
        if (!hasCJK(name) && isLast && options.breadcrumbLabel) {
          name = options.breadcrumbLabel;
        }
        return {
          '@type': 'ListItem',
          position: index + 1,
          name,
          item: item.link ? base + item.link : undefined,
        };
      });
      scripts.push({
        type: 'application/ld+json',
        innerHTML: JSON.stringify({
          '@context': 'https://schema.org',
          '@type': 'BreadcrumbList',
          itemListElement: breadcrumbItems,
        }),
      });
    }

    (h.script as Record<string, unknown>[]).push(...scripts);

    return h;
  });

  useHead(head);
};
