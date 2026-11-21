interface ApiResponse<T> {
  data: T
}

interface SitemapResource {
  slug: string
  updated_at?: string
}

function mapSitemapEntries(resources: SitemapResource[], prefix: string) {
  return resources
    .filter(resource => resource.slug?.trim())
    .map(resource => ({
      loc: `/${prefix}/${resource.slug}`,
      lastmod: resource.updated_at,
    }))
}

export default defineSitemapEventHandler(async () => {
  const [projects, pages, cases, categories] = await Promise.all([
    $fetch<ApiResponse<SitemapResource[]>>('/api/v1/projects?status=1'),
    $fetch<ApiResponse<SitemapResource[]>>('/api/v1/pages'),
    $fetch<ApiResponse<SitemapResource[]>>('/api/v1/cases'),
    $fetch<ApiResponse<SitemapResource[]>>('/api/v1/project-categories').catch(() => {
      console.warn('[sitemap] Failed to load project categories')
      return { data: [] }
    }),
  ])

  return [
    ...mapSitemapEntries(projects.data, 'projects'),
    ...mapSitemapEntries(pages.data, 'pages'),
    ...mapSitemapEntries(cases.data, 'case'),
    ...mapSitemapEntries(categories.data, 'projects/category'),
  ]
})
