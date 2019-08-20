package routes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

var Status = make(map[string]Response)

type Response struct {
	ID           string `json:"id"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Status       string `json:"status"`
	DownloadType string `json:"download_type"`
	Files 		 map[string]string `json:"files"`
}

type Request struct{
	ID           string `json:"id"`
}

func DownloadStatus(w http.ResponseWriter, r *http.Request){
	var req Request

	respond, _ := ioutil.ReadAll(r.Body)

	json.Unmarshal(respond,&req)

	v,ok:= Status[req.ID]

	if(ok) {
		b, _ := json.Marshal(v)
		w.Write(b)
	}else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Invalid ID"))
	}
}



