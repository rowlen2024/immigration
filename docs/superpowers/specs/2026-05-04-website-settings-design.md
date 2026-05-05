# 网站设置 (Website Settings) 功能设计

## 概述

在管理后台新增"网站设置"菜单，允许管理员和编辑者配置全站级参数（网站标题、Logo、SEO/GEO 元数据、联系方式、页脚信息、第三方代码等），替代当前 Header/Footer/useSeo 中的硬编码值。

## 数据模型

### 复用 `home_configs` 表

沿用现有 key-value JSON 模式，新增 `site` config key。所有设置存储在一个 JSON 对象中，与其他 config key（如 `home`）共存。

```json
// home_configs 表 config_key = "site" 的 config_value:
{
  "site_name": "MyGo移民",
  "site_logo": "",
  "site_favicon": "",
  "seo_title": "{site_name} | 专业投资移民服务",
  "seo_description": "MyGo 专注美国EB-5、香港CIES及巴拿马投资移民...",
  "seo_keywords": "投资移民,EB-5,香港移民,巴拿马移民",
  "og_image": "",
  "canonical_base": "https://www.mygo-immigration.com",
  "organization_name": "MyGo Immigration Consulting Ltd",
  "organization_description": "",
  "organization_logo": "",
  "organization_url": "https://www.mygo-immigration.com",
  "same_as": [],
  "contact_phone": "400-xxx-xxxx",
  "contact_email": "info@mygo-immigration.com",
  "contact_address": "上海市浦东新区陆家嘴金融中心",
  "contact_wechat": "MyGo_Immigration",
  "ga_tracking_id": "",
  "baidu_tongji_id": "",
  "custom_head_code": "",
  "custom_body_code": "",
  "copyright_text": "© {year} {site_name}. All rights reserved.",
  "icp_number": ""
}
```

### 字段定义

| 字段 | 类型 | 说明 | 组 |
|------|------|------|-----|
| `site_name` | string | 网站名称 | 基础信息 |
| `site_logo` | string | Logo 图片 URL | 基础信息 |
| `site_favicon` | string | Favicon URL | 基础信息 |
| `seo_title` | string | 首页标题模板，`{site_name}` 占位 | SEO |
| `seo_description` | string | 全局 Meta Description | SEO |
| `seo_keywords` | string | 全局 Meta Keywords | SEO |
| `og_image` | string | 默认 OG 分享图 URL | SEO |
| `canonical_base` | string | 首选域名 base URL | SEO |
| `organization_name` | string | Schema.org Organization 法律实体名称 | GEO |
| `organization_description` | string | 机构描述（知识图谱摘要） | GEO |
| `organization_logo` | string | 知识图谱 Logo URL（≥112×112px） | GEO |
| `organization_url` | string | 官网 URL（实体认证） | GEO |
| `same_as` | string[] | 社交媒体链接数组（linkedin/小红书/公众号等） | GEO |
| `contact_phone` | string | 客服电话 | 联系方式 |
| `contact_email` | string | 客服邮箱 | 联系方式 |
| `contact_address` | string | 公司地址 | 联系方式 |
| `contact_wechat` | string | 企业微信 | 联系方式 |
| `ga_tracking_id` | string | Google Analytics ID | 第三方代码 |
| `baidu_tongji_id` | string | 百度统计 ID | 第三方代码 |
| `custom_head_code` | string | 自定义 `<head>` 内代码 | 第三方代码 |
| `custom_body_code` | string | 自定义 `</body>` 前代码 | 第三方代码 |
| `copyright_text` | string | 版权文字，`{year}` / `{site_name}` 占位 | 页脚 |
| `icp_number` | string | ICP 备案号 | 页脚 |

## API 设计

### GET `/api/v1/site-config`（公开）

返回网站设置。无需认证。

```json
// Response: { code: 200, data: { ... } }
```

### PUT `/api/v1/admin/site-config`（需认证）

更新网站设置。

- 权限：`content:write`
- Body: 完整的 site config JSON（全量替换）
- Response: `{ code: 200, data: null }`

## 前端设计

### 管理页面 `/admin/settings`

- 布局：`admin.vue`
- 路由鉴权：`auth.ts` 中间件
- 权限要求：有 `content:write` 权限的用户可见（通过 auth store 的 role 判断）

### 分组表单

页面按 6 个分组组织，使用折叠面板（ElCollapse）或卡片（ElCard）+ 锚点导航：

1. **基础信息** — site_name, site_logo, site_favicon
2. **SEO 设置** — seo_title 到 canonical_base
3. **GEO / 结构化数据** — organization_* 到 same_as
4. **联系方式** — contact_phone 到 contact_wechat
5. **第三方代码** — ga_tracking_id 到 custom_body_code
6. **页脚设置** — copyright_text, icp_number

### 交互细节

