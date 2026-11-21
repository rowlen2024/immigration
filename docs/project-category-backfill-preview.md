# 二期项目类目线上回填预演

2026 年 11 月 21 日 17:27（北京时间），使用北极星移民官网公开接口只读取数。本文记录待审核的回填方案。

数据来自三项线上公开接口。 [页眉导航](https://www.northstarvisa.com/api/v1/navigation?position=header) · [页脚导航](https://www.northstarvisa.com/api/v1/navigation?position=footer) · [已发布项目](https://www.northstarvisa.com/api/v1/projects?status=1)

## 预演结果

拟建 **8 个类目**，拟写入 **38 条项目归属**，覆盖线上全部 **35 个已发布项目**。其中 35 条来自七个业务类目，另 3 条来自「热门项目」。已核对项目 ID、发布状态和 slug，类目内无重复，项目直链与现网地址一致。

| 类目 | 类目页 | 现有导航 ID | 导航显示位置 | 归属数 |
| --- | --- | ---: | --- | ---: |
| 热门项目 | `/projects/category/hot` | 1 | 页眉 | 3 |
| 投资移民 | `/projects/category/investment` | 2 | 页眉和页脚 | 8 |
| 购房移民 | `/projects/category/property` | 92 | 页眉和页脚 | 2 |
| 创业移民 | `/projects/category/entrepreneur` | 45 | 页眉和页脚 | 5 |
| 捐款入籍 | `/projects/category/citizenship-by-donation` | 54 | 页眉和页脚 | 7 |
| 技术/人才/留学移民 | `/projects/category/talent-and-study` | 69 | 页眉 | 6 |
| 长期签证 | `/projects/category/long-stay-visas` | 77 | 页眉 | 4 |
| 跨境服务 | `/projects/category/cross-border-services` | 86 | 页眉 | 3 |

八个栏目目前都是一级 `custom` 导航，目标链接为空。切换时按上表原地绑定，通过类目 slug 查得 `category_id`，再把 `link_type` 改为 `category`。公开链接由后端按 slug 生成。

页眉 10 个一级入口和 53 个可见节点沿用原有导航 ID、层级和排序。移民项目总览 `106`、关于我们 `5` 及全部快捷链接沿用现有目标。页脚四栏沿用 21 条项目直链，其中投资移民有 7 条，类目页根据页眉候选收录 8 个项目。

## 项目归属明细

类目内排序取现网页眉的快捷链接顺序，以 10 为步长。「现有子导航 ID」用于核对快捷链接，项目归属写入使用「项目 ID」。

### 热门项目

| 排序 | 项目 ID | 现有子导航 ID | 项目名称 | 现有详情路径 |
| ---: | ---: | ---: | --- | --- |
| 10 | 1 | 107 | 美国EB-5投资移民 | `/projects/eb5` |
| 20 | 3 | 26 | 巴拿马合格投资者永居 | `/projects/panama` |
| 30 | 2 | 27 | 香港新资本投资者入境计划 | `/projects/cies` |

### 投资移民

| 排序 | 项目 ID | 现有子导航 ID | 项目名称 | 现有详情路径 |
| ---: | ---: | ---: | --- | --- |
| 10 | 1 | 29 | 美国EB-5投资移民 | `/projects/eb5` |
| 20 | 2 | 31 | 香港新资本投资者入境计划 | `/projects/cies` |
| 30 | 6 | 33 | 马耳他永居 | `/projects/malta-permanent-residence` |
| 40 | 7 | 35 | 新加坡自雇PIC | `/projects/singapore-pic` |
| 50 | 8 | 36 | 新加坡家族办公室 | `/projects/singapore-family-office` |
| 60 | 12 | 104 | 土耳其投资入籍 | `/projects/turkey-citizenship` |
| 70 | 9 | 38 | 匈牙利投资移民 | `/projects/hungary-investment` |
| 80 | 5 | 85 | 葡萄牙黄金签证 | `/projects/portugal-fund` |

### 购房移民

| 排序 | 项目 ID | 现有子导航 ID | 项目名称 | 现有详情路径 |
| ---: | ---: | ---: | --- | --- |
| 10 | 3 | 95 | 巴拿马合格投资者永居 | `/projects/panama` |
| 20 | 11 | 96 | 希腊25万欧起购房移民 | `/projects/greece-property` |

### 创业移民

| 排序 | 项目 ID | 现有子导航 ID | 项目名称 | 现有详情路径 |
| ---: | ---: | ---: | --- | --- |
| 10 | 13 | 47 | 美国EB-1C高管移民 | `/projects/usa-eb1c` |
| 20 | 14 | 49 | 日本经营管理签证 | `/projects/japan-business-manager` |
| 30 | 15 | 51 | 新加坡自雇EP | `/projects/singapore-self-employed-ep` |
| 40 | 16 | 53 | 创新创始人签证（英国） | `/projects/uk-innovator-founder` |
| 50 | 35 | 101 | 阿联酋创业移民 | `/projects/a-lian-qiu-2-nian-qi-tou-zi-ren-ju-liu-qian-zheng` |

### 捐款入籍

| 排序 | 项目 ID | 现有子导航 ID | 项目名称 | 现有详情路径 |
| ---: | ---: | ---: | --- | --- |
| 10 | 17 | 56 | 多米尼克捐款入籍 | `/projects/dominica-citizenship` |
| 20 | 18 | 58 | 圣基茨捐款入籍 | `/projects/st-kitts-citizenship` |
| 30 | 19 | 60 | 圣卢西亚捐款入籍 | `/projects/st-lucia-citizenship` |
| 40 | 20 | 62 | 格林纳达捐款入籍 | `/projects/grenada-citizenship` |
| 50 | 21 | 64 | 安提瓜捐款入籍 | `/projects/antigua-citizenship` |
| 60 | 22 | 66 | 瓦努阿图捐款入籍 | `/projects/vanuatu-citizenship` |
| 70 | 23 | 68 | 圣多美捐款入籍 | `/projects/sao-tome-citizenship` |

### 技术/人才/留学移民

| 排序 | 项目 ID | 现有子导航 ID | 项目名称 | 现有详情路径 |
| ---: | ---: | ---: | --- | --- |
| 10 | 24 | 71 | EB-1A杰出人才移民 | `/projects/usa-eb1a` |
| 20 | 25 | 72 | NIW国家利益豁免 | `/projects/usa-niw` |
| 30 | 26 | 74 | 中国香港专才计划 | `/projects/hongkong-professional` |
| 40 | 27 | 75 | 中国香港优才计划 | `/projects/hongkong-quality-migrant` |
| 50 | 28 | 76 | 香港高等院校进修 | `/projects/hongkong-study` |
| 60 | 36 | 102 | 中国香港高端人才 | `/projects/xiang-gang-gao-duan-ren-cai` |

### 长期签证

| 排序 | 项目 ID | 现有子导航 ID | 项目名称 | 现有详情路径 |
| ---: | ---: | ---: | --- | --- |
| 10 | 29 | 79 | 美国L-1签证 | `/projects/usa-l1a` |
| 20 | 37 | 105 | 美国 O1 杰出人才签证 | `/projects/mei-guo-o1-jie-chu-ren-cai-qian-zheng` |
| 30 | 30 | 81 | 西班牙非盈利居留 | `/projects/spain-non-lucrative` |
| 40 | 31 | 83 | 泰国精英签证（泰国尊荣卡） | `/projects/thailand-elite` |

### 跨境服务

| 排序 | 项目 ID | 现有子导航 ID | 项目名称 | 现有详情路径 |
| ---: | ---: | ---: | --- | --- |
| 10 | 32 | 87 | 美国远程结婚 | `/projects/mei-guo-yuan-cheng-jie-hun` |
| 20 | 33 | 89 | 香港银行公司账户开户 | `/projects/xiang-gang-yin-hang-gong-si-zhang-hu-kai-hu` |
| 30 | 34 | 90 | 马来西亚税号（TIN） | `/projects/ma-lai-xi-ya-shui-hao-（-t-i-n-）` |

## 应用前核对

- 使用后台数据复核导航 ID、旧链接类型、显示位置、父子关系及项目 ID、slug、发布状态。公开快照无法覆盖隐藏导航、草稿和下架项目。这些记录由运营按需确认归属。
- 优先审核「投资移民」里的新加坡自雇 PIC（项目 ID `7`）和家族办公室（项目 ID `8`）。预演沿用了现网页眉归属。
- 核对类目 slug 是否已被占用，再以草稿准备八个类目和归属。录入介绍并验证类目页后，再发布类目和绑定导航。正式绑定前备份八条导航原值，便于回滚。
- 回滚时先恢复八条导航的旧链接类型和目标，再按需回退表结构。
