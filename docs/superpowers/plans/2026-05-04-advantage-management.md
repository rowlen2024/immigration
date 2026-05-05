# 优势管理增强 — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add section title/subtitle to advantage management, SVG icon picker (100 Lucide icons), and unify three admin card UI styles (hero slides, project showcase, advantage).

**Architecture:** Backend adds `IconType` field to `AdvantageItem` and new `AdvantageSectionConfig` struct; `Get()` loads `advantage_section` from DB. Frontend adds `lucideIcons.ts` data file and `IconPicker.vue` component; admin `homepage.vue` gains section fields + icon picker; `index.vue` reads title/subtitle from API and renders SVG icons.

**Tech Stack:** Go (Gin + GORM), Nuxt 3 (Vue 3 + Element Plus), MySQL

---

### Task 1: Backend — Add data structures and update Get()

**Files:**
- Modify: `backend/internal/service/home_config_svc.go`

- [ ] **Step 1: Add IconType to AdvantageItem**

At line 28 (inside `AdvantageItem` struct), add `IconType` field:

```go
type AdvantageItem struct {
	Icon        string `json:"icon"`
	IconType    string `json:"icon_type"` // "lucide" for svg icons, empty for legacy emoji
	Title       string `json:"title"`
	Description string `json:"description"`
}
```

- [ ] **Step 2: Add AdvantageSectionConfig struct**

After `AdvantageItem` (after line 31), insert:

```go
// AdvantageSectionConfig holds the advantage section title and subtitle.
type AdvantageSectionConfig struct {
	SectionTitle    string `json:"section_title"`
	SectionSubtitle string `json:"section_subtitle"`
}
```

- [ ] **Step 3: Add AdvantageSection to HomeConfigData**

In `HomeConfigData` struct (around line 43), add before `ProjectShowcase`:

```go
type HomeConfigData struct {
	HeroSlides        []HeroSlide              `json:"hero_slides"`
	AdvantageItems    []AdvantageItem          `json:"advantage_items"`
	AdvantageSection  *AdvantageSectionConfig  `json:"advantage_section"`
	ProjectShowcase   *ProjectShowcaseConfig   `json:"project_showcase"`
}
```

- [ ] **Step 4: Update Get() to load advantage_section**

In `Get()` method (after the `advantage_items` block, before `project_showcase`, around line 57), add:

```go
	if advSecCfg, err := s.repo.FindByKey("advantage_section"); err == nil {
		var asc AdvantageSectionConfig
		if err := json.Unmarshal(advSecCfg.ConfigValue, &asc); err == nil {
			data.AdvantageSection = &asc
		}
	}
```

- [ ] **Step 5: Run backend tests to verify**

Run: `cd backend && go test ./internal/service/... -v`
Expected: All existing tests pass.

- [ ] **Step 6: Commit**

```bash
git add backend/internal/service/home_config_svc.go
git commit -m "feat: add AdvantageSection config and IconType field to AdvantageItem"
```

---

### Task 2: Frontend — Create Lucide icon data file

**Files:**
- Create: `frontend/composables/lucideIcons.ts`

- [ ] **Step 1: Create the file with 100 icons**

Write `frontend/composables/lucideIcons.ts`:

