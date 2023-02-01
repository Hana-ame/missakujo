package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/bigkevmcd/go-configparser"
	"github.com/valyala/fastjson"
)

func main() {

	// s := []byte(`[{"a":"1"},{"a":"2"},{"a":"3"}]`)
	// fmt.Printf("%s\n", fastjson.GetString(s, "0", "a"))
	// fmt.Printf("%s\n", fastjson.GetString(s, "1", "a"))
	// fmt.Printf("%s\n", fastjson.GetString(s, "2", "a"))
	// fmt.Printf("%s\n", fastjson.GetBytes(s, "0"))

	p, err := configparser.NewConfigParserFromFile("missakujo.ini")
	if err != nil {
		log.Println(`You must have a file named 'missakujo.ini' to set the configuration`, err)
		return
	}
	host, err := p.Get("config", "host")
	if err != nil {
		log.Println(`You must have a 'host' key in the configuration`, err)
		return
	}
	user, err := p.Get("config", "user")
	if err != nil {
		log.Println(`You must have a 'user' key in the configuration`, err)
		return
	}
	token, err := p.Get("config", "token")
	if err != nil {
		log.Println(`You must have a 'token' key in the configuration`, err)
		return
	}

	sinceRaw, err := p.Get("config", "since")
	if err != nil {
		log.Println(`'since' key is not found in the configuration`, err)
	}
	untilRaw, err := p.Get("config", "until")
	if err != nil {
		log.Println(`'until' key is not found in the configuration`, err)
	}

	renotesLessThanRaw, err := p.Get("config", "renotes_less_than")
	if err != nil {
		log.Println(`'renotes_less_than' key is not found in the configuration`, err)
	}
	renotesLessThan, err := strconv.Atoi(renotesLessThanRaw)
	if err != nil {
		log.Println(`something wrong with 'renote_less_than', use default value 99999 (it will delete all notes)`, err)
		renotesLessThan = 99999
	}

	// reactionsLessThanRaw, err := p.Get("config", "reactions_less_than")
	// if err != nil {
	// 	log.Println(`'reactions_less_than' key is not found in the configuration`, err)
	// }
	// reactionsLessThan, err := strconv.Atoi(reactionsLessThanRaw)
	// if err != nil {
	// 	log.Println(`something wrong with 'reactions_less_than', use default value 99999 (it will delete all notes)`, err)
	// 	reactionsLessThan = 99999
	// }

	deleteRenote, err := p.Get("config", "delete_renote")
	if err != nil {
		log.Println(`'delete_renote' key is not found in the configuration`, err)
	}
	deleteReply, err := p.Get("config", "delete_reply")
	if err != nil {
		log.Println(`'delete_reply' key is not found in the configuration`, err)
	}

	// since := sinceRaw
	// until := untilRaw
	const timeForm = "2006-01-02 15:04:05"

	since, err := time.Parse(timeForm, sinceRaw)
	if err != nil {
		log.Println(err)
		return
	}
	until, err := time.Parse(timeForm, untilRaw)
	if err != nil {
		log.Println(err)
		return
	}

	sinceFmt := since.Format(timeForm)
	untilFmt := until.Format(timeForm)

	fmt.Println("host : " + host)
	fmt.Println("user : " + user)
	fmt.Println("token : " + token)
	fmt.Println(sinceFmt + " to " + untilFmt)
	fmt.Printf("delete renotes less than : %d\n", renotesLessThan)
	// fmt.Printf("delete reactions less than : %d\n", reactionsLessThan)
	fmt.Printf("delete renote? : %t\n", deleteRenote == `true`)
	fmt.Printf("delete reply? : %t\n", deleteReply == `true`)

	sinceInt := since.UnixMilli()
	untilInt := until.UnixMilli()

	listNotes, _, deleteNotes := getApis(host, token, user)

	for sinceInt < untilInt {
		payload, err := listNotes(sinceInt, untilInt)

		if err != nil {
			log.Println(err)
			return
		}

		var p fastjson.Parser
		v, err := p.Parse(string(payload))
		if err != nil {
			log.Fatal(err)
		}

		i := 0
		for i < 10 {
			no := strconv.Itoa(i)
			id := v.GetStringBytes(no, "id")
			log.Println(string(id))
			if string(id) == "" {
				return
			}
			renoteCount := v.GetInt(no, "renoteCount")
			createdAtRaw := v.GetStringBytes(no, "createdAt")
			renoteId := v.GetStringBytes(no, "renoteId")
			replyId := v.GetStringBytes(no, "replyId")

			createdAt, err := time.Parse(time.RFC3339, string(createdAtRaw))
			if err != nil {
				log.Println(err)
			}
			if untilInt >= createdAt.UnixMilli() {
				untilInt = createdAt.UnixMilli() - 1
			}

			if renoteCount < renotesLessThan &&
				((string(renoteId) == "" && deleteRenote != `true`) || deleteRenote == `true`) &&
				((string(replyId) == "" && deleteReply != `true`) || deleteReply == `true`) {
				log.Println("deleting "+no+" createdAt", createdAt)
				payload, err := deleteNotes(string(id))
				for err != nil || fastjson.GetString(payload, "error", "message") != "" {
					time.Sleep(time.Second * 2)
					payload, err = deleteNotes(string(id))
				}
				time.Sleep(time.Second * 2)
			}

			i++
		}
	}

}

