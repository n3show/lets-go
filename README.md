# Let's Go

A mobile-friendly, static "Go by Example"-style reference site, generated
from plain Go source files.

## How it works

- `examples/NN-slug.go` — one Go file per topic. `// comment` lines become
  the prose explanation; everything else becomes the code block. Files are
  ordered by the `NN-` numeric prefix.
- `main.go` — the generator. Parses every file in `examples/`, renders them
  through an HTML template, and writes static pages into `docs/`.
- `static/style.css` / `static/app.js` — the site's styling and the small
  bit of JS needed for the mobile nav drawer and "Copy" buttons.
- `docs/` — the generated, ready-to-deploy static site (this is what you'd
  point GitHub Pages, Netlify, or any static host at).

## Adding a new topic

1. Create `examples/26-your-topic.go` (next number in sequence).
2. Write it like a normal commented Go file — comments become prose,
   code becomes the snippet.
3. Regenerate:

   ```
   go run main.go
   ```

That's it — the new page, its nav entry, and prev/next links are built
automatically.

## Running locally

```
go run main.go              # builds docs/
cd docs && python3 -m http.server 8000
```

Then open http://localhost:8000

## Deploying

`docs/` is fully static — drop it on GitHub Pages (point Pages at the
`docs/` folder), Netlify, Vercel, or any static file host. No build step
needed on the server side; you just regenerate `docs/` locally and commit/
upload it.

## Currently included topics

Hello World, Variables, Constants, For, If/Else, Switch, Arrays, Slices,
Maps, Range, Functions, Multiple Return Values, Variadic Functions,
Closures, Recursion, Pointers, Structs, Methods, Interfaces, Errors,
Goroutines, Channels, Select, WaitGroups, Defer.
