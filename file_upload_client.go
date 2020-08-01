package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

// func main() {
// 	postFile("./xxx.txt", "http://localhost:8080/upload", "uploadfile")
// }

func postFile(filename string, targetURL string, formKey string) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile(formKey, filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}

	//打开文件句柄操作
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return err
	}
	defer fh.Close()

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	// 或者
	//
	//  req, err := http.NewRequest("GET","http://golang.org",nil)
	// req.Close = true
	//or do this:
	//req.Header.Add("Connection", "close")
	// resp, err := http.DefaultClient.Do(req)
	//
	// or
	//
	// tr := &http.Transport{DisableKeepAlives: true}
	// client := &http.Client{Transport: tr}
	// resp, err := client.Get("http://golang.org")
	//
	resp, err := http.Post(targetURL, contentType, bodyBuf)
	if err != nil {
		return err
	}
	defer resp.Body.Close() // 要在 判断 err 之后 防止 resp 为空报错
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(resp.Status)
	fmt.Println(string(respBody))
	return nil
}
