# Admin UI Redesign — Design Spec

**日期**: 2026-05-04
**产品**: MyGo Immigration 后台管理系统
**风格**: 现代极简 (Modern Minimal) — slate 冷色调 + 金色点缀

---

## 1. 色彩系统

| Token | 值 | 用途 |
|-------|-----|------|
| `--color-bg-app` | `#f8fafc` | 页面底色 (slate-50) |
| `--color-bg-surface` | `#ffffff` | 卡片/表格/表单背景 |
| `--color-bg-sidebar` | `#0f172a` | 侧边栏底色 (slate-900) |
| `--color-primary` | `#0f172a` | 主色 — 按钮、链接、活跃态 (slate-900) |
| `--color-accent` | `#e2a83e` | 强调色 — 核心 CTA、重要状态标识 |
| `--color-accent-soft` | `#fef9c3` | 强调色浅底 — 高亮提示 (yellow-100) |
| `--color-border` | `#e2e8f0` | 边框/分割线 (slate-200) |
| `--color-border-light` | `#f1f5f9` | 微弱分割线 (slate-100) |
| `--color-text` | `#0f172a` | 主文字 (slate-900) |
| `--color-text-secondary` | `#64748b` | 辅助文字 (slate-500) |
| `--color-text-muted` | `#94a3b8` | 弱化文字 (slate-400) |
| `--color-success` | `#16a34a` | 成功/已发布 |
| `--color-warning` | `#d97706` | 警告 |
| `--color-danger` | `#dc2626` | 危险/删除 |
| `--color-info` | `#3b82f6` | 信息/提示 |

**语义色保留方案**: 成功绿、警告橙、危险红沿用 Tailwind 语义色值。删除按钮统一用红色 outline 而非实心。

---

## 2. 字体系统

| 用途 | 字体 | 字号 | 字重 |
|------|------|------|------|
| 页面标题 | PingFang SC | 22px | 600 |
| 区块标题 | PingFang SC | 16px | 600 |
| 正文 | PingFang SC | 14px | 400 |
| 辅助文字 | PingFang SC | 12px | 400 |
| 统计数字 | Tabular Nums | 28px | 700 |

现有 `--font-sans` 变量保持不变 (`PingFang SC, Microsoft YaHei, Helvetica Neue, sans-serif`)。

---

## 3. 图标系统

统一使用 **Lucide** 图标集 (lucide-vue-next)，全体 stroke-width: 1.5, size: 18px。

侧边栏导航图标映射：
- 控制台 → `LayoutDashboard`
- 项目管理 → `FolderKanban`
- 页面管理 → `FileText`
- 首页配置 → `Home`
- 导航管理 → `Navigation`
- FAQ 管理 → `MessageCircleQuestion`
- 案例管理 → `Users`
- 咨询管理 → `MessageSquare`
- 用户管理 → `Shield`
- 媒体库 → `Image`
- 网站设置 → `Settings`

---

## 4. 侧边栏 (admin.vue)

### 结构
- 顶部：Logo 区 (MyGo 图标 + "MyGo 管理后台" 文字)
- 中部：导航菜单 (图标 + 文字)，当前页高亮 + 左侧 3px 竖条指示器
- 底部：用户信息区 (头像 + 角色名)，悬停展开下拉菜单 (个人设置/退出)

### 交互
- 桌面端宽度 220px，可折叠至 56px（仅图标 + tooltip）
- 移动端保持现有 slide-over 模式
- 悬浮态：`rgba(255,255,255,0.08)` 背景
- 活跃态：`rgba(255,255,255,0.16)` 背景 + 左侧 3px 白色竖条
- 折叠时图标大小保持 20px，居中显示

---

## 5. 顶栏 (Topbar)

精简为一行：
- 左侧：汉堡菜单折叠按钮
- 右侧：通知铃铛 (预留) + 用户头像下拉 + "访问网站" 外链

高度 56px (比当前 60px 更紧凑)，白色背景，底部 1px 边框 `--color-border`。

---

## 6. 控制台 (Dashboard)

### 统计卡片区
- 4 列 grid, gap 16px
- 每张卡片：图标(彩色圆角方块) + 标签 + 数字 + 环比变化
- 鼠标悬浮：轻微上浮 2px + 阴影增强

### 内容区 (2 栏布局)
- 左栏 (2/3)：最近咨询列表 — 名称 + 项目 + 相对时间
- 右栏 (1/3)：快捷操作 — 4 个操作按钮 (新建项目、新建页面、新建 FAQ、查看咨询)
- 快捷操作按钮使用卡片式排列，带图标

---

## 7. 列表页统一模板 (Projects/Pages/FAQs/Cases/Users/Leads)

### 页面结构
```
[页面标题]                    [新建按钮]
[搜索框] [状态筛选]
[列表区域 — 无边框分割线式]
[分页栏 — 右对齐]
```

### 列表行
- 行高 52px，左右 padding 0 (内容靠边)
- 行间分割线：`--color-border-light` (slate-100)
- 悬浮高亮：`--color-bg-app` (slate-50)
- 操作列：文字按钮 (编辑/删除)，不使用 el-button 外框
- 状态标签：圆角 pill 样式，小尺寸

