package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"sort"

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
	stats := make(map[string]interface{})
	stats["UL"] = int(totalSpeed["UL"])
	stats["DL"] = int(totalSpeed["DL"])

	totalPeers := 0
	totalSeeds := 0
	for _, v := range g.Client.Torrents() {
		totalPeers += v.Stats().TotalPeers
		totalSeeds += v.Stats().ActivePeers
	}
	stats["Peers"] = totalPeers
	stats["Seeds"] = totalSeeds

	stats["Trackers"] = g.getAllTrackers()
	stats["TrackersMap"] = g.getAllTrackersMap()
	stats["TrackersNo"] = len(g.getAllTrackers())
	stats["States"] = g.getAllStates()
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
	r.ParseForm()

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
				Trackers: g.getTrackers(v),
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
func (g *Gourmet) getTrackers(t *torrent.Torrent) []string {
	ret := []string{}
	data := make(map[string]bool)

	for _, val := range reflect.ValueOf(t.Metainfo().AnnounceList.DistinctValues()).MapKeys() {
		data[val.String()] = true
	}

	for _, v := range reflect.ValueOf(data).MapKeys() {
		ret = append(ret, v.String())
	}
	sort.Strings(ret)
	return ret
}
func (g *Gourmet) getAllTrackersMap() map[string]int {
	data := make(map[string]int)

	for _, v := range g.Client.Torrents() {

		for _, val := range reflect.ValueOf(v.Metainfo().AnnounceList.DistinctValues()).MapKeys() {
			data[val.String()]++
		}

	}

	return data
}
func (g *Gourmet) getAllStates() map[string]int {

	data := make(map[string]int)
	data["Stopped"] = 0
	data["Downloading"] = 0
	data["Seeding"] = 0

	for _, v := range g.Client.Torrents() {
		data[v.Status]++
	}

	return data
}

func (g *Gourmet) getAllTrackers() []string {
	ret := []string{}
	data := make(map[string]bool)
	for _, v := range g.Client.Torrents() {

		for _, val := range reflect.ValueOf(v.Metainfo().AnnounceList.DistinctValues()).MapKeys() {
			data[val.String()] = true
		}

	}
	for _, v := range reflect.ValueOf(data).MapKeys() {
		ret = append(ret, v.String())
	}
	sort.Strings(ret)
	return ret
}