```typescript
export interface IconDef {
  name: string
  category: string
  content: string
}

export const iconCategories = ['全部', '商务', '金融', '法律', '教育', '通信', '通用']

export const lucideIcons: IconDef[] = [
  // === 商务 (Business) ===
  { name: 'briefcase', category: '商务', content: '<rect width="20" height="14" x="2" y="7" rx="2" ry="2"/><path d="M16 21V5a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v16"/>' },
  { name: 'building', category: '商务', content: '<rect width="16" height="20" x="4" y="2" rx="2" ry="2"/><path d="M9 22v-4h6v4M8 6h.01M16 6h.01M8 10h.01M16 10h.01M8 14h.01M16 14h.01M12 10h.01M12 14h.01"/>' },
  { name: 'building-2', category: '商务', content: '<path d="M6 22V2h12v20M6 12H4a2 2 0 0 0-2 2v6a2 2 0 0 0 2 2M18 9h2a2 2 0 0 1 2 2v9a2 2 0 0 1-2 2"/><path d="M10 6h4M10 10h4M10 14h4M10 18h4"/>' },
  { name: 'users', category: '商务', content: '<path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M22 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75"/>' },
  { name: 'user-plus', category: '商务', content: '<path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><line x1="19" x2="19" y1="8" y2="14"/><line x1="22" x2="16" y1="11" y2="11"/>' },
  { name: 'user-check', category: '商务', content: '<path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><polyline points="16 11 18 13 22 9"/>' },
  { name: 'target', category: '商务', content: '<circle cx="12" cy="12" r="10"/><circle cx="12" cy="12" r="6"/><circle cx="12" cy="12" r="2"/>' },
  { name: 'trending-up', category: '商务', content: '<polyline points="22 7 13.5 15.5 8.5 10.5 2 17"/><polyline points="16 7 22 7 22 13"/>' },
  { name: 'bar-chart', category: '商务', content: '<line x1="12" x2="12" y1="20" y2="10"/><line x1="18" x2="18" y1="20" y2="4"/><line x1="6" x2="6" y1="20" y2="16"/>' },
  { name: 'pie-chart', category: '商务', content: '<path d="M21.21 15.89A10 10 0 1 1 8 2.83"/><path d="M22 12A10 10 0 0 0 12 2v10z"/>' },
  { name: 'presentation', category: '商务', content: '<path d="M2 3h20"/><path d="M21 3v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V3"/><path d="m7 21 5-5 5 5"/>' },
  { name: 'clipboard', category: '商务', content: '<rect width="16" height="20" x="4" y="2" rx="2" ry="2"/><path d="M9 2v4h6V2"/>' },
  { name: 'file-text', category: '商务', content: '<path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z"/><polyline points="14 2 14 8 20 8"/><line x1="16" x2="8" y1="13" y2="13"/><line x1="16" x2="8" y1="17" y2="17"/><line x1="10" x2="8" y1="9" y2="9"/>' },
  { name: 'folder', category: '商务', content: '<path d="M4 20h16a2 2 0 0 0 2-2V8a2 2 0 0 0-2-2h-7.93a2 2 0 0 1-1.66-.9l-.82-1.2A2 2 0 0 0 7.93 3H4a2 2 0 0 0-2 2v13c0 1.1.9 2 2 2z"/>' },
  { name: 'mail', category: '商务', content: '<rect width="20" height="16" x="2" y="4" rx="2"/><path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7"/>' },
  { name: 'phone', category: '商务', content: '<path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"/>' },
  { name: 'calendar', category: '商务', content: '<rect width="18" height="18" x="3" y="4" rx="2" ry="2"/><line x1="16" x2="16" y1="2" y2="6"/><line x1="8" x2="8" y1="2" y2="6"/><line x1="3" x2="21" y1="10" y2="10"/>' },
  { name: 'clock', category: '商务', content: '<circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/>' },
  { name: 'map-pin', category: '商务', content: '<path d="M20 10c0 6-8 12-8 12s-8-6-8-12a8 8 0 0 1 16 0z"/><circle cx="12" cy="10" r="3"/>' },
  { name: 'globe', category: '商务', content: '<circle cx="12" cy="12" r="10"/><path d="M12 2a14.5 14.5 0 0 0 0 20 14.5 14.5 0 0 0 0-20M2 12h20"/>' },
  { name: 'award', category: '商务', content: '<circle cx="12" cy="8" r="6"/><path d="M15.477 12.89 17 22l-5-3-5 3 1.523-9.11"/>' },
  { name: 'handshake', category: '商务', content: '<path d="m11 17 4-4M14 17l4-4M9 10l-2 2"/><path d="M2 14.5V12a4 4 0 0 1 4-4h1.5M22 14.5V12a4 4 0 0 0-4-4h-1.5"/><path d="M20 18H4a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2z"/>' },

  // === 金融 (Finance) ===
  { name: 'banknote', category: '金融', content: '<rect width="20" height="12" x="2" y="6" rx="2"/><circle cx="12" cy="12" r="2"/><path d="M6 12h.01M18 12h.01"/>' },
  { name: 'dollar-sign', category: '金融', content: '<line x1="12" x2="12" y1="2" y2="22"/><path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"/>' },
  { name: 'credit-card', category: '金融', content: '<rect width="20" height="14" x="2" y="5" rx="2"/><line x1="2" x2="22" y1="10" y2="10"/>' },
  { name: 'wallet', category: '金融', content: '<path d="M19 7V4a1 1 0 0 0-1-1H5a2 2 0 0 0 0 4h15a1 1 0 0 1 1 1v4h-3a2 2 0 0 0 0 4h3a1 1 0 0 0 1-1v-2a1 1 0 0 0-1-1"/><path d="M3 5v14a2 2 0 0 0 2 2h15a1 1 0 0 0 1-1v-4"/>' },
  { name: 'landmark', category: '金融', content: '<line x1="3" x2="21" y1="22" y2="22"/><line x1="6" x2="6" y1="18" y2="11"/><line x1="10" x2="10" y1="18" y2="11"/><line x1="14" x2="14" y1="18" y2="11"/><line x1="18" x2="18" y1="18" y2="11"/><polygon points="12 2 20 7 4 7"/>' },
  { name: 'piggy-bank', category: '金融', content: '<path d="M19 5c-1.5 0-2.8 1.1-3.5 2.5H7.5C6.7 6.1 5.4 5 3.9 5 2.6 5 1.5 6 1.5 7.2c0 .5.2 1 .7 1.4-.3.5-.5 1-.5 1.5 0 1.2.9 2.3 2 2.8V19h2v-2h2v2h2v-2h2v2h2v-6.1c1.1-.5 2-1.6 2-2.8 0-.5-.2-1-.5-1.5.5-.4.7-.9.7-1.4 0-1.2-1.1-2.2-2.5-2.2z"/><path d="M12 11h.01"/>' },
  { name: 'receipt', category: '金融', content: '<path d="M4 2v20l2-1 2 1 2-1 2 1 2-1 2 1 2-1 2 1V2l-2 1-2-1-2 1-2-1-2 1-2-1-2 1-2-1Z"/><path d="M8 7h8M8 11h8M8 15h5"/>' },
  { name: 'calculator', category: '金融', content: '<rect width="16" height="20" x="4" y="2" rx="2"/><line x1="8" x2="16" y1="6" y2="6"/><line x1="16" x2="16" y1="14" y2="18"/><path d="M16 10h.01M12 10h.01M8 10h.01M12 14h.01M8 14h.01M12 18h.01M8 18h.01"/>' },
  { name: 'scale', category: '金融', content: '<path d="m16 16 3-8 3 8c-.87.65-1.92 1-3 1s-2.13-.35-3-1ZM2 16l3-8 3 8c-.87.65-1.92 1-3 1s-2.13-.35-3-1Z"/><path d="M12 22V8M6 22V2h12v20"/>' },
  { name: 'gem', category: '金融', content: '<polygon points="6 3 12 9 18 3"/><polygon points="4 8 12 22 20 8"/>' },
  { name: 'coins', category: '金融', content: '<circle cx="8" cy="8" r="6"/><path d="M18.09 10.37A6 6 0 1 1 10.34 18"/><path d="M7 6h1v4"/><path d="m16.71 13.88.71.71-2.83 2.83"/>' },
  { name: 'arrow-up-down', category: '金融', content: '<path d="m21 16-4 4-4-4"/><path d="M17 20V4"/><path d="m3 8 4-4 4 4"/><path d="M7 4v16"/>' },
  { name: 'percent', category: '金融', content: '<line x1="19" x2="5" y1="5" y2="19"/><circle cx="6.5" cy="6.5" r="2.5"/><circle cx="17.5" cy="17.5" r="2.5"/>' },
  { name: 'chart-line', category: '金融', content: '<path d="M3 3v18h18"/><path d="m19 9-5 5-4-4-3 3"/>' },

  // === 法律 (Legal) ===
  { name: 'shield', category: '法律', content: '<path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>' },
  { name: 'shield-check', category: '法律', content: '<path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/><path d="m9 12 2 2 4-4"/>' },
  { name: 'shield-alert', category: '法律', content: '<path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/><path d="M12 8v4M12 16h.01"/>' },
  { name: 'gavel', category: '法律', content: '<path d="m14.5 12.5-8 8a2.119 2.119 0 1 1-3-3l8-8"/><path d="m16 16 6-6"/><path d="m8 8 6-6"/><path d="m9 7 8 8"/><path d="m21 11-8-8"/>' },
  { name: 'file-signature', category: '法律', content: '<path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z"/><polyline points="14 2 14 8 20 8"/><path d="m10 18 2-2-1-1-2 2v1h1z"/><path d="M12 12 9 9"/>' },
  { name: 'stamp', category: '法律', content: '<path d="M5 21h14"/><path d="M6 3h12l-2 8H8l-2-8z"/><path d="M12 11v10"/>' },
  { name: 'book-open', category: '法律', content: '<path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"/><path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"/>' },
  { name: 'book-marked', category: '法律', content: '<path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1 0-5H20"/><polyline points="10 2 10 10 13 7 16 10 16 2"/>' },
  { name: 'scroll-text', category: '法律', content: '<path d="M8 21h12a2 2 0 0 0 2-2v-2H10v2a2 2 0 1 1-4 0V5a2 2 0 1 0-4 0v3h4"/><path d="M19 17V5a2 2 0 0 0-2-2H4"/><path d="M15 8h-5M15 12h-5"/>' },
  { name: 'scale-3d', category: '法律', content: '<path d="M5 21h14"/><path d="M12 3v18"/><path d="m5 9 7-4 7 4M5 15l7 4 7-4"/>' },

  // === 教育 (Education) ===
  { name: 'graduation-cap', category: '教育', content: '<path d="M22 10v6M2 10l10-5 10 5-10 5z"/><path d="M6 12v5c0 2 3 3 6 3s6-1 6-3v-5"/>' },
  { name: 'book', category: '教育', content: '<path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1 0-5H20"/>' },
  { name: 'library', category: '教育', content: '<path d="m16 6 4 14"/><path d="M12 6v14"/><path d="M8 8v12"/><path d="M4 4v16"/>' },
  { name: 'pencil', category: '教育', content: '<path d="M17 3a2.85 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z"/><path d="m15 5 4 4"/>' },
  { name: 'lightbulb', category: '教育', content: '<path d="M15 14c.2-1 .7-1.7 1.5-2.5 1-.9 1.5-2.2 1.5-3.5A6 6 0 0 0 6 8c0 1 .2 2.2 1.5 3.5.7.7 1.3 1.5 1.5 2.5"/><path d="M9 18h6"/><path d="M10 22h4"/>' },
  { name: 'brain', category: '教育', content: '<path d="M9.5 2A2.5 2.5 0 0 1 12 4.5v15a2.5 2.5 0 0 1-4.96.44 2.5 2.5 0 0 1-2.96-3.08 3 3 0 0 1-.34-5.58 2.5 2.5 0 0 1 1.32-4.24 2.5 2.5 0 0 1 1.98-3A2.5 2.5 0 0 1 9.5 2Z"/><path d="M14.5 2A2.5 2.5 0 0 0 12 4.5v15a2.5 2.5 0 0 0 4.96.44 2.5 2.5 0 0 0 2.96-3.08 3 3 0 0 0 .34-5.58 2.5 2.5 0 0 0-1.32-4.24 2.5 2.5 0 0 0-1.98-3A2.5 2.5 0 0 0 14.5 2Z"/>' },
  { name: 'microscope', category: '教育', content: '<path d="M6 22h12"/><path d="M12 22v-6"/><path d="M12 2v11"/><path d="M19 2H5"/><circle cx="12" cy="8" r="3"/>' },
  { name: 'flask', category: '教育', content: '<path d="M9 3h6M10 3v4l-4.5 8.5A2 2 0 0 0 7.2 18h9.6a2 2 0 0 0 1.7-2.5L14 7V3"/>' },
  { name: 'school', category: '教育', content: '<path d="m4 6 8-4 8 4"/><path d="m18 10 4 2v8a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2v-8l4-2"/><path d="M14 22v-4a2 2 0 0 0-2-2v0a2 2 0 0 0-2 2v4"/><path d="M18 5v17M6 5v17"/><circle cx="12" cy="9" r="2"/>' },
  { name: 'compass', category: '教育', content: '<circle cx="12" cy="12" r="10"/><polygon points="16.24 7.76 14.12 14.12 7.76 16.24 9.88 9.88 16.24 7.76"/>' },
  { name: 'palette', category: '教育', content: '<circle cx="13.5" cy="6.5" r=".5"/><circle cx="17.5" cy="10.5" r=".5"/><circle cx="8.5" cy="7.5" r=".5"/><circle cx="6.5" cy="12.5" r=".5"/><path d="M12 2C6.5 2 2 6.5 2 12s4.5 10 10 10c.926 0 1.648-.746 1.648-1.688 0-.437-.18-.835-.437-1.125-.29-.289-.438-.652-.438-1.125a1.64 1.64 0 0 1 1.668-1.668h1.996c3.051 0 5.555-2.503 5.555-5.554C21.965 6.012 17.461 2 12 2z"/>' },

  // === 通信 (Communication) ===
  { name: 'message-circle', category: '通信', content: '<path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"/>' },
  { name: 'messages-square', category: '通信', content: '<path d="M14 9a2 2 0 0 1-2 2H6l-4 4V4c0-1.1.9-2 2-2h8a2 2 0 0 1 2 2v5z"/><path d="M18 9h2a2 2 0 0 1 2 2v11l-4-4h-6a2 2 0 0 1-2-2v-1"/>' },
  { name: 'send', category: '通信', content: '<path d="m22 2-7 20-4-9-9-4Z"/><path d="M22 2 11 13"/>' },
  { name: 'thumbs-up', category: '通信', content: '<path d="M7 11v8a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1v-7a1 1 0 0 1 1-1h3a4 4 0 0 0 4-4V6a2 2 0 0 1 4 0v5h3a2 2 0 0 1 2 2l-1 5a2 3 0 0 1-2 2h-7a3 3 0 0 1-3-3"/>' },
  { name: 'star', category: '通信', content: '<polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/>' },
  { name: 'heart', category: '通信', content: '<path d="M20.42 4.58a5.4 5.4 0 0 0-7.65 0l-.77.78-.77-.78a5.4 5.4 0 0 0-7.65 0C1.46 6.7 1.33 10.28 4 13l8 8 8-8c2.67-2.72 2.54-6.3.42-8.42z"/>' },
  { name: 'bell', category: '通信', content: '<path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"/><path d="M13.73 21a2 2 0 0 1-3.46 0"/>' },
  { name: 'at-sign', category: '通信', content: '<circle cx="12" cy="12" r="4"/><path d="M16 8v5a3 3 0 0 0 6 0v-1a10 10 0 1 0-3.92 7.94"/>' },
  { name: 'link', category: '通信', content: '<path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"/><path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"/>' },
  { name: 'rss', category: '通信', content: '<path d="M4 11a9 9 0 0 1 9 9"/><path d="M4 4a16 16 0 0 1 16 16"/><circle cx="5" cy="19" r="1"/>' },
  { name: 'share-2', category: '通信', content: '<circle cx="18" cy="5" r="3"/><circle cx="6" cy="12" r="3"/><circle cx="18" cy="19" r="3"/><line x1="8.59" x2="15.42" y1="13.51" y2="17.49"/><line x1="15.41" x2="8.59" y1="6.51" y2="10.49"/>' },

  // === 通用 (General) ===
  { name: 'search', category: '通用', content: '<circle cx="11" cy="11" r="8"/><path d="m21 21-4.3-4.3"/>' },
  { name: 'check', category: '通用', content: '<polyline points="20 6 9 17 4 12"/>' },
  { name: 'x', category: '通用', content: '<path d="M18 6 6 18"/><path d="m6 6 12 12"/>' },
  { name: 'plus', category: '通用', content: '<path d="M5 12h14"/><path d="M12 5v14"/>' },
  { name: 'minus', category: '通用', content: '<path d="M5 12h14"/>' },
  { name: 'chevron-right', category: '通用', content: '<path d="m9 18 6-6-6-6"/>' },
  { name: 'chevron-left', category: '通用', content: '<path d="m15 18-6-6 6-6"/>' },
  { name: 'arrow-right', category: '通用', content: '<path d="M5 12h14"/><path d="m12 5 7 7-7 7"/>' },
  { name: 'menu', category: '通用', content: '<line x1="4" x2="20" y1="12" y2="12"/><line x1="4" x2="20" y1="6" y2="6"/><line x1="4" x2="20" y1="18" y2="18"/>' },
  { name: 'home', category: '通用', content: '<path d="m3 9 9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/><polyline points="9 22 9 12 15 12 15 22"/>' },
  { name: 'info', category: '通用', content: '<circle cx="12" cy="12" r="10"/><path d="M12 16v-4M12 8h.01"/>' },
  { name: 'help-circle', category: '通用', content: '<circle cx="12" cy="12" r="10"/><path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"/><path d="M12 17h.01"/>' },
  { name: 'settings', category: '通用', content: '<path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z"/><circle cx="12" cy="12" r="3"/>' },
  { name: 'edit', category: '通用', content: '<path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>' },
  { name: 'trash-2', category: '通用', content: '<path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/><line x1="10" x2="10" y1="11" y2="17"/><line x1="14" x2="14" y1="11" y2="17"/>' },
  { name: 'copy', category: '通用', content: '<rect width="14" height="14" x="8" y="8" rx="2" ry="2"/><path d="M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2"/>' },
  { name: 'download', category: '通用', content: '<path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" x2="12" y1="15" y2="3"/>' },
  { name: 'upload', category: '通用', content: '<path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" x2="12" y1="3" y2="15"/>' },
  { name: 'external-link', category: '通用', content: '<path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/><polyline points="15 3 21 3 21 9"/><line x1="10" x2="21" y1="14" y2="3"/>' },
  { name: 'eye', category: '通用', content: '<path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/>' },
  { name: 'eye-off', category: '通用', content: '<path d="M9.88 9.88a3 3 0 1 0 4.24 4.24"/><path d="M10.73 5.08A10.43 10.43 0 0 1 12 5c7 0 10 7 10 7a13.16 13.16 0 0 1-1.67 2.68"/><path d="M6.61 6.61A13.526 13.526 0 0 0 2 12s3 7 10 7a9.74 9.74 0 0 0 5.39-1.61"/><line x1="2" x2="22" y1="2" y2="22"/>' },
  { name: 'lock', category: '通用', content: '<rect width="18" height="11" x="3" y="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/>' },
  { name: 'unlock', category: '通用', content: '<rect width="18" height="11" x="3" y="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 9.9-1"/>' },
  { name: 'key', category: '通用', content: '<circle cx="7.5" cy="15.5" r="5.5"/><path d="m21 2-9.6 9.6"/><path d="m15.5 7.5 3 3L22 7l-3-3"/>' },
  { name: 'filter', category: '通用', content: '<polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3"/>' },
  { name: 'zap', category: '通用', content: '<polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/>' },
  { name: 'flag', category: '通用', content: '<path d="M4 15s1-1 4-1 5 2 8 2 4-1 4-1V3s-1 1-4 1-5-2-8-2-4 1-4 1z"/><line x1="4" x2="4" y1="22" y2="15"/>' },
  { name: 'tag', category: '通用', content: '<path d="M12 2H2v10l9.17 9.17a2 2 0 0 0 2.83 0l7-7a2 2 0 0 0 0-2.83L12 2Z"/><path d="M7 7h.01"/>' },

  // === 更多实用图标 ===
  { name: 'route', category: '商务', content: '<circle cx="6" cy="19" r="3"/><path d="M9 19h8.5a3.5 3.5 0 0 0 0-7h-11a3.5 3.5 0 0 1 0-7H15"/><circle cx="18" cy="5" r="3"/>' },
  { name: 'plane', category: '商务', content: '<path d="M17.8 19.2 16 11l3.5-3.5C21 6 21.5 4 21 3c-1-.5-3 0-4.5 1.5L13 8 4.8 6.2c-.5-.1-.9.1-1.1.5l-.3.5c-.2.5-.1 1 .3 1.3L9 12l-2 3H4l-1 1 3 2 2 3 1-1v-3l3-2 3.5 5.3c.3.4.8.5 1.3.3l.5-.2c.4-.3.6-.7.5-1.2z"/>' },
  { name: 'rocket', category: '商务', content: '<path d="M4.5 16.5c-1.5 1.26-2 5-2 5s3.74-.5 5-2c.71-.84.7-2.13-.09-2.91a2.18 2.18 0 0 0-2.91-.09z"/><path d="m12 15-3-3a22 22 0 0 1 2-3.95A12.88 12.88 0 0 1 22 2c0 2.72-.78 7.5-6 11a22.35 22.35 0 0 1-4 2z"/><path d="M9 12H4s.55-3.03 2-4c1.62-1.08 5 0 5 0"/><path d="M12 15v5s3.03-.55 4-2c1.08-1.62 0-5 0-5"/>' },
  { name: 'smile', category: '通信', content: '<circle cx="12" cy="12" r="10"/><path d="M8 14s1.5 2 4 2 4-2 4-2"/><line x1="9" x2="9.01" y1="9" y2="9"/><line x1="15" x2="15.01" y1="9" y2="9"/>' },
  { name: 'badge-check', category: '商务', content: '<path d="M3.85 8.62a4 4 0 0 1 4.78-4.77 4 4 0 0 1 6.74 0 4 4 0 0 1 4.78 4.78 4 4 0 0 1 0 6.74 4 4 0 0 1-4.77 4.78 4 4 0 0 1-6.75 0 4 4 0 0 1-4.78-4.77 4 4 0 0 1 0-6.76Z"/><path d="m9 12 2 2 4-4"/>' },
  { name: 'umbrella', category: '通用', content: '<path d="M22 12a10.06 10.06 1 0 0-20 0Z"/><path d="M12 12v8a2 2 0 0 0 4 0"/><path d="M12 2v1"/>' },
  { name: 'crown', category: '商务', content: '<path d="m2 4 3 12h14l3-12-6 7-4-7-4 7-6-7z"/><path d="M3 20h18"/>' },
  { name: 'trophy', category: '商务', content: '<path d="M6 9H4.5a2.5 2.5 0 0 1 0-5H6"/><path d="M18 9h1.5a2.5 2.5 0 0 0 0-5H18"/><path d="M4 22h16"/><path d="M10 22v-3.34a6 6 0 0 1-5-5.9V9h14v3.76a6 6 0 0 1-5 5.9V22h-4z"/>' },
  { name: 'palmtree', category: '通用', content: '<path d="M13 8c0-2.76-2.24-5-5-5S3 5.24 3 8h2l1-4 1 4h2"/><path d="M11 13c0-2.76-2.24-5-5-5s-5 2.24-5 5h2l1-4 1 4h2"/><path d="M9 18c0-2.76-2.24-5-5-5s-5 2.24-5 5h2l1-4 1 4h2"/><path d="M22 22H2"/>' },
]

export function getIconByName(name: string): IconDef | undefined {
  return lucideIcons.find((i) => i.name === name)
}

export function getIconsByCategory(category: string): IconDef[] {
  if (category === '全部') return lucideIcons
  return lucideIcons.filter((i) => i.category === category)
}

export function searchIcons(query: string): IconDef[] {
  const q = query.toLowerCase()
  return lucideIcons.filter((i) => i.name.includes(q))
}
```

