package main

import (
	"bytes"
	"fmt"
	"github.com/chrislusf/seaweedfs/weed/operation"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
)

var wg sync.WaitGroup

func main() {
	assignResult, err := operation.Assign("192.168.100.118:9333", grpc.WithInsecure(), &operation.VolumeAssignRequest{
		Count: 1,
	})
	if err != nil {
		fmt.Println(err)
	}
	//data := make([]byte, 1024)
	//rand.Read(data)
	data, _ := os.Open("./main.go")
	all, err := ioutil.ReadAll(data)
	defer data.Close()
	fmt.Println(assignResult)
	targetUrl := fmt.Sprintf("http://%s/%s", assignResult.PublicUrl, assignResult.Fid)
	fmt.Println(targetUrl)
	res, err := operation.UploadData(targetUrl, "", false, all, false, "application/txt", nil, assignResult.Auth)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func uploadFile(client *http.Client, r io.Reader) int {
	defer wg.Done()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", uuid.NewV4().String())
	if err != nil {
		fmt.Println(err)
	}
	_, err = io.Copy(part, r)
	if err != nil {
		fmt.Println(err)
	}
	writer.Close()
	request, err := http.NewRequest("POST", "http://192.168.100.118:9333/submit", body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(request)
	all, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(all))
	return len(all)
}
