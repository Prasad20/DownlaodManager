package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var Mapping map[string]DownloadType

type DownloadType interface {
	getFiles()
}

type Serial struct {
	Type string   `json:"type"`
	Urls []string `json:"urls"`
	status string
}

type Concurrent struct {
	Type string   `json:"type"`
	Urls []string `json:"urls"`
	status string
}

func (s Serial) getFiles(w http.ResponseWriter, r *http.Request){

	response, _ := ioutil.ReadAll(r.Body)

	json.Unmarshal(response,&s)

	filepath := ""

	for i:=0;i<len(s.Urls);i++ {
		name := strings.Split(s.Urls[i],"/")
		fmt.Println(name)
		if err := filesDownload(s.Urls[i],filepath+name[len(name)-1]); err != nil {
			panic(err)
			return
		}
	}

	s.status = "Downloaded"
}

func (c Concurrent) getFiles(w http.ResponseWriter, r *http.Request){

	response, _ := ioutil.ReadAll(r.Body)

	json.Unmarshal(response,&c)

	filepath := ""

	for i:=0;i<len(c.Urls);i++ {
		name := strings.Split(c.Urls[i],"/")
		fmt.Println(name)
		if err := filesDownload(c.Urls[i],filepath+name[len(name)-1]); err != nil {
			panic(err)
			return
		}
	}

	c.status = "Downloaded"
}

func filesDownload(fileUrl string, filepath string) error {

	fmt.Println(filepath)

	resp, err := http.Get(fileUrl)

	if err !=nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func SerialDownload(w http.ResponseWriter, r *http.Request){
	var s Serial
	s.status = "Failed"
	s.getFiles(w,r)
}

func ConDownload(w http.ResponseWriter, r *http.Request){
	var c Concurrent
	c.status = "Failed"
	c.getFiles(w,r)
}