- [ ] **Step 2: Verify file compiles**

Run: `cd frontend && npx nuxi typecheck`
Expected: No new type errors.

- [ ] **Step 3: Commit**

```bash
git add frontend/composables/lucideIcons.ts
git commit -m "feat: add 100 Lucide SVG icons data file with category support"
```

---

### Task 3: Frontend — Create IconPicker component

**Files:**
- Create: `frontend/components/admin/IconPicker.vue`

- [ ] **Step 1: Write IconPicker.vue**

```vue
<template>
  <div class="icon-picker">
    <div class="icon-picker__trigger" @click="dialogVisible = true">
      <div v-if="modelValue" class="icon-picker__preview">
        <svg
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="icon-picker__svg"
          v-html="selectedIcon?.content || ''"
        />
      </div>
      <span v-else class="icon-picker__placeholder">选择图标</span>
    </div>

    <el-dialog
      v-model="dialogVisible"
      title="选择图标"
      width="620px"
      destroy-on-close
    >
      <div class="icon-picker__dialog">
        <div class="icon-picker__tabs">
          <span
            v-for="cat in iconCategories"
            :key="cat"
            class="icon-picker__tab"
            :class="{ active: activeCategory === cat }"
            @click="activeCategory = cat"
          >{{ cat }}</span>
        </div>
        <el-input
          v-model="searchQuery"
          placeholder="搜索图标..."
          clearable
          class="icon-picker__search"
        />
        <div class="icon-picker__grid">
          <div
            v-for="icon in filteredIcons"
            :key="icon.name"
            class="icon-picker__item"
            :class="{ selected: modelValue === icon.name }"
            :title="icon.name"
            @click="select(icon.name)"
          >
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
              class="icon-picker__item-svg"
              v-html="icon.content"
            />
          </div>
        </div>
        <div v-if="modelValue" class="icon-picker__selected-name">
          已选: <strong>{{ modelValue }}</strong>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { lucideIcons, iconCategories, searchIcons, getIconsByCategory, getIconByName } from '~/composables/lucideIcons'

const props = defineProps<{ modelValue: string }>()
const emit = defineEmits<{ 'update:modelValue': [value: string] }>()

const dialogVisible = ref(false)
const activeCategory = ref('全部')
const searchQuery = ref('')

const selectedIcon = computed(() => getIconByName(props.modelValue))

const filteredIcons = computed(() => {
  if (searchQuery.value) return searchIcons(searchQuery.value)
  return getIconsByCategory(activeCategory.value)
})

function select(name: string) {
  emit('update:modelValue', name)
  dialogVisible.value = false
}
</script>

<style scoped>
.icon-picker__trigger {
  width: 44px;
  height: 44px;
  border: 2px dashed #dcdfe6;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: border-color 0.2s;
}
.icon-picker__trigger:hover {
  border-color: #c8963e;
}
.icon-picker__preview {
  width: 36px;
  height: 36px;
  background: rgba(26, 58, 92, 0.08);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}
.icon-picker__svg {
  width: 20px;
  height: 20px;
  color: #c8963e;
}
.icon-picker__placeholder {
  font-size: 11px;
  color: #c0c4cc;
}
.icon-picker__dialog {
  max-height: 460px;
  overflow-y: auto;
}
.icon-picker__tabs {
  display: flex;
  gap: 4px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}
.icon-picker__tab {
  padding: 4px 12px;
  font-size: 12px;
  border-radius: 4px;
  cursor: pointer;
  color: #606266;
  background: #f5f7fa;
  transition: all 0.2s;
}
.icon-picker__tab.active {
  background: #1a3a5c;
  color: #fff;
}
.icon-picker__tab:hover:not(.active) {
  background: #e8eaed;
}
.icon-picker__search {
  margin-bottom: 12px;
}
.icon-picker__grid {
  display: grid;
  grid-template-columns: repeat(8, 1fr);
  gap: 6px;
}
.icon-picker__item {
  aspect-ratio: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  cursor: pointer;
  border: 2px solid transparent;
  transition: all 0.15s;
}
.icon-picker__item:hover {
  background: rgba(26, 58, 92, 0.06);
}
.icon-picker__item.selected {
  border-color: #c8963e;
  background: rgba(200, 150, 62, 0.1);
}
.icon-picker__item-svg {
  width: 22px;
  height: 22px;
  color: #606266;
}
.icon-picker__item.selected .icon-picker__item-svg {
  color: #c8963e;
}
.icon-picker__selected-name {
  margin-top: 12px;
  padding-top: 10px;
  border-top: 1px solid #ebeef5;
  font-size: 13px;
  color: #909399;
}
</style>
```

