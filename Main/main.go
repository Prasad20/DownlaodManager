package main

import (
	"github.com/Prasad20/DownlaodManager/routes"
	"log"
	"net/http"
)


func homePage(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(`{"message": "ok"}`))
}

func handleRequests() {

	http.HandleFunc("/Health", homePage)
	http.HandleFunc("/Downloads", routes.Select)
	http.HandleFunc("/Downloads/",routes.DownloadStatus)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {
	handleRequests()
}
