package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"mygo-immigration/backend/internal/config"
	"mygo-immigration/backend/internal/database"
	"mygo-immigration/backend/internal/service"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var force bool

func init() {
	generateCmd.Flags().BoolVar(&force, "force", false, "强制重新生成所有变体，即使文件已存在")
	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:   "generate-variants",
	Short: "为所有被引用的图片按场景生成对应规格的缩略图变体",
	Long:  "扫描各实体表、富文本内容和首页配置中的 /uploads/ 图片 URL，根据所属场景（律师3:4、项目16:9、案例16:9、评价1:1、资讯3:2、轮播2.4:1、通用）生成对应尺寸的 JPEG 变体。使用 --force 强制覆盖已有变体。",
	Run:   runGenerateVariants,
}

// column queries mapped to their upload context
var columnContexts = []struct {
	table, col, context string
}{
	{"projects", "cover_image", "project"},
	{"cases", "photo_url", "case"},
	{"pages", "cover_image", "page-cover"},
	{"lawyers", "photo_url", "lawyer"},
	{"testimonials", "avatar_url", "testimonial"},
}

func runGenerateVariants(cmd *cobra.Command, args []string) {
	cfg := config.Load()

	db, err := database.InitMySQL(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// url → context
	urlCtx := collectAllURLsWithContext(db)

	// Deduplicate: filename → {savePath, context}
	// If an image appears in multiple contexts, prefer the most specific (non-general) one.
	type entry struct {
		path    string
		context string
	}
	seen := make(map[string]entry) // filename → entry
	for u, ctx := range urlCtx {
		u = strings.TrimSpace(u)
		if u == "" || !strings.Contains(u, "/uploads/") {
			continue
		}
		name := filepath.Base(u)
		savePath := filepath.Join("./uploads", name)
		if _, err := os.Stat(savePath); os.IsNotExist(err) {
			continue
		}
		existing, ok := seen[name]
		if ok && existing.context != "general" {
			continue // keep the more specific context
		}
		seen[name] = entry{path: savePath, context: ctx}
	}

	fmt.Printf("Found %d unique image files on disk\n\n", len(seen))

	ok, skipped, failed := 0, 0, 0

	for name, e := range seen {
		ext := filepath.Ext(name)
		baseName := name[:len(name)-len(ext)]

		// Check if variants already exist
		specs := contextSpecsForContext(e.context)
		allFilesExist := true
		for _, spec := range specs {
			vf := filepath.Join("./uploads", baseName+"_"+spec.Name+".jpg")
			if _, err := os.Stat(vf); os.IsNotExist(err) {
				allFilesExist = false
				break
			}
		}

		if !force && allFilesExist {
			fmt.Printf("SKIP %-45s context=%-16s all variants exist\n", name, e.context)
			skipped++
			continue
		}

		variants, err := service.GenerateVariantsForContext(e.path, baseName, e.context)
		if err != nil {
			fmt.Printf("FAIL %-45s context=%-16s %v\n", name, e.context, err)
			failed++
			continue
		}

		// Update media table variants if a record exists
		db.Table("media").Where("filename = ?", name).Update("variants", toJSON(variants))

		n := len(variants)
		if force && allFilesExist {
			fmt.Printf("OVERWRITE %-40s context=%-16s → %d variants\n", name, e.context, n)
		} else {
			fmt.Printf("OK   %-45s context=%-16s → %d variants\n", name, e.context, n)
		}
		ok++
	}

	fmt.Printf("\nDone: %d ok, %d skipped, %d failed\n", ok, skipped, failed)
}

// contextSpecsForContext returns the variant specs for a given context (used for file-existence checks).
func contextSpecsForContext(ctx string) []service.VariantSpec {
	specs, ok := service.GetContextSpecs(ctx)
	if !ok {
		specs, _ = service.GetContextSpecs("general")
	}
	return specs
}

// collectAllURLsWithContext gathers every /uploads/ URL with its context.
func collectAllURLsWithContext(db *gorm.DB) map[string]string {
	result := make(map[string]string)

	// 1. Entity columns — each has a known context
	var rows []string
	for _, q := range columnContexts {
		rows = rows[:0]
		db.Table(q.table).Select(q.col).Where(q.col+" LIKE ?", "/uploads/%").Find(&rows)
		for _, u := range rows {
			u = strings.TrimSpace(u)
			if u != "" {
				result[u] = q.context
			}
		}
	}

	// 2. Rich text <img> tags — context is "general" (inline content images)
	var contents []string
	db.Table("cases").Select("content").Where("content LIKE ?", "%/uploads/%").Find(&contents)
	for _, u := range extractImgSrcs(contents) {
		if _, exists := result[u]; !exists {
			result[u] = "general"
		}
	}
	contents = nil
	db.Table("pages").Select("content").Where("content LIKE ?", "%/uploads/%").Find(&contents)
	for _, u := range extractImgSrcs(contents) {
		if _, exists := result[u]; !exists {
			result[u] = "general"
		}
	}

	// 3. Home config JSON — try to detect context from surrounding keys, fallback to "general"
	var configs []string
	db.Table("home_configs").Select("config_value").Find(&configs)
	for u, ctx := range extractUploadPathsWithHints(configs) {
		if _, exists := result[u]; !exists {
			result[u] = ctx
		}
	}

	return result
}

func extractImgSrcs(contents []string) []string {
	var urls []string
	for _, c := range contents {
		for {
			start := strings.Index(c, "<img")
			if start == -1 {
				break
			}
			c = c[start:]
			end := strings.Index(c, ">")
			if end == -1 {
				break
			}
			tag := c[:end+1]
			c = c[end+1:]

			srcStart := strings.Index(tag, "src=")
			if srcStart == -1 {
				continue
			}
			tag = tag[srcStart+4:]
			quote := tag[0]
			tag = tag[1:]
			srcEnd := strings.IndexByte(tag, quote)
			if srcEnd == -1 {
				continue
			}
			src := tag[:srcEnd]
			if strings.Contains(src, "/uploads/") {
				urls = append(urls, src)
			}
		}
	}
	return urls
}

// extractUploadPathsWithHints extracts /uploads/ paths from JSON and
// attempts to detect context from JSON keys like "image", "slide", "hero".
func extractUploadPathsWithHints(configs []string) map[string]string {
	result := make(map[string]string)
	for _, raw := range configs {
		detectContextFromJSON(raw, result)
	}
	return result
}

// detectContextFromJSON scans JSON text for /uploads/ URLs and guesses
// context based on nearby JSON key names.
func detectContextFromJSON(raw string, out map[string]string) {
	for {
		idx := strings.Index(raw, "/uploads/")
		if idx == -1 {
			break
		}

		raw = raw[idx:]
		end := strings.IndexAny(raw, "\"'\\s,}]")
		if end == -1 {
			end = len(raw)
		}
		url := raw[:end]
		raw = raw[1:]

		if _, exists := out[url]; exists {
			continue
		}

		// Scan backwards from the URL position for the nearest JSON key
		searchLen := 500
		if idx < searchLen {
			searchLen = idx
		}
		preceding := raw[idx-searchLen : idx]
		ctx := guessContextFromJSON(preceding)
		out[url] = ctx
	}
}

func guessContextFromJSON(preceding string) string {
	// Check the nearest JSON key before the /uploads/ URL
	// Pattern: "key": "/uploads/..." or "key":"/uploads/..."
	lastKey := ""
	for {
		idx := strings.LastIndex(preceding, "\"")
		if idx == -1 {
			break
		}
		preceding = preceding[:idx]
		prevIdx := strings.LastIndex(preceding, "\"")
		if prevIdx == -1 {
			break
		}
		lastKey = strings.ToLower(preceding[prevIdx+1 : idx])
		break
	}

	switch lastKey {
	case "image", "slide_image", "hero_image", "background_image":
		return "homepage-slide"
	case "og_image":
		return "og-image"
	case "logo", "site_logo", "organization_logo":
		return "general"
	case "favicon", "site_favicon":
		return "favicon"
	case "qr", "wechat", "contact_wechat", "contact_wechat_mp", "contact_wechat_video":
		return "qr-code"
	case "avatar", "avatar_url":
		return "testimonial"
	case "photo", "photo_url", "cover", "cover_image":
		return "project" // most common default for covers
	default:
		return "general"
	}
}

func toJSON(v map[string]string) string {
	if len(v) == 0 {
		return "null"
	}
	parts := make([]string, 0, len(v))
	for k, val := range v {
		parts = append(parts, fmt.Sprintf(`"%s":"%s"`, k, val))
	}
	return "{" + strings.Join(parts, ",") + "}"
}
