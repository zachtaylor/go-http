package handler

import (
	"html/template"
	"net/http"
)

const goGetTpl = `<html>
	<meta name="go-import" content="{{.Host}}/{{.Package}} git git://{{.Host}}/{{.Package}}">
	<meta name="go-source" content="{{.Host}}/{{.Package}} git://{{.Host}}/{{.Package}} https://{{.Host}}/{{.Package}}/tree/master{/dir} https://{{.Host}}/{{.Package}}/tree/master{/dir}/{file}#L{line}">
</html>`

type goGetData struct {
	Host    string
	Package string
}

// GoGet returns a Handler that writes data for the go tool to find code using go get
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
