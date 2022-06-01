package sites

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type PirateBayMetadata struct {
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

func PirateBayTrackers() []string {

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

	return trackers
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

func SearchTorrent(search string) []Metadata {

	searchUrl := fmt.Sprintf("https://pirate-proxy.club/newapi/q.php?q=%s&cat=", search)
	fmt.Println(searchUrl)

	response, _ := http.Get(searchUrl)
	bytes, _ := ioutil.ReadAll(response.Body)

	var pirateBayMetadata []PirateBayMetadata
	json.Unmarshal(bytes, &pirateBayMetadata)

	var metadata []Metadata

	for i := range pirateBayMetadata {

		pirateBay := pirateBayMetadata[i]
		sizeFloat, _ := strconv.ParseFloat(pirateBay.Size, 64)
		size := getSizeString(sizeFloat)
		metadata = append(metadata, Metadata{Name: pirateBay.Name, Hash: pirateBay.Info_hash, Seeders: pirateBay.Seeders, Size: size})
	}

	return metadata
}
