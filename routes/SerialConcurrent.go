package routes

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

type DownloadType interface {
	Files(*Response)
	StartResponse() Response
}

type Selection struct {
	Type string   `json:"type"`
	Urls []string `json:"urls"`
}

type Serial struct {
	Urls []string `json:"urls"`
}

type Concurrent struct {
	Urls []string `json:"urls"`
	channel chan string
	bound int
}

func GenerateID(w http.ResponseWriter) Response{
	var resp Response

	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}

	resp.ID = strings.Trim(string(out), "\n")

	b, _ := json.Marshal(resp.ID)
	w.Write(b)

	return resp
}


func SerialDownload(w http.ResponseWriter,Urls []string){
	var s Serial

	s.Urls = Urls

	resp := GenerateID(w)

	s.StartResponse(&resp)

	Status[resp.ID] = resp
}

func ConDownload(w http.ResponseWriter,Urls []string){
	var c Concurrent

	c.Urls = Urls

	resp := GenerateID(w)
	resp.Status = "Downloading"
	resp.DownloadType = "Concurrent"

	go c.StartResponse(&resp)

	Status[resp.ID] = resp
}

func Select (w http.ResponseWriter, r *http.Request){

	var sel Selection

	respond, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(respond,&sel)

	switch(sel.Type){
	case "serial": SerialDownload(w,sel.Urls)
	case "concurrent": ConDownload(w,sel.Urls)
	}
}


