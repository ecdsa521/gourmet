package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/ecdsa521/torrent"
	"github.com/julienschmidt/httprouter"
)

func (g *Gourmet) apiDetails(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	data := make(map[string]interface{})

	b, _ := json.Marshal(data)
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)

}

func (g *Gourmet) apiStats(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	stats := map[string]int{}
	stats["UL"] = int(totalSpeed["UL"])
	stats["DL"] = int(totalSpeed["DL"])
	stats["Peers"] = 0
	stats["Seeds"] = 0

	for _, v := range g.Client.Torrents() {
		stats["Peers"] += v.Stats().TotalPeers
		stats["Seeds"] += v.Stats().ActivePeers
	}

	b, _ := json.Marshal(stats)
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

func (g *Gourmet) apiRemove(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	go func() {
		fmt.Printf("Want to remove %s\n", r.FormValue("hash"))
		t, succ := g.getTorrent(r.FormValue("hash"))
		if succ {
			<-t.GotInfo()
			t.Drop()
		}
	}()

	b, _ := json.Marshal("ok: " + ps.ByName("hash"))
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}
func (g *Gourmet) apiAnnounce(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	r.ParseForm()
	go func() {
		fmt.Printf("Want to announce %s\n", r.FormValue("hash"))
		t, succ := g.getTorrent(r.FormValue("hash"))
		if succ {
			<-t.GotInfo()
			t.Announce()

		}
	}()

	b, _ := json.Marshal("ok: " + ps.ByName("hash"))
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}
func (g *Gourmet) apiStopDL(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	go func() {
		fmt.Printf("Want to stop %s\n", r.FormValue("hash"))
		t, succ := g.getTorrent(r.FormValue("hash"))
		if succ {
			<-t.GotInfo()
			//t.Announce()
			t.Close()
			t.SetStatus("Stopped")
		}
	}()

	b, _ := json.Marshal("ok: " + ps.ByName("hash"))
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}
func (g *Gourmet) apiStartDL(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	go func() {

		fmt.Printf("Want to start %s\n", r.FormValue("hash"))
		t, succ := g.getTorrent(r.FormValue("hash"))
		if succ {
			<-t.GotInfo()

			t.Reopen()
			t.DownloadAll()
			t.SetStatus("Started")
		}
	}()

	b, _ := json.Marshal("ok: " + ps.ByName("hash"))
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

func (g *Gourmet) apiList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	//	list := g.fakeList(1000)
	data := []GEntry{}
	//data := g.fakeList(5000)
	for _, v := range g.Client.Torrents() {
		var hex string = v.InfoHash().HexString()
		g.speedCalcDL(v)
		g.speedCalcUL(v)
		if v.Info() != nil {
			//	fmt.Printf("%s: %v\n", v.InfoHash().HexString(), v.Activity())
			s := v.Activity()
			if s["closed"] {
				v.SetStatus("Stopped")
			}
			if s["seeding"] {
				v.SetStatus("Seeding")
			}
			if s["needData"] {
				v.SetStatus("Downloading")
			}
			data = append(data, GEntry{
				Name:     v.Name(),
				Hash:     v.InfoHash().HexString(),
				Size:     v.Length(),
				Done:     v.BytesCompleted(),
				Peers:    v.Stats().TotalPeers,
				Seeds:    v.Stats().ActivePeers,
				Uploaded: v.Stats().BytesWritten,
				Trackers: v.Metainfo().AnnounceList,
				Status:   v.Status,
				Activity: v.Activity(),
				UL:       ulSpeedCalc[hex].lastSpeed,
				DL:       dlSpeedCalc[hex].lastSpeed,
			})
		}

	}
	b, _ := json.Marshal(data)
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}
func (g *Gourmet) apiAddMagnet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()

	t, err := g.Client.AddMagnet(r.FormValue("magnet"))
	if err == nil {
		if r.FormValue("autostart") == "on" {
			<-t.GotInfo()
			t.Reopen()
			t.DownloadAll()
			t.SetStatus("Downloading")

		} else {
			t.SetStatus("Stopped")
		}
	}
	fmt.Printf("Adding magnet: %s\n", r.FormValue("magnet"))
	b, _ := json.Marshal(fmt.Sprintf("ok: %d", len(g.Client.Torrents())))
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}
func (g *Gourmet) getTorrent(hash string) (*torrent.Torrent, bool) {

	var ih metainfo.Hash
	hex.Decode(ih[:], []byte(hash))

	return g.Client.Torrent(ih)

}
