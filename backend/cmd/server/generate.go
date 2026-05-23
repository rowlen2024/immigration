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
	Short: "为所有被引用的图片生成缩略图变体",
	Long:  "扫描所有实体表、富文本内容和首页配置中的 /uploads/ 图片 URL，为每个存在的图片文件生成 thumb/sm/md/lg 四种尺寸 JPEG 变体。使用 --force 强制覆盖已有变体。",
	Run:   runGenerateVariants,
}

func runGenerateVariants(cmd *cobra.Command, args []string) {
	cfg := config.Load()

	db, err := database.InitMySQL(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	urls := collectAllURLs(db)

	// Deduplicate and resolve to disk files
	seen := make(map[string]string) // filename → full path
	for _, u := range urls {
		u = strings.TrimSpace(u)
		if u == "" || !strings.Contains(u, "/uploads/") {
			continue
		}
		name := filepath.Base(u)
		if _, ok := seen[name]; ok {
			continue
		}
		savePath := filepath.Join("./uploads", name)
		if _, err := os.Stat(savePath); os.IsNotExist(err) {
			continue
		}
		seen[name] = savePath
	}

	fmt.Printf("Found %d unique image files on disk\n\n", len(seen))

	var svc service.MediaService
	ok, skipped, failed := 0, 0, 0
	variantNames := []string{"thumb", "sm", "md", "lg"}

	for name, savePath := range seen {
		ext := filepath.Ext(name)
		baseName := name[:len(name)-len(ext)]

		allFilesExist := true
		for _, vn := range variantNames {
			vf := filepath.Join("./uploads", baseName+"_"+vn+".jpg")
			if _, err := os.Stat(vf); os.IsNotExist(err) {
				allFilesExist = false
				break
			}
		}

		if !force && allFilesExist {
			fmt.Printf("SKIP %s: all variants exist\n", name)
			skipped++
			continue
		}

		variants, err := svc.GenerateVariants(savePath, baseName)
		if err != nil {
			fmt.Printf("FAIL %s: %v\n", name, err)
			failed++
			continue
		}

		// Update media table variants if a record exists
		db.Table("media").Where("filename = ?", name).Update("variants", toJSON(variants))

		fmt.Printf("OK   %s → %d variants\n", name, len(variants))
		ok++
	}

	fmt.Printf("\nDone: %d ok, %d skipped, %d failed\n", ok, skipped, failed)
}

// collectAllURLs gathers every /uploads/ URL from entity columns,
// rich text content, and home config.
func collectAllURLs(db *gorm.DB) []string {
	var urls []string
	var rows []string

	// Direct column URLs
	queries := []struct{ table, col string }{
		{"projects", "cover_image"},
		{"cases", "photo_url"},
		{"pages", "cover_image"},
		{"lawyers", "photo_url"},
		{"testimonials", "avatar_url"},
	}
	for _, q := range queries {
		db.Table(q.table).Select(q.col).Where(q.col+" LIKE ?", "/uploads/%").Find(&rows)
		urls = append(urls, rows...)
	}

	// Rich text content <img> tags
	var contents []string
	db.Table("cases").Select("content").Where("content LIKE ?", "%/uploads/%").Find(&contents)
	urls = append(urls, extractImgSrcs(contents)...)
	contents = nil
	db.Table("pages").Select("content").Where("content LIKE ?", "%/uploads/%").Find(&contents)
	urls = append(urls, extractImgSrcs(contents)...)

	// Home config JSON
	var configs []string
	db.Table("home_configs").Select("config_value").Find(&configs)
	urls = append(urls, extractUploadPaths(configs)...)

	return urls
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

func extractUploadPaths(configs []string) []string {
	var urls []string
	for _, raw := range configs {
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
			urls = append(urls, raw[:end])
			raw = raw[1:]
		}
	}
	return urls
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
