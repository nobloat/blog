# Hello blog

When I decided to start writing about technical topics, I needed a simple way to publish blog posts. The requirements were straightforward: write in Markdown, generate static HTML, support RSS feeds, and keep everything version-controlled in Git. Most importantly, I wanted something I could understand completely, maintain easily, and that wouldn't break on major upgrades.

## Existing tools were overkill

The static site generator ecosystem is vast, but most solutions come with significant complexity:

[**Hugo**](https://gohugo.io/) is fast and popular, but it's a massive framework with hundreds of features I'd never use. The configuration files, theme system, and plugin ecosystem add layers of abstraction that make it hard to understand what's actually happening. When something breaks or I want to customize behavior, it is too complex for me to reason about all this. And I have had issues in the past where untouched static websites suddenly would no longer build after one or two years.

[**Jekyll**](https://jekyllrb.com/) requires Ruby, Bundler, and a complex gem ecosystem. Every time I'd want to update or deploy, I'd need to ensure the right Ruby version, manage gem dependencies, and deal with potential version conflicts. The overhead of maintaining a Ruby environment just to generate static HTML felt excessive.

[**Next.js**](https://nextjs.org/) and other React-based static site generators are powerful, but they bring the entire JavaScript ecosystem with them. Node modules, build tools, transpilation, and the constant churn of the npm ecosystem—all for what is essentially text processing and template rendering.

Even simpler tools like [**Zola**](https://www.getzola.org/) or [**11ty**](https://www.11ty.dev/) still require learning their specific conventions, configuration formats, and template languages. They're better than the heavyweights, but they're still frameworks with their own abstractions.

What I needed was clear: write Markdown files, run a simple command, get HTML. No configuration files, no theme system, no plugin architecture. Everything should be in Git, work with text editors, and require no setup beyond having Go installed.

None of the existing solutions met these requirements. They either required complex setup, had too many dependencies, introduced unnecessary abstractions, or were too opinionated about structure. Plus this might be a fun project for a sunny afternoon in the park. 

The implementation consists of two Go files: `main.go` (core functionality) and `data.go` (site configuration), with no external dependencies beyond the standard library. It reads Markdown files, converts them to HTML, generates an index page, creates an RSS feed, and outputs everything to a `public/` directory. The entire codebase is **under 400 lines** and does exactly what I need, nothing more.

## How it works

The blog generator follows a simple workflow:

### 1. Content Structure

Posts are Markdown files in the `articles/` directory, named with a date prefix: `YYYY-MM-DD-title.md`. The date prefix serves two purposes: it provides the publication date for sorting and RSS feeds, and it makes chronological organization obvious when browsing files.

```
articles/
  2025-07-01-hello-blog.md
  2025-12-03-x-platform-translation-system.md
```

The first line of each Markdown file is treated as the title (a `# Heading`), and the rest is converted to HTML content.

### 2. Markdown "Parsing"

The markdown parser is intentionally minimal. It handles:
- Headings (`#`, `##`, `###`)
- Paragraphs
- Lists (`- item`)
- Inline formatting (bold, italic, code, links, images)
- Code blocks with syntax highlighting classes
- Automatic anchor generation for `##` headings

The parser is a simple line-by-line regex parser, building HTML output. It might cover all markdown, but for what I am using it is more than enough. No external markdown libraries

```go
func parseMarkdown(input string) (content string, title string) {
    lines := strings.Split(input, "\n")
    // Process each line, convert to HTML
    // Extract title from first # heading
    // Return HTML content and title
}
```

### 3. Post Loading and Sorting

The `loadPosts()` function scans the articles directory, reads each `.md` file, parses the date from the filename prefix, converts Markdown to HTML, and sorts posts by date in descending order (newest first).

```go
func loadPosts(dir string) []Post {
    files, _ := os.ReadDir(dir)
    var posts []Post
    for _, f := range files {
        if strings.HasSuffix(f.Name(), ".md") {
            // Parse date from filename: YYYY-MM-DD-title.md
            dateStr := f.Name()[:10]
            postDate, err := time.Parse("2006-01-02", dateStr)
            // Convert markdown to HTML
        }
    }
    // Sort by date, newest first
    sort.Slice(posts, func(i, j int) bool {
        return posts[i].Date.After(posts[j].Date)
    })
    return posts
}
```

If a file doesn't match the expected format, it logs a warning and skips it, ensuring only properly formatted posts are included.

### 4. HTML Generation

The generator creates three types of HTML:

**Index Page** (`index.html`): Lists all posts with links, plus a links section for external resources.

**Post Pages** (`YYYY-MM-DD-title.html`): Individual post pages with navigation back to the index.

**RSS Feed** (`feed.xml`): Standard Atom feed for RSS readers.

All HTML is generated using Go's `html/template` package, which is part of the standard library. Templates are read from simple HTML files (`index.html` and `article.html`) that use Go's template syntax—no complex template system, just straightforward HTML with template variables.

### 5. Configuration

Site metadata is stored in `data.go` as a simple Go struct. This includes the site title, slogan, base URL, links, projects, and tools that appear on the index page. The configuration is just a variable declaration—no YAML, no JSON, no complex config parsing. To change the site title or add a link, you edit `data.go` directly.

```go
var config = Config{
    Title:   "][ nobloat.org",
    Slogan:  "pragmatic software minimalism",
    BaseURL: "https://nobloat.org",
    Links: map[string]string{
        "Choosing boring technology": "https://boringtechnology.club/",
        // ...
    },
    // ...
}
```

This configuration is used throughout the generator: for page titles, RSS feed metadata, sitemap URLs, and populating the index page with links and projects. Separating configuration from logic makes it easy to update site metadata without touching the core generation code.

### 6. File Watching (Optional)

For development, the generator includes a `--watch` flag that uses the `fsnotify` package to monitor the articles directory, CSS file, template files, and the generator itself. When any file changes, it automatically rebuilds the site.

```bash
go run main.go --watch
```

When you modify content, templates, or CSS, changes are detected immediately and the site rebuilds automatically—no manual intervention needed. Edit a post, see it update. Modify the HTML templates, get instant feedback. Change the stylesheet, see the new styles applied. 

This is the only external dependency ([github.com/fsnotify/fsnotify](https://github.com/fsnotify/fsnotify)), and it's only needed for the watch feature. The core build functionality requires no external packages.


## Conclusion

This blog generator does exactly what I need: converts Markdown to HTML, generates an index and RSS feed, and outputs static files. It's under 400 lines of code, uses only the go standard library for core functionality, and I understand every part of it.

It might not be suitable for someone who needs complex features like tags, categories, pagination, or theme systems. But for a simple technical blog, it's perfect. It fits the "nobloat" philosophy.

The entire codebase is very small, making it easy to read, modify, and maintain.

## Links / References

- [github.com/nobloat/blog](https://github.com/nobloat/blog)
- [nobloat.org](https://nobloat.org)
- [GitHub: nobloat](https://github.com/nobloat)
