package handler

import (
	"html/template"
	"net/http"
)

var goGetTpl = `<html>
	<meta name="go-import" content="{{.Host}}/{{.Package}} git https://{{.Host}}/{{.Package}}">
</html>`

type goGetData struct {
	Host    string
	Package string
}

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
