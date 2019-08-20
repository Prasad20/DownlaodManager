package routes

import (
	"encoding/json"
//	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)


type DownloadType interface {
	getFiles(http.ResponseWriter, *http.Request, *Response)
}

type Serial struct {
	Type string   `json:"type"`
	Urls []string `json:"urls"`
}

type Concurrent struct {
	Type string   `json:"type"`
	Urls []string `json:"urls"`
}

func (s Serial) getFiles(w http.ResponseWriter, r *http.Request, resp *Response){

	respond, _ := ioutil.ReadAll(r.Body)

	json.Unmarshal(respond,&s)

	Files := make(map[string]string)

	channel := make(chan string)

	for _,i:=range(s.Urls){
		go filesDownload(i,Files,resp,channel)
		Files[i] = <-channel
	}

	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}

	resp.ID = strings.Trim(string(out), "\n")

	resp.EndTime = time.Now()
	resp.Files = Files

	if(resp.Status=="Downloading") {
		resp.Status = "Successful"
	}

	Status[resp.ID] = *resp

	b, err := json.Marshal(resp)
	w.Write(b)
}


func (c Concurrent) getFiles(w http.ResponseWriter, r *http.Request, resp *Response){

	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}

	resp.ID = strings.Trim(string(out), "\n")

	respond, _ := ioutil.ReadAll(r.Body)

	json.Unmarshal(respond,&c)

	Files := make(map[string]string)

	b, err := json.Marshal(resp)

	w.Write(b)

	channel := make(chan string)

	for _,i:=range(c.Urls){
		go filesDownload(i,Files,resp,channel)
		Files[i] = <-channel
	}

	resp.EndTime = time.Now()
	resp.Files = Files

	if(resp.Status=="Downloading") {
		resp.Status = "Successful"
	}

	Status[resp.ID] = *resp
}

func filesDownload(fileUrl string,Files map[string]string, resp *Response,channel chan string) {

	res, err := http.Get(fileUrl)

	if err !=nil {
		resp.Status="Failed"
		channel<-"Failed"
		return
	}
	defer res.Body.Close()

	name := strings.Split(fileUrl,"/")

	filepath := "/tmp/"+name[len(name)-1]

	output, err := os.Create(filepath)
	if err != nil {
		resp.Status="Failed"
		channel<-"Failed"
		return
	}
	defer output.Close()

	_, err = io.Copy(output, res.Body)
	if err != nil {
		resp.Status="Failed"
		channel<-"Failed"
		return
	}

	 channel<-filepath
}

func SerialDownload(w http.ResponseWriter, r *http.Request){
	var s Serial

	var resp Response

	resp.StartTime = time.Now()
	resp.DownloadType = "Serial"
	resp.Status = "Downloading"

	s.getFiles(w,r,&resp)
}

func ConDownload(w http.ResponseWriter, r *http.Request){
	var c Concurrent

	var resp Response

	resp.StartTime = time.Now()
	resp.DownloadType = "Concurrent"
	resp.Status = "Downloading"

	c.getFiles(w,r,&resp)
}

