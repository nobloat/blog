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
	"unicode"

	"github.com/fsnotify/fsnotify"
)

type Post struct {
	Title   string
	Slug    string
	Date    time.Time
	Content template.HTML
}

type Config struct {
	Title    string
	BaseURL  string
	Links    map[string]string
	Projects map[string]string
}

var config = Config{
	Title:   "][ nobloat.org",
	BaseURL: "https://nobloat.org",
	Links: map[string]string{
		"GitHub":               "https://github.com/nobloat",
		"nobloat.org":          "https://nobloat.org",
		"zeitkapsl.eu":         "https://zeitkapsl.eu",
		"hardcode.at":          "https://hardcode.at",
		"spiessknafl.at/peter": "https://spiessknafl.at/peter",
		"RSS Feed":             "./feed.xml",
	},
	Projects: map[string]string{
		"[nobloat/css](https://github.com/nobloat/css)":             "modular vanilla CSS3 components",
		"[nobloat/bare-jvm](https://github.com/nobloat/bare-jvm)":   "[baremessages](https://baremessages.org/) implementation for the JVM",
		"[cinemast/dbolve](https://github.com/cinemast/dbolve)":     "very minimalistic database migration tool for golang projects",
		"[nobloat/tinyviper](https://github.com/nobloat/tinyviper)": "very minimalistic alternative to the famous [spf13/viper](https://github.com/spf13/viper) configuration library",
	},
}

func sanitizeAnchor(input string) string {
	var out strings.Builder
	for _, r := range input {
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r):
			out.WriteRune(r)
		case r == '-' || r == '_' || r == '.' || r == '~':
			out.WriteRune(r)
		default:
			out.WriteRune('-') // replace unsafe characters with dash
		}
	}
	return strings.ToLower(out.String())
}

func buildSite() {
	posts := loadPosts("articles")
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

	watchPaths := []string{"articles", "style.css", "main.go", "index.html", "article.html"}
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

var (
	// Inline formatting
	codeRe   = regexp.MustCompile("`([^`\n]+)`")
	boldRe   = regexp.MustCompile(`\*\*(.+?)\*\*`)
	italicRe = regexp.MustCompile(`\*(.+?)\*`)
	strikeRe = regexp.MustCompile(`~~(.+?)~~`)
	imageRe  = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
	linkRe   = regexp.MustCompile(`\[([^\]]*)\]\(([^)]+)\)`)
)

func formatInline(text string) string {
	text = html.EscapeString(text)
	text = imageRe.ReplaceAllString(text, `<figure><img src="$2" alt="$1"><figcaption>$1</figcaption></figure>`)
	text = linkRe.ReplaceAllString(text, `<a href="$2">$1</a>`)
	text = codeRe.ReplaceAllString(text, "<code>$1</code>")
	text = boldRe.ReplaceAllString(text, "<strong>$1</strong>")
	text = strikeRe.ReplaceAllString(text, "<del>$1</del>")
	text = italicRe.ReplaceAllString(text, "<em>$1</em>")
	return text
}

func parseMarkdown(input string) (content string, title string) {
	lines := strings.Split(input, "\n")
	var out strings.Builder
	inList := false
	inCode := false
	codeLang := ""

	// First line is title
	if len(lines) > 0 && strings.HasPrefix(lines[0], "# ") {
		title = strings.TrimPrefix(lines[0], "# ")
	}

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
				out.WriteString(fmt.Sprintf("<pre><code class=\"language-%s\">", codeLang))
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
			id := sanitizeAnchor(strings.TrimPrefix(line, "## "))
			out.WriteString("<h2 id=\"" + id + "\"><a href=\"#" + id + "\">" + formatInline(strings.TrimPrefix(line, "## ")) + "</a></h2>\n")
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
		out.WriteString("</code></pre>\n")
	}

	return out.String(), title
}

func writeIfChanged(path string, content []byte) error {

	fmt.Println("writing:", path)
	return os.WriteFile(path, content, 0644)
}

var funcMap = template.FuncMap{
	"md2html": formatInline,
	"safeHTML": func(s string) template.HTML {
		return template.HTML(s)
	},
}

func generateIndex(posts []Post) {
	tpl, err := os.ReadFile("index.html")
	if err != nil {
		panic(err)
	}
	tmpl := template.Must(template.New("index").Funcs(funcMap).Parse(string(tpl)))
	var buf bytes.Buffer
	tmpl.Execute(&buf, map[string]any{"Title": config.Title, "Posts": posts, "Links": config.Links, "Projects": config.Projects})
	_ = writeIfChanged("public/index.html", buf.Bytes())
}

func generatePosts(posts []Post) {
	tpl, err := os.ReadFile("article.html")
	if err != nil {
		panic(err)
	}
	tmpl := template.Must(template.New("post").Funcs(funcMap).Parse(string(tpl)))
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
