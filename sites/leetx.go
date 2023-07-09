package sites

import (
	"fmt"
	"net/url"

	"github.com/gocolly/colly"
	"github.com/symon991/pirate/config"
)

type LeetxSearch struct{}

func (*LeetxSearch) Search(search string) ([]Metadata, error) {

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

	c.Visit(fmt.Sprintf(config.ReadConfig().Sites.LeetxUrlTemplate+"/search/%s/1/", url.QueryEscape(search)))
	return metadata, nil
}

func (*LeetxSearch) GetMagnet(metadata Metadata) string {

	c := colly.NewCollector()
	var result string

	c.OnHTML(".torrentdown1", func(e *colly.HTMLElement) {
		result = e.Attr("href")
	})

	c.Visit(fmt.Sprintf(config.ReadConfig().Sites.LeetxUrlTemplate+"%s", metadata.Hash))

	return result
}

func (*LeetxSearch) GetName() string {
	return "1337x"
}
