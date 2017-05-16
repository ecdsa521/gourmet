package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/time/rate"

	yaml "gopkg.in/yaml.v2"

	"github.com/ecdsa521/torrent"
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
	Name     string
	Hash     string
	Path     string
	Size     int64
	Done     int64
	Peers    int
	Seeds    int
	Comment  string
	Piece    int //piece size
	UL       int64
	DL       int64
	Uploaded int64
	Status   string
	Activity map[string]bool
	Files    []*torrent.File
	PeerList []*torrent.Peer
	Trackers []string
}

type speed struct {
	lastTime  int64
	lastSize  int64
	lastSpeed int64
}

var ulSpeedCalc map[string]speed
var dlSpeedCalc map[string]speed
var totalSpeed map[string]int64

/*
func (g *Gourmet) fakeList(num int) []GEntry {
	tmp, _ := ioutil.ReadFile("/usr/share/dict/words")
	g._wordlist = strings.Split(string(tmp), "\n")

	data := []GEntry{}
	for i := 0; i <= num; i++ {
		seeds, peers, size, done := rand.Intn(500), rand.Intn(500), rand.Int63n(50000), rand.Int63n(50000)
		if done > size {
			done = size / (rand.Int63n(4) + 1)
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
*/

//Start the webserver and setup routes
func (g *Gourmet) Start() {
	g.router = httprouter.New()
	g.router.GET("/", g.listPage)
	g.router.GET("/config", g.configPage)
	g.router.GET("/api/list", g.apiList)
	g.router.GET("/api/start", g.apiStartDL)
	g.router.GET("/api/stop", g.apiStopDL)
	g.router.GET("/api/remove", g.apiRemove)
	g.router.GET("/api/stats", g.apiStats)
	g.router.GET("/api/announce", g.apiAnnounce)
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

func (g *Gourmet) speedCalcDL(v *torrent.Torrent) {
	var hex string = v.InfoHash().HexString()
	if _, ok := dlSpeedCalc[hex]; ok {
		lastSize, lastTime := dlSpeedCalc[hex].lastSize, dlSpeedCalc[hex].lastTime
		lastSpeed := (v.BytesCompleted() - lastSize) / ((time.Now().Unix() - lastTime) + 1)
		if lastSpeed < 5000 { //do not show network pings
			lastSpeed = 0
		}
		dlSpeedCalc[hex] = speed{
			lastSize:  v.BytesCompleted(),
			lastTime:  time.Now().Unix(),
			lastSpeed: lastSpeed,
		}
	} else {
		dlSpeedCalc[hex] = speed{
			lastSize:  v.BytesCompleted(),
			lastTime:  time.Now().Unix(),
			lastSpeed: 0,
		}
	}
	totalSpeed["DL"] = 0
	for _, v := range dlSpeedCalc {
		totalSpeed["DL"] += v.lastSpeed
	}
}
func (g *Gourmet) speedCalcUL(v *torrent.Torrent) {
	var hex string = v.InfoHash().HexString()
	if _, ok := ulSpeedCalc[hex]; ok {
		lastSize, lastTime := ulSpeedCalc[hex].lastSize, ulSpeedCalc[hex].lastTime
		lastSpeed := (v.Stats().BytesWritten - lastSize) / ((time.Now().Unix() - lastTime) + 1)
		if lastSpeed < 5000 { //do not show network pings
			lastSpeed = 0
		}
		ulSpeedCalc[hex] = speed{
			lastSize:  v.Stats().BytesWritten,
			lastTime:  time.Now().Unix(),
			lastSpeed: lastSpeed,
		}
	} else {
		ulSpeedCalc[hex] = speed{
			lastSize:  v.Stats().BytesWritten,
			lastTime:  time.Now().Unix(),
			lastSpeed: 0,
		}
	}
	totalSpeed["UL"] = 0
	for _, v := range ulSpeedCalc {
		totalSpeed["UL"] += v.lastSpeed
	}
}

func main() {
	ulSpeedCalc = map[string]speed{}
	dlSpeedCalc = map[string]speed{}
	totalSpeed = make(map[string]int64)
	gourmet := Gourmet{}
	text, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	yaml.Unmarshal([]byte(text), &gourmet.Config)

	upLimit := gourmet.Config["UL"].(int)
	downLimit := gourmet.Config["DL"].(int)
	gourmet.ClientConfig = &torrent.Config{}
	gourmet.ClientConfig.Seed = true
	if upLimit > 0 {
		gourmet.ClientConfig.UploadRateLimiter = rate.NewLimiter(rate.Limit(upLimit), upLimit*2)
	}
	if downLimit > 0 {
		gourmet.ClientConfig.DownloadRateLimiter = rate.NewLimiter(rate.Limit(downLimit), downLimit*2)
	}
	gourmet.ClientConfig.Seed = true
	rand.Seed(time.Now().UnixNano())
	gourmet.ClientConfig.ListenAddr = fmt.Sprintf(":%d", rand.Intn(65530)+1)

	gourmet.Client, err = torrent.NewClient(gourmet.ClientConfig)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Starting listener on port %s\nPeer ID: %x\n", gourmet.ClientConfig.ListenAddr, gourmet.Client.PeerID())
	trackers := [][]string{
		{"udp://tracker.opentrackr.org:1337"},
		{"udp://tracker.coppersurfer.tk:6969"},
		{"udp://tracker.leechers-paradise.org:6969"},
		{"udp://zer0day.ch:1337"},
		{"udp://explodie.org:6969"},
	}

	a, _ := gourmet.Client.AddMagnet("magnet:?xt=urn:btih:6REDNTETZGFY7FH2WLNO5QHXS4MBDIQD")
	a.AddTrackers(trackers)
	b, _ := gourmet.Client.AddMagnet("magnet:?xt=urn:btih:LEDGO2NZVVBNULSQQYI4GPL4ISALHBL3")
	b.AddTrackers(trackers)
	b.AddTrackers([][]string{{
		"udp://test.tracker.please.ignore",
	}})
	a.SetStatus("Seeding")
	b.SetStatus("Seeding")
	go gourmet.Start()
	select {}
}
