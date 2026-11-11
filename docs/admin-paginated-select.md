# 后台分页多选组件

## 目标与职责

`frontend/components/admin/PaginatedMultiSelect.vue` 用于从大量远程选项中搜索、翻页和多选。默认每页 20 条，首次展开下拉框才发起请求，不循环加载全部分页。

- 组件负责关键词输入、分页、加载状态、失败重试和过期响应隔离。
- 业务页面通过 `loadOptions` 提供当前页选项与总数，负责后台接口、鉴权、业务筛选和字段映射。
- 后端继续负责 SQL 分页及业务权限。组件不改变接口契约、不转换标识、不清洗关键词、不重复校验业务规则。
- 当前首个使用场景为项目管理的“添加资讯”，通过原有 `/admin/pages/options` 接口检索已发布新闻，不涉及数据迁移。

## 使用契约

类型定义位于 `frontend/types/paginated-select.ts`。

| 属性 | 契约 |
| --- | --- |
| `v-model` | 原始标识数组，支持字符串或数字；跨页和搜索保留选择 |
| `load-options` | 接收 `{ page, perPage, keyword }`，返回 `Promise<{ items: [{ value, label }], total }>` |
| `page-size` | 默认 20，调用方应遵守目标接口允许的分页范围 |
| `placeholder` | 默认“输入关键词搜索...” |
| `empty-text` | 默认“未找到匹配项” |
| `disabled` | 默认 `false`，禁止修改选择 |

`total` 是当前关键词筛选后的总数，不能使用当前页条数代替。组件仅保存当前页选项，已选项标题使用 Element Plus 的选项缓存保留。关键词防抖沿用选择器行为，关键词变化时回到第一页。

```vue
<PaginatedMultiSelect
  v-model="selectedIds"
  :load-options="loadOptions"
  placeholder="输入标题搜索新闻页面..."
/>
```

```ts
import type { PaginatedSelectRequest, PaginatedSelectResult } from '~/types/paginated-select';

/** 传给业务保存接口的原始新闻标识。 */
const selectedIds = ref<number[]>([]);

/** 使用原有后台客户端查询一页已发布新闻。 */
const loadOptions = async (request: PaginatedSelectRequest): Promise<PaginatedSelectResult<number>> => {
  /** 自动携带鉴权信息并解包分页信封的后台客户端。 */
  const api = useApi();
  /** 后台返回的当前页数据。 */
  const data = await api<{ items: { id: number; title: string }[]; total: number }>('/admin/pages/options', {
    query: {
      page: request.page,
      per_page: request.perPage,
      title: request.keyword,
      page_type: 'news',
      status: 'published',
    },
  });
  return {
    items: data.items.map(item => ({
      value: item.id,
      label: item.title,
    })),
    total: data.total,
  };
};
```

## 异常与生命周期

- 请求期间显示加载状态并禁用分页，避免继续选择旧列表。
- 请求失败时展示“加载失败 / 重试”，保留已有选择；重试当前页及当前关键词。
- 新搜索或翻页发起后，较早的请求结果不再写入状态。
- 组件销毁后，未完成请求不再写入状态。弹窗使用 `destroy-on-close` 时，每次重新打开从初始查询开始；父页面负责清空当次临时选择。
- 取消或关闭弹窗不调用业务保存接口；确认添加继续使用现有新闻关联接口。

## 验收数据与场景

浏览器验收使用 625 篇模拟已发布新闻，ID 为 1～625，标题为“分页验收新闻 N”，每页 20 篇，共 32 页，末页 5 篇。模拟响应沿用后台 `data + pagination` 信封，不接触实际数据库。

1. 未展开选择器时不请求选项；首次展开仅请求 `page=1&per_page=20`。
2. 翻至第 2 页显示 ID 21～40，末页显示 ID 621～625。
3. 从第 1、2 页各选一篇，再搜索“620”并选择 ID 620；三个标题和标识均保留，保存提交三者。
4. 从后续页输入新标题时回到第 1 页；清空搜索后恢复全量分页，选择仍保留。
5. 不存在的关键词显示空结果，确认按钮按已选数量启用或禁用。
6. 先发出的搜索延迟返回，后发搜索先完成时，界面始终保留后发查询结果。
7. 模拟请求失败，显示重试入口；重试仍使用失败时的页码和关键词。
8. 请求未完成时关闭弹窗再打开，旧请求不污染新一轮选择。
9. 取消不写入关联接口；重新打开清空临时选择。

构建与类型检查命令：在 `frontend/` 执行 `npm run build`、`npx nuxi typecheck`。浏览器模拟验收用于确认前端交互，实际数据库及部署环境需另行联调。

### 本次验证结果（2026-11-10）

- Nuxt 生产构建通过，`git diff --check` 通过。
- 通过 Chrome DevTools MCP 在真实项目管理页注入上述模拟接口数据，完成全部九项浏览器验收；实际提交的跨页选择为 `[1, 21, 620]`，所有选项请求均为每页 20 条，并保留已发布新闻筛选条件。
- `npx nuxi typecheck` 未通过：现有 `frontend/pages/admin/settings.vue` 存在 `idx`、`field` 未定义错误；本次修改文件未产生类型诊断。
- 尚未连接实际 MySQL 或部署环境验证。
