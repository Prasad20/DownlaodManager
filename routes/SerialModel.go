package routes

import (
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func (s Serial) fileDownload(fileUrl string,Files map[string]string, resp *Response) {
	res, err := http.Get(fileUrl)

	if err !=nil {
		resp.Status="Failed"
		return
	}
	defer res.Body.Close()

	name := strings.Split(fileUrl,"/")

	filepath := "/tmp/"+name[len(name)-1]

	output, err := os.Create(filepath)
	if err != nil {
		resp.Status="Failed"
		return
	}
	defer output.Close()

	_, err = io.Copy(output, res.Body)
	if err != nil {
		resp.Status="Failed"
		return
	}

	Files[fileUrl] = filepath
}

func (s Serial) Files(resp *Response){
	resp.Files = make(map[string]string)
	for _,link:=range(s.Urls){
		_,err := http.Get(link)
		if(err!=nil){
			resp.Files[link] = "Invalid Request"
			resp.Status="Failed"
			continue
		}
		s.fileDownload(link,resp.Files,resp)
	}
}

func (s Serial) StartResponse(resp *Response) {

	resp.Status = "Downloading"
	resp.DownloadType = "Serial"

	resp.StartTime = time.Now()
	s.Files(resp)
	resp.EndTime = time.Now()

	if(resp.Status!="Failed") {
		resp.Status = "Successful"
	}
}