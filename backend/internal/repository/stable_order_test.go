package repository

import (
	"os"
	"strings"
	"testing"
)

func TestRepositoryListQueriesUseStableOrder(t *testing.T) {
	tests := []struct {
		name string
		file string
		want string
	}{
		{name: "pages", file: "page_repo.go", want: `Order("sort_order asc, id asc")`},
		{name: "page search", file: "page_repo.go", want: `Order("sort_order asc, id asc").`},
		{name: "projects", file: "project_repo.go", want: `Order("sort_order asc, id asc")`},
		{name: "cases", file: "case_repo.go", want: `Order("sort_order asc, id asc")`},
		{name: "faqs", file: "faq_repo.go", want: `Order("faqs.sort_order asc, faqs.id asc")`},
		{name: "faq search", file: "faq_repo.go", want: `Order("sort_order asc, id asc").`},
		{name: "faq distinct projects", file: "faq_repo.go", want: `Order("projects.sort_order asc, projects.id asc").`},
		{name: "testimonials", file: "testimonial_repo.go", want: `Order("sort_order asc, id asc")`},
		{name: "leads", file: "lead_repo.go", want: `Order("created_at desc, id desc")`},
		{name: "media", file: "media_repo.go", want: `Order("created_at desc, id desc")`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := os.ReadFile(tt.file)
			if err != nil {
				t.Fatalf("read %s: %v", tt.file, err)
			}
			if !strings.Contains(string(content), tt.want) {
				t.Fatalf("%s does not contain stable order %q", tt.file, tt.want)
			}
		})
	}
}
