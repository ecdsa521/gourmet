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

	"golang.org/x/time/rate"

	yaml "gopkg.in/yaml.v2"

	"github.com/anacrolix/torrent"
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
		seeds, peers, size, done := rand.Intn(500), rand.Intn(500), rand.Intn(50000), rand.Intn(50000)
		if done > size {
			done = size / (rand.Intn(4) + 1)
		}
		data = append(data, GEntry{
			Name:  g.fakeName(),
			Size:  size,
			Done:  done,
			Seeds: seeds,
			Peers: peers,
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
	g.router.GET("/api/start/:hash", g.apiStartDL)
	g.router.GET("/api/add/magnet", g.apiAddMagnet)
	g.router.ServeFiles("/static/*filepath", http.Dir("static"))
	http.ListenAndServe(fmt.Sprintf(":%d", g.Config["Port"]), g.router)
}

func (g *Gourmet) listPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	g.genTemplate(w, r, "list.html", map[string]interface{}{
		"Config":  g.Config,
		"Navbar":  g.navbar("torrents"),
		"Toolbox": g.toolbox(),
		"Title":   "torrent",
		"Content": "",
		"Footer":  g.footer(),
	})
}
func (g *Gourmet) configPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	g.genTemplate(w, r, "config.html", map[string]interface{}{
		"Config":  g.Config,
		"Navbar":  g.navbar("config"),
		"Toolbox": g.toolbox(),
		"Title":   "torrent",
		"Content": "",
		"Footer":  g.footer(),
	})
}
func (g *Gourmet) apiAddMagnet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()

	g.Client.AddMagnet(r.FormValue("magnet"))
	fmt.Printf("Adding magnet: %s\n", r.FormValue("magnet"))
	b, _ := json.Marshal(fmt.Sprintf("ok: %d", len(g.Client.Torrents())))
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

func (g *Gourmet) apiStartDL(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	go func() {

		for _, t := range g.Client.Torrents() {
			if t.InfoHash().HexString() == ps.ByName("hash") {
				fmt.Printf("Want to start %s\n", ps.ByName("hash"))
				<-t.GotInfo()
				t.DownloadAll()

			} else {
				fmt.Printf("Skipping %s\n", t.InfoHash().HexString())
			}
		}
	}()

	b, _ := json.Marshal("ok: " + ps.ByName("hash"))
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}
func (g *Gourmet) apiList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//	list := g.fakeList(1000)
	data := []GEntry{}
	for _, v := range g.Client.Torrents() {
		data = append(data, GEntry{
			Name:  v.Name(),
			Hash:  fmt.Sprintf("%x", v.InfoHash().AsString()),
			Size:  int(v.Length()),
			Done:  int(v.BytesCompleted()),
			Peers: v.Stats().TotalPeers,
			Seeds: v.Stats().ActivePeers,
		})
	}
	b, _ := json.Marshal(data)
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

	upLimit := gourmet.Config["UL"].(int)
	downLimit := gourmet.Config["DL"].(int)
	gourmet.ClientConfig = &torrent.Config{}
	if upLimit > 0 {
		gourmet.ClientConfig.UploadRateLimiter = rate.NewLimiter(rate.Limit(upLimit), upLimit*2)
	}
	if downLimit > 0 {
		gourmet.ClientConfig.DownloadRateLimiter = rate.NewLimiter(rate.Limit(downLimit), downLimit*2)
	}
	rand.Seed(time.Now().UnixNano())
	gourmet.ClientConfig.ListenAddr = fmt.Sprintf(":%d", rand.Intn(65530)+1)
	gourmet.Client, err = torrent.NewClient(gourmet.ClientConfig)
	if err != nil {
		panic(err)
	}
	gourmet.Client.AddMagnet("magnet:?xt=urn:btih:6REDNTETZGFY7FH2WLNO5QHXS4MBDIQD")
	gourmet.Client.AddMagnet("magnet:?xt=urn:btih:LEDGO2NZVVBNULSQQYI4GPL4ISALHBL3")

	gourmet.Start()
}