- [ ] **Step 2: Verify typecheck**

Run: `cd frontend && npx nuxi typecheck`
Expected: No new type errors.

- [ ] **Step 3: Commit**

```bash
git add frontend/components/admin/IconPicker.vue
git commit -m "feat: add IconPicker component with Lucide SVG grid, category tabs, and search"
```

---

### Task 4: Frontend — Update admin homepage.vue

**Files:**
- Modify: `frontend/pages/admin/homepage.vue`

This task has three parts: (A) unified hero slides card, (B) advantage section title/subtitle + icon picker, (C) unified styles.

- [ ] **Step 1: Move "新增 Slide" button from card header to content area**

Replace the hero slides card template (lines 10-29) — move button from `#header` slot to after the list:

Change from:
```vue
<el-card class="config-card">
  <template #header>
    <div class="card-header">
      <h3 class="card-title">轮播管理</h3>
      <el-button type="primary" size="small" @click="openAddSlide">新增 Slide</el-button>
    </div>
  </template>
  <div v-if="heroSlides.length === 0" ...>
  ...
</el-card>
```

To:
```vue
<el-card class="config-card">
  <template #header>
    <h3 class="card-title">轮播管理</h3>
  </template>
  <div v-if="heroSlides.length === 0" class="empty-hint">暂无轮播，点击下方按钮添加。</div>
  <div v-else class="config-list">
    <div v-for="(slide, i) in heroSlides" :key="i" class="config-item">
      <img v-if="slide.image" :src="slide.image" class="slide-thumb" />
      <span v-else class="slide-label">(无图片)</span>
      <div class="config-item-actions">
        <el-button size="small" :disabled="i === 0" @click="moveSlide(i, -1)">↑</el-button>
        <el-button size="small" :disabled="i === heroSlides.length - 1" @click="moveSlide(i, 1)">↓</el-button>
        <el-button size="small" @click="openEditSlide(i)">编辑</el-button>
        <el-button size="small" type="danger" @click="removeSlide(i)">删除</el-button>
      </div>
    </div>
  </div>
  <div class="config-list-actions">
    <el-button type="primary" size="small" @click="openAddSlide">新增 Slide</el-button>
  </div>
  <div class="card-footer">
    <el-button type="primary" :loading="slideSaving" @click="saveSlides">保存轮播</el-button>
  </div>
</el-card>
```

