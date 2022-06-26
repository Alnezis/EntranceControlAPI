package face

import (
	"EntranceControlAPI/api"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type Result struct {
	Count int `json:"count"`
	Faces []struct {
		Dist float64 `json:"dist,omitempty"`
		ID   string  `json:"id,omitempty"`
	} `json:"faces,omitempty"`
}

func AddFace(id, url string) (*string, error) { //curl -X POST -F "file=@j.jpeg" http://localhost:8080/faces?id=person1

	fileName := "add_" + api.RandString(10) + "_" + id + ".jpg"
	dir := "./images/"

	response, e := http.Get(url)
	api.CheckErrInfo(e, "get url")
	defer response.Body.Close()

	//open a file for writing
	file, err := os.Create(dir + fileName)
	if err != nil {

	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println("Success!")

	rsp, err := call(fmt.Sprintf("http://176.9.84.19:8080/faces?id=%s", id), "POST", dir, fileName)
	if err != nil {
		return nil, err
	}
	resp_body, _ := ioutil.ReadAll(rsp.Body)

	fmt.Println(string(resp_body))

	if response.StatusCode == 400 {
		return nil, err
	}
	var res []string
	err3 := json.Unmarshal(resp_body, &res)
	if err3 != nil {
		fmt.Println("whoops:", err3)
		return nil, err
		//outputs: whoops: <nil>
	}
	fmt.Println("faces")
	fmt.Println(res)
	f := "https://alnezis.riznex.ru:1337/images/" + fileName
	return &f, nil
}

func CheckFace(dir, fileName string) (*Result, error) {
	rsp, err := call("http://176.9.84.19:8080/", "POST", dir, fileName)
	if err != nil {
		return nil, err
	}

	resp_body, _ := ioutil.ReadAll(rsp.Body)

	var res Result
	err3 := json.Unmarshal(resp_body, &res)
	if err3 != nil {
		fmt.Println("whoops:", err3)
		//outputs: whoops: <nil>
	}

	fmt.Println(res)

	fmt.Println(fmt.Sprintf("Count: %d", res.Count))
	for _, face := range res.Faces {
		fmt.Println(fmt.Sprintf("face: " + face.ID))
	}
	return &res, nil
}

func call(urlPath, method, dir, fileName string) (*http.Response, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(dir + fileName)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		return nil, err
	}
	writer.Close()
	defer file.Close()
	req, err := http.NewRequest(method, urlPath, bytes.NewReader(body.Bytes()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rsp, _ := client.Do(req)
	fmt.Println(rsp.StatusCode)
	if rsp.StatusCode != http.StatusOK {
		log.Printf("Request failed with response code: %d", rsp.StatusCode)
	}
	return rsp, nil
}