func getApis(
	host string,
	token string,
	userId string,
) (
	func(since, until int64) ([]byte, error),
	func(noteId string) ([]byte, error),
	func(noteId string) ([]byte, error),
) {
	listNotesEndpoint := `https://` + host + `/api/users/notes`
	showNotesEndpoint := `https://` + host + `/api/notes/show`
	deleteNotesEndpoint := `https://` + host + `/api/notes/delete`
	deleteNotes := func(noteId string) ([]byte, error) {
		data := make(map[string]any)
		data["noteId"] = noteId
		data["i"] = token
		body, err := json.Marshal(data)
		if err != nil {
			log.Println(`Marshal Error: `, err)
			return nil, err
		}
		return fetchHandler(deleteNotesEndpoint, body)
	}
	showNotes := func(noteId string) ([]byte, error) {
		data := make(map[string]any)
		data["noteId"] = noteId
		data["i"] = token
		body, err := json.Marshal(data)
		if err != nil {
			log.Println(`Marshal Error: `, err)
			return nil, err
		}
		return fetchHandler(showNotesEndpoint, body)
	}
	listNotes := func(since, until int64) ([]byte, error) {
		data := make(map[string]any)
		data["userId"] = userId
		data["sinceDate"] = since
		data["untilDate"] = until
		data["i"] = token
		body, err := json.Marshal(data)
		if err != nil {
			log.Println(`Marshal Error: `, err)
			return nil, err
		}
		return fetchHandler(listNotesEndpoint, body)
	}

	return listNotes, showNotes, deleteNotes
}

var Client *http.Client = &http.Client{}

func fetchHandler(url string, body []byte) ([]byte, error) {
	log.Println("fetch", url, string(body))
	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewReader(body),
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := Client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	reader := getPlainTextReader(resp.Body, resp.Header.Get("Content-Encoding"))
	payload, err := io.ReadAll(reader)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("fetched", string(payload))
	return payload, nil
}

func getPlainTextReader(body io.ReadCloser, encoding string) io.ReadCloser {
	switch encoding {
	case "gzip":
		reader, err := gzip.NewReader(body)
		if err != nil {
			log.Println("error decoding gzip response", reader)
			log.Println("will return raw body")
			return body
		}
		return reader
	case "br":
		reader := brotli.NewReader(body)
		if reader == nil {
			log.Println("error decoding br response", reader)
			log.Println("will return raw body")
			return body
		}
		return io.NopCloser(reader)
	default:
		return body
	}
}

func xor(a, b bool) bool {
	if a == b {
		return false
	} else {
		return true
	}
}
