// main.go
package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"html"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Post struct {
	Title   string
	Slug    string
	Date    time.Time
	Content template.HTML
}

type Config struct {
	Title   string
	BaseURL string
	Links   map[string]string
}

var config = Config{
	Title:   "nobloat.org",
	BaseURL: "https://nobloat.org",
	Links: map[string]string{
		"GitHub":       "https://github.com/nobloat",
		"nobloat.org":  "https://nobloat.org",
		"zeitkapsl.eu": "https://zeitkapsl.eu",
	},
}

func buildSite() {
	posts := loadPosts("content")
	os.MkdirAll("public", 0755)
	copyStaticAssets()
	generateIndex(posts)
	generatePosts(posts)
	generateSitemap(posts)
	generateFeed(posts)
	fmt.Println("Build complete.")
}

func main() {
	watch := flag.Bool("watch", false, "Rebuild site on file changes")
	flag.Parse()

	buildSite()

	if *watch {
		fmt.Println("Watching for changes...")
		watchFiles()
	}
}

func watchFiles() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	watchPaths := []string{"content", "style.css", "main.go"}
	for _, path := range watchPaths {
		if err := watcher.Add(path); err != nil {
			log.Println("watch error:", err)
		}
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			fmt.Println("Changed:", event.Name)
			buildSite()
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("watch error:", err)
		}
	}
}

func loadPosts(dir string) []Post {
	files, _ := os.ReadDir(dir)
	var posts []Post
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".md") {
			path := filepath.Join(dir, f.Name())
			stat, _ := os.Stat(path)
			data, _ := os.ReadFile(path)
			content, title := parseMarkdown(string(data))
			slug := strings.TrimSuffix(f.Name(), ".md")
			posts = append(posts, Post{
				Title:   title,
				Slug:    slug,
				Date:    stat.ModTime(),
				Content: template.HTML(content),
			})
		}
	}
	return posts
}

func parseMarkdown(input string) (content string, title string) {
	lines := strings.Split(input, "\n")
	var out strings.Builder
	inList := false
	inCode := false
	codeLang := ""

	// Inline formatting
	codeRe := regexp.MustCompile("`([^`\n]+)`")
	boldRe := regexp.MustCompile(`\*\*(.+?)\*\*`)
	italicRe := regexp.MustCompile(`\*(.+?)\*`)
	strikeRe := regexp.MustCompile(`~~(.+?)~~`)

	formatInline := func(text string) string {
		text = html.EscapeString(text)
		text = codeRe.ReplaceAllString(text, "<code>$1</code>")
		text = boldRe.ReplaceAllString(text, "<strong>$1</strong>")
		text = strikeRe.ReplaceAllString(text, "<del>$1</del>")
		text = italicRe.ReplaceAllString(text, "<em>$1</em>")
		return text
	}

	title = lines[0][2:]

	for _, raw := range lines {
		line := strings.TrimSpace(raw)

		// Handle code block start/end
		if strings.HasPrefix(line, "```") {
			if inCode {
				out.WriteString("</code></pre>\n")
				inCode = false
				continue
			}
			inCode = true
			codeLang = strings.TrimSpace(strings.TrimPrefix(line, "```"))
			if codeLang == "" {
				out.WriteString("<pre><code>\n")
			} else {
				out.WriteString(fmt.Sprintf("<pre><code class=\"language-%s\">\n", codeLang))
			}
			continue
		}

		if inCode {
			out.WriteString(html.EscapeString(raw) + "\n")
			continue
		}

		// End list if needed
		if inList && line == "" {
			out.WriteString("</ul>\n")
			inList = false
			continue
		}

		switch {
		case strings.HasPrefix(line, "# "):
			if inList {
				out.WriteString("</ul>\n")
				inList = false
			}
			out.WriteString("<h1>" + formatInline(strings.TrimPrefix(line, "# ")) + "</h1>\n")
		case strings.HasPrefix(line, "## "):
			if inList {
				out.WriteString("</ul>\n")
				inList = false
			}
			out.WriteString("<h2>" + formatInline(strings.TrimPrefix(line, "## ")) + "</h2>\n")
		case strings.HasPrefix(line, "### "):
			if inList {
				out.WriteString("</ul>\n")
				inList = false
			}
			out.WriteString("<h3>" + formatInline(strings.TrimPrefix(line, "### ")) + "</h3>\n")
		case strings.HasPrefix(line, "- "):
			if !inList {
				out.WriteString("<ul>\n")
				inList = true
			}
			out.WriteString("<li>" + formatInline(strings.TrimPrefix(line, "- ")) + "</li>\n")
		case line == "":
			if inList {
				out.WriteString("</ul>\n")
				inList = false
			}
		default:
			if inList {
				out.WriteString("</ul>\n")
				inList = false
			}
			out.WriteString("<p>" + formatInline(line) + "</p>\n")
		}
	}

	if inList {
		out.WriteString("</ul>\n")
	}
	if inCode {
		out.WriteString("</code></pre>\n") // ensure it's closed
	}

	return out.String(), title
}

