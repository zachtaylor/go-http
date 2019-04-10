package handlers

import (
	"html/template"
	"net/http"
)

const goGetTpl = `<html>
	<meta name="go-import" content="{{.Host}}/{{.Package}} git https://{{.Host}}/{{.Package}}">
	<meta name="go-source" content="{{.Host}}/{{.Package}} https://{{.Host}}/{{.Package}} https://{{.Host}}/{{.Package}}/tree/master{/dir} https://{{.Host}}/{{.Package}}/tree/master{/dir}/{file}#L{line}">
</html>`

type goGetData struct {
	Host    string
	Package string
}

// GoGet returns a Handler that prints data for the go tool to find go code
// for any package requested with git at https at the host name
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
