interface ProjectChoice {
  slug: string
  title: string
}

export function useProjectComparison() {
  const selection = useState<ProjectChoice[]>('public:comparison-selection', () => [])
  const isFull = computed(() => selection.value.length >= 2)
  const isSelected = (slug: string) => selection.value.some(project => project.slug === slug)
  const compareTo = computed(() => ({
    path: '/compare',
    query: { a: selection.value[0]?.slug, b: selection.value[1]?.slug },
  }))

  function remove(slug: string) {
    selection.value = selection.value.filter(project => project.slug !== slug)
  }

  function toggle(project: ProjectChoice) {
    if (isSelected(project.slug)) remove(project.slug)
    else if (!isFull.value) selection.value = [...selection.value, project]
  }

  return { selection, isFull, isSelected, compareTo, remove, toggle }
}
