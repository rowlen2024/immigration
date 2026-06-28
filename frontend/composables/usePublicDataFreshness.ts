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

const hasNuxtPayloadData = (nuxtApp: ReturnType<typeof useNuxtApp>, dataKey: string) => {
  const payloadData = nuxtApp.payload?.data as Record<string, unknown> | undefined
  const staticData = nuxtApp.static?.data as Record<string, unknown> | undefined

  return Object.prototype.hasOwnProperty.call(payloadData || {}, dataKey)
    || Object.prototype.hasOwnProperty.call(staticData || {}, dataKey)
}

const scheduleFreshnessCheck = (
  nuxtApp: ReturnType<typeof useNuxtApp>,
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
    inflight = flushFreshnessQueue(nuxtApp, force, versions, refreshing).finally(() => {
      inflight = null
    })
  }, 0)
}

const flushFreshnessQueue = async (
  nuxtApp: ReturnType<typeof useNuxtApp>,
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
    const versionKeysWithFreshValue = versionKeys.filter((key) => freshVersions[key])
    const changedDataKeys = entries
      .filter((entry) => {
        const next = freshVersions[entry.versionKey]
        if (!next || refreshing.value[entry.dataKey]) return false
        const previous = versions.value[entry.versionKey]
        return force
          || (previous !== undefined && previous !== next)
          || (previous === undefined && hasNuxtPayloadData(nuxtApp, entry.dataKey))
      })
      .map((entry) => entry.dataKey)

    if (changedDataKeys.length === 0) {
      for (const key of versionKeysWithFreshValue) versions.value[key] = freshVersions[key]
      return
    }

    const uniqueDataKeys = Array.from(new Set(changedDataKeys))
    for (const key of uniqueDataKeys) refreshing.value[key] = true
    try {
      await refreshNuxtData(uniqueDataKeys)
      for (const key of versionKeysWithFreshValue) versions.value[key] = freshVersions[key]
    } finally {
      for (const key of uniqueDataKeys) refreshing.value[key] = false
    }
  } catch {
    // Freshness checks should never break page rendering.
  }
}

export const usePublicDataFreshness = (keys: PublicFreshnessKeys) => {
  const nuxtApp = useNuxtApp()
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
    scheduleFreshnessCheck(nuxtApp, entries, route.fullPath, route.query.fresh === '1', versions, refreshing, checked)
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
