package main

import (
	"fmt"
	"html"
	"net/http"
	"net/http/cgi"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	probDir = "/math.fent.uk/problems" /* inside /var/www chroot */
)

func basetop(title string) string {
	return `
<!doctype html>
<html>
<head>
<meta charset="utf-8">
<title>` + html.EscapeString(title) + `</title>
<link rel="stylesheet" href="/static/style.css">
</head>
<body>
`
}

func basebot() string {
	return `
<footer>math.fent.uk</footer>
</body>
</html>
`
}

func todayid() string {
	location, err := time.LoadLocation("Europe/London")
	if err != nil {
		location = time.UTC
	}
	now := time.Now().In(location)
	return now.Format("020106") /* DDMMYY */
}

func serveindex(w http.ResponseWriter) {
	ents, err := os.ReadDir(probDir)
	if err != nil {
		http.Error(w, "cannot read problems", 500)
		return
	}

	var files []string
	for _, e := range ents {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".html") {
			files = append(files, name)
		}
	}

	sort.Slice(files, func(i, j int) bool {
		/* newest first */
		return files[i] > files[j]
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, basetop("Problems"))
	fmt.Fprint(w, "<h3>Problems</h3>\n<ul>\n")
	for _, f := range files {
		fmt.Fprintf(w, `<li><a href="/problems/%s">%s</a></li>`+"\n",
			html.EscapeString(f),
			html.EscapeString(strings.TrimSuffix(f, ".html")),
		)
	}
	fmt.Fprint(w, "</ul>\n")
	fmt.Fprint(w, basebot())
}

func serveproblem(w http.ResponseWriter, file string) {
	/* allow only simple filenames */
	if strings.Contains(file, "/") || strings.Contains(file, "..") {
		http.NotFound(w, nil)
		return
	}
	if !strings.HasSuffix(file, ".html") {
		http.NotFound(w, nil)
		return
	}

	path := filepath.Join(probDir, file)
	b, err := os.ReadFile(path)
	if err != nil {
		http.NotFound(w, nil)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, basetop("Daily Math Puzzles"))
	fmt.Fprint(w, string(b))
	fmt.Fprint(w, basebot())
}

func handler(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path

	switch {
	case p == "/":
		http.Redirect(w, r, "/problems/"+todayid()+".html", http.StatusFound)
		return

	case p == "/problems" || p == "/problems/":
		serveindex(w)
		return

	case strings.HasPrefix(p, "/problems/"):
		serveproblem(w, strings.TrimPrefix(p, "/problems/"))
		return

	default:
		http.NotFound(w, r)
		return
	}
}

func main() {
	cgi.Serve(http.HandlerFunc(handler))
}

