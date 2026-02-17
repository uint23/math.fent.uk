package main

import (
	"bytes"
	"fmt"
	"html"
	"net/http"
	"net/http/cgi"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/yuin/goldmark"
)

const (
	probDir = "/math.fent.uk/problems" /* inside /var/www chroot */
)

var md = goldmark.New()

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

func rendermarkdown(src []byte) (string, bool) {
	var out bytes.Buffer

	if err := md.Convert(src, &out); err != nil {
		return "", false
	}

	return out.String(), true
}

func splitmd(mdsrc string) (prob string, sol string) {
	lines := strings.Split(mdsrc, "\n")
	for i := 0; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "### Solution" {
			prob = strings.Join(lines[:i], "\n")
			sol = strings.Join(lines[i+1:], "\n")
			return
		}
	}

	return mdsrc, ""
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
		if strings.HasSuffix(name, ".md") {
			files = append(files, strings.TrimSuffix(name, ".md")+".html")
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

	mdfile := strings.TrimSuffix(file, ".html") + ".md"
	path := filepath.Join(probDir, mdfile)
	b, err := os.ReadFile(path)
	if err != nil {
		http.NotFound(w, nil)
		return
	}

	probmd, solmd := splitmd(string(b))

	probhtml, ok := rendermarkdown([]byte(probmd))
	if !ok {
		http.Error(w, "markdown error", 500)
		return
	}

	solhtml := ""
	if strings.TrimSpace(solmd) != "" {
		solhtml, ok = rendermarkdown([]byte(solmd))
		if !ok {
			http.Error(w, "markdown error", 500)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, basetop("Daily Math Puzzles"))
	fmt.Fprint(w, probhtml)
	if solhtml != "" {
		fmt.Fprint(w, "<details>\n<summary>Solution</summary>\n")
		fmt.Fprint(w, solhtml)
		fmt.Fprint(w, "\n</details>\n")
	}
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

