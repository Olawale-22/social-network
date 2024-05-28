package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"social-network/backend/pkg/db/sqlite"
	"social-network/backend/pkg/server"

	"github.com/gorilla/mux"
)

func main() {
	sqlite.Init()
	sqlite.InitMigration()
	hub := server.NewHub()
	go hub.Run()
	r := mux.NewRouter()

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server.ServeWs(hub, w, r)
	})

	r.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		queryValues := r.URL.Query()
		studentID := queryValues.Get("studentID")
		HighDee, _ := strconv.Atoi(studentID)
		getDataHandler(w, HighDee, r)
	}).Methods("GET")

	// Gestionnaire POST pour la route /upload
	r.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		fmt.Println("ENDPOINT")
		fmt.Println(r.FormValue("privacy"))

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Erreur de lecture du corps de la requête", http.StatusInternalServerError)
			return
		}

		fmt.Println("BODY: ", body)
		file, handler, err := r.FormFile("image")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		fmt.Println("FORMFILE")
		fmt.Println(file)

		// Créer un nouveau fichier sur le serveur
		newFile, err := os.Create("./backend/pkg/server/upload/" + handler.Filename)
		if err != nil {
			http.Error(w, "Could not create file on server", http.StatusInternalServerError)
			return
		}
		defer newFile.Close()

		fmt.Println("CREATE")

		// Copier le contenu du fichier téléchargé dans le fichier sur le serveur
		_, err = io.Copy(newFile, file)
		if err != nil {
			http.Error(w, "Could not copy file to server", http.StatusInternalServerError)
			return
		}

		fmt.Println("COPY")

		fmt.Println("image: ", handler.Filename)
		http.ServeFile(w, r, "./backend/pkg/server/upload/"+handler.Filename)

		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{"message": "File uploaded successfully"}
		json.NewEncoder(w).Encode(response)
		fmt.Println("SUCCESS")
	}).Methods("POST")

	r.PathPrefix("/upload/").Handler(http.StripPrefix("/upload/", http.FileServer(http.Dir("./backend/pkg/server/upload/"))))

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	fmt.Println("Server started listenning on :8080")
}

func getDataHandler(w http.ResponseWriter, input int, r *http.Request) {
	data, err := server.FetchDataFromDB(input)
	fmt.Println("DATA FROM API: ", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set CORS headers to allow requests from specific origins
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Remove the trailing slash
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
