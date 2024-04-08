package torrents

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
)

type LeetxSearch struct {
	BaseSearch
}

func (s *LeetxSearch) SearchPreset(preset string) ([]Metadata, error) {

	c := colly.NewCollector()

	var metadata []Metadata

	c.OnHTML("table tbody tr", func(e *colly.HTMLElement) {
		metadata = append(metadata, Metadata{
			Name:    e.ChildText(".name a:nth-of-type(2)"),
			Hash:    e.ChildAttr(".name a:nth-of-type(2)", "href"),
			Seeders: e.ChildText(".seeds"),
			Size:    e.ChildText(".size"),
		})
	})

	url := s.ConfigHandler.Config.Sites.LeetxUrlTemplate + "/top-100"

	err := c.Visit(url)
	if err != nil {
		return nil, fmt.Errorf("visiting url %s: %w", url, err)
	}

	return metadata, nil
}

func (s *LeetxSearch) Search(search string) ([]Metadata, error) {

	return s.SearchWithPage(search, 1)
}

func (s *LeetxSearch) SearchWithPage(search string, page uint64) ([]Metadata, error) {

	c := colly.NewCollector()

	var metadata []Metadata

	c.OnHTML("table tbody tr", func(e *colly.HTMLElement) {
		metadata = append(metadata, Metadata{
			Name:    e.ChildText(".name a:nth-of-type(2)"),
			Hash:    e.ChildAttr(".name a:nth-of-type(2)", "href"),
			Seeders: e.ChildText(".seeds"),
			Size:    strings.Split(e.ChildText(".size"), "B")[0] + "B",
		})
	})

	url := fmt.Sprintf(s.ConfigHandler.Config.Sites.LeetxUrlTemplate+"/search/%s/%d/", url.QueryEscape(search), page)

	err := c.Visit(url)
	if err != nil {
		return nil, fmt.Errorf("visiting url %s: %w", url, err)
	}
	return metadata, nil
}

func (s *LeetxSearch) GetMagnet(metadata Metadata) (string, error) {

	c := colly.NewCollector()
	var result string

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if e.Text == "Magnet Download" {
			result = e.Attr("href")
		}
	})

	url := fmt.Sprintf(s.ConfigHandler.Config.Sites.LeetxUrlTemplate+"%s", metadata.Hash)

	err := c.Visit(url)

	if err != nil {
		return "", fmt.Errorf("visiting url %s: %w", url, err)
	}

	return result, nil
}

func (*LeetxSearch) GetName() string {
	return "1337x"
}
