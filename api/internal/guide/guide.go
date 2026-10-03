package guide

import (
	"embed"
	"errors"
	"io/fs"
	"path"
	"strings"

	"fuda/internal/markdown"
)

//go:embed pages/*.md
var pages embed.FS

var ErrNotFound = errors.New("guide page not found")

type Page struct {
	Slug  string `json:"slug"`
	Title string `json:"title"`
}

type Rendered struct {
	Page
	HTML string `json:"html"`
}

func List() []Page {
	entries, _ := fs.ReadDir(pages, "pages")
	out := make([]Page, 0, len(entries))
	for _, e := range entries {
		content, err := pages.ReadFile("pages/" + e.Name())
		if err != nil {
			continue
		}
		slug := slugOf(e.Name())
		out = append(out, Page{Slug: slug, Title: markdown.Title(content, slug)})
	}
	return out
}

func Render(slug string) (Rendered, error) {
	entries, _ := fs.ReadDir(pages, "pages")
	for _, e := range entries {
		if slugOf(e.Name()) != slug {
			continue
		}
		content, err := pages.ReadFile("pages/" + e.Name())
		if err != nil {
			return Rendered{}, err
		}
		html, err := markdown.Render(content, "guide/"+e.Name(), markdown.Links{})
		if err != nil {
			return Rendered{}, err
		}
		return Rendered{Page: Page{Slug: slug, Title: markdown.Title(content, slug)}, HTML: html}, nil
	}
	return Rendered{}, ErrNotFound
}

func slugOf(name string) string {
	name = strings.TrimSuffix(path.Base(name), ".md")
	if _, rest, ok := strings.Cut(name, "-"); ok && len(name) > 3 && name[0] >= '0' && name[0] <= '9' {
		return rest
	}
	return name
}
