# Findings — MyGo 移民

**关联任务:** task_plan.md

---

## 图片变体系统 (2026-05-23)

### 发现的问题

#### C1: `aspect-ratio` CSS 与固定 height 冲突
- **位置:** `.card-image`, `.case-image`, `.lawyer-photo`
- **问题:** `height: 200px` + `aspect-ratio: 16/9` → 浏览器从高度反算宽度 355px，父容器 366px，产生 11px 右侧间隙
- **修复:** 移除 `.card-image`、`.case-image`、`.lawyer-photo` 的 `aspect-ratio`，它们已有 height/flex-basis 约束

#### C2: CSS background-image 无回退机制
- **位置:** hero 轮播、项目详情 hero
- **问题:** `getVariantUrl()` 直接替换为变体 URL，变体 404 后背景静默空白
- **修复:** 改为绝对定位 `<ResponsiveImage>` 元素，利用 `@error` 事件自动回退

#### C3: `chai2010/webp` 使用 CGO
- **问题:** v1.4.0 实际上依赖 CGO，与 Dockerfile `CGO_ENABLED=0` 不兼容
- **修复:** 用 `golang.org/x/image/webp`（纯 Go 解码）+ `jpeg.Encode`（输出），变体全部 JPEG

#### C4: 前端无法直接拿到 `media.variants` JSON
- **问题:** 项目/案例/律师等实体表存的是 URL 字符串，没有关联到 media 表的 variants
- **修复:** 变体命名规则确定（`{base}_thumb.jpg`），前端纯字符串计算变体 URL

### 设计决策

#### D5: 变体输出 JPEG 而非 WebP
- **理由:** Go 生态无成熟的纯 Go WebP 编码器。JPEG Q=80 对缩略图场景效果足够，文件体积和兼容性最均衡

#### D6: cobra 单二进制多子命令
- **理由:** 不额外编译 CLI 二进制，`./server` 正常启动服务，`./server generate-variants` 执行批量变体。cobra 框架方便以后扩展更多子命令

#### D7: ImageInput 预览用 sm 变体
- **理由:** 管理后台表单中 120px 宽的预览区不需要全分辨率原图，sm(400×300) 足够清晰
