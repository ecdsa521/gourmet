package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/julienschmidt/httprouter"
)

//GFile represents one file from torrent
type GFile struct {
	Name string
	Path string
	Size int
	Done int
}

//GEntry represents one torrent in list
type GEntry struct {
	Name    string
	Hash    string
	Path    string
	Size    int
	Done    int
	Peers   int
	Seeds   int
	Comment string
	Piece   int //piece size

	Files    []GFile
	Trackers []string
}

func (g *Gourmet) fakeList(num int) []GEntry {
	tmp, _ := ioutil.ReadFile("/usr/share/dict/words")
	g._wordlist = strings.Split(string(tmp), "\n")

	data := []GEntry{}
	for i := 0; i <= num; i++ {
		data = append(data, GEntry{
			Name:  g.fakeName(),
			Size:  rand.Intn(1000000),
			Done:  rand.Intn(1000000),
			Seeds: rand.Intn(500),
			Peers: rand.Intn(500),
			Hash:  fmt.Sprintf("%x", sha1.Sum([]byte(g.fakeName()))),
			Path:  "C:/Windows/system32/" + g.fakeName(),
		})
	}
	return data
}
func (g *Gourmet) fakeName() string {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(len(g._wordlist))
	return strings.TrimSpace(g._wordlist[r])
}

//Start the webserver and setup routes
func (g *Gourmet) Start() {
	g.router = httprouter.New()
	g.router.GET("/", g.listPage)
	g.router.GET("/config", g.configPage)
	g.router.GET("/api/list", g.apiList)

	g.router.ServeFiles("/static/*filepath", http.Dir("static"))
	http.ListenAndServe(fmt.Sprintf(":%d", g.Config["Port"]), g.router)
}

func (g *Gourmet) listPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	g.genTemplate(w, r, "list.html", map[string]interface{}{
		"Config":  g.Config,
		"Navbar":  g.navbar("list"),
		"Title":   "some title",
		"Content": "some content",
	})
}
func (g *Gourmet) configPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	g.genTemplate(w, r, "config.html", map[string]interface{}{
		"Config":  g.Config,
		"Navbar":  g.navbar("config"),
		"Title":   "some title",
		"Content": "some content",
	})
}

func (g *Gourmet) apiList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	list := g.fakeList(1000)
	b, _ := json.Marshal(list)
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}
func main() {
	gourmet := Gourmet{}
	text, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	yaml.Unmarshal([]byte(text), &gourmet.Config)
	gourmet.Start()
}
