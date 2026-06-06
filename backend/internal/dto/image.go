package dto

// ImageVariantInfo 描述一个图片变体的 URL 和宽度。
type ImageVariantInfo struct {
	URL   string `json:"url"`
	Width int    `json:"width"`
}
