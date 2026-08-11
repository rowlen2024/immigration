/** 面包屑项同时供页面展示与结构化数据使用。 */
export interface BreadcrumbItem {
  /** 面包屑显示名称。 */
  label: string
  /** 站内页面路径。 */
  href: string
}

/** 所有公开页面面包屑的固定起点。 */
export const HOME_BREADCRUMB: BreadcrumbItem = {
  label: '首页',
  href: '/',
}
