/** 分页选择器支持的原始标识类型，组件不转换标识格式。 */
export type PaginatedSelectValue = string | number;

/** 分页选择器传给业务加载函数的查询条件。 */
export interface PaginatedSelectRequest {
  /** 从 1 开始的页码。 */
  page: number;
  /** 本次需要的选项条数。 */
  perPage: number;
  /** 用户输入的原始关键词，由业务调用方映射到接口检索字段。 */
  keyword: string;
}

/** 分页选择器展示的单个选项。 */
export interface PaginatedSelectOption<T extends PaginatedSelectValue> {
  /** 选项的原始业务标识。 */
  value: T;
  /** 选项展示文本。 */
  label: string;
}

/** 业务加载函数返回的当前页选项与筛选后的总数。 */
export interface PaginatedSelectResult<T extends PaginatedSelectValue> {
  /** 当前页选项，组件不会继续预取其他分页。 */
  items: PaginatedSelectOption<T>[];
  /** 当前关键词对应的全部选项数量。 */
  total: number;
}
