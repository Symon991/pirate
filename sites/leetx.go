package sites

import (
	"fmt"
	"net/url"

	"github.com/gocolly/colly"
	"github.com/symon991/pirate/config"
)

type LeetxSearch struct{}

func (s *LeetxSearch) Search(search string) ([]Metadata, error) {

	return s.SearchWithPage(search, 1)
}

func (*LeetxSearch) SearchWithPage(search string, page uint64) ([]Metadata, error) {

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

	url := fmt.Sprintf(config.GetConfig().Sites.LeetxUrlTemplate+"/search/%s/%d/", url.QueryEscape(search), page)

	err := c.Visit(url)
	if err != nil {
		return nil, fmt.Errorf("visiting url %s: %w", url, err)
	}
	return metadata, nil
}

func (*LeetxSearch) GetMagnet(metadata Metadata) (string, error) {

	c := colly.NewCollector()
	var result string

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if e.Text == "Magnet Download" {
			result = e.Attr("href")
		}
	})

	url := fmt.Sprintf(config.GetConfig().Sites.LeetxUrlTemplate+"%s", metadata.Hash)

	err := c.Visit(url)

	if err != nil {
		return "", fmt.Errorf("visiting url %s: %w", url, err)
	}

	return result, nil
}

func (*LeetxSearch) GetName() string {
	return "1337x"
}
