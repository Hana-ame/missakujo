package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Hana-ame/missakujo/backend"
	"github.com/Hana-ame/missakujo/utils"
)

func main() {
	// app := backend.App()

	// app.Listen(":3000")

	http.HandleFunc("/delete", handleDelete)
	http.HandleFunc("/webfinger/", handleWebfinger)

	err := http.ListenAndServe("127.23.0.1:8080", nil)
	fmt.Println(err)
}

func handleWebfinger(w http.ResponseWriter, r *http.Request) {
	arr := strings.Split(r.URL.Path, "/")
	acct := arr[len(arr)-1]
	userId, err := utils.ResolveUser(acct)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, err.Error())
		return
	}

	w.Write([]byte(fmt.Sprintf(`{"userId":"%s"}`, userId)))
}

func handleDelete(w http.ResponseWriter, r *http.Request) {

	resJson, err := getRequsetJson(r) //;fmt.Println(resJson)
	if err != nil {
		// log.Fatal(err)
		w.WriteHeader(400)
		fmt.Fprintf(w, err.Error())
		return
	}

	rs := backend.Wrapper(&resJson)

	w.Write([]byte(rs))
}

func getRequsetJson(r *http.Request) (req backend.DelReqCtx, err error) {

	resBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(resBody, &req)
	if err != nil {
		return
	}
	return
}
