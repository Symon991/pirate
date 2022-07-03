package sites

import (
	"fmt"
	"sort"
	"strconv"
)

const nyaaUrlTemplate = "https://nyaa.si/?page=rss&q=%s"
const opensubtitlesUrlTemplate = "https://www.opensubtitles.org/en/search/sublanguageid-%s/moviename-%s/rss_2_00"
const pirateBayUrlTemplate = "https://pirate-proxy.club/newapi/q.php?q=%s&cat="

type Metadata struct {
	Name     string
	Hash     string
	Seeders  string
	Size     string
	Category string
}

func GetMagnet(metadata Metadata, trackers []string) string {

	trackerString := ""
	for a := range trackers {
		trackerString += fmt.Sprintf("&tr=%s", trackers[a])
	}
	return fmt.Sprintf("magnet://?xt=urn:btih:%s&dn=%s%s", metadata.Hash, metadata.Name, trackerString)
}

func PrintMetadata(metadata []Metadata) {

	for a := range metadata {
		fmt.Printf("%d - %s - %s - %s\n", a, metadata[a].Name, metadata[a].Seeders, metadata[a].Size)
	}
}

func SortMetadata(metadata []Metadata) {

	sort.Slice(metadata, func(p, q int) bool {
		intP, _ := strconv.ParseInt(metadata[p].Seeders, 10, 32)
		intQ, _ := strconv.ParseInt(metadata[q].Seeders, 10, 32)
		return intP > intQ
	})
}
