package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Link struct {
	HRef       string             `json:"href,omitempty"`
	Type       string             `json:"type,omitempty"`
	Rel        string             `json:"rel,omitempty"`
	Template   string             `json:"template,omitempty"`
	Properties map[string]*string `json:"properties,omitempty"`
	Titles     map[string]string  `json:"titles,omitempty"`
}

type Resource struct {
	Subject    string            `json:"subject,omitempty"`
	Aliases    []string          `json:"aliases,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
	Links      []Link            `json:"links"`
}

func ResolveUser(name string) (string, error) {
	name = strings.TrimPrefix(name, "@")
	arr := strings.Split(name, "@")
	if len(arr) != 2 {
		return "", fmt.Errorf("should be somewhat like user@example.com")
	}
	username, host := arr[0], arr[1]

	resp, err := http.Get("https://" + host + "/.well-known/webfinger?resource=acct:" + username + "@" + host)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res *Resource
	json.NewDecoder(resp.Body).Decode(&res)

	log.Println(res)

	for _, link := range res.Links {
		if link.Rel == "self" {
			href := link.HRef
			arr := strings.Split(href, "/")
			r := arr[len(arr)-1]
			return r, nil
		}
	}

	return "", fmt.Errorf("no self link found")

}