- [ ] **Step 2: Add section title/subtitle fields to advantage card**

Replace the advantage card template (lines 78-101) — add section fields before the list:

```vue
<el-card class="config-card">
  <template #header>
    <h3 class="card-title">优势管理</h3>
  </template>
  <el-form label-width="100px" class="section-form">
    <el-form-item label="区域标题">
      <el-input v-model="advantageSection.section_title" placeholder="为什么选择 MyGo移民？" />
    </el-form-item>
    <el-form-item label="区域副标题">
      <el-input v-model="advantageSection.section_subtitle" placeholder="专业服务，值得信赖" />
    </el-form-item>
  </el-form>
  <div v-if="advantageItems.length === 0" class="empty-hint">暂无优势项，点击下方按钮添加。</div>
  <div v-else class="config-list">
    <div v-for="(item, i) in advantageItems" :key="i" class="config-item">
      <div class="adv-icon-preview">
        <svg
          v-if="item.icon_type === 'lucide'"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="adv-icon-svg"
          v-html="getIconContent(item.icon)"
        />
        <span v-else class="adv-icon-emoji">{{ item.icon }}</span>
      </div>
      <div class="config-item-info">
        <strong>{{ item.title }}</strong>
        <span class="config-item-desc">{{ item.description }}</span>
      </div>
      <div class="config-item-actions">
        <el-button size="small" :disabled="i === 0" @click="moveAdv(i, -1)">↑</el-button>
        <el-button size="small" :disabled="i === advantageItems.length - 1" @click="moveAdv(i, 1)">↓</el-button>
        <el-button size="small" @click="openEditAdv(i)">编辑</el-button>
        <el-button size="small" type="danger" @click="removeAdv(i)">删除</el-button>
      </div>
    </div>
  </div>
  <div class="config-list-actions">
    <el-button type="primary" size="small" @click="openAddAdv">新增优势项</el-button>
  </div>
  <div class="card-footer">
    <el-button type="primary" :loading="advSaving" @click="saveAdvantages">保存优势设置</el-button>
  </div>
</el-card>
```

