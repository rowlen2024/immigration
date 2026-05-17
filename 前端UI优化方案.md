# MyGo 移民官网前台 UI 优化方案

---

## 一、现状分析

### 1.1 技术架构

| 维度 | 现状 |
|------|------|
| 框架 | Nuxt 3 (SPA 模式, ssr: false) |
| CSS 方案 | **纯手写 CSS，无任何框架**（无 Tailwind/Bootstrap/UnoCSS） |
| UI 库 | Element Plus 已安装但**前台仅使用了 `el-rate`**（星级评分） |
| 图标 | 自定义 lucideIcons.ts（128个 SVG 图标硬编码），通过 `v-html` 渲染 |
| 样式组织 | 3 个全局 CSS 文件 + 每个 .vue 文件内 `<style scoped>` 块 |
| 字体 | **无 Web Font**，仅使用系统字体栈 |

### 1.2 视觉风格

**当前风格定位：深蓝海军 + 金色奢华风**

```
主色: #15294D (深蓝) / #0F172A (更深蓝)
强调色: #C8963E (金色)
背景: #FFFFFF (卡片) / #F8F9FB (区域底色)
文字: #1A1A2E (标题) / #555770 (正文) / #8E8EA0 (辅助)
```

**设计参考：** 类似高端律所/金融机构风格，强调专业感和信任感。

### 1.3 色彩体系问题

当前 CSS 变量存在 **"双重体系"**：

```css
/* 规范命名 (variables.css :root 前段) */
--color-primary: #0f172a;
--color-accent: #e2a83e;
--color-text: #0f172a;

/* 遗留别名 (variables.css :root 后段) */
--primary: #15294D;       /* 与 --color-primary 不同值! */
--accent: #C8963E;        /* 与 --color-accent 不同值! */
--text-primary: #1A1A2E;  /* 与 --color-text 不同值! */
```

**问题：** 两套变量指向不同色值，页面混用导致色彩不统一。例如：
- `admin.css` 覆盖 Element Plus 时用的是 `--color-primary`
- 大部分前台页面 scoped style 用的是 `--primary` / `--accent`
- `projects/[slug].vue` 的优势卡片同时使用了 `--color-text` 和 `--text-primary`

### 1.4 排版体系问题

**中文字体回退不可靠：**

```css
--font-serif: Georgia, 'Noto Serif CJK SC', 'Songti SC', serif;
```

`Georgia` 是英文字体，**不包含中文字符**。Windows 系统通常不安装 `Noto Serif CJK SC` 和 `Songti SC`，导致 Serif 标题在 Windows 上实际渲染为系统默认宋体或无衬线字体，与设计意图不符。

**字体大小缺乏层级：**
- Hero 标题：无明确规范，各页不一致
- 区块标题：32px（`.section-header h2`）
- 页面标题：36px
- 正文：16px（body 默认）
- 缺少系统化的 Type Scale（如 12/14/16/18/24/32/40/48）

### 1.5 布局结构

```
┌──────────────────────────────┐
│  Header (fixed, 64px)        │
├──────────────────────────────┤
│                              │
│  Main Content                │
│  max-width: 1200px           │
│                              │
├──────────────────────────────┤
│  Footer (深色背景)            │
└──────────────────────────────┘
         ↑
  ContactSidebar (fixed, right)
```

### 1.6 页面清单

| 路由 | 页面 | 状态 |
|------|------|------|
| `/` | 首页（Hero轮播+信任条+项目卡片+案例+评价+律师+优势+CTA） | ~900行单文件 |
| `/projects/:slug` | 项目详情（Tab导航+概览+条件+费用+流程+优势+案例+FAQ） | 复杂多段式 |
| `/compare` | 项目对比（下拉选择器+对比表） | 功能齐全 |
| `/compare/:a-vs-:b` | 对比直达页 | 含硬编码降级数据 |
| `/cases` | 成功案例列表 | 网格布局 |
| `/case/:slug` | 案例详情 | CMS内容渲染 |
| `/contact` | 联系表单 | 表单+侧边栏 |
| `/faq` | FAQ列表 | 项目筛选+手风琴 |
| `/pages/:slug` | 动态CMS页 | 3种模板 |

---

## 二、UI Bug 清单

### 2.1 🔴 严重 — 功能性 Bug

**Bug #1: Windows 下衬线标题不生效**
- **位置:** `variables.css:53` — `--font-serif` 字体栈
- **原因:** `Georgia` 无中文字符，`Noto Serif CJK SC` / `Songti SC` 在 Windows 上默认不存在
- **影响:** 首页 Hero 标题、所有 `.section-header h2` 的衬线效果在 Windows 上完全失效，标题与正文无字体差异
- **修复:** 引入 Google Fonts（Noto Serif SC）或在字体栈前添加 Windows 可用的中文字体

