package markdown

import (
	"strings"
	"testing"
)

func TestRenderRewritesLinks(t *testing.T) {
	links := Links{
		DocsRoot:  "docs",
		Docs:      map[string]bool{"docs/board/TASKS.md": true, "docs/runbooks/deploy.md": true},
		TaskPaths: map[string]string{"docs/board/tasks/SS-2-x.md": "SS-2"},
		CodeURL:   func(p string) string { return "https://git.example/blob/develop/" + p },
	}
	source := strings.Join([]string{
		"[rules](../TASKS.md)",
		"[deploy](../../runbooks/deploy.md#steps)",
		"[task](SS-2-x.md)",
		"[code](../../../api/main.go)",
		"[gone](../missing.md)",
		"[web](https://example.com)",
		"[outside](../../../../etc/passwd)",
	}, "\n\n")

	html, err := Render([]byte(source), "docs/board/tasks/SS-1.md", links)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`<a href="/docs/board/TASKS.md">rules</a>`,
		`<a href="/docs/runbooks/deploy.md#steps">deploy</a>`,
		`<a href="/?task=SS-2">task</a>`,
		`<a href="https://git.example/blob/develop/api/main.go" target="_blank" rel="noreferrer">code</a>`,
		`<p>gone</p>`,
		`<a href="https://example.com" target="_blank" rel="noreferrer">web</a>`,
		`<p>outside</p>`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in\n%s", want, html)
		}
	}
}

func TestRenderEscapesRawHTMLAndKeepsMermaid(t *testing.T) {
	html, err := Render([]byte("<script>alert(1)</script>\n\n```mermaid\nflowchart LR\n  A --> B\n```\n"), "docs/x.md", Links{})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(html, "<script>") {
		t.Error("raw HTML must not pass through")
	}
	if !strings.Contains(html, `<code class="language-mermaid">`) {
		t.Errorf("mermaid block not marked: %s", html)
	}
}

func TestRenderWithoutCodeURLDropsCodeLinks(t *testing.T) {
	html, _ := Render([]byte("[code](../main.go)"), "docs/x.md", Links{})
	if strings.Contains(html, "<a ") {
		t.Errorf("code link kept without a host: %s", html)
	}
}

func TestTitle(t *testing.T) {
	if got := Title([]byte("intro\n# The title\n## sub"), "fallback"); got != "The title" {
		t.Errorf("got %q", got)
	}
	if got := Title([]byte("no heading"), "fallback"); got != "fallback" {
		t.Errorf("got %q", got)
	}
}