- [ ] **Step 3: Update project showcase card header to match**

Replace the project showcase card header (line 33):

```vue
<template #header><h3 class="card-title">项目展示区</h3></template>
```

And add `config-list` / `config-item` classes to the featured list (lines 43-69), wrapping the featured items:

```vue
<div v-if="projectShowcase.featured_slugs.length === 0" class="empty-hint">
  未选择精选项目，首页将展示全部项目。
</div>
<div v-else class="config-list">
  <div v-for="(slug, i) in projectShowcase.featured_slugs" :key="slug" class="config-item">
    <span class="config-item-name">{{ getProjectTitle(slug) }}</span>
    <div class="config-item-actions">
      <el-button size="small" :disabled="i === 0" @click="moveFeatured(i, -1)">↑</el-button>
      <el-button size="small" :disabled="i === projectShowcase.featured_slugs.length - 1" @click="moveFeatured(i, 1)">↓</el-button>
      <el-button size="small" type="danger" @click="removeFeatured(i)">移除</el-button>
    </div>
  </div>
</div>
```

- [ ] **Step 4: Update TypeScript — add advantageSection ref, import lucideIcons**

In the `<script setup>` block (after line 165):

Add import:
```typescript
import { getIconByName } from '~/composables/lucideIcons'
```

Add `advantageSection` ref (after `advantageItems` ref, around line 197):
```typescript
const advantageSection = ref<{ section_title: string; section_subtitle: string }>({
  section_title: '',
  section_subtitle: '',
});
```

