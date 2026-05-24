package service

import "github.com/microcosm-cc/bluemonday"

var HTMLSanitizer = func() *bluemonday.Policy {
	p := bluemonday.NewPolicy()
	p.AllowElements("h1", "h2", "h3", "h4", "h5", "h6", "p", "br", "hr",
		"ul", "ol", "li", "blockquote", "pre", "code", "strong", "em", "u", "s",
		"a", "img", "table", "thead", "tbody", "tr", "th", "td",
		"div", "span", "iframe", "video", "source")
	p.AllowAttrs("src", "alt", "title", "width", "height").OnElements("img")
	p.AllowAttrs("href", "title", "target", "rel").OnElements("a")
	p.AllowAttrs("src", "frameborder", "allowfullscreen").OnElements("iframe")
	p.AllowAttrs("src", "controls", "width", "height").OnElements("video", "source")
	p.AllowAttrs("style", "class").OnElements("span", "div", "td", "th")
	p.AllowAttrs("class").OnElements("table", "thead", "tbody", "tr", "img", "a")
	p.AllowStyles("color", "background-color", "text-align").OnElements("span", "td", "th")
	p.AllowURLSchemes("http", "https", "mailto")
	p.AllowRelativeURLs(true)
	p.RequireNoFollowOnLinks(true)
	return p
}()
