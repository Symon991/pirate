package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
)

type Metadata struct {
	Id        int
	Name      string
	Info_hash string
	Leechers  string
	Seeders   string
	Num_files string
	Size      string
	Username  string
	Added     string
	Status    string
	Category  string
	Imdb      string
}

func getSizeString(size float64) string {

	if size > 1024*1024*1024 {
		return fmt.Sprintf("%f GB", size/1024.0/1024.0/1024.0)
	}

	if size > 1024*1024 {
		return fmt.Sprintf("%f MB", size/1024.0/1024.0)
	}

	if size > 1024 {
		return fmt.Sprintf("%f KB", size/1024.0)
	}

	return fmt.Sprintf("%f Bytes", size)
}

func getMagnet(metadata Metadata, trackers []string) string {

	trackerString := ""
	for a := range trackers {
		trackerString += fmt.Sprintf("&tr=%s", trackers[a])
	}
	return fmt.Sprintf("magnet:?xt=urn:btih:%s&dn=%s%s", metadata.Info_hash, metadata.Name, trackerString)
}

func addToRemote(remote string, magnet string, category string) {

	values := url.Values{"urls": {magnet}}

	if len(category) > 0 {
		values.Add("category", category)
		values.Add("autoTMM", "true")
	}

	response, err := http.PostForm(fmt.Sprintf("http://%s/api/v2/torrents/add", remote), values)

	if err != nil {
		fmt.Printf("Errore post: %s\n", err.Error())
		return
	}

	body, _ := ioutil.ReadAll(response.Body)

	fmt.Println(string(body))
}

func searchTorrent(search string) []Metadata {

	searchUrl := fmt.Sprintf("https://pirate-proxy.club/newapi/q.php?q=%s&cat=", search)
	fmt.Println(searchUrl)

	response, _ := http.Get(searchUrl)
	bytes, _ := ioutil.ReadAll(response.Body)

	var metadata []Metadata
	json.Unmarshal(bytes, &metadata)

	return metadata
}

func main() {

	trackers := []string{
		"udp://tracker.coppersurfer.tk:6969/announce",
		"udp://tracker.openbittorrent.com:6969/announce",
		"udp://9.rarbg.to:2710/announce",
		"udp://9.rarbg.me:2780/announce",
		"udp://9.rarbg.to:2730/announce",
		"udp://tracker.opentrackr.org:1337",
		"http://p4p.arenabg.com:1337/announce",
		"udp://tracker.torrent.eu.org:451/announce",
		"udp://tracker.tiny-vps.com:6969/announce",
		"udp://open.stealth.si:80/announce",
	}

	var search string
	var first bool
	var searchOnly bool
	var remote string
	var category string
	flag.StringVar(&search, "s", "", "Search string")
	flag.StringVar(&remote, "add", "", "Remote qBittorent")
	flag.StringVar(&category, "c", "", "Category")
	flag.BoolVar(&first, "f", false, "Get first")
	flag.BoolVar(&searchOnly, "o", false, "SearchOnly")
	flag.Parse()

	metadata := searchTorrent(search)

	sort.Slice(metadata, func(p, q int) bool {
		intP, _ := strconv.ParseInt(metadata[p].Seeders, 10, 32)
		intQ, _ := strconv.ParseInt(metadata[q].Seeders, 10, 32)
		return intP > intQ
	})

	for a := range metadata {
		size, _ := strconv.ParseFloat(metadata[a].Size, 32)
		fmt.Printf("%d - %s - %s - %s\n", a, metadata[a].Name, metadata[a].Seeders, getSizeString(size))
	}

	if searchOnly {
		return
	}

	index := 0

	if !first {
		fmt.Printf("Pick torrent: ")
		fmt.Scanf("%d", &index)
	}

	magnet := getMagnet(metadata[index], trackers)

	if len(remote) > 0 {
		addToRemote(remote, magnet, category)
	} else {
		fmt.Println(magnet)
	}
}
