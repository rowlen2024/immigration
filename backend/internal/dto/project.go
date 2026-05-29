package dto

// DashboardStats holds admin dashboard statistics.
type DashboardStats struct {
	TotalProjects int64   `json:"totalProjects"`
	TotalPages    int64   `json:"totalPages"`
	TotalLeads    int64   `json:"totalLeads"`
	TotalCases    int64   `json:"totalCases"`
	UnreadLeads   int64   `json:"unreadLeads"`
	Trends        []Trend `json:"trends"`
}

// Trend represents a single metric's month-over-month trend.
type Trend struct {
	Key       string `json:"key"`
	Direction string `json:"direction"`
	Percent   int    `json:"percent"`
	Label     string `json:"label"`
}

// CompareRow represents a single comparison row.
type CompareRow struct {
	Label  string     `json:"label"`
	Values []string   `json:"values"`
	Items  [][]string `json:"items,omitempty"`
}

// CompareResult holds the full comparison output.
type CompareResult struct {
	Projects []CompareProject `json:"projects"`
	Rows     []CompareRow     `json:"rows"`
}

// CompareProject holds minimal project info for the comparison header.
type CompareProject struct {
	Title string `json:"title"`
	Slug  string `json:"slug"`
}
