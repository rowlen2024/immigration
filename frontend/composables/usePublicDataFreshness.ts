type PublicVersionsResponse = {
  data?: Record<string, string>
}

type PublicFreshnessKey = string | {
  versionKey: string
  dataKey?: string
}

type PublicFreshnessKeys = PublicFreshnessKey[] | Ref<PublicFreshnessKey[]> | (() => PublicFreshnessKey[])

type FreshnessEntry = {
  versionKey: string
  dataKey: string
}

let pendingEntries = new Map<string, FreshnessEntry>()
let flushTimer: ReturnType<typeof setTimeout> | null = null
let inflight: Promise<void> | null = null
const CHECK_DEDUPE_MS = 5_000

const resolveFreshnessEntries = (keys: PublicFreshnessKeys): FreshnessEntry[] => {
  const value = typeof keys === 'function' ? keys() : unref(keys)
  const entries = (value || [])
    .filter(Boolean)
    .map((item) => {
      if (typeof item === 'string') {
        return { versionKey: item, dataKey: item }
      }
      return { versionKey: item.versionKey, dataKey: item.dataKey || item.versionKey }
    })
    .filter((item) => item.versionKey && item.dataKey)

  return Array.from(new Map(entries.map((item) => [`${item.versionKey}:${item.dataKey}`, item])).values())
}

const scheduleFreshnessCheck = (
  entries: FreshnessEntry[],
  routeKey: string,
  force: boolean,
  versions: Ref<Record<string, string>>,
  refreshing: Ref<Record<string, boolean>>,
  checked: Ref<Record<string, number>>,
) => {
  const now = Date.now()
  for (const entry of entries) {
    const checkKey = `${routeKey}:${entry.versionKey}:${entry.dataKey}`
    const lastCheckedAt = checked.value[checkKey] || 0
    if (!force && now - lastCheckedAt < CHECK_DEDUPE_MS) continue
    pendingEntries.set(`${entry.versionKey}:${entry.dataKey}`, entry)
    checked.value[checkKey] = now
  }

  if (pendingEntries.size === 0 || flushTimer) return

  flushTimer = setTimeout(() => {
    flushTimer = null
    inflight = flushFreshnessQueue(force, versions, refreshing).finally(() => {
      inflight = null
    })
  }, 0)
}

const flushFreshnessQueue = async (
  force: boolean,
  versions: Ref<Record<string, string>>,
  refreshing: Ref<Record<string, boolean>>,
) => {
  if (inflight) await inflight

  const entries = Array.from(pendingEntries.values())
  pendingEntries.clear()
  if (entries.length === 0) return

  const versionKeys = Array.from(new Set(entries.map((entry) => entry.versionKey)))
  try {
    const response = await $fetch<PublicVersionsResponse>('/api/v1/public-versions', {
      query: { keys: versionKeys.join(',') },
    })
    const freshVersions = response?.data || {}
    const changedDataKeys = entries
      .filter((entry) => {
        const next = freshVersions[entry.versionKey]
        if (!next || refreshing.value[entry.dataKey]) return false
        const previous = versions.value[entry.versionKey]
        return force || (previous !== undefined && previous !== next)
      })
      .map((entry) => entry.dataKey)

    for (const key of versionKeys) {
      if (freshVersions[key]) versions.value[key] = freshVersions[key]
    }
    if (changedDataKeys.length === 0) return

    const uniqueDataKeys = Array.from(new Set(changedDataKeys))
    for (const key of uniqueDataKeys) refreshing.value[key] = true
    try {
      await refreshNuxtData(uniqueDataKeys)
    } finally {
      for (const key of uniqueDataKeys) refreshing.value[key] = false
    }
  } catch {
    // Freshness checks should never break page rendering.
  }
}

export const usePublicDataFreshness = (keys: PublicFreshnessKeys) => {
  const route = useRoute()
  const versions = useState<Record<string, string>>('public-data-versions', () => ({}))
  const refreshing = useState<Record<string, boolean>>('public-data-refreshing', () => ({}))
  const checked = useState<Record<string, number>>('public-data-checked', () => ({}))

  if (import.meta.server) {
    onServerPrefetch(async () => {
      const entries = resolveFreshnessEntries(keys)
      const versionKeys = Array.from(new Set(entries.map((entry) => entry.versionKey)))
      const missingKeys = versionKeys.filter((key) => !versions.value[key])
      if (missingKeys.length === 0) return

      try {
        const response = await $fetch<PublicVersionsResponse>('/api/v1/public-versions', {
          query: { keys: missingKeys.join(',') },
        })
        const freshVersions = response?.data || {}
        for (const key of missingKeys) {
          if (freshVersions[key]) versions.value[key] = freshVersions[key]
        }
      } catch {
        // Freshness baselines should never break SSR rendering.
      }
    })
    return
  }

  const check = () => {
    const entries = resolveFreshnessEntries(keys)
    if (entries.length === 0) return
    scheduleFreshnessCheck(entries, route.fullPath, route.query.fresh === '1', versions, refreshing, checked)
  }

  const handleFocus = () => check()
  const handleVisibilityChange = () => {
    if (document.visibilityState === 'visible') check()
  }

  onMounted(() => {
    check()
    window.addEventListener('focus', handleFocus)
    document.addEventListener('visibilitychange', handleVisibilityChange)
  })

  onUnmounted(() => {
    window.removeEventListener('focus', handleFocus)
    document.removeEventListener('visibilitychange', handleVisibilityChange)
  })
  watch(() => resolveFreshnessEntries(keys).map((entry) => `${entry.versionKey}:${entry.dataKey}`).join(','), check)
  watch(() => route.fullPath, check)
}