Update `AdvantageItem` interface (line 178) to add `icon_type`:
```typescript
interface AdvantageItem {
  icon: string;
  icon_type: string;
  title: string;
  description: string;
}
```

Update `load()` function — add `advantage_section` to API response type and populate `advantageSection`:
```typescript
const [config, projects] = await Promise.all([
  api<{
    hero_slides: HeroSlide[];
    advantage_items: AdvantageItem[];
    advantage_section: { section_title: string; section_subtitle: string } | null;
    project_showcase: ProjectShowcase | null;
  }>('/admin/home-config'),
  api<{ items: ProjectOption[] }>('/projects'),
]);

if (config) {
  heroSlides.value = config.hero_slides || [];
  advantageItems.value = config.advantage_items || [];
  if (config.advantage_section) {
    advantageSection.value = config.advantage_section;
  }
  if (config.project_showcase) {
    projectShowcase.value = config.project_showcase;
  }
}
```

Add helper function `getIconContent`:
```typescript
function getIconContent(name: string): string {
  return getIconByName(name)?.content || ''
}
```

- [ ] **Step 5: Update advantage edit dialog — replace icon text input with IconPicker**

Replace the icon form item in the advantage dialog (lines 147-149):

```vue
<el-form-item label="图标">
  <IconPicker v-model="advForm.icon" />
</el-form-item>
```

Update `advForm` initialization to include `icon_type`:
```typescript
const advForm = ref<AdvantageItem>({ icon: '', icon_type: 'lucide', title: '', description: '' });
```

Update `openAddAdv()`:
```typescript
function openAddAdv() {
  advEditIndex.value = -1;
  advForm.value = { icon: '', icon_type: 'lucide', title: '', description: '' };
  advDialogVisible.value = true;
}
```

- [ ] **Step 6: Update saveAdvantages to include advantage_section**

```typescript
async function saveAdvantages() {
  advSaving.value = true;
  try {
    const api = useApi();
    await api('/admin/home-config', {
      method: 'PUT',
      body: {
        advantage_items: advantageItems.value,
        advantage_section: advantageSection.value,
      },
    });
    ElMessage.success('优势设置已保存');
  } catch {
    ElMessage.error('保存失败');
  } finally {
    advSaving.value = false;
  }
}
```

- [ ] **Step 7: Replace scoped styles with unified classes**

Replace the scoped `<style>` section. Remove old classes: `.slide-item`, `.slide-label`, `.slide-actions`, `.adv-item`, `.adv-icon`, `.adv-info`, `.adv-desc`, `.adv-actions`, `.featured-area`, `.featured-row`, `.featured-name`, `.featured-actions`.

Add unified classes:

```css
.page-header { margin-bottom: 24px; }
.page-title { font-size: 22px; font-weight: 600; }
.config-body { max-width: 900px; }
.config-card { margin-bottom: 20px; }
.card-title { margin: 0; font-size: 16px; font-weight: 600; }
.empty-hint {
  color: #909399;
  font-size: 14px;
  padding: 16px 0;
  text-align: center;
}

/* Unified list styles for all three cards */
.config-list {
  display: flex;
  flex-direction: column;
}
.config-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 0;
  border-bottom: 1px solid #ebeef5;
}
.config-item:last-child { border-bottom: none; }
.config-item-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}
.config-item-info strong { font-size: 14px; }
.config-item-desc { font-size: 12px; color: #909399; }
.config-item-name { font-size: 14px; flex: 1; }
.config-item-actions {
  display: flex;
  gap: 6px;
  flex-shrink: 0;
}
.config-list-actions {
  text-align: right;
  padding-top: 8px;
}

/* Slide thumbnail */
.slide-thumb { width: 120px; height: 68px; object-fit: cover; border-radius: 4px; flex-shrink: 0; }
.slide-label { font-size: 14px; color: #909399; flex: 1; }

/* Section form (title/subtitle) */
.section-form {
  padding-bottom: 8px;
  margin-bottom: 8px;
  border-bottom: 1px solid #ebeef5;
}

/* Advantage icon preview in list */
.adv-icon-preview {
  width: 36px;
  height: 36px;
  background: rgba(26, 58, 92, 0.08);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.adv-icon-svg {
  width: 18px;
  height: 18px;
  color: #c8963e;
}
.adv-icon-emoji { font-size: 18px; }

/* Card footer */
.card-footer { text-align: center; padding-top: 16px; }

/* Add project select */
.add-project-select { margin-top: 8px; width: 100%; }
```

