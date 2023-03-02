package helpep

import (
	"html/template"
	"net/http"
	"bytes"
)

const tplPath = "templates/"

func Templating(w http.ResponseWriter, tmplName string, layout string, args ...any) {
	buf := &bytes.Buffer{}

	tmpl, err := template.New(layout).ParseFiles(tplPath+tmplName+".html", tplPath+layout+".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(args) == 0 {
		err = tmpl.Execute(buf, nil)
	} else {
		err = tmpl.Execute(buf, args[0])
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf.WriteTo(w)
}