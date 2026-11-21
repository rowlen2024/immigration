export function projectRegionLabel(country?: string | null): string {
  const label = country?.trim() ?? ''
  return label === '中国（香港特别行政区）' ? '中国香港' : label
}
