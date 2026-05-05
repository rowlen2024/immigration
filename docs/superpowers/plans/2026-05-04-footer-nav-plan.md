# 页脚导航动态化 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Footer 导航从硬编码改为通过 API 动态获取，后台统一管理

**Architecture:** Navigation 模型新增 `display_position` 枚举字段（header/footer/both），公开 API 通过 query param `?position=` 过滤，前端 Footer 调用新接口渲染

**Tech Stack:** Go/Gin/GORM (backend), Nuxt 3/Vue 3/Element Plus (frontend), MySQL

---

### Task 1: Database Migration

**Files:**
- Create: `database/migrations/000017_add_navigation_display_position.up.sql`
- Create: `database/migrations/000017_add_navigation_display_position.down.sql`

- [ ] **Step 1: Create up migration**

`database/migrations/000017_add_navigation_display_position.up.sql`:
```sql
SET NAMES utf8mb4;

ALTER TABLE `navigations`
    ADD COLUMN `display_position` VARCHAR(16) NOT NULL DEFAULT 'header' AFTER `status`;

ALTER TABLE `navigations`
    ADD INDEX `idx_nav_display_position` (`display_position`);
```

- [ ] **Step 2: Create down migration**

`database/migrations/000017_add_navigation_display_position.down.sql`:
```sql
SET NAMES utf8mb4;

ALTER TABLE `navigations` DROP INDEX `idx_nav_display_position`;
ALTER TABLE `navigations` DROP COLUMN `display_position`;
```

- [ ] **Step 3: Run migration against local MySQL**

```bash
docker exec -i mygo-mysql mysql -uroot -proot mygo < database/migrations/000017_add_navigation_display_position.up.sql
```
Expected: no errors, column added.

- [ ] **Step 4: Commit**

```bash
git add database/migrations/000017_add_navigation_display_position.up.sql database/migrations/000017_add_navigation_display_position.down.sql
git commit -m "feat: add display_position column to navigations"
```

---

### Task 2: Model - Add DisplayPosition Field

**Files:**
- Modify: `backend/internal/model/navigation.go`

- [ ] **Step 1: Add field to Navigation struct**

In `backend/internal/model/navigation.go`, after `Status` line (line 18), insert:

```go
DisplayPosition string         `gorm:"size:16;not null;default:'header'" json:"display_position"`
```

Result struct should have these fields in order:
```go
type Navigation struct {
	ID              uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Label           string         `gorm:"size:255;not null" json:"label"`
	Link            *string        `gorm:"size:512;default:null" json:"link"`
	LinkType        string         `gorm:"size:32;not null;default:'custom'" json:"link_type"`
	ProjectID       *uint64        `gorm:"index;default:null" json:"project_id"`
	PageID          *uint64        `gorm:"index;default:null" json:"page_id"`
	ParentID        *uint64        `gorm:"index;default:null" json:"parent_id"`
	SortOrder       int            `gorm:"not null;default:0" json:"sort_order"`
	Status          bool           `gorm:"not null;default:1" json:"status"`
	DisplayPosition string         `gorm:"size:16;not null;default:'header'" json:"display_position"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	Children        []Navigation   `gorm:"foreignKey:ParentID" json:"children"`
	Project         *Project       `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Page            *Page          `gorm:"foreignKey:PageID" json:"page,omitempty"`
}
```

- [ ] **Step 2: Verify compilation**

```bash
cd backend && go build ./...
```
Expected: no errors.

- [ ] **Step 3: Commit**

```bash
git add backend/internal/model/navigation.go
git commit -m "feat: add DisplayPosition field to Navigation model"
```

---

### Task 3: Repository - Add FindAllActiveByPosition

**Files:**
- Modify: `backend/internal/repository/interfaces.go`
- Modify: `backend/internal/repository/nav_repo.go`

- [ ] **Step 1: Add method to interface**

In `backend/internal/repository/interfaces.go`, inside `NavigationRepository` interface, add after `FindAllActive()`:

