package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/russross/blackfriday"
)

//Gourmet ...
type Gourmet struct {
	Config    map[string]interface{}
	router    *httprouter.Router
	_wordlist []string
}

func eq(d1 interface{}, d2 interface{}, s string, e string) string {
	if d1 == d2 {
		return s
	}
	return e

}
func neq(d1 interface{}, d2 interface{}, s string, e string) string {
	if d1 != d2 {
		return s
	}

	return e

}
func markDowner(args ...interface{}) template.HTML {
	s := blackfriday.MarkdownCommon([]byte(fmt.Sprintf("%s", args...)))
	return template.HTML(s)
}

func (g *Gourmet) genTemplate(w http.ResponseWriter, r *http.Request, file string, data map[string]interface{}) {
	var t *template.Template
	files := []string{
		fmt.Sprintf("theme/%s/layout.html", g.Config["Theme"]),
		fmt.Sprintf("theme/%s/navbar.html", g.Config["Theme"]),
		fmt.Sprintf("theme/%s/sidebar.html", g.Config["Theme"]),
		fmt.Sprintf("theme/%s/%s", g.Config["Theme"], file),
	}

	t = template.Must(template.New("layout.html").Funcs(template.FuncMap{
		"markDown": markDowner}).ParseFiles(files...))

	t.ExecuteTemplate(w, "layout", data)

}

func (g *Gourmet) navbar(active string) []map[string]interface{} {

	data := []map[string]interface{}{
		{"Name": "List", "Href": "/", "Class": eq("list", active, "active", "")},
		{"Name": "Config", "Href": "/config", "Class": eq("config", active, "active", "")},
	}

	return data
}