**Bug #2: CSS 变量色值不一致**
- **位置:** `variables.css` — 两套变量体系
- **表现:**
  - `--color-primary: #0f172a` vs `--primary: #15294D`
  - `--color-accent: #e2a83e` vs `--accent: #C8963E`
- **影响:** 页面不同区域使用不同变量，呈现略有差异的主色和强调色
- **修复:** 统一为单一变量体系，废弃旧别名或统一色值后保留别名

**Bug #3: 首页轮播重复 CSS**
- **位置:** `index.vue` 媒体查询内
- **表现:** `carousel-viewport` 的 margin/padding 覆盖出现两次（复制粘贴残留）
- **影响:** 代码冗余，不影响渲染（后一条覆盖前一条）

### 2.2 🟡 中等 — 体验问题

**Bug #4: 案例卡片代码重复**
- **位置:** `index.vue` / `cases.vue` / `projects/[slug].vue`
- **表现:** 案例卡片的 HTML 结构和 CSS 样式在 3 处重复，细微差异
- **影响:** 修改案例卡片样式需同步 3 处，容易遗漏

**Bug #5: 评价轮播代码重复**
- **位置:** `index.vue` / `projects/[slug].vue`
- **表现:** 轮播的 offset/page/gap/autoplay 逻辑完全重复
- **影响:** 同上，维护成本高

**Bug #6: 联系侧边栏在移动端的遮挡**
- **位置:** `ContactSidebar.vue` — `fixed, right: 24px, bottom: 94px`
- **表现:** 移动端可能与内容重叠，`z-index` 未在全局层级的明确定义
- **影响:** 小屏幕可能遮挡表单按钮或底部内容

**Bug #7: 自定义滚动条过细**
- **位置:** `global.css:99-108`
- **表现:** `width: 6px`，深色轨道 + 金色滑块，在浅色页面区域（如卡片内）可能不协调
- **影响:** 可用性问题 — 6px 宽度在触摸屏上几乎无法操作

**Bug #8: 加载状态不统一**
- **位置:** 各页面
- **表现:** 有的显示 "加载中..."，有的用骨架屏（skeleton），有的无加载状态
- **影响:** 用户体验不一致

### 2.3 🟢 轻微 — 视觉瑕疵

**Bug #9: Emoji 作为结构图标**
- **位置:** `Header.vue` 中的 `&#9662;`（▾）下拉箭头
- **表现:** 使用 HTML 实体作为导航下拉指示器
- **影响:** 不同平台渲染效果不一致

**Bug #10: focus-visible 样式引用了不一致的变量**
- **位置:** `global.css:87`
- **表现:** `outline: 2px solid var(--accent)` 使用的是旧 `--accent` 而非 `--color-accent`
- **影响:** 焦点环颜色与其他强调色不一致

**Bug #11: 对比页 select 下拉框样式未统一**
- **位置:** `compare.vue`
- **表现:** 原生 `<select>` 元素在不同浏览器/平台上样式差异大
- **影响:** 视觉不统一，尤其在移动端

---

## 三、可选优化方案

### 方案 A：渐进优化（保守）🛡️

**适合：** 希望快速见效、风险最低

**核心思路：** 修复 Bug + 小幅视觉微调，不动整体架构

| 改动项 | 说明 |
|--------|------|
| 统一 CSS 变量 | 合并两套变量，统一色值，废弃旧别名 |
| 引入 Web Font | 添加 Noto Serif SC（Google Fonts），修复衬线标题 |
| 提取案例卡片组件 | 创建 `components/CaseCard.vue`，消除重复 |
| 提取评价轮播组件 | 创建 `components/TestimonialCarousel.vue` |
| 统一加载状态 | 全部改为骨架屏 |
| 修复 Emoji 下拉箭头 | 改为 SVG 图标 |
| 增加移动端间距 | ContactSidebar 位置自适应 |

**工作量：** 1-2 天
**风险：** 极低

---

### 方案 B：视觉升级（适中）⭐ 推荐

**适合：** 希望显著提升品牌质感和转化率

**核心思路：** 在方案 A 基础上 + 品牌视觉重塑 + 交互优化

#### B1. 色彩升级

将现有"深蓝+金色"升级为更精致的调色板：

