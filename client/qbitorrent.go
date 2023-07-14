package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func AddToRemote(remote string, magnet string, category string, authCookie string) error {

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
	body, _ := io.ReadAll(response.Body)

	fmt.Printf("API Response: %s", string(body))
	return nil
}

func LogInRemote(remote string, username string, password string) (string, error) {

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
