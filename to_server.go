package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"missakujo/utils"

	"github.com/andybalholm/brotli"
	"github.com/google/uuid"
	"github.com/valyala/fastjson"
)

// 里面是屎山，我不调了。
// 返回sessionID，用来查看文件用
func Wrapper(ctx *DelReqCtx) string {
	u, err := uuid.NewRandom()
	if err != nil {
		return err.Error()
	}
	sessionID := u.String()

	since, err := time.Parse(timeForm, ctx.Since)
	if err != nil {
		return err.Error()
	}
	until, err := time.Parse(timeForm, ctx.Until)
	if err != nil {
		return err.Error()
	}

	offset := ctx.TimeOffset

	sinceInt := since.UnixMilli()
	untilInt := until.UnixMilli()

	sinceInt -= int64(offset) * 1000
	untilInt -= int64(offset) * 1000

	go HandleDelete(
		ctx.Host,
		ctx.User,
		ctx.Token,
		sinceInt, untilInt,
		ctx.RenoteLessThan,
		ctx.DeleteRenote, ctx.DeleteReply,
		sessionID,
	)

	return sessionID
}

func HandleDelete(
	host string,
	user string,
	token string,
	sinceInt, untilInt int64,
	renotesLessThan int,
	deleteRenote, deleteReply string,
	sessionId string,
) {

	log, err := utils.NewFileWriter(sessionId + ".txt")
	if err != nil {
		log.Log(err.Error())
		return
	}
	defer log.Close()

	// sinceInt := since.UnixMilli()
	// untilInt := until.UnixMilli()

	listNotes, _, deleteNotes := getApis(host, token, user, log)

	for sinceInt < untilInt {
		payload, err := listNotes(sinceInt, untilInt)

		if err != nil {
			log.Println(err)
			return
		}

		var p fastjson.Parser
		v, err := p.Parse(string(payload))
		if err != nil {
			log.Println(err)
			return
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
				((string(renoteId) == "" && deleteRenote != `true`) ||
					deleteRenote == `true`) &&
				((string(replyId) == "" && deleteReply != `true`) ||
					deleteReply == `true`) {

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
	log *utils.FileWriter,
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
		return fetchHandler(deleteNotesEndpoint, body, log)
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
		return fetchHandler(showNotesEndpoint, body, log)
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
		return fetchHandler(listNotesEndpoint, body, log)
	}

	return listNotes, showNotes, deleteNotes
}

var Client *http.Client = &http.Client{}

func fetchHandler(url string, body []byte, log *utils.FileWriter) ([]byte, error) {
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
