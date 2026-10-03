package guide

import (
	"strings"
	"testing"
)

func TestPagesListInOrderWithTitles(t *testing.T) {
	pages := List()
	if len(pages) < 8 || pages[0].Slug != "overview" || pages[0].Title == "" {
		t.Fatalf("pages: %+v", pages)
	}
	for _, p := range pages {
		if _, err := Render(p.Slug); err != nil {
			t.Errorf("%s: %v", p.Slug, err)
		}
	}
}

func TestGuideLinksStayInsideTheGuide(t *testing.T) {
	page, err := Render("overview")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(page.HTML, `href="/guide/setup-a-repo"`) {
		t.Errorf("guide link missing: %s", page.HTML)
	}
	if _, err := Render("nope"); err == nil {
		t.Error("unknown page must fail")
	}
}
