package main

import (
	"fmt"
	"net/http"
	"net/http/cgi"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Fprint(w, `
<!doctype html>
<html>
<head>
<meta charset="utf-8">
<title>mathfent</title>
<link rel="stylesheet" href="/static/style.css">
</head>
<body>

<h3>Daily Math Puzzles</h3>

<p>
Two equal circles, both of radius 5 rest on a line and are
touching each other. A square sits between them as shown.
Find <b><i>x</i></b>.
</p>

<img src="/static/example.png" alt="diagram">

<details>
<summary>Solution</summary>
<p><i>x</i> = 2</strong>.</p>
</details>

<footer>
footer
</footer>

</body>
</html>
`)
}

func main() {
	cgi.Serve(http.HandlerFunc(handler))
}

