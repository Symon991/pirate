package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/symon991/pirate/client"
	"github.com/symon991/pirate/config"
	"github.com/symon991/pirate/sites"
)

func main() {

	torrentCmd := flag.NewFlagSet("torrent", flag.ExitOnError)
	configCmd := flag.NewFlagSet("config", flag.ExitOnError)
	subtitleCmd := flag.NewFlagSet("subtitle", flag.ExitOnError)
	var err error

	switch os.Args[1] {
	case "torrent":
		err = handleTorrent(torrentCmd, os.Args)
	case "config":
		handleConfig(configCmd, os.Args)
	case "subtitle":
		handleSubtitle(subtitleCmd, os.Args)
	}

	if err != nil {
		fmt.Println(err)
	}
}

func handleTorrent(flags *flag.FlagSet, args []string) error {

	var search string
	var first bool
	var searchOnly bool
	var remote string
	var category string
	var site string
	var page uint64

	flags.StringVar(&search, "s", "", "Search string")
	flags.StringVar(&remote, "add", "", "qBittorrent Remote")
	flags.StringVar(&category, "c", "", "qBittorrent Category")
	flags.BoolVar(&first, "f", false, "Non-interactive mode, automatically selects first result")
	flags.BoolVar(&searchOnly, "o", false, "Search Only")
	flags.StringVar(&site, "t", "piratebay", "Site")
	flags.Uint64Var(&page, "p", 1, "Page")
	flags.Parse(args[2:])

	searchSite := sites.GetSearch(site)
	metadata, err := searchSite.SearchWithPage(search, page)
	if err != nil {
		return fmt.Errorf("handleTorrent: %s", err)
	}

	if len(metadata) == 0 {
		fmt.Printf("No results.\n")
		return nil
	}

	sites.PrintMetadata(metadata)

	if searchOnly {
		return nil
	}

	index := 0

	if !first {
		fmt.Printf("Pick torrent: ")
		fmt.Scanf("%d", &index)
	}

	magnet := searchSite.GetMagnet(metadata[index])

	if len(remote) > 0 {
		var authCookie string
		remoteConfig := config.GetRemote(remote)

		if remoteConfig.UserName != "" {
			auth, err := client.LogInRemote(remoteConfig.Url, remoteConfig.UserName, remoteConfig.Password)
			if err != nil {
				fmt.Println(err.Error())
			}
			authCookie = auth
		}

		if err := client.AddToRemote(remoteConfig.Url, magnet, category, authCookie); err != nil {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println(magnet)
	}
	return nil
}

func handleConfig(flags *flag.FlagSet, args []string) {

	var url string
	var name string
	var subtitleDir string
	var username string
	var password string

	flags.StringVar(&url, "url", "", "qBittorrent Remote url")
	flags.StringVar(&name, "name", "", "qBittorrent Remote name")
	flags.StringVar(&subtitleDir, "subtitleDir", "", "Subtitle download directory")
	flags.StringVar(&username, "username", "", "Username for auth")
	flags.StringVar(&password, "password", "", "Password for auth")
	flags.Parse(args[2:])

	userConfig := config.ReadConfig()

	if name != "" && url != "" {
		userConfig.Remotes = append(userConfig.Remotes, config.Remote{Url: url, Name: name, UserName: username, Password: password})
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

	if len(opensubs) == 0 {
		fmt.Printf("No results.\n")
		return
	}

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
