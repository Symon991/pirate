package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/symon991/pirate/client"
	"github.com/symon991/pirate/config"
	"github.com/symon991/pirate/sites"
)

func main() {

	_, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("error loading configuration: %s", err.Error())
	}

	torrentCmd := flag.NewFlagSet("torrent", flag.ExitOnError)
	configCmd := flag.NewFlagSet("config", flag.ExitOnError)
	subtitleCmd := flag.NewFlagSet("subtitle", flag.ExitOnError)

	switch os.Args[1] {
	case "torrent":
		err = handleTorrent(torrentCmd, os.Args)
	case "config":
		err = handleConfig(configCmd, os.Args)
	case "subtitle":
		err = handleSubtitle(subtitleCmd, os.Args)
	}

	if err != nil {
		fmt.Printf("\n%s\n", err)
	}
}

func handleTorrent(flags *flag.FlagSet, args []string) error {

	var search string
	var remote string
	var category string
	var site string

	var index int64
	var page uint64 = 1
	var metadata []sites.Metadata
	var err error

	flags.StringVar(&search, "s", "", "Search string")
	flags.StringVar(&remote, "add", "", "qBittorrent Remote")
	flags.StringVar(&category, "c", "", "qBittorrent Category")
	flags.StringVar(&site, "t", "leetx", "Site")
	flags.Parse(args[2:])

	searchSite := sites.GetSearch(site)

	fmt.Printf("\nSearching \"%s\" on %s...", search, site)

main:
	for {

		metadata, err = searchSite.SearchWithPage(search, page)
		if err != nil {
			return fmt.Errorf("handleTorrent: %s", err)
		}

		fmt.Printf("\n\nShowing result for page %d:\n\n", page)

		if len(metadata) == 0 {
			fmt.Printf("No results.\n")
			return nil
		}

		sites.PrintMetadata(metadata)

		var choice string

		fmt.Printf("\nPick torrent or move to the [n]ext or [p]revious page: ")
		fmt.Scanf("%s", &choice)

		switch choice {
		case "n":
			page = page + 1
		case "p":
			page = page - 1
		default:
			index, err = strconv.ParseInt(choice, 10, 64)
			if err == nil {
				break main
			}
		}

		if page < 1 {
			page = 1
		}
	}

	magnet, err := searchSite.GetMagnet(metadata[index])
	if err != nil {
		return fmt.Errorf("get magnet: %w", err)
	}

	if len(remote) > 0 {

		fmt.Printf("\nAdding torrent to remote %s with category %s...", remote, category)

		var authCookie string
		remoteConfig, _ := config.GetRemote(remote)

		if remoteConfig.UserName != "" {
			auth, err := client.LogInRemote(remoteConfig.Url, remoteConfig.UserName, remoteConfig.Password)
			if err != nil {
				return fmt.Errorf("error adding torrent to remote: %w", err)
			}
			authCookie = auth
		}

		if err := client.AddToRemote(remoteConfig.Url, magnet, category, authCookie); err != nil {
			return fmt.Errorf("error adding torrent to remote: %w", err)
		}

	} else {

		fmt.Printf("\nMagnet link for %d - %s: \n\n%s\n\n", index, metadata[index].Name, magnet)
	}
	return nil
}

func handleConfig(flags *flag.FlagSet, args []string) error {

	var url string
	var name string
	var subtitleDir string
	var username string
	var password string
	var list bool

	flags.StringVar(&url, "url", "", "qBittorrent Remote url")
	flags.StringVar(&name, "name", "", "qBittorrent Remote name")
	flags.StringVar(&subtitleDir, "subtitleDir", "", "Subtitle download directory")
	flags.StringVar(&username, "username", "", "Username for auth")
	flags.StringVar(&password, "password", "", "Password for auth")
	flags.BoolVar(&list, "list", false, "List current remotes and subtitle directory")
	flags.Parse(args[2:])

	userConfig := config.GetConfig()

	if list {

		fmt.Printf("\nListing current remotes and subtitle directory:\n")

		for _, remote := range userConfig.Remotes {

			fmt.Printf("\n- Name: %s, Url: %s, Username: %s", remote.Name, remote.Url, remote.UserName)
		}

	} else {

		if name != "" && url != "" {
			userConfig.Remotes = append(userConfig.Remotes, config.Remote{Url: url, Name: name, UserName: username, Password: password})
		}

		if subtitleDir != "" {
			userConfig.SubtitleDir = subtitleDir
		}

		config.WriteConfig()
	}

	fmt.Printf("\n\n")

	return nil
}

func handleSubtitle(flags *flag.FlagSet, args []string) error {

	var search string
	var language string
	var index uint64

	flags.StringVar(&search, "s", "", "Search string")
	flags.StringVar(&language, "l", "ita", "Subtitle language (eng, ita)")
	flags.Parse(args[2:])

	fmt.Printf("\nSearching \"%s\" on %s...\n\n", search, "OpenSubtitles")

	opensubs := sites.SearchOpensubs(search, language)

	if len(opensubs) == 0 {
		fmt.Printf("No results.\n")
		return nil
	}

	for i := range opensubs {
		fmt.Printf("%d - %s - %s\n", i, opensubs[i].Title, opensubs[i].Release)
		fmt.Printf("\t%s - %s\n", opensubs[i].Link[0].Url, opensubs[i].Format)
	}

	fmt.Printf("\nPick subtitle: ")
	fmt.Scanf("%d", &index)

	fmt.Printf("\nDownloading subtitle into %s...", config.GetSubtitleDir())

	sites.DownloadSubtitle(config.GetSubtitleDir(), opensubs[index].Link[0].Url)

	fmt.Printf(" Done\n\n")

	return nil
}
