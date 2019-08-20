package main

import (
	"github.com/Prasad20/DownlaodManager/routes"
	"log"
	"net/http"
)


func homePage(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(`{"message": "Hello"}`))
}

func handleRequests() {

	http.HandleFunc("/Health", homePage)
	http.HandleFunc("/DownloadSerial", routes.SerialDownload)
	http.HandleFunc("/DownloadConcurrent", routes.ConDownload)
	http.HandleFunc("/Status",routes.DownloadStatus)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {
	handleRequests()
}
