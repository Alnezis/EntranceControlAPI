package main

import (
	"EntranceControlAPI/api"
	"EntranceControlAPI/app"
	"EntranceControlAPI/controllers"
	"EntranceControlAPI/firebase"
	"EntranceControlAPI/token"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

func main() {
	_cors := cors.New(cors.Options{
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	})

	router := mux.NewRouter()

	//{"number_phone": "79384846350"}

	//image
	router.HandleFunc("/uploadImage", controllers.UploadFile).Methods("POST")
	router.HandleFunc("/checkFace", controllers.CheckFace).Methods("POST")

	router.PathPrefix("/images").Handler(http.StripPrefix("/images", http.FileServer(http.Dir("images/")))).Methods("GET")

	router.HandleFunc("/user/get", token.IsAuthorized(controllers.GetUser)).Methods("GET")

	router.HandleFunc("/appointments", controllers.Appointments).Methods("GET")

	go firebase.Demon()

	cert := "/etc/letsencrypt/live/alnezis.riznex.ru/fullchain.pem"
	key := "/etc/letsencrypt/live/alnezis.riznex.ru/privkey.pem"
	if _, err := os.Stat(cert); err != nil {
		if os.IsNotExist(err) {
			log.Println("no ssl")
			handler := _cors.Handler(router)
			err := http.ListenAndServe(fmt.Sprintf(":%d", app.CFG.Port), handler)
			if err != nil {
				log.Println(err)
			}
			return
		}
	}
	log.Println("yes ssl")
	handler := _cors.Handler(router)
	err := http.ListenAndServeTLS(fmt.Sprintf(":%d", app.CFG.Port), cert, key, handler)
	if err != nil {
		api.CheckErrInfo(err, "ListenAndServeTLS")
		//	return
	}
}
