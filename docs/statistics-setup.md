# 网站统计 ID 获取指南

后台"站点设置"页面需要填入两个统计 ID 来启用网站流量分析。

## Google Analytics ID（格式：`G-XXXXXXXXXX`）

1. 打开 [analytics.google.com](https://analytics.google.com)，用 Google 账号登录
2. 点击 **"开始衡量"** → 填写账号名称 → 创建媒体资源，平台选择 **"网络"**
3. 输入网站域名和名称，点击 **"创建数据流"**
4. 创建后在详情页顶部获取 **衡量 ID**（`G-` 开头）

## 百度统计 ID（格式：32 位字符串）

1. 打开 [tongji.baidu.com](https://tongji.baidu.com)，用百度账号登录
2. 点击 **"新增网站"**，填写域名后完成添加
3. 添加成功后在"代码获取"页面，JS 脚本 URL 中 `hm.js?` 后面的字符串即为统计 ID：
   ```
   https://hm.baidu.com/hm.js?xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
                                         ↑ 这一串就是 ID ↑
   ```

拿到两个 ID 后，填入后台 **站点设置** 页面的对应字段，前台所有页面即会自动加载统计脚本。
