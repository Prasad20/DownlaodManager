package routes

import (
	"encoding/json"
	"net/http"
	"strings"
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

func parse(path string) string{
	id := strings.Split(path, string('/'))

	return id[len(id)-1]
}

func DownloadStatus(w http.ResponseWriter, r *http.Request){
	req := parse(r.URL.Path)

	v,ok:= Status[req]

	if(ok) {
		b, _ := json.Marshal(v)
		w.Write(b)
	}else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Invalid ID"))
	}
}