- 每个参数标签旁显示 `?` 图标（ElTooltip），鼠标悬停显示该参数的用途说明
- 顶栏或底部固定"保存"按钮（提交全量 JSON）
- 保存成功后显示 ElMessage 提示
- Logo/Favicon 字段 使用 media picker（可选，初期先支持输入 URL）
- `same_as` 使用可动态增减的输入列表
- `custom_head_code` / `custom_body_code` 使用 textarea（monospace 字体）

### Tooltip 说明文案映射

每个字段的 tips 说明：

```
site_name: "网站在导航栏、浏览器标题栏等位置显示的名称"
site_logo: "网站主 Logo 图片地址，显示在页头导航栏"
site_favicon: "浏览器标签页和书签栏显示的小图标，建议 32×32px"
seo_title: "搜索引擎结果页显示的标题模板。{site_name} 会被替换为网站名称"
seo_description: "搜索引擎结果页显示的页面描述，建议 120-160 字，各内页无独立描述时使用此值"
seo_keywords: "页面关键词，用逗号分隔。现代搜索引擎权重已降低，但仍建议填写"
og_image: "在微信、Facebook 等社交平台分享链接时显示的默认预览图"
canonical_base: "网站首选访问域名，用于规范搜索引擎索引（如 https://www.example.com）"
organization_name: "向 Google 知识图谱和 AI 搜索声明的企业法律实体全称"
organization_description: "用于生成 AI 搜索摘要和知识面板的企业介绍，建议包含成立时间、核心业务和规模"
organization_logo: "Google 知识面板展示的 Logo，建议 112×112px 以上高清 PNG"
organization_url: "声明机构的官方网址，用于多平台交叉验证网站真实性"
same_as: "与官网关联的社交媒体链接（LinkedIn、公众号、小红书等）。AI 搜索引擎通过双向链接验证实体可信度"
contact_phone: "网站底部和结构化数据中展示的客服电话"
contact_email: "网站底部和结构化数据中展示的客服邮箱"
contact_address: "公司办公地址，显示在网站底部"
contact_wechat: "企业微信号，显示在网站底部联系方式区域"
ga_tracking_id: "Google Analytics 4 衡量 ID（格式：G-XXXXXXXX），用于网站流量分析"
baidu_tongji_id: "百度统计站点 ID，用于中国国内流量分析"
custom_head_code: "插入到每个页面 <head> 标签内的自定义代码（meta 验证标签、第三方脚本等）"
custom_body_code: "插入到每个页面 </body> 闭合标签前的自定义代码"
copyright_text: "页脚版权声明。{year} 动态替换为当前年份，{site_name} 替换为网站名称"
icp_number: "ICP 备案号（中国大陆运营网站必需），如 沪ICP备XXXXXXXX号"
```

## 消费者改造

### Header.vue
- `logo-text` 改为从 site config 读取 `site_name`
- 如果 `site_logo` 有值，显示 `<img>` 替代纯文本

### Footer.vue
- 联系方式（phone/email/address/wechat）从 site config 读取
- 版权声明从 `copyright_text` 读取，替换 `{year}` 和 `{site_name}`

### useSeo.ts
- 默认标题后缀从 `site_name` 读取（不再硬编码"MyGo移民"）
- 注入全局 OG 标签和 canonical

### default.vue / app.vue
- `<head>` 中注入 `custom_head_code`
- 页面末尾注入 `custom_body_code`
- 注入 GA / 百度统计脚本（如果已填写）
- 注入 Organization JSON-LD schema（如果 GEO 参数已填写）
- 设置 `<link rel="icon">` 为 `site_favicon`

## 完整参数列表总结

共 22 个参数，分 6 组。所有参数均为可选（后端提供默认空值，前端渐进增强）。

## 实现范围

### 后端
1. 复用 `HomeConfigService`，新增 `GetSiteConfig()` / `UpdateSiteConfig()` 方法
2. 新增公开路由 `GET /api/v1/site-config`
3. 新增管理路由 `PUT /api/v1/admin/site-config`（`content:write`）

### 前端
1. 新增 `useSiteConfig` composable（fetch + 响应式缓存）
2. 新增 `/admin/settings` 页面（分组表单 + Tooltip）
3. 修改 `Header.vue` 使用 site config
4. 修改 `Footer.vue` 使用 site config
5. 修改 `useSeo.ts` 使用 site config
6. 修改 `default.vue` 注入结构化数据
7. 在 `admin.vue` 侧边栏新增"网站设置"菜单项
8. `app.vue` / 全局注入 favicon、第三方代码、JSON-LD

## 不纳入本次范围

- 不拆分设置项到独立的 `site_configs` 表
- 不做设置历史版本/回滚
- 不做多站点/多语言设置
- Logo 上传暂用现有媒体库 + URL 填入方式，不做专门的 Logo 上传裁剪
