package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Resource struct {
	XMLName xml.Name `xml:resource"`
	Items   []Item   `xml:"item"`
}

type Item struct {
	Key   string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

func main() {
	var dir = "translations"
	var _url = "https://darkorbit-22.bpsecure.com/spacemap/templates/"

	os.RemoveAll(filepath.Join(".", dir))
	os.MkdirAll(filepath.Join(".", dir), os.ModePerm)

	languages := []string{"bg", "cs", "da", "de", "el", "en", "es", "fi", "fr", "hu", "it", "ja", "nl", "no", "pl", "pt", "ro", "ru", "sk", "sv", "tr"}

	resources := []string{"resource_eic", "resource_inventory", "resource_achievement", "resource_chat", "resource_loadingScreen", "resource_items", "flashres", "resource_quest"}

	for i := 0; i < len(languages); i++ {
		os.MkdirAll(filepath.Join(".", dir, languages[i]), os.ModePerm)

		for j := 0; j < len(resources); j++ {
			filename := fmt.Sprint(resources[j], ".php")
			trans := fmt.Sprint(resources[j], ".xml")
			file, _ := os.Create(filepath.Join(".", dir, languages[i], filename))
			defer file.Close()

			u, _ := url.Parse(_url)
			u.Path = path.Join(u.Path, languages[i], trans)

			resp, _ := http.Get(u.String())
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)

			v := new(Resource)
			xml.Unmarshal([]byte(body), &v)

			file.Write([]byte("<?php\n\nreturn [\n"))

			for _, item := range v.Items {
				line := []byte(fmt.Sprintf("	\"%s\" => \"%s\",\n", item.Key, sanitize(item.Value)))
				file.Write(line)
			}

			file.Write([]byte("];\n"))
		}
	}
}

func sanitize(text string) []byte {
	result := strings.Replace(text, `\`, `\\`, -1)
	result = strings.Replace(result, `"`, `\"`, -1)

	return []byte(fmt.Sprint(result))
}
