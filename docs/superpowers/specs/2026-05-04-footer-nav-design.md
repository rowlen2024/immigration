# 页脚导航动态化设计

## 目标

将 Footer 组件从硬编码改为从后台导航管理动态获取，实现页脚链接可通过后台统一管理。

## 方案：Navigation 模型新增 `display_position` 字段

使用单一枚举字段控制导航项在哪些位置显示，替代两个独立布尔值。

枚举值: `header`（仅头部）、`footer`（仅页脚）、`both`（两处都显示）

## 数据库变更

新增 migration `000015_add_navigation_display_position`:

```sql
-- up
ALTER TABLE navigations
  ADD COLUMN display_position VARCHAR(16) NOT NULL DEFAULT 'header'
  AFTER status;
ALTER TABLE navigations ADD INDEX idx_nav_display_position (display_position);

-- down (separate file)
ALTER TABLE navigations DROP INDEX idx_nav_display_position;
ALTER TABLE navigations DROP COLUMN display_position;
```

现有导航项默认为 `header`，保持头部导航不受影响。管理员后续在后台为需要出现在页脚的导航项修改此字段。

## 后端变更

### Model (`internal/model/navigation.go`)

Navigation 结构体新增:
```go
DisplayPosition string `gorm:"size:16;not null;default:'header'" json:"display_position"`
```

### Repository (`internal/repository/nav_repo.go`)

新增方法，按位置筛选启用的导航:
```go
func (r *NavRepo) FindAllActiveByPosition(position string) ([]model.Navigation, error)
```
查询逻辑: `WHERE status=1 AND display_position IN (position, 'both')`

### Service (`internal/service/nav_svc.go`)

- `GetTree(position string)` — 新增 `position` 参数，替换原有的无参版本
- 内部调用 `FindAllActiveByPosition(position)`

### Handler (`internal/handler/nav_handler.go`)

`GetNavigation` 接收可选 query param `position`:
- 默认值 `header`（向后兼容）
- 可以传 `footer` 获取页脚导航

### Router

无需新增路由。沿用现有公开路由:
```
GET /api/v1/navigation?position=footer
```

## 前端变更

### Footer.vue

- 调用 `GET /api/v1/navigation?position=footer` 获取页脚导航树
- 顶层节点渲染为列标题（`footer-heading`），子节点渲染为列内链接
- 品牌列（左侧）和联系方式列（右侧）保持从 `siteConfig` 获取，不通过导航系统
- 保留硬编码 fallback，API 失败时降级

渲染逻辑示例:
```
API 返回树: [
  { label: "移民项目", children: [{label:"美国EB-5", link:"/projects/eb5"}, ...] },
  { label: "关于我们", children: [{label:"公司简介", link:"/pages/about"}, ...] },
]

渲染:
  列: 移民项目          列: 关于我们
      美国EB-5              公司简介
      香港投资              成功案例
      巴拿马购房            常见问题
      项目对比              联系我们
```

### Header.vue

- 请求改为 `GET /api/v1/navigation?position=header`
- 其余逻辑不变（树形结构渲染桌面端下拉菜单 + 移动端可展开菜单）

### Admin navigation.vue

表单新增 `display_position` 字段:
```
[radio group]
  ○ 仅头部    ○ 仅页脚    ○ 头部+页脚
```

## 影响范围

| 层 | 文件 | 变更类型 |
|---|------|---------|
| DB | `database/migrations/000015_*` | 新增 migration |
| Model | `internal/model/navigation.go` | 新增字段 |
| Repo | `internal/repository/nav_repo.go` | 新增方法 |
| Repo | `internal/repository/interfaces.go` | 接口新增方法 |
| Service | `internal/service/nav_svc.go` | GetTree 加参数 |
| Handler | `internal/handler/nav_handler.go` | 接收 query param |
| Frontend | `components/global/Footer.vue` | 重写数据获取逻辑 |
| Frontend | `components/global/Header.vue` | 请求 URL 微调 |
| Frontend | `pages/admin/navigation.vue` | 表单新增字段 |

## 风险

- 现有 Header 依赖 `GET /api/v1/navigation` 无参数返回全部激活导航；改为默认 `position=header` 后行为一致（默认值就是 header），向后兼容
- 需要执行 migration 后才能启动后端，否则新字段不存在