### 搜索/筛选栏
- 搜索框左侧，300px 宽，圆角 8px
- 状态筛选下拉紧随其后
- 整体与列表内容用 16px 间距分隔

---

## 8. 抽屉表单 (Drawer)

### 容器
- 宽度 560px (桌面)，移动端全屏
- 从右侧滑入，背景遮罩 `rgba(0,0,0,0.3)`
- 动画 duration 250ms，ease-out

### 内部结构
```
[标题栏: 标题 + 关闭按钮]        ← 固定顶部
[表单内容区: 可滚动]              ← 弹性填充
[操作栏: 取消 + 保存]            ← 固定底部
```

### 表单字段
- label 在上方 (top)，字号 13px，颜色 `--color-text-secondary`
- input 高度 38px，圆角 6px，边框 `--color-border`
- focus 态：边框 `--color-primary`，ring 2px `rgba(15,23,42,0.1)`
- 字段间距 20px

---

## 9. 首页配置页 (homepage.vue)

### 布局改为 Tab 切换
```
[页面标题: 首页配置]
[Tab: 轮播管理 | 项目展示区 | 优势管理]
[当前 Tab 内容区]
```

### 轮播管理 Tab
- 列表项：缩略图 (120×68) + 标题 + 上下移动/编辑/删除按钮
- 顶部 "新增 Slide" 按钮
- 底部 "保存轮播" 按钮（全宽，居中）

### 项目展示区 Tab
- 表单区 (标题/副标题) + 精选项目列表
- 项目列表项：名称 + 上下移动/移除
- 下拉选择框添加新项目
- 底部 "保存" 按钮

### 优势管理 Tab
- 表单区 (标题/副标题) + 优势项列表
- 优势项：图标预览 + 标题/描述 + 上下移动/编辑/删除
- 底部 "新增优势项" + "保存优势设置"

---

## 10. 其他页面微调

### 导航管理 (navigation.vue)
- 树节点行内操作按钮改为图标按钮 (Pencil, Plus, Trash2)
- 链接类型标签保持但缩小
- "已隐藏" 标签颜色改灰

### 媒体库 (media.vue)
- 网格卡片改为 4 列 → 2 列 (平板) → 1 列 (手机)
- 增加图片预览大图 (点击放大)
- 文件信息排版优化

### 网站设置 (settings.vue)
- 卡片间距缩小到 16px
- 字段标签宽度统一 160px
- "?" 提示按钮改为灰色圆形，hover 变深

### 登录页 (login.vue)
- 背景渐变保持现有风格
- 卡片增加 Logo 区
- 输入框圆角加大到 8px
- 按钮改为 slate-900 色

---

## 11. 空状态 & 加载态

### 空状态
每个列表页提供语义化空状态：
```
[图标: 文件夹空]
暂无数据
[新建按钮]
```

### 加载态
- 列表加载：灰色占位条闪烁 (skeleton)，不使用全局 spinner
- 提交中：按钮 loading spinner + 禁用态
- 全局刷新：仅顶部进度条 (NProgress 风格)

---

## 12. 响应式断点

| 断点 | 宽度 | 适配 |
|------|------|------|
| 手机 | < 768px | 侧边栏 overlay, 统计卡片 2 列, 抽屉全屏 |
| 平板 | 768–1024px | 侧边栏折叠, 统计卡片 3 列 |
| 桌面 | > 1024px | 侧边栏展开 220px, 统计卡片 4 列 |

---

## 13. 实施范围

### Phase 1 (核心框架)
- `variables.css` — 色彩变量全面更新
- `admin.css` — 共享样式重写
- `admin.vue` (layout) — 侧边栏+顶栏重构

### Phase 2 (列表页模板)
- `projects.vue` — 列表+搜索+抽屉表单
- `pages.vue` — 同上模板
- `faqs.vue` — 同上模板
- `cases.vue` — 同上模板
- `users.vue` — 同上模板
- `leads.vue` — 同上模板

### Phase 3 (特殊页面)
- `index.vue` — 控制台重构
- `homepage.vue` — Tab 切换重构
- `navigation.vue` — 图标按钮优化
- `media.vue` — 网格优化
- `settings.vue` — 排版优化
- `login.vue` — 视觉升级

---

## 14. 技术实现备注

- **图标**: 使用项目已有的 `composables/lucideIcons.ts` (内联 SVG)，无需额外依赖
- **Element Plus**: 保留使用，通过 CSS 变量覆盖默认主题色
- **抽屉**: 使用 Element Plus `el-drawer` 组件替代 `el-dialog`
- **骨架屏**: 使用 Element Plus `el-skeleton` 组件
- **分页**: 保持 `el-pagination`，仅调整间距样式
- **不引入新依赖** (图表库等)，保持简洁

---

## 15. 设计检查清单

- [ ] 所有文本对比度 ≥ 4.5:1
- [ ] 可点击区域最小 44×44px
- [ ] 动画时长 150–300ms
- [ ] 无 emoji 作为 UI 图标
- [ ] 所有图标来自 Lucide 统一图标集
- [ ] 响应式：375px / 768px / 1024px / 1440px
- [ ] prefers-reduced-motion 尊重