```css
/* 方案 B 推荐色彩 */
--color-primary:       #0F1B2D;  /* 更深邃的海军蓝 */
--color-primary-light: #1A2F4A;  /* 悬浮/卡片蓝 */
--color-accent:        #C8963E;  /* 保留金色，提高饱和度 */
--color-accent-hover:  #D4A84B;  /* 悬浮金 */
--color-accent-soft:   #FBF5E8;  /* 金色浅底 */
--color-gold-gradient: linear-gradient(135deg, #C8963E, #E2B86B);

/* 新增渐变 */
--gradient-hero:    linear-gradient(135deg, #060E18 0%, #0F1B2D 40%, #152036 70%, #0A1628 100%);
--gradient-cta:     linear-gradient(135deg, #C8963E, #D4A84B);
--gradient-card:    linear-gradient(180deg, rgba(200,150,62,0.05), transparent);
```

**变化点：**
- 主色略微加深，增加高级感
- 金色渐变替代纯色按钮
- 新增金色浅底用于卡片/区块

#### B2. 排版升级

```css
/* 引入 Web 字体 */
@import url('https://fonts.googleapis.com/css2?family=Noto+Serif+SC:wght@400;600;700&display=swap');
@import url('https://fonts.googleapis.com/css2?family=Noto+Sans+SC:wght@300;400;500;700&display=swap');

--font-sans:  'Noto Sans SC', 'PingFang SC', 'Microsoft YaHei', sans-serif;
--font-serif: 'Noto Serif SC', 'Georgia', 'Songti SC', serif;

/* Type Scale */
--text-xs:   0.75rem;   /* 12px */
--text-sm:   0.875rem;  /* 14px */
--text-base: 1rem;      /* 16px */
--text-lg:   1.125rem;  /* 18px */
--text-xl:   1.375rem;  /* 22px */
--text-2xl:  1.75rem;   /* 28px */
--text-3xl:  2.25rem;   /* 36px */
--text-4xl:  3rem;      /* 48px */
```

#### B3. 首页重构

**Hero 区域优化：**
- 新增渐变光晕动效（ambient light blobs）
- 标题改为逐字淡入动画
- CTA 按钮增加微光效果（shimmer）
- 轮播指示器改为进度条样式

**信任条优化：**
- 数字滚动动画保留，增加金色分隔线
- 背景改为半透明深蓝 + 毛玻璃效果

**项目卡片优化：**
- 卡片增加悬浮上浮效果（translateY + shadow）
- 图片区域增加金色渐变叠加
- 特性标签改为 pill 样式

#### B4. 细节微交互

| 元素 | 当前 | 优化后 |
|------|------|--------|
| 按钮悬浮 | color 切换 | 微上浮 + 阴影扩散 |
| 卡片悬浮 | 无 | translateY(-4px) + shadow-lg |
| 导航悬浮 | mega panel 显示 | 增加淡入动画（200ms） |
| Tab 切换 | instant | 滑动指示器（300ms） |
| FAQ 展开 | 无动画 | 高度过渡 + 图标旋转 |
| 页面滚动 | 无 | 滚动渐显 (reveal 已有，需确保覆盖全部区域) |

#### B5. 移动端优化

- 统一断点系统：`375 / 768 / 1024 / 1280 / 1440`
- 首页 Hero 移动端减小高度、增大文字
- 项目详情 Tab 导航移动端改为横向滚动
- 联系侧边栏移动端改为底部固定条

**工作量：** 3-5 天
**风险：** 低
**推荐理由：** 在保持现有功能和架构的前提下，最大化视觉提升效果，投入产出比最优

---

### 方案 C：品牌重塑（激进）🚀

**适合：** 希望全面升级，引入设计系统

**核心思路：** 引入 Tailwind CSS + 设计令牌系统，全面重构

| 改动项 | 说明 |
|--------|------|
| 安装 Tailwind CSS | 替换手写 CSS，使用 utility class |
| 建立 Design Token | 统一 color/spacing/typography/shadow 令牌 |
| 全面组件化 | 提取 Button/Card/Badge/Section 等基础组件 |
| 暗色模式 | 支持 prefers-color-scheme |
| 国际化准备 | 提取所有硬编码中文为 i18n key |
| 动画系统 | 引入 GSAP/Framer Motion 替代手写 CSS animation |

**工作量：** 2-3 周
**风险：** 中高（大量代码重写，可能引入回归问题）

---

## 四、方案对比

| 维度 | 方案 A (保守) | 方案 B (适中) ⭐ | 方案 C (激进) |
|------|:--:|:--:|:--:|
| Bug 修复 | ✅ | ✅ | ✅ |
| 视觉提升 | 🟡 小幅 | ✅ 显著 | ✅ 全面 |
| 代码质量 | 🟡 部分提取组件 | ✅ 组件化 | ✅ 完全组件化 |
| 性能影响 | 无 | 极小（字体加载） | 可能降低（框架层） |
| 工作量 | 1-2天 | 3-5天 | 2-3周 |
| 回归风险 | 极低 | 低 | 中高 |
| 暗色模式 | ❌ | ❌ | ✅ |
| 维护性提升 | 🟡 | ✅ | ✅✅ |

