export function stripHtml(html: string, maxLen = 80): string {
  if (!html) return ''
  return html.replace(/<[^>]+>/g, '').replace(/&nbsp;/g, ' ').slice(0, maxLen)
}