func writeIfChanged(path string, content []byte) error {
	existing, err := os.ReadFile(path)
	if err == nil && bytes.Equal(existing, content) {
		fmt.Println("unchanged:", path)
		return nil
	}
	fmt.Println("writing:", path)
	return os.WriteFile(path, content, 0644)
}

func generateIndex(posts []Post) {
	tmpl := template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html>
<head>
<title>{{.Title}}</title>
<link rel="stylesheet" href="style.css">
</head>
<body>
<h1>{{.Title}}</h1>
<ul>
{{range .Posts}}<li><a href="{{.Slug}}.html">{{.Title}}</a></li>{{end}}
</ul>
<h2>Links</h2>
<ul>
{{range $name, $url := .Links}}<li><a href="{{$url}}">{{$name}}</a></li>{{end}}
</ul>
</body>
</html>`))
	var buf bytes.Buffer
	tmpl.Execute(&buf, map[string]any{"Title": config.Title, "Posts": posts, "Links": config.Links})
	_ = writeIfChanged("public/index.html", buf.Bytes())
}

func generatePosts(posts []Post) {
	tmpl := template.Must(template.New("post").Parse(`
<!DOCTYPE html>
<html>
<head>
<title>{{.Title}}</title>
<link rel="stylesheet" href="style.css">
</head>
<body>
<article>{{.Content}}</article>
<a href="index.html">Back to home</a>
</body>
</html>`))
	for _, post := range posts {
		var buf bytes.Buffer
		tmpl.Execute(&buf, post)
		_ = writeIfChanged("public/"+post.Slug+".html", buf.Bytes())
	}
}

func copyStaticAssets() {
	input, err := os.ReadFile("style.css")
	if err == nil {
		_ = writeIfChanged("public/style.css", input)
	}
}

func generateSitemap(posts []Post) {
	type URL struct {
		Loc     string `xml:"loc"`
		LastMod string `xml:"lastmod"`
	}
	type Urlset struct {
		XMLName xml.Name `xml:"urlset"`
		Xmlns   string   `xml:"xmlns,attr"`
		URLs    []URL    `xml:"url"`
	}
	var urls []URL
	for _, post := range posts {
		urls = append(urls, URL{
			Loc:     config.BaseURL + "/" + post.Slug + ".html",
			LastMod: post.Date.Format("2006-01-02"),
		})
	}
	urls = append(urls, URL{Loc: config.BaseURL + "/index.html", LastMod: time.Now().Format("2006-01-02")})
	data, _ := xml.MarshalIndent(Urlset{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}, "", "  ")
	_ = writeIfChanged("public/sitemap.xml", []byte(xml.Header+string(data)))
}

func generateFeed(posts []Post) {
	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8" ?>
<feed xmlns="http://www.w3.org/2005/Atom">
`)
	buf.WriteString(fmt.Sprintf("<title>%s</title>\n", config.Title))
	buf.WriteString(fmt.Sprintf("<link href=\"%s/feed.xml\" rel=\"self\" />\n", config.BaseURL))
	buf.WriteString(fmt.Sprintf("<updated>%s</updated>\n", time.Now().Format(time.RFC3339)))
	for _, post := range posts {
		buf.WriteString("<entry>\n")
		buf.WriteString(fmt.Sprintf("<title>%s</title>\n", post.Title))
		buf.WriteString(fmt.Sprintf("<link href=\"%s/%s.html\"/>\n", config.BaseURL, post.Slug))
		buf.WriteString(fmt.Sprintf("<updated>%s</updated>\n", post.Date.Format(time.RFC3339)))
		buf.WriteString(fmt.Sprintf("<id>%s/%s.html</id>\n", config.BaseURL, post.Slug))
		buf.WriteString("</entry>\n")
	}
	buf.WriteString("</feed>")
	_ = writeIfChanged("public/feed.xml", buf.Bytes())
}