- [ ] **Step 8: Verify typecheck**

Run: `cd frontend && npx nuxi typecheck`
Expected: No type errors.

- [ ] **Step 9: Commit**

```bash
git add frontend/pages/admin/homepage.vue
git commit -m "feat: add advantage section title/subtitle, IconPicker integration, and unify three admin card styles"
```

---

### Task 5: Frontend — Update index.vue homepage rendering

**Files:**
- Modify: `frontend/pages/index.vue`

- [ ] **Step 1: Replace hardcoded advantage title/subtitle with API data**

At lines 73-75, change:
```vue
<h2 class="section-title">为什么选择MyGo移民</h2>
<p class="section-subtitle">专业服务，值得信赖</p>
```

To:
```vue
<h2 class="section-title">{{ advantageTitle }}</h2>
<p class="section-subtitle">{{ advantageSubtitle }}</p>
```

- [ ] **Step 2: Add advantage section data from API**

In the `<script setup>` block (after `showcaseConfig` computed, around line 192), add:

```typescript
const advantageSection = computed(() => {
  if (homeConfig.value) {
    const config = homeConfig.value as unknown as Record<string, unknown>;
    const data = config.data as Record<string, unknown> | undefined;
    if (data && data.advantage_section) {
      return data.advantage_section as { section_title?: string; section_subtitle?: string };
    }
  }
  return null;
});

const advantageTitle = computed(() => advantageSection.value?.section_title || '为什么选择MyGo移民');
const advantageSubtitle = computed(() => advantageSection.value?.section_subtitle || '专业服务，值得信赖');
```

- [ ] **Step 3: Replace emoji icon rendering with SVG + fallback**

In the advantages computed (lines 265-300), update the map to pass `icon_type`:

```typescript
const advantages = computed(() => {
  if (homeConfig.value) {
    const config = homeConfig.value as unknown as Record<string, unknown>;
    const data = config.data as Record<string, unknown> | undefined;
    if (data && Array.isArray(data.advantage_items)) {
      return (data.advantage_items as Array<Record<string, string>>).map((a) => ({
        icon: a.icon || '',
        iconType: a.icon_type || '',
        title: a.title || '',
        description: a.description || '',
      }));
    }
  }

  return [
    { icon: '\u{1F3C6}', iconType: '', title: '10年行业经验', description: '深耕投资移民领域，拥有丰富的成功案例和行业资源' },
    { icon: '\u{1F465}', iconType: '', title: '专业顾问团队', description: '资深移民律师、前移民官组成的一流专业团队' },
    { icon: '\u{1F512}', iconType: '', title: '100%成功率', description: '严格的审核流程确保申请质量，保持行业领先成功率' },
    { icon: '\u{1F4DE}', iconType: '', title: '一站式服务', description: '从方案定制到成功获批，全程跟踪服务让您无忧' },
  ];
});
```

- [ ] **Step 4: Update template — icon rendering with SVG**

At line 79, replace:
```vue
<div class="advantage-icon">{{ adv.icon }}</div>
```

With:
```vue
<div class="advantage-icon">
  <svg
    v-if="adv.iconType === 'lucide'"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    stroke-width="2"
    stroke-linecap="round"
    stroke-linejoin="round"
    class="advantage-svg"
    v-html="getIconContent(adv.icon)"
  />
  <span v-else class="advantage-emoji">{{ adv.icon }}</span>
</div>
```

- [ ] **Step 5: Add import and helper function**

In `<script setup>` (after the `useSeo` call at line 102):

```typescript
import { getIconByName } from '~/composables/lucideIcons'

function getIconContent(name: string): string {
  return getIconByName(name)?.content || ''
}
```

- [ ] **Step 6: Update CSS for SVG icon (matching design mockup)**

Replace `.advantage-icon` CSS (line 587-589):

```css
.advantage-icon {
  width: 56px;
  height: 56px;
  margin: 0 auto 16px;
  background: rgba(26, 58, 92, 0.08);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}
.advantage-svg {
  width: 28px;
  height: 28px;
  color: #c8963e;
}
.advantage-emoji {
  font-size: 28px;
  line-height: 1;
}
```

- [ ] **Step 7: Verify frontend compiles**

Run: `cd frontend && npx nuxi typecheck`
Expected: No type errors.

- [ ] **Step 8: Commit**

```bash
git add frontend/pages/index.vue
git commit -m "feat: render advantage title/subtitle from API, SVG icon with emoji fallback"
```

---

### Task 6: Backend — Run tests and verify end-to-end

**Files:**
- Test: `backend/internal/service/` (existing test files)
- Test: `backend/internal/handler/` (existing test files)

- [ ] **Step 1: Run backend tests**

```bash
cd backend && go test ./... -v
```

Expected: All tests pass.

- [ ] **Step 2: Verify the Go server starts correctly**

```bash
cd backend && go build ./cmd/server/
```

Expected: Build succeeds with no errors.

- [ ] **Step 3: Start dev server and smoke test**

Start backend and frontend, then verify:
- Visit `/admin/homepage` — three cards render with unified style
- Add an advantage item with a Lucide icon — icon picker works
- Save — data persists via API
- Visit `/` — advantages render with SVG icons and section title/subtitle

- [ ] **Step 4: Final commit if any fixes were needed**

```bash
git add -A
git commit -m "chore: final adjustments and verification"
```
