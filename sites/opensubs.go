package sites

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/symon991/pirate/config"
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

func SearchOpensubs(search string, language string) []OpensubsItem {

	searchUrl := fmt.Sprintf(config.ReadConfig().Sites.OpensubtitlesUrlTemplate, language, search)
	fmt.Println(searchUrl)

	response, _ := http.Get(searchUrl)
	bytes, _ := io.ReadAll(response.Body)

	var opensubs Opensubs
	xml.Unmarshal(bytes, &opensubs)

	items := opensubs.Items

	regex := regexp.MustCompile(`(?:Released as: ([^;]*);[\s\w]*)?Format: ([^;]*);`)

	for i := range items {
		matches := regex.FindStringSubmatch(items[i].Description)
		items[i].Release = matches[1]
		items[i].Format = matches[2]
	}

	return items
}

func DownloadSubtitle(path string, url string) {

	var filename string
	response, _ := http.Get(url)

	contentDisposition := response.Header.Get("content-disposition")
	fmt.Sscanf(contentDisposition, "attachment; filename=%s", &filename)
	filename = strings.Replace(filename, "\"", "", -1)

	dir := filepath.Join(path, filename)

	os.MkdirAll(path, os.ModeDir)
	file, error := os.Create(dir)
	if error != nil {
		fmt.Println(error)
	}

	io.Copy(file, response.Body)
}
