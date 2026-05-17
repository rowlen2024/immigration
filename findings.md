# Findings — MyGo 移民前台 UI 优化

**关联任务:** task_plan.md

---

## 发现的 Bug

### B1: Windows 下衬线标题不生效
- **根因:** `--font-serif` 字体栈以 `Georgia` 为首选，Georgia 不含中文字符，Windows 不支持 `Noto Serif CJK SC` / `Songti SC`
- **修复:** 引入 Google Fonts `Noto Serif SC` 作为首选，`Georgia` 仅作英文回退

### B2: CSS 变量双重体系
- **根因:** `variables.css` 同时存在 `--color-*` 规范和 `--primary`/`--accent` 等遗留别名，且色值不同（如 `--color-primary: #0f172a` vs `--primary: #15294D`）
- **修复:** 统一为单一色值体系，遗留别名通过 `var()` 指向规范变量

### B3: 首页轮播 CSS 重复
- **位置:** `index.vue` 767px 断点内 `.carousel-viewport` 规则出现两次，第一条被第二条覆盖
- **修复:** 提取到 TestimonialCarousel 组件后自然消除

### B4: Emoji 作为结构图标
- **位置:** Header.vue 导航下拉箭头 `&#9662;`，FAQAccordion `+/−`
- **修复:** 全部替换为内联 SVG 图标，配合过渡动画

### B5: 案例卡片 3 处重复
- **范围:** `index.vue` / `cases.vue` / `projects/[slug].vue`
- **修复:** 提取统一 CaseCard 组件

### B6: 评价轮播 2 处重复
- **范围:** `index.vue` / `projects/[slug].vue`（核心逻辑 100% 相同）
- **修复:** 提取统一 TestimonialCarousel 组件，清理死代码（testimonialPrev/Next/detailNext 从未调用）

---

## 设计决策

### D1: 金色按钮文字用白色而非深色
- **理由:** 金色渐变底色上，白色文字 + `text-shadow` 在保证可读性的同时更显高级；深色文字在金色上偏暗淡

### D2: 评价卡片去掉引号，改用视觉层次区分
- **理由:** 文字引号 `"..."` 效果单调，缺乏社交证明氛围。新方案用金色左竖线 + 暖色渐变背景 + 星级居中 + 已签约验证标签来营造"评论区"质感

### D3: FAQ 手风琴用 SVG 十字图标替代 `+/−` 文字
- **理由:** 文字符号在不同平台渲染不一致，SVG 图标配合旋转动画（+ 变 ×）更精致

### D4: 保留遗留变量别名
- **理由:** 大量现有代码引用 `--primary` / `--accent` / `--text-primary` 等旧变量名，保留别名使旧代码无需改动即可使用新色值，降低回归风险