```go
FindAllActiveByPosition(position string) ([]model.Navigation, error)
```

- [ ] **Step 2: Implement in NavRepo**

In `backend/internal/repository/nav_repo.go`, add after `FindAllActive()`:

```go
func (r *NavRepo) FindAllActiveByPosition(position string) ([]model.Navigation, error) {
	var items []model.Navigation
	err := r.db.Where("status = 1 AND deleted_at IS NULL AND display_position IN ?", []string{position, "both"}).
		Order("sort_order asc, id asc").
		Find(&items).Error
	return items, err
}
```

- [ ] **Step 3: Verify compilation**

```bash
cd backend && go build ./...
```
Expected: no errors.

- [ ] **Step 4: Commit**

```bash
git add backend/internal/repository/interfaces.go backend/internal/repository/nav_repo.go
git commit -m "feat: add FindAllActiveByPosition to NavRepo"
```

---

### Task 4: Service - Update GetTree with Position Parameter

**Files:**
- Modify: `backend/internal/service/nav_svc.go`

- [ ] **Step 1: Replace GetTree and GetAdminTree signatures**

In `backend/internal/service/nav_svc.go`, change `GetTree` (line 18-26):

```go
func (s *NavService) GetTree(position string) ([]model.Navigation, error) {
	items, err := s.repo.FindAllActiveByPosition(position)
	if err != nil {
		return nil, err
	}
	tree := buildTree(items, nil)
	s.fillLinks(tree)
	return tree, nil
}
```

And keep `GetAdminTree` unchanged (it returns ALL items for admin management, not position-filtered).

- [ ] **Step 2: Verify compilation**

```bash
cd backend && go build ./...
```
Expected: no errors (handler will break until Task 5, but we fix that next).

- [ ] **Step 3: Commit**

```bash
git add backend/internal/service/nav_svc.go
git commit -m "feat: add position parameter to NavService.GetTree"
```

---

### Task 5: Handler - Update GetNavigation with Query Param

**Files:**
- Modify: `backend/internal/handler/nav_handler.go`

- [ ] **Step 1: Update GetNavigation to accept position query param**

In `backend/internal/handler/nav_handler.go`, replace `GetNavigation` (lines 12-19):

```go
func (h *Handler) GetNavigation(c *gin.Context) {
	position := c.DefaultQuery("position", "header")
	if position != "header" && position != "footer" {
		c.JSON(http.StatusBadRequest, dto.Error(400, "position must be 'header' or 'footer'"))
		return
	}
	tree, err := h.svc.Nav.GetTree(position)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(tree))
}
```

- [ ] **Step 2: Verify compilation**

```bash
cd backend && go build ./...
```
Expected: no errors.

- [ ] **Step 3: Run existing tests**

```bash
cd backend && go test ./... -count=1 2>&1 | tail -20
```
Expected: all tests pass (existing nav-related tests may need update if they call GetTree — check and fix if needed).

- [ ] **Step 4: Commit**

```bash
git add backend/internal/handler/nav_handler.go
git commit -m "feat: add position query param to GET /navigation"
```

---

### Task 6: Admin Frontend - Add Display Position to Navigation Form

**Files:**
- Modify: `frontend/pages/admin/navigation.vue`

- [ ] **Step 1: Add display_position to defaultForm**

In the `defaultForm` function, add after `status: true,` (around line 208):

```ts
display_position: 'header',
```

Full `defaultForm` return:
```ts
const defaultForm = (): {
  label: string;
  link: string;
  link_type: string;
  project_id: number | null;
  page_id: number | null;
  parent_id: number | null;
  sort_order: number;
  status: boolean;
  display_position: string;
} => ({
  label: '',
  link: '',
  link_type: 'custom',
  project_id: null,
  page_id: null,
  parent_id: null,
  sort_order: 0,
  status: true,
  display_position: 'header',
});
```

- [ ] **Step 2: Add display_position to form type**

