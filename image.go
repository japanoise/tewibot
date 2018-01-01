package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func fetchImage(tag string) string {
	if Global.DanbooruLogin == "" || Global.DanbooruAPIKey == "" {
		return `Error: Can't do API requests without both a Danbooru Login & API key. https://danbooru.donmai.us/wiki_pages/43568`
	}
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go
	body := strings.NewReader(`tags=` + tag + ` rating:safe&limit=1&random=true`)
	req, err := http.NewRequest("GET", "https://danbooru.donmai.us/posts.json", body)
	if err != nil {
		return err.Error()
	}

	req.SetBasicAuth(Global.DanbooruLogin, Global.DanbooruAPIKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := netClient.Do(req)
	if err != nil {
		return err.Error()
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err.Error()
		}
		if len(body) < 3 {
			return "Nobody here but us chickens!"
		}
		return imageLinkForJson(body)
	} else {
		return resp.Status
	}
}

// PLEASE DO NOT USE THIS FUNCTION ANYWHERE ELSE, IT'S TERRIBLE
// Assumes:
// - if it's passed a list, it's a singleton list NOT AN EMPTY LIST
// - json data isn't empty
// These are huge assumptions. Don't blame me if you pass json that you've not checked these against into this function.
// If you really want to do something like what the bot does, tweak the method above.
func imageLinkForJson(b []byte) string {
	// Very hacky. I cba to do the whole danbooru API, so we're just throwing this in as an unstructured JSON object.
	if b[0] == '[' {
		// xtreme hack: unlistify it. Makes the rest a little easier.
		// We assume it's a singleton list.
		b[0] = ' '
		b[len(b)-1] = ' '
	}
	var obj interface{}
	err := json.Unmarshal(b, &obj)
	if err != nil {
		return err.Error()
	}
	objmap := obj.(map[string]interface{})
	fileurl := objmap["file_url"]
	if fileurl == nil {
		return "Malformed json data (nothing for file_url found)"
	}
	switch fut := fileurl.(type) {
	case string:
		return "https://danbooru.donmai.us" + fut
	default:
		return "Malformed json data (wrong type for file_url; was expecting string)"
	}
}
