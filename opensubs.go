package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type OpensubsItem struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Release     string
	Format      string
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

	items := opensubs.Items[1:10]

	regex := regexp.MustCompile("(?:Released as: ([^;]*);[\\s\\w]*)?Format: ([^;]*);")

	for i := range items {
		matches := regex.FindStringSubmatch(items[i].Description)
		items[i].Release = matches[1]
		items[i].Format = matches[2]
	}

	return items
}

func downloadSubtitle(path string, url string) {

	var filename string
	response, _ := http.Get(url)

	contentDisposition := response.Header.Get("content-disposition")
	fmt.Sscanf(contentDisposition, "attachment; filename=%s", &filename)
	filename = strings.Replace(filename, "\"", "", -1)

	dir := path + "\\" + filename

	os.MkdirAll(path, os.ModeDir)
	file, error := os.Create(dir)
	if error != nil {
		fmt.Println(error)
	}

	io.Copy(file, response.Body)
}
