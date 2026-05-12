package config

// CompareField defines a single dimension available for project comparison.
type CompareField struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	From  string `json:"from"`
}

// CompareFields is the list of all available comparison dimensions.
var CompareFields = []CompareField{
	{Key: "overview_text", Label: "项目介绍", From: "project"},
	{Key: "investment_amount", Label: "投资金额", From: "project"},
	{Key: "processing_period", Label: "办理周期", From: "project"},
	{Key: "target_crowd", Label: "适合人群", From: "project"},
	{Key: "country", Label: "国家", From: "project"},
	{Key: "costs_total", Label: "费用总计", From: "project"},
	{Key: "requirements_count", Label: "申请条件", From: "requirements"},
	{Key: "timeline_steps", Label: "流程步骤数", From: "timeline"},
	{Key: "tagline", Label: "标语", From: "project"},
	{Key: "policy_title", Label: "政策标题", From: "project"},
}
