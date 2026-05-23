# Task Plan — MyGo 移民

**更新时间:** 2026-05-23

---

## Phase 6: 后端图片自动缩略图 + 压缩 ✅

**目标:** 上传时自动生成多尺寸变体 + >5MB 压缩

- [x] 选型：`disintegration/imaging` + `golang.org/x/image/webp`（纯 Go，兼容 CGO_ENABLED=0）
- [x] 新建 `service/media_thumbnail.go` — `GenerateVariants()` + `CompressIfLarge()`
- [x] `model/media.go` 新增 `Variants datatypes.JSON` 字段
- [x] `handler/media_handler.go` 上传流程接入：落盘 → 压缩(>5MB) → 生成变体 → 存 DB
- [x] `service/media_svc.go` CleanupUnused 同步删除变体文件
- [x] 迁移 `000024_add_media_variants.up.sql`

**变体规格:** thumb(200×200 裁切) / sm(400×300 Fit) / md(800×450 Fit) / lg(1920×800 Fit) — 全部 JPEG Q=80

---

## Phase 7: 前端图片变体接入 ✅

**目标:** 前端按场景加载对应尺寸变体，旧图自动回退原 URL

- [x] 新建 `utils/image.ts` — `getVariantUrl(url, variant)` 纯字符串计算
- [x] 新建 `components/ResponsiveImage.vue` — 变体 404 自动回退原图
- [x] 共享组件：CaseCard(sm), LawyerCarousel(sm), TestimonialCarousel(thumb), MediaPicker(sm/md), ImageInput(sm)
- [x] 公众页面：index.vue(Hero lg, 卡片 sm, 横幅 lg), projects/[slug].vue(Hero lg, 新闻 thumb)
- [x] 管理后台：cases/lawyers/projects/homepage/media 表格缩略图(thumb), 预览(md)

---

## Phase 8: Hero 背景回退修复 + CLI 命令行工具 ✅

**问题:** CSS background-image 无 @error 回退机制，变体 404 后背景空白

- [x] index.vue hero slide: `:style backgroundImage` → `<ResponsiveImage class="hero-slide-bg">`
- [x] projects/[slug].vue detail hero: 同上，加独立 gradient overlay div
- [x] 删除未使用的 `getVariantUrl` import

**CLI 工具:** `./server generate-variants` 为历史图片批量生成变体
- [x] cobra 三文件结构：main.go / root.go / generate.go
- [x] 移除 Dockerfile 多余 binary，单二进制多子命令

---

## 待办

- [ ] 部署后端 → 运行 `docker exec <容器> /app/server generate-variants` 为历史图片生成变体
- [ ] 后续 PR: 前端升级使用 `<picture>` / `srcset` 进一步提升性能
