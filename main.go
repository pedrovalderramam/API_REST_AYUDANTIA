package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Alumno struct {
	Nombre     string  `json:"nombre"`
	Nota1      float32 `json:"nota1"`
	Nota2      float32 `json:"nota2"`
	Nota3      float32 `json:"nota3"`
	Nota4      float32 `json:"nota4"`
	Promedio   float32 `json:"promedio"`
	Siatuacion string  `json:"situacion"`
}

var alumnoStore = make(map[string]Alumno)
var id int

func main() {
	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("api/alumnos", GetAlumnoHandler).Methods("GET")
	r.HandleFunc("api/alumnos", PostAlumnoHandler).Methods("POST")
	r.HandleFunc("api/alumnos/{id}", PutAlumnoHandler).Methods("PUT")
	r.HandleFunc("api/alumnos/{id}", DeleteAlumnoHandler).Methods("DELETE")

	server := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Println("Listening http://localhost:8080 ...")
	server.ListenAndServe()

}
