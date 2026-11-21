# 项目地区导航第一期发布与验收

第一期交付页眉按项目地区分组、页脚直链迁移和项目后台地区选择。地区来源是项目关联的字典，导航继续保存项目链接。类目页和项目多类目归属属于后续版本。

## 开发与验收

- `backend/database/migrations/000030_add_project_regions.up.sql` 新增地区字典及项目关联，回填明确匹配的旧地区。
- `/api/v1/project-regions?ids=1,2` 批量返回已发布项目的地区，独立路由可避免占用项目详情的 slug。`Header.vue` 使用现有导航中的项目 ID 关联地区，服务端生成地区标题和项目链接。页脚原有分栏继续从现有导航读取直链。
- `frontend/pages/admin/projects.vue` 使用后台地区列表。既有 `country` 和 `flag_emoji` 字段仍供旧客户端使用。
- `backend/cmd/server/navigation_regions_command.go` 提供导航数据的预演、应用和回滚。用实际页眉、页脚目标集合核对迁移结果。
- 发布前运行 `go build ./...`、`go test ./...`、`npx nuxi typecheck` 和 `npm run build`。检查首页原始 HTML、桌面及移动菜单、地区选择、旧筛选及页脚。

## 生产发布顺序

1. 记录当前镜像版本，备份受影响的项目和导航数据。导出 `/api/v1/navigation?position=header` 和 `/api/v1/navigation?position=footer` 的当前结果。
2. 上线后端并执行地区迁移。用下面的查询检查缺少地区关联的项目，包含已下架的项目。通过后台为未匹配项确认归属。「中国香港高端人才」的旧地区为“中国”，需要按单个项目确认。

   ```sql
   SELECT id, slug, country, status
   FROM projects
   WHERE deleted_at IS NULL AND region_id IS NULL
   ORDER BY status DESC, id;
   ```

3. 上线能读取旧树及扁平导航的前端。在导航管理下新增顶级「移民项目」，链接类型为自定义，路径为 `/projects`，位置为页眉，排在项目分组前。先确认桌面及手机都能直接点进总览。
4. 在已确认指向本次目标库的后端环境生成迁移预演。以下候选 ID 来自线上公开快照，应用前要和当前后台节点及预演里的原值核对。使用容器里的 `./server` 时，将命令前缀 `go run ./cmd/server` 替换为 `./server`。将报告保存在容器外或挂载目录，以便回滚。

   ```powershell
   $regionNodeIds = '22,23,24,28,30,32,34,103,37,84,93,94,46,48,50,52,100,55,57,59,61,63,65,67,70,73,78,80,82'
   go run ./cmd/server flatten-navigation-regions --ids $regionNodeIds > navigation-region-plan.json
   ```

5. 核对预演后应用。工具会重新校验预演内容，拒绝应用已过期的报告。

   ```powershell
   go run ./cmd/server flatten-navigation-regions --apply navigation-region-plan.json
   ```

6. 等待或按线上方式刷新 SSR 缓存，再抓取首份 HTML 和公开导航。检查页脚仍有原 21 条项目直链，热门项目仍按原顺序展示 EB-5、巴拿马、香港。检查同一地区的 L-1、O-1 都可到达，香港的别名归到同一地区。

## 回滚

使用保存的预演报告恢复受影响的导航记录。若工具发现上线后编辑，核对冲突后再处理。新增的「移民项目」入口在后台单独回滚。地区表和项目地区列可以与旧版本共存，回滚展示时可保留已确认的地区数据。

```powershell
go run ./cmd/server flatten-navigation-regions --rollback navigation-region-plan.json
```