---

## 五、我的推荐：方案 B

作为精通移民行业的产品经理，我推荐方案 B，理由如下：

1. **移民行业的核心是信任** — 方案 B 的色彩升级让品牌更显高端可信，而非激进改变让老客户感到陌生
2. **转化率优先** — 微交互优化（按钮光影/卡片悬浮/滚动渐显）直接提升 CTA 点击意愿
3. **投入产出比最优** — 3-5天完成，不影响现有功能迭代计划
4. **Windows 体验修复** — 字体问题是当前最严重的视觉缺陷，方案 B 彻底解决

---

## 六、方案 B 详细实施计划

### 阶段 1：基础设施（Day 1）

| 任务 | 文件 | 说明 |
|------|------|------|
| 1.1 统一 CSS 变量 | `variables.css` | 合并两套变量，定义最终色值 |
| 1.2 引入 Web Fonts | `variables.css` + `nuxt.config.ts` | 添加 Noto Serif SC + Noto Sans SC |
| 1.3 建立 Type Scale | `variables.css` | 添加 `--text-*` 系列变量 |
| 1.4 统一阴影系统 | `variables.css` | 定义完整 shadow 层级 |
| 1.5 全局替换变量引用 | 所有 .vue 文件 | `--primary` → `--color-primary` 等 |

### 阶段 2：组件提取与优化（Day 2-3）

| 任务 | 文件 | 说明 |
|------|------|------|
| 2.1 提取 CaseCard 组件 | 新建 `components/CaseCard.vue` | 消除 3 处重复 |
| 2.2 提取 TestimonialCarousel | 新建 `components/TestimonialCarousel.vue` | 消除 2 处重复 |
| 2.3 优化 Header | `Header.vue` | Emoji→SVG, 添加过渡动画 |
| 2.4 优化 Footer | `Footer.vue` | 移动端适配优化 |
| 2.5 优化 ContactSidebar | `ContactSidebar.vue` | 移动端底部固定条 |

### 阶段 3：核心页面视觉升级（Day 3-4）

| 任务 | 文件 | 说明 |
|------|------|------|
| 3.1 首页 Hero 升级 | `index.vue` | 渐变光晕+逐字动画+shimmer按钮 |
| 3.2 首页项目卡片升级 | `index.vue` | 悬浮效果+金色叠加+pill标签 |
| 3.3 首页信任条升级 | `index.vue` | 毛玻璃+金色分隔线 |
| 3.4 项目详情页优化 | `projects/[slug].vue` | Tab动画+卡片悬浮+色彩统一 |
| 3.5 对比页优化 | `compare.vue` + `compare/[a]-vs-[b].vue` | select下拉美化+对比表视觉增强 |
| 3.6 案例列表页优化 | `cases.vue` | 使用 CaseCard 组件 |
| 3.7 联系页优化 | `contact.vue` | 表单焦点状态+提交按钮动画 |
| 3.8 FAQ 页优化 | `faq.vue` | 展开动画+筛选按钮悬浮 |

### 阶段 4：动效与打磨（Day 4-5）

| 任务 | 文件 | 说明 |
|------|------|------|
| 4.1 统一过渡时间 | 全局 | 所有动画 200-300ms |
| 4.2 滚动渐显检查 | 各页面 | 确保所有区块有 reveal |
| 4.3 移动端全面测试 | 全局 | 375/768/1024/1440 断点 |
| 4.4 性能验证 | 全局 | Lighthouse + 视觉回归 |
| 4.5 编译验证 | 全局 | `npx nuxi typecheck` + `npm run build` |

---

## 七、待确认事项

在开始实施前，请确认以下问题：

1. **方案选择**：您倾向于方案 A、B 还是 C？或者需要在 B 的基础上调整？
2. **字体选择**：方案 B 推荐 Google Fonts（Noto Serif SC + Noto Sans SC），是否需要考虑国内 CDN 加速（如字体文件自托管）？
3. **金色调整**：当前金色偏沉稳（#C8963E），方案 B 建议略微提亮（增加渐变）。是否需要保持现有金色色值？
4. **首页重构范围**：首页改动影响最大，是否需要保持现有布局完全不变，只改视觉？
5. **移动端优先级**：移动端访问占比大概多少？是否需要专门优化移动端体验？

---

*本方案由 UI/UX Pro Max 设计系统辅助生成，结合 MyGo 移民网站实际代码分析完成。*
*设计系统推荐：Trust & Authority 风格 + Enterprise Gateway 布局模式*
