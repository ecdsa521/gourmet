package main

import (
	"fmt"
	"html/template"
	"math/rand"
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
func (g *Gourmet) toolbox() []map[string]interface{} {

	data := []map[string]interface{}{
		{"Name": "Add", "Icon": "plus-sign", "Func": "tfAdd"},
		{"Name": "Del", "Icon": "minus-sign", "Func": "tfDel"},
		{"Name": "Magnet", "Icon": "magnet", "Func": "tfMagnet"},
		{"Sep": true},
		{"Name": "Start", "Icon": "play", "Func": "tfStart"},
		{"Name": "Stop", "Icon": "stop", "Func": "tfStop"},
		{"Sep": true},
	}
	return data
}

func (g *Gourmet) navbar(active string) []map[string]interface{} {

	data := []map[string]interface{}{
		{"Name": "List", "Href": "/", "Class": eq("torrents", active, "active", "")},
		{"Name": "Config", "Href": "/config", "Class": eq("config", active, "active", "")},
		{"Name": "File manager", "Href": "/fileman", "Class": eq("fileman", active, "active", "")},
	}

	return data
}

func (g *Gourmet) footer() map[string]interface{} {
	return map[string]interface{}{
		"UL": rand.Intn(2500),
		"DL": rand.Intn(6500),
	}
}