Update `form` initialization type implicitly via defaultForm — no extra step needed. But ensure the `form` reactive has the right initial values by checking it calls `defaultForm()`.

- [ ] **Step 3: Add form field in template**

In template, after the `<el-row>` containing status switch (around line 137), add:

```html
<el-form-item label="显示位置" prop="display_position">
  <el-radio-group v-model="form.display_position">
    <el-radio value="header">仅头部</el-radio>
    <el-radio value="footer">仅页脚</el-radio>
    <el-radio value="both">头部+页脚</el-radio>
  </el-radio-group>
</el-form-item>
```

- [ ] **Step 4: Update openEdit to copy display_position**

In `openEdit`, add to `Object.assign` (around line 324):

```ts
display_position: row.display_position || 'header',
```

- [ ] **Step 5: Update handleSave to include display_position**

In `handleSave`, add to the `body` object (around line 348):

```ts
display_position: form.display_position,
```

- [ ] **Step 6: Display position tag in tree view (optional but nice)**

In the tree node template, after the link_type tag (around line 27), add:

```html
<el-tag
  v-if="data.display_position && data.display_position !== 'header'"
  :type="data.display_position === 'both' ? 'success' : 'warning'"
  size="small"
  effect="plain"
>
  {{ data.display_position === 'both' ? '头部+页脚' : '仅页脚' }}
</el-tag>
```

- [ ] **Step 7: Type check**

```bash
cd frontend && npx nuxi typecheck
```
Expected: no errors.

- [ ] **Step 8: Commit**

```bash
git add frontend/pages/admin/navigation.vue
git commit -m "feat: add display_position field to admin navigation form"
```

---

### Task 7: Footer.vue - Dynamic Navigation

**Files:**
- Modify: `frontend/components/global/Footer.vue`

- [ ] **Step 1: Replace script section**

Replace `<script setup>` block entirely:

```vue
<script setup lang="ts">
const { siteConfig, fetch: fetchSiteConfig } = useSiteConfig();

interface NavItem {
  id: number;
  label: string;
  link: string;
  children: NavItem[];
  status: boolean;
}

const FALLBACK_FOOTER: NavItem[] = [
  {
    id: 1, label: '移民项目', link: '', status: true,
    children: [
      { id: 11, label: '美国EB-5投资移民', link: '/projects/eb5', children: [], status: true },
      { id: 12, label: '香港投资移民', link: '/projects/cies', children: [], status: true },
      { id: 13, label: '巴拿马购房移民', link: '/projects/panama', children: [], status: true },
      { id: 14, label: '项目对比', link: '/compare', children: [], status: true },
    ],
  },
  {
    id: 2, label: '关于我们', link: '', status: true,
    children: [
      { id: 21, label: '公司简介', link: '/pages/about', children: [], status: true },
      { id: 22, label: '成功案例', link: '/cases', children: [], status: true },
      { id: 23, label: '常见问题', link: '/faq', children: [], status: true },
      { id: 24, label: '联系我们', link: '/pages/contact', children: [], status: true },
    ],
  },
];

const footerNav = ref<NavItem[]>([]);

const fetchFooterNav = async () => {
  try {
    const api = useApi();
    const data = await api<NavItem[]>('/navigation?position=footer');
    if (data && (data as NavItem[]).length > 0) {
      footerNav.value = data as NavItem[];
    } else {
      footerNav.value = FALLBACK_FOOTER;
    }
  } catch {
    footerNav.value = FALLBACK_FOOTER;
  }
};

const copyrightText = computed(() => {
  const template = siteConfig.value?.copyright_text || '© {year} {site_name}. All rights reserved.';
  return template
    .replace('{year}', String(new Date().getFullYear()))
    .replace('{site_name}', siteConfig.value?.site_name || 'MyGo移民');
});

onMounted(() => {
  fetchFooterNav();
  fetchSiteConfig();
});
</script>
```

- [ ] **Step 2: Replace template — hardcoded nav columns with dynamic rendering**

