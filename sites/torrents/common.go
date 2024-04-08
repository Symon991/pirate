package torrents

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/symon991/pirate/config"
)

type Metadata struct {
	Name     string
	Hash     string
	Seeders  string
	Size     string
	Category string
}

type Search interface {
	Search(search string) ([]Metadata, error)
	SearchPreset(preset string) ([]Metadata, error)
	SearchWithPage(search string, page uint64) ([]Metadata, error)
	GetMagnet(metadata Metadata) (string, error)
	PrintMetadata(metadata []Metadata)
	SortMetadata(metadata []Metadata)
	GetName() string
}

type BaseSearch struct {
	ConfigHandler config.ConfigHandler
}

func GetSearch(site string, configHandler *config.ConfigHandler) Search {

	switch site {
	case "piratebay":
		return &PirateBaySearch{
			BaseSearch: BaseSearch{
				ConfigHandler: *configHandler,
			},
		}
	case "nyaa":
		return &NyaaSearch{
			BaseSearch: BaseSearch{
				ConfigHandler: *configHandler,
			},
		}
	case "leetx":
		return &LeetxSearch{
			BaseSearch: BaseSearch{
				ConfigHandler: *configHandler,
			},
		}
	}
	return nil
}

func (s *BaseSearch) GetMagnet(metadata Metadata, trackers []string) (string, error) {

	trackerString := ""
	for a := range trackers {
		trackerString += fmt.Sprintf("&tr=%s", trackers[a])
	}
	return fmt.Sprintf("magnet:?xt=urn:btih:%s&dn=%s%s", metadata.Hash, metadata.Name, trackerString), nil
}

func (s *BaseSearch) PrintMetadata(metadata []Metadata) {

	for a := range metadata {
		fmt.Printf("%d - %s - %s - %s\n", a, metadata[a].Name, metadata[a].Seeders, metadata[a].Size)
	}
}

func (s *BaseSearch) SortMetadata(metadata []Metadata) {

	sort.Slice(metadata, func(p, q int) bool {
		intP, _ := strconv.ParseInt(metadata[p].Seeders, 10, 32)
		intQ, _ := strconv.ParseInt(metadata[q].Seeders, 10, 32)
		return intP > intQ
	})
}
