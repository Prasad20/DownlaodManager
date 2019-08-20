package routes

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
)

func (c Concurrent) fileDownload(fileUrl string, Files map[string]string, resp *Response) {
	res, err := http.Get(fileUrl)

	if err != nil {
		resp.Status = "Failed"
		Files[fileUrl] = "Invalid Request"
		return
	}
	defer res.Body.Close()

	name := strings.Split(fileUrl, "/")

	filepath := "/tmp/" + name[len(name)-1]

	output, err := os.Create(filepath)
	if err != nil {
		resp.Status = "Failed"
		Files[fileUrl] = "Invalid Request"
		return
	}
	defer output.Close()

	_, err = io.Copy(output, res.Body)
	if err != nil {
		resp.Status = "Failed"
		Files[fileUrl] = "Invalid Request"
		return
	}

	Files[fileUrl] = filepath
}

func (c Concurrent) Files(resp *Response) {
	c.channel = make(chan string)
	resp.Files = make(map[string]string)

	for i := 0; i < c.bound; i++ {
		go func() {
			for {
				url := <-c.channel
				_, err := http.Get(url)
				if err != nil {
					resp.Files[url] = "Invalid Request"
					resp.Status = "Failed"
					return
				}
				c.fileDownload(url, resp.Files, resp)
			}
		}()
	}

	for _, link := range c.Urls {
		c.channel <- link
	}

	if resp.Status != "Failed" {
		resp.Status = "Successful"
	}

	resp.EndTime = time.Now()
}

func (c Concurrent) StartResponse(resp *Response) {

	c.bound = int(math.Min(7, float64(len(c.Urls))))

	fmt.Println(c.bound)

	resp.StartTime = time.Now()
	c.Files(resp)

	Status[resp.ID] = *resp
}