Replace the middle two `<div class="footer-column">` sections (移民项目 and 关于我们, lines 12-30) with:

```vue
<div v-for="col in footerNav" :key="col.id" class="footer-column">
  <h3 class="footer-heading">{{ col.label }}</h3>
  <ul class="footer-links">
    <li v-for="child in col.children" :key="child.id">
      <NuxtLink :to="child.link" class="footer-link">{{ child.label }}</NuxtLink>
    </li>
  </ul>
</div>
```

Final template structure should be:
```vue
<template>
  <footer class="site-footer">
    <div class="footer-container">
      <!-- Brand column (unchanged) -->
      <div class="footer-column footer-brand">
        ...
      </div>

      <!-- Dynamic nav columns -->
      <div v-for="col in footerNav" :key="col.id" class="footer-column">
        <h3 class="footer-heading">{{ col.label }}</h3>
        <ul class="footer-links">
          <li v-for="child in col.children" :key="child.id">
            <NuxtLink :to="child.link" class="footer-link">{{ child.label }}</NuxtLink>
          </li>
        </ul>
      </div>

      <!-- Contact column (unchanged) -->
      <div class="footer-column">
        ...
      </div>
    </div>
    <!-- bottom bar unchanged -->
  </footer>
</template>
```

- [ ] **Step 3: Type check**

```bash
cd frontend && npx nuxi typecheck
```
Expected: no errors.

- [ ] **Step 4: Commit**

```bash
git add frontend/components/global/Footer.vue
git commit -m "feat: make footer navigation dynamic via API"
```

---

### Task 8: Header.vue - Explicit Position Parameter

**Files:**
- Modify: `frontend/components/global/Header.vue`

- [ ] **Step 1: Update fetchNav to request header position**

In `frontend/components/global/Header.vue`, in `fetchNav` function (line 114), change:

```ts
const data = await api<NavItem[]>('/navigation');
```

To:

```ts
const data = await api<NavItem[]>('/navigation?position=header');
```

- [ ] **Step 2: Type check**

```bash
cd frontend && npx nuxi typecheck
```
Expected: no errors.

- [ ] **Step 3: Commit**

```bash
git add frontend/components/global/Header.vue
git commit -m "feat: add explicit position=header to Header nav API call"
```

---

### Task 9: Integration Verification

- [ ] **Step 1: Start backend and verify API**

```bash
cd backend && go run ./cmd/server &
sleep 3
curl -s http://localhost:8080/api/v1/navigation | head -c 200
curl -s "http://localhost:8080/api/v1/navigation?position=footer" | head -c 200
```

Expected: both return JSON arrays with `{code:0, data:[...]}`. The first returns header nav, the second returns empty (no footer items set yet).

- [ ] **Step 2: Set some nav items to footer via DB**

```bash
docker exec -i mygo-mysql mysql -uroot -proot mygo -e "UPDATE navigations SET display_position='both' WHERE id IN (1,2,3,4,5);"
```

- [ ] **Step 3: Re-verify footer API returns data**

```bash
curl -s "http://localhost:8080/api/v1/navigation?position=footer" | python3 -m json.tool 2>/dev/null || curl -s "http://localhost:8080/api/v1/navigation?position=footer"
```

Expected: returns tree with items that have display_position='both' or 'footer'.

- [ ] **Step 4: Start frontend dev server**

```bash
cd frontend && npm run dev &
```

- [ ] **Step 5: Visual check**

Open `http://localhost:3000`, verify:
- Header navigation still works correctly
- Footer shows dynamic columns from API
- Footer fallback works when API is unavailable (stop backend to test)

- [ ] **Step 6: Reset test data**

```bash
docker exec -i mygo-mysql mysql -uroot -proot mygo -e "UPDATE navigations SET display_position='header' WHERE id IN (1,2,3,4,5);"
```

- [ ] **Step 7: Commit any remaining changes**

```bash
git status
git diff
```
