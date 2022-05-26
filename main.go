package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
)

type Metadata struct {
	Name     string
	Hash     string
	Seeders  string
	Size     string
	Category string
}

func getMagnet(metadata Metadata, trackers []string) string {

	trackerString := ""
	for a := range trackers {
		trackerString += fmt.Sprintf("&tr=%s", trackers[a])
	}
	return fmt.Sprintf("magnet:?xt=urn:btih:%s&dn=%s%s", metadata.Hash, metadata.Name, trackerString)
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

func printMetadata(metadata []Metadata) {

	for a := range metadata {
		fmt.Printf("%d - %s - %s - %s\n", a, metadata[a].Name, metadata[a].Seeders, metadata[a].Size)
	}
}

func sortMetadata(metadata []Metadata) {

	sort.Slice(metadata, func(p, q int) bool {
		intP, _ := strconv.ParseInt(metadata[p].Seeders, 10, 32)
		intQ, _ := strconv.ParseInt(metadata[q].Seeders, 10, 32)
		return intP > intQ
	})
}

func main() {

	var search string
	var first bool
	var searchOnly bool
	var remote string
	var category string
	var site string

	flag.StringVar(&search, "s", "", "Search string")
	flag.StringVar(&remote, "add", "", "qBittorrent Remote")
	flag.StringVar(&category, "c", "", "qBittorrent Category")
	flag.BoolVar(&first, "f", false, "Non-interactive mode, automatically selects first result")
	flag.BoolVar(&searchOnly, "o", false, "Search Only")
	flag.StringVar(&site, "t", "piratebay", "Site")
	flag.Parse()

	var metadata []Metadata
	var trackers []string

	switch site {
	case "nyaa":
		metadata = searchNyaa(search)
		trackers = nyaaTrackers()
	case "piratebay":
		metadata = searchTorrent(search)
		trackers = pirateBayTrackers()
	default:
		metadata = searchTorrent(search)
		trackers = pirateBayTrackers()
	}

	sortMetadata(metadata)
	printMetadata(metadata)

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
