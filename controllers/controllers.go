package controllers

import (
	"EntranceControlAPI/api"
	"EntranceControlAPI/face"
	"EntranceControlAPI/firebase"
	"EntranceControlAPI/user"
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Response struct {
	Result interface{} `json:"result"`
	Error  *Error      `json:"error"`
}

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type code struct {
	NumberPhone string `json:"number_phone"`
}

type checkCode struct {
	UserID int `json:"user_id"`
	Code   string
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	id := r.Header["Id"][0]

	json.NewEncoder(w).Encode(&Response{
		Result: user.GetUser(id),
	})
}

func Appointments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	json.NewEncoder(w).Encode(&Response{
		Result: firebase.Appointments(),
	})
}

func CheckFace(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var imgByte []byte
	err := json.NewDecoder(r.Body).Decode(&imgByte)
	if err != nil {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: err.Error(),
				Code:    0,
			},
		})
		return
	}

	img, _, err := image.Decode(bytes.NewReader(imgByte))
	if err != nil {
		log.Fatalln(err)
	}
	fileName := api.RandString(15) + ".jpeg"

	directory := "./data/photos/"

	out, _ := os.Create(directory + fileName)
	defer out.Close()
	var opts jpeg.Options
	opts.Quality = 100
	err = jpeg.Encode(out, img, &opts)
	//jpeg.Encode(out, img, nil)

	if err != nil {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: err.Error(),
				Code:    0,
			},
		})
		return
	}

	result, err := face.CheckFace(directory, fileName)
	if err != nil {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: err.Error(),
				Code:    0,
			},
		})
		return
	}

	if result.Count == 0 {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: "Ошибка №1, лица на фото не обнаружено.",
				Code:    1,
			},
		})
		return
	}

	if len(result.Faces) == 0 {
		json.NewEncoder(w).Encode(&Response{
			Result: user.User{
				ID:       "",
				PhotoURL: "https://alnezis.riznex.ru:1337/images/add-user.png",
				FIO:      "Нет данных.",
			},
		})
		return
	}
	id := result.Faces[0].ID

	json.NewEncoder(w).Encode(&Response{
		Result: user.GetUser(id),
	})
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	//  Ensure our file does not exceed 5MB
	r.Body = http.MaxBytesReader(w, r.Body, 5*1024*1024)

	file, handler, err := r.FormFile("image")

	// Capture any errors that may arise
	if err != nil {
		fmt.Fprintf(w, "Error getting the file")
		fmt.Println(err)
		return
	}

	defer file.Close()

	fmt.Printf("Uploaded file name: %+v\n", handler.Filename)
	fmt.Printf("Uploaded file size %+v\n", handler.Size)
	fmt.Printf("File mime type %+v\n", handler.Header)

	// Get the file content type and access the file extension
	fileType := strings.Split(handler.Header.Get("Content-Type"), "/")[1]

	// Create the temporary file name
	fileName := fmt.Sprintf("upload-*.%s", fileType)
	// Create a temporary file with a dir folder
	tempFile, err := ioutil.TempFile("images", fileName)

	if err != nil {
		fmt.Println(err)
	}

	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	tempFile.Write(fileBytes)
	fmt.Fprintf(w, "Successfully uploaded file")
}
