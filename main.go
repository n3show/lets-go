package main

import (
	"bufio"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Segment is a prose block, a code block, or a terminal/shell block.
type Segment struct {
	IsCode     bool
	IsTerminal bool
	Text       string // raw text (prose: plain; code: will be escaped+highlighted in template)
}

type Example struct {
	Slug     string // url-friendly id, e.g. "for"
	Title    string // display title, e.g. "For"
	Number   int    // order
	Segments []Segment
}

type NavItem struct {
	Slug   string
	Title  string
	Active bool
}

type PageData struct {
	Title      string
	Nav        []NavItem
	Example    *Example
	IsHome     bool
	PrevSlug   string
	PrevTitle  string
	NextSlug   string
	NextTitle  string
	TotalCount int
}

var fileNameRe = regexp.MustCompile(`^(\d+)-(.+)\.go$`)

func titleFromSlug(slug string) string {
	parts := strings.Split(slug, "-")
	for i, p := range parts {
		if p == "" {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	return strings.Join(parts, " ")
}

func parseExampleFile(path string) (*Example, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	base := filepath.Base(path)
	m := fileNameRe.FindStringSubmatch(base)
	if m == nil {
		return nil, fmt.Errorf("bad filename: %s", base)
	}
	num, _ := strconv.Atoi(m[1])
	slug := m[2]

	var segments []Segment
	var curIsCode bool
	var curLines []string
	first := true

	flush := func() {
		if len(curLines) == 0 {
			return
		}
		text := strings.Join(curLines, "\n")
		text = strings.TrimRight(text, "\n")
		if text != "" {
			segments = append(segments, Segment{IsCode: curIsCode, Text: text})
		}
		curLines = nil
	}

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		isCommentLine := strings.HasPrefix(trimmed, "// ") || trimmed == "//"

		if isCommentLine {
			content := strings.TrimPrefix(trimmed, "//")
			content = strings.TrimPrefix(content, " ")
			if first || curIsCode {
				flush()
				curIsCode = false
			}
			first = false
			curLines = append(curLines, content)
		} else {
			if first || !curIsCode {
				flush()
				curIsCode = true
			}
			first = false
			curLines = append(curLines, line)
		}
	}
	flush()

	// trim leading/trailing blank code lines within each code segment,
	// and flag prose segments that are actually shell sessions (start with "$ ").
	for i := range segments {
		if segments[i].IsCode {
			lines := strings.Split(segments[i].Text, "\n")
			start, end := 0, len(lines)
			for start < end && strings.TrimSpace(lines[start]) == "" {
				start++
			}
			for end > start && strings.TrimSpace(lines[end-1]) == "" {
				end--
			}
			segments[i].Text = strings.Join(lines[start:end], "\n")
		} else if strings.HasPrefix(strings.TrimSpace(segments[i].Text), "$ ") {
			segments[i].IsTerminal = true
		}
	}

	return &Example{
		Slug:     slug,
		Title:    titleFromSlug(slug),
		Number:   num,
		Segments: segments,
	}, nil
}

const layoutTmpl = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{if .IsHome}}Let's Go — Go examples, one idea at a time{{else}}{{.Title}} · Let's Go{{end}}</title>
<link rel="stylesheet" href="static/style.css">
<script>
(function () {
  var saved = null;
  try { saved = localStorage.getItem('lg-theme'); } catch (e) {}
  var theme = saved || (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light');
  document.documentElement.setAttribute('data-theme', theme);
})();
</script>
</head>
<body>
<div class="shell">
  <header class="topbar">
    <button class="nav-toggle" id="navToggle" aria-label="Toggle navigation" aria-expanded="false">
      <span></span><span></span><span></span>
    </button>
    <a class="brand" href="index.html">
      <span class="brand-mark">&gt;_</span>
      <span class="brand-text">Let's Go</span>
    </a>
    <span class="topbar-count">{{.TotalCount}} examples</span>
    <button class="theme-toggle" id="themeToggle" aria-label="Toggle color theme" type="button">
      <svg class="icon-sun" viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="4"/><path d="M12 2v2M12 20v2M4.93 4.93l1.41 1.41M17.66 17.66l1.41 1.41M2 12h2M20 12h2M4.93 19.07l1.41-1.41M17.66 6.34l1.41-1.41"/></svg>
      <svg class="icon-moon" viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79Z"/></svg>
    </button>
  </header>

  <div class="layout">
    <nav class="sidebar" id="sidebar">
      <ol class="toc">
      {{range .Nav}}
        <li class="{{if .Active}}is-active{{end}}"><a href="{{.Slug}}.html">{{.Title}}</a></li>
      {{end}}
      </ol>
    </nav>
    <div class="scrim" id="scrim"></div>

    <main class="content">
    {{if .IsHome}}
      <section class="hero">
        <p class="eyebrow">A field guide, not a manual</p>
        <h1>Let's Go</h1>
        <p class="lede">Small, runnable Go programs, one idea each. Read the comment, read the code, run it in your head, move on. Built to work just as well on a phone on the metro as on a desktop.</p>
        <a class="cta" href="{{(index .Nav 0).Slug}}.html">Start with {{(index .Nav 0).Title}} →</a>
      </section>
      <ol class="index-grid">
      {{range .Nav}}
        <li><a href="{{.Slug}}.html"><span class="idx-title">{{.Title}}</span></a></li>
      {{end}}
      </ol>
    {{else}}
      <article class="example">
        <header class="example-head">
          <p class="eyebrow">Example {{printf "%02d" .Example.Number}} of {{.TotalCount}}</p>
          <h1>{{.Example.Title}}</h1>
        </header>
        {{range .Example.Segments}}
          {{if .IsCode}}
            <pre class="code-block"><button class="copy-btn" type="button">Copy</button><code>{{.Text | safeHTML}}</code></pre>
          {{else if .IsTerminal}}
            <pre class="terminal-block">{{.Text | safeHTML}}</pre>
          {{else}}
            <div class="prose">{{.Text | safeHTML}}</div>
          {{end}}
        {{end}}
        <nav class="pager">
          {{if .PrevSlug}}<a class="pager-prev" href="{{.PrevSlug}}.html">← {{.PrevTitle}}</a>{{else}}<span></span>{{end}}
          {{if .NextSlug}}<a class="pager-next" href="{{.NextSlug}}.html">{{.NextTitle}} →</a>{{end}}
        </nav>
      </article>
    {{end}}
    </main>
  </div>
</div>
<script src="static/app.js"></script>
</body>
</html>
`

func proseToHTML(text string) template.HTML {
	lines := strings.Split(text, "\n")
	var b strings.Builder
	for i, l := range lines {
		trimmed := strings.TrimSpace(l)
		if trimmed == "" {
			b.WriteString("<br><br>")
			continue
		}
		b.WriteString(template.HTMLEscapeString(trimmed))
		if i < len(lines)-1 {
			b.WriteString("<br>")
		}
	}
	return template.HTML(b.String())
}

func terminalToHTML(text string) template.HTML {
	lines := strings.Split(text, "\n")
	var b strings.Builder
	for i, l := range lines {
		trimmed := strings.TrimSpace(l)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "$ ") {
			rest := strings.TrimPrefix(trimmed, "$ ")
			b.WriteString(`<span class="term-line"><span class="term-prompt">$</span> <span class="term-cmd">`)
			b.WriteString(template.HTMLEscapeString(rest))
			b.WriteString(`</span></span>`)
		} else {
			b.WriteString(`<span class="term-line term-out">`)
			b.WriteString(template.HTMLEscapeString(trimmed))
			b.WriteString(`</span>`)
		}
		if i < len(lines)-1 {
			b.WriteString("\n")
		}
	}
	return template.HTML(b.String())
}

func main() {
	files, err := filepath.Glob("examples/*.go")
	if err != nil || len(files) == 0 {
		fmt.Println("no example files found:", err)
		os.Exit(1)
	}
	sort.Strings(files)

	var examples []*Example
	for _, fp := range files {
		ex, err := parseExampleFile(fp)
		if err != nil {
			fmt.Println("skip", fp, err)
			continue
		}
		examples = append(examples, ex)
	}
	sort.Slice(examples, func(i, j int) bool { return examples[i].Number < examples[j].Number })

	var nav []NavItem
	for _, ex := range examples {
		nav = append(nav, NavItem{Slug: ex.Slug, Title: ex.Title})
	}

	funcMap := template.FuncMap{
		"safeHTML": func(s string) template.HTML { return template.HTML(s) },
	}
	tmpl := template.Must(template.New("layout").Funcs(funcMap).Parse(layoutTmpl))

	os.MkdirAll("docs/static", 0755)

	// copy static assets
	copyFile("static/style.css", "docs/static/style.css")
	copyFile("static/app.js", "docs/static/app.js")

	// render index
	{
		navCopy := make([]NavItem, len(nav))
		copy(navCopy, nav)
		data := PageData{Title: "Let's Go", Nav: navCopy, IsHome: true, TotalCount: len(examples)}
		writePage(tmpl, "docs/index.html", data)
	}

	// render each example
	for i, ex := range examples {
		navCopy := make([]NavItem, len(nav))
		copy(navCopy, nav)
		navCopy[i].Active = true

		// build prose-rendered segments
		renderedSegs := make([]Segment, len(ex.Segments))
		copy(renderedSegs, ex.Segments)
		exCopy := *ex
		exCopy.Segments = renderedSegs

		data := PageData{
			Title:      ex.Title,
			Nav:        navCopy,
			Example:    &exCopy,
			TotalCount: len(examples),
		}
		if i > 0 {
			data.PrevSlug = examples[i-1].Slug
			data.PrevTitle = examples[i-1].Title
		}
		if i < len(examples)-1 {
			data.NextSlug = examples[i+1].Slug
			data.NextTitle = examples[i+1].Title
		}
		writePageExample(tmpl, fmt.Sprintf("docs/%s.html", ex.Slug), data)
	}

	fmt.Printf("built %d pages into docs/\n", len(examples)+1)
}

func writePage(tmpl *template.Template, path string, data PageData) {
	f, err := os.Create(path)
	if err != nil {
		fmt.Println("error creating", path, err)
		return
	}
	defer f.Close()
	if err := tmpl.Execute(f, data); err != nil {
		fmt.Println("error executing template for", path, err)
	}
}

func writePageExample(tmpl *template.Template, path string, data PageData) {
	for i, seg := range data.Example.Segments {
		if seg.IsCode {
			data.Example.Segments[i].Text = template.HTMLEscapeString(seg.Text)
		} else if seg.IsTerminal {
			data.Example.Segments[i].Text = string(terminalToHTML(seg.Text))
		} else {
			data.Example.Segments[i].Text = string(proseToHTML(seg.Text))
		}
	}
	f, err := os.Create(path)
	if err != nil {
		fmt.Println("error creating", path, err)
		return
	}
	defer f.Close()
	if err := tmpl.Execute(f, data); err != nil {
		fmt.Println("error executing template for", path, err)
	}
}

func copyFile(src, dst string) {
	b, err := os.ReadFile(src)
	if err != nil {
		fmt.Println("warn: could not read", src, err)
		return
	}
	os.WriteFile(dst, b, 0644)
}
