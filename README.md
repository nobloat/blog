# nobloat blog

A tiny static site generator that powers [nobloat.org](https://nobloat.org). It converts Markdown files in `articles/` into HTML pages, an index, a sitemap, and an Atom feed.

- Under 400 lines of Go code, with only the standard library plus `fsnotify` (for optional watch mode)
- Write posts as `YYYY-MM-DD-title.md` files; the filename prefix doubles as the publish date and the first heading becomes the page title
- Markdown "parser" that supports headings, paragraphs, lists, inline formatting, fenced code blocks with language classes, block quotes, and automatic anchors for `##` sections
- Plain HTML templates (`index.html`, `article.html`) and a single `style.css`.
- Configuration is regular Go code in `data.go`, keeping links, tools, and metadata under version control without extra or parsing

## Build & Run
1. Install Go
2. Generate the site once:
   ```bash
   go run .
   ```
   The static output is written to `public/` (e.g. open `public/index.html` in a browser).
3. Rebuild automatically while editing posts, templates, or CSS:
   ```bash
   go run . --watch
   ```
4. Deploy however you like. A simple `rsync -av --delete public/ user@host:html/blog/` keeps a remote target in sync.

For convenience you can also run `make` (build once) or `make dev` (watch mode).

To publish a new post, drop a Markdown file into `articles/`, run the build, and commit the generated `public/` files.
