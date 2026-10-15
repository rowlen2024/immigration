# 联系咨询邮件通知

## 目标与范围

用户通过 `POST /api/v1/leads` 提交联系我们表单，咨询记录保存成功后，异步向运营通知邮箱发送一封邮件。接口仍按原契约返回咨询记录，邮件结果不影响提交结果。

## 配置契约

网站设置新增 `notification_email`，与 `contact_email` 客服邮箱独立。后台通过现有 `GET/PUT /api/v1/admin/site-config` 读取、保存，公开网站设置接口沿用同一配置结构。

- `notification_email` 默认为空，为空时不发送。
- 设置保存入口校验非空值为邮箱地址，保留输入原值。
- 配置保存在 `home_configs` 的 `site` JSON 中，无需数据库迁移；旧配置读取后新字段为空。
- SMTP 连接参数来自后端环境变量，不进入网站设置。

| 环境变量 | 含义 | 默认值 |
| --- | --- | --- |
| `SMTP_HOST` | SMTP 服务地址 | `smtp.exmail.qq.com` |
| `SMTP_PORT` | 隐式 TLS 端口 | `465` |
| `SMTP_USERNAME` | 发件邮箱及认证账号 | 空，部署时配置 |
| `SMTP_PASSWORD` | 邮箱密码或邮箱服务要求的 SMTP 授权码 | 空，部署时配置 |

发件账号配置为 `luolian@mixlinker.com`。密码仅通过部署环境注入，不写入源码、示例、日志或文档。Docker Compose 将上述变量传给后端容器。SMTP 连接始终启用 TLS 并验证服务器证书，连接与后续 SMTP 操作均受超时约束。

## 职责与流程

1. Handler 维持现有表单绑定及响应行为。
2. `LeadService.Create` 完成现有校验、净化及数据库保存；保存失败直接返回，不启动通知。
3. 保存成功后复制咨询记录，由独立 goroutine 调用通知器，不持有请求上下文或 Gin Context。
4. 通知器读取当前网站设置；通知邮箱为空则结束，否则组装纯文本邮件。
5. 邮件包含咨询编号、姓名、电话、邮箱、意向项目、留言和提交时间。用户输入只进入正文，中文主题及正文使用 MIME 编码。
6. 读取配置或 SMTP 发送失败只记录咨询编号及错误，不重试，不回滚已保存记录。

邮件发送器和配置读取器通过窄接口注入通知器，通知器作为 `LeadService` 的必需依赖注入；测试始终提供替身。通知器不重复校验已由入口保证的表单和邮箱字段。

## 一致性与交付语义

异步任务为进程内尽力通知，进程退出时尚未完成的任务可能丢失。每次成功提交触发一次尝试，不增加通知状态、补偿任务或持久化队列。清空通知邮箱可停止后续邮件，SMTP 配置变更需要重启后端生效。

## 验证方案

- CLI：配置单元测试、本地 TLS SMTP 协议测试、通知器测试、咨询异步行为及失败日志测试、设置字段保存与绑定测试。
- 测试覆盖旧配置、空邮箱、配置解析失败、邮件发送失败、数据库保存失败、中文和多行正文、异步不阻塞及咨询快照。
- 回归：`go test ./...`、`go build ./...`、`go vet ./...`、`npx nuxi typecheck` 和前端构建；浏览器检查后台字段及保存请求。
- Skill：Context7 查询 Go SMTP/TLS API，按完成前验证和代码评审流程检查证据、差异及校验责任边界。
- 测试使用本地 SMTP 服务或替身，不向真实运营邮箱发送测试邮件。

## 本次验证记录

- `go build ./...`、`go vet ./...` 通过。
- `go test -race ./internal/config ./internal/mailer ./internal/service ./internal/handler -count=1` 通过。
- `go test ./... -count=1` 中，仅已有的 `TestPageRepoFindRelatedBySlugMySQL` 因本机 `localhost:3307` 无 MySQL 服务而失败；跳过该用例后其余测试通过。依赖 `MIGRATION_TEST_DSN` 的迁移集成测试按现有约定跳过。
- `npm run build` 通过，存在已有依赖弃用及包体积警告。
- `npx nuxi typecheck` 在 `settings.vue` 原有模板中报告 `idx/field` 作用域错误；对修改前设置页复验得到相同错误。
- `docker compose config --quiet` 通过；实际部署仍需提供 SMTP 账号和密码。
- 浏览器使用隔离上下文及内存接口替身，验证独立读取、修改、清空运营邮箱、重新进入页面后回显，以及 1440px 桌面和 390px 移动视口的新字段布局；未改动真实网站设置。
- 已对腾讯 SMTP 端点完成 TLS 与 EHLO 只读握手，确认其支持 `AUTH PLAIN`；未验证真实账号认证或实际收件。

部署时配置上述 SMTP 环境变量，启动数据库和后端，并在网站设置中填写运营通知邮箱，即可启用后续咨询通知。
