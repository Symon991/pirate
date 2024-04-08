package torrents

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type PirateBayMetadata struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Info_hash string `json:"info_hash"`
	Leechers  string `json:"leechers"`
	Seeders   string `json:"seeders"`
	Num_files string `json:"num_files"`
	Size      string `json:"size"`
	Username  string `json:"username"`
	Added     string `json:"added"`
	Status    string `json:"status"`
	Category  string `json:"category"`
	Imdb      string `json:"imdb"`
}

type PirateBaySearch struct {
	BaseSearch
}

func (s *PirateBaySearch) SearchPreset(preset string) ([]Metadata, error) {

	return make([]Metadata, 0), nil
}

func (p PirateBaySearch) Search(search string) ([]Metadata, error) {

	return p.SearchWithPage(search, 1)
}

func (p PirateBaySearch) SearchWithPage(search string, page uint64) ([]Metadata, error) {

	searchUrl := fmt.Sprintf(p.ConfigHandler.Config.Sites.PirateBayUrlTemplate, search)
	fmt.Println(searchUrl)

	response, err := http.Get(searchUrl)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("api get: %w", err)
	}

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("api read response: %w", err)
	}

	var pirateBayMetadata []PirateBayMetadata
	err = json.Unmarshal(bytes, &pirateBayMetadata)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	var metadata []Metadata

	for i := range pirateBayMetadata {

		pirateBay := pirateBayMetadata[i]
		sizeFloat, err := strconv.ParseFloat(pirateBay.Size, 64)
		if err != nil {
			return nil, fmt.Errorf("size conversion index %d: %w", i, err)
		}
		size := getSizeString(sizeFloat)
		metadata = append(metadata, Metadata{Name: pirateBay.Name, Hash: pirateBay.Info_hash, Seeders: pirateBay.Seeders, Size: size})
	}

	return metadata, nil
}

func (p PirateBaySearch) GetName() string {
	return "Pirate Bay"
}

func (p PirateBaySearch) GetMagnet(metadata Metadata) (string, error) {

	return p.BaseSearch.GetMagnet(metadata, pirateBayTrackers())
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

func pirateBayTrackers() []string {

	return []string{
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
}
