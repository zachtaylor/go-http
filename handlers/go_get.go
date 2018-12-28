package handlers

import (
	"html/template"
	"net/http"
)

const goGetTpl = `<html>
	<meta name="go-import" content="{{.Host}}/{{.Package}} git https://{{.Host}}/{{.Package}}">
</html>`

type goGetData struct {
	Host    string
	Package string
}

// GoGet returns a Handler that prints the minimum go-import data to
// satisfy the `go get` tool for its' requested package
func GoGet(host string) http.Handler {
	t := template.Must(template.New("").Parse(goGetTpl))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pkg := r.RequestURI[1 : len(r.RequestURI)-len("?go-get=1")]
		t.Execute(w, goGetData{
			Host:    host,
			Package: pkg,
		})
	})
}
