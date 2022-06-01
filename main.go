package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"pirate/config"
	"pirate/sites"
	"sort"
	"strconv"
)

func getMagnet(metadata sites.Metadata, trackers []string) string {

	trackerString := ""
	for a := range trackers {
		trackerString += fmt.Sprintf("&tr=%s", trackers[a])
	}
	return fmt.Sprintf("magnet:?xt=urn:btih:%s&dn=%s%s", metadata.Hash, metadata.Name, trackerString)
}

func addToRemote(remote string, magnet string, category string) {

	values := url.Values{"urls": {magnet}}

	if len(category) > 0 {
		values.Add("category", category)
		values.Add("autoTMM", "true")
	}

	response, err := http.PostForm(fmt.Sprintf("http://%s/api/v2/torrents/add", remote), values)

	if err != nil {
		fmt.Printf("Errore post: %s\n", err.Error())
		return
	}

	body, _ := ioutil.ReadAll(response.Body)

	fmt.Println(string(body))
}

func printMetadata(metadata []sites.Metadata) {

	for a := range metadata {
		fmt.Printf("%d - %s - %s - %s\n", a, metadata[a].Name, metadata[a].Seeders, metadata[a].Size)
	}
}

func sortMetadata(metadata []sites.Metadata) {

	sort.Slice(metadata, func(p, q int) bool {
		intP, _ := strconv.ParseInt(metadata[p].Seeders, 10, 32)
		intQ, _ := strconv.ParseInt(metadata[q].Seeders, 10, 32)
		return intP > intQ
	})
}

func main() {

	torrentCmd := flag.NewFlagSet("torrent", flag.ExitOnError)
	configCmd := flag.NewFlagSet("config", flag.ExitOnError)
	subtitleCmd := flag.NewFlagSet("subtitle", flag.ExitOnError)

	flag.NewFlagSet("config", flag.ExitOnError)

	switch os.Args[1] {
	case "torrent":
		handleTorrent(torrentCmd, os.Args)
	case "config":
		handleConfig(configCmd, os.Args)
	case "subtitle":
		handleSubtitle(subtitleCmd, os.Args)
	}

}

func handleTorrent(flags *flag.FlagSet, args []string) {

	var search string
	var first bool
	var searchOnly bool
	var remote string
	var category string
	var site string

	flags.StringVar(&search, "s", "", "Search string")
	flags.StringVar(&remote, "add", "", "qBittorrent Remote")
	flags.StringVar(&category, "c", "", "qBittorrent Category")
	flags.BoolVar(&first, "f", false, "Non-interactive mode, automatically selects first result")
	flags.BoolVar(&searchOnly, "o", false, "Search Only")
	flags.StringVar(&site, "t", "piratebay", "Site")
	flags.Parse(args[2:])

	var metadata []sites.Metadata
	var trackers []string

	switch site {
	case "nyaa":
		metadata = sites.SearchNyaa(search)
		trackers = sites.NyaaTrackers()
	case "piratebay":
		metadata = sites.SearchTorrent(search)
		trackers = sites.PirateBayTrackers()
	default:
		metadata = sites.SearchTorrent(search)
		trackers = sites.PirateBayTrackers()
	}

	sortMetadata(metadata)
	printMetadata(metadata)

	if searchOnly {
		return
	}

	index := 0

	if !first {
		fmt.Printf("Pick torrent: ")
		fmt.Scanf("%d", &index)
	}

	magnet := getMagnet(metadata[index], trackers)

	if len(remote) > 0 {
		remoteConfig := config.GetRemote(remote)
		addToRemote(remoteConfig.Url, magnet, category)
	} else {
		fmt.Println(magnet)
	}
}

func handleConfig(flags *flag.FlagSet, args []string) {

	var url string
	var name string
	var subtitleDir string

	flags.StringVar(&url, "url", "", "Search string")
	flags.StringVar(&name, "name", "", "qBittorrent Remote")
	flags.StringVar(&subtitleDir, "subtitleDir", "", "Subtitle directory")
	flags.Parse(args[2:])

	userConfig := config.ReadConfig()

	if name != "" && url != "" {
		userConfig.Remotes = append(userConfig.Remotes, config.Remote{Url: url, Name: name})
	}

	if subtitleDir != "" {
		userConfig.SubtitleDir = subtitleDir
	}

	config.WriteConfig(userConfig)
}

func handleSubtitle(flags *flag.FlagSet, args []string) {

	var search string
	var language string
	var first bool
	flags.StringVar(&search, "s", "", "Search string")
	flags.StringVar(&language, "l", "", "Subtitle language (eng, ita)")
	flags.BoolVar(&first, "f", false, "Non-interactive mode, automatically selects first result")
	flags.Parse(args[2:])

	opensubs := sites.SearchOpensubs(search, language)

	for i := range opensubs {
		fmt.Printf("%d - %s - %s\n", i, opensubs[i].Title, opensubs[i].Link[0].Url)
		fmt.Printf("\t%s - %s\n", opensubs[i].Release, opensubs[i].Format)
	}

	index := 0

	if !first {
		fmt.Printf("Pick subtitle: ")
		fmt.Scanf("%d", &index)
	}

	userConfig := config.ReadConfig()

	sites.DownloadSubtitle(userConfig.SubtitleDir, opensubs[index].Link[0].Url)
}
