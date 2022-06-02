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
)

func getMagnet(metadata sites.Metadata, trackers []string) string {

	trackerString := ""
	for a := range trackers {
		trackerString += fmt.Sprintf("&tr=%s", trackers[a])
	}
	return fmt.Sprintf("magnet:?xt=urn:btih:%s&dn=%s%s", metadata.Hash, metadata.Name, trackerString)
}

func addToRemote(remote string, magnet string, category string) error {

	values := url.Values{"urls": {magnet}}

	if len(category) > 0 {
		values.Add("category", category)
		values.Add("autoTMM", "true")
	}

	response, err := http.PostForm(fmt.Sprintf("http://%s/api/v2/torrents/add", remote), values)

	if err != nil {
		return fmt.Errorf("error adding torrent to remote: %s", err.Error())
	}

	body, _ := ioutil.ReadAll(response.Body)

	fmt.Println(string(body))
	return nil
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

	sites.SortMetadata(metadata)
	sites.PrintMetadata(metadata)

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
		if error := addToRemote(remoteConfig.Url, magnet, category); error != nil {
			fmt.Println(error)
		}
	} else {
		fmt.Println(magnet)
	}
}

func handleConfig(flags *flag.FlagSet, args []string) {

	var url string
	var name string
	var subtitleDir string

	flags.StringVar(&url, "url", "", "qBittorrent Remote url")
	flags.StringVar(&name, "name", "", "qBittorrent Remote name")
	flags.StringVar(&subtitleDir, "subtitleDir", "", "Subtitle download directory")
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

	sites.DownloadSubtitle(config.GetSubtitleDir(), opensubs[index].Link[0].Url)
}
