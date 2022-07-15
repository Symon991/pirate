package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/symon991/pirate/config"
	"github.com/symon991/pirate/sites"
)

func addToRemote(remote string, magnet string, category string, authCookie string) error {

	values := url.Values{"urls": {magnet}}

	if len(category) > 0 {
		values.Add("category", category)
		values.Add("autoTMM", "true")
	}

	request, err := http.NewRequest("POST", fmt.Sprintf("http://%s/api/v2/torrents/add", remote), strings.NewReader(values.Encode()))

	if err != nil {
		return fmt.Errorf("error creating request: %s", err.Error())
	}

	if authCookie != "" {
		request.AddCookie(&http.Cookie{Name: "SID", Value: authCookie})
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return fmt.Errorf("error adding torrent to remote: %s", err.Error())
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	fmt.Println(string(body))
	return nil
}

func logInRemote(remote string, username string, password string) (string, error) {

	values := url.Values{"username": {username}, "password": {password}}

	response, err := http.PostForm(fmt.Sprintf("http://%s/api/v2/auth/login", remote), values)

	if err != nil {
		return "", fmt.Errorf("error login remote: %s", err.Error())
	}

	for _, cookie := range response.Cookies() {
		if cookie.Name == "SID" {
			return cookie.Value, nil
		}
	}

	return "", fmt.Errorf("error login remote: cookie wasn't set")
}

func main() {

	torrentCmd := flag.NewFlagSet("torrent", flag.ExitOnError)
	configCmd := flag.NewFlagSet("config", flag.ExitOnError)
	subtitleCmd := flag.NewFlagSet("subtitle", flag.ExitOnError)

	switch os.Args[1] {
	case "torrent":
		handleTorrent(torrentCmd, os.Args)
	case "config":
		handleConfig(configCmd, os.Args)
	case "subtitle":
		handleSubtitle(subtitleCmd, os.Args)
	}

}

func handleTorrent(flags *flag.FlagSet, args []string) error {

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
	var err error
	err = nil

	switch site {
	case "nyaa":
		metadata = sites.SearchNyaa(search)
		trackers = sites.NyaaTrackers()
	case "piratebay":
		metadata, err = sites.SearchTorrent(search)
		trackers = sites.PirateBayTrackers()
	default:
		metadata, err = sites.SearchTorrent(search)
		trackers = sites.PirateBayTrackers()
	}

	if err != nil {
		return fmt.Errorf("handleTorrent: %s", err)
	}

	if len(metadata) == 0 {
		fmt.Printf("No results.\n")
		return nil
	}

	sites.SortMetadata(metadata)
	sites.PrintMetadata(metadata)

	if searchOnly {
		return nil
	}

	index := 0

	if !first {
		fmt.Printf("Pick torrent: ")
		fmt.Scanf("%d", &index)
	}

	magnet := sites.GetMagnet(metadata[index], trackers)

	if len(remote) > 0 {
		var authCookie string
		remoteConfig := config.GetRemote(remote)

		if remoteConfig.UserName != "" {
			auth, err := logInRemote(remoteConfig.Url, remoteConfig.UserName, remoteConfig.Password)
			if err != nil {
				fmt.Println(err.Error())
			}
			authCookie = auth
		}

		if err := addToRemote(remoteConfig.Url, magnet, category, authCookie); err != nil {
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
