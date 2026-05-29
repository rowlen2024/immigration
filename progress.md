# Progress Log — MyGo 移民

**关联计划:** task_plan.md

---

## Session 3 — 2026-05-23 — 图片变体系统 + 前端接入

### 分析阶段
- 浏览器打开首页 http://localhost:3000/ 和 http://47.113.121.245/
- Explore agent 完整扫描后端 media 上传/存储/清理流程
- Explore agent 研究 Go 图片处理库，确定 `imaging` + `x/image/webp`
- Explore agent 扫描前端全部 32 处图片加载点，分类为"应用变体"和"保持原 URL"

### 实施阶段 — Phase 6 后端
- `go get imaging` + `go get chai2010/webp` → 编译失败（CGO）→ 换 `x/image/webp` 解码 + JPEG 编码输出
- 新建 `service/media_thumbnail.go`: GenerateVariants + CompressIfLarge
- `model/media.go`: 新增 `Variants datatypes.JSON`
- `handler/media_handler.go`: 落盘 → 压缩(>5MB) → 生成 4 变体 → 写 DB
- `service/media_svc.go`: CleanupUnused 遍历 variants JSON 同步删文件
- 迁移 `000024_add_media_variants.up.sql`
- `go build ./...` 通过，全部测试 PASS（唯一失败 TestAuthHandler 预存问题）

### 实施阶段 — Phase 7 前端
- 新建 `utils/image.ts` (getVariantUrl) + `components/ResponsiveImage.vue`
- 共享组件: CaseCard(sm), LawyerCarousel(sm), TestimonialCarousel(thumb), MediaPicker(sm/md), ImageInput(sm)
- 公众页面: index.vue(Hero lg/卡片 sm/横幅 lg), projects/[slug].vue(Hero lg/新闻 thumb)
- 管理后台: cases/lawyers/projects/homepage/media → thumb; media preview → md
- 不改动: Logo, QR 码, Favicon, OG Image, RichEditor, Settings ImageInput
- `npx nuxi typecheck` 无新增错误
- 浏览器验证: Network 显示变体 URL 正确发起 → 404 → @error 回退原图

### 实施阶段 — Phase 8 修复 + CLI
- Hero 背景 `background-image` → `<ResponsiveImage>` 绝对定位
- 删除旧 `cmd/generate_variants/` 目录
- cobra 重构: `cmd/server/main.go`(入口) + `root.go`(服务) + `generate.go`(子命令)
- Dockerfile 还原为单二进制
- `go build ./...` 通过, `go run ./cmd/server --help` 正常

### 当前状态
- 后端 + 前端代码就绪，等待部署
- 部署后需执行: `docker exec <容器> /app/server generate-variants` 为 163 个历史图片生成变体
- 旧图片通过 ResponsiveImage @error 自动回退，零白屏过渡
