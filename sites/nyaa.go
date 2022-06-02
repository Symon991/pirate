package sites

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Item struct {
	Title    string `xml:"title"`
	Hash     string `xml:"infoHash"`
	Seeders  string `xml:"seeders"`
	Category string `xml:"category"`
	Size     string `xml:"size"`
}

type Nyaa struct {
	Items []Item `xml:"channel>item"`
}

func NyaaTrackers() []string {

	trackers := []string{
		"http://nyaa.tracker.wf:7777/announce",
		"udp://open.stealth.si:80/announce",
		"udp://tracker.opentrackr.org:1337/announce",
		"udp://exodus.desync.com:6969/announce",
		"udp://tracker.torrent.eu.org:451/announce",
	}

	return trackers
}

func SearchNyaa(search string) []Metadata {

	searchUrl := fmt.Sprintf(nyaaUrlTemplate, search)
	fmt.Println(searchUrl)

	response, _ := http.Get(searchUrl)
	bytes, _ := ioutil.ReadAll(response.Body)

	var nyaa Nyaa
	xml.Unmarshal(bytes, &nyaa)

	var metadata []Metadata

	for i := range nyaa.Items {

		item := nyaa.Items[i]
		metadata = append(metadata, Metadata{Name: item.Title, Hash: item.Hash, Seeders: item.Seeders, Size: item.Size})
	}

	return metadata
}
