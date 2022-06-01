package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

type OpensubsItem struct {
	Title       string      `xml:"title"`
	Description string      `xml:"description"`
	Link        []Enclosure `xml:"enclosure"`
}

type Opensubs struct {
	Items []OpensubsItem `xml:"channel>item"`
}

type Enclosure struct {
	Url string `xml:"url,attr"`
}

func searchOpensubs(search string, language string) []OpensubsItem {

	searchUrl := fmt.Sprintf("https://www.opensubtitles.org/en/search/sublanguageid-%s/moviename-%s/rss_2_00", language, search)
	fmt.Println(searchUrl)

	response, _ := http.Get(searchUrl)
	bytes, _ := ioutil.ReadAll(response.Body)

	var opensubs Opensubs
	xml.Unmarshal(bytes, &opensubs)

	return opensubs.Items
}
