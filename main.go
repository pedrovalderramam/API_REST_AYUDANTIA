package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Alumno struct {
	Nombre    string    `json:"nombre"`
	Nota1     string    `json:"nota1"`
	Nota2     string    `json:"nota2"`
	Nota3     string    `json:"nota3"`
	Nota4     string    `json:"nota4"`
	Promedio  string    `json:"promedio"`
	Situacion string    `json:"situacion"`
	CreatedAt time.Time `json:"create_at"`
}

//Función que agrega promedio y situación
func promedio(nota1 float64, nota2 float64, nota3 float64, nota4 float64) float64 {

	if nota1 <= 1.0 || nota1 >= 7.0 {
		return 0.0
	}
	if nota2 <= 1.0 || nota2 >= 7.0 {
		return 0.0
	}
	if nota3 <= 1.0 || nota3 >= 7.0 {
		return 0.0
	}
	if nota4 <= 1.0 || nota4 >= 7.0 {
		return 0.0
	}

	var prom = (nota1 + nota2 + nota3 + nota4) / 4
	return prom
}

var alumnoStore = make(map[string]Alumno)
var id int

//GetAlumnoHandler - GET - /api/alumnos
func GetAlumnoHandler(w http.ResponseWriter, r *http.Request) {
	var alumnos []Alumno

	for _, V := range alumnoStore {

		var n1, n2, n3, n4 float64

		n1, _ = strconv.ParseFloat(V.Nota1, 64)
		n2, _ = strconv.ParseFloat(V.Nota2, 64)
		n3, _ = strconv.ParseFloat(V.Nota3, 64)
		n4, _ = strconv.ParseFloat(V.Nota4, 64)

		var prom = promedio(n1, n2, n3, n4)

		if prom == 0.00 {
			V.Promedio = "Error al ingresar promedio, notas invalidas"
		} else {
			V.Promedio = fmt.Sprintf("%.2f", prom)
			if V.Promedio >= "4.0" {
				V.Situacion = "APROBADO"
			} else {
				V.Situacion = "REPROBADO"
			}
		}

		alumnos = append(alumnos, V)

	}

	w.Header().Set("Content-Type", "appliacation/json")
	j, err := json.Marshal(alumnos)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

//PostAlumnoHandler - POST - /api/alumnos
func PostAlumnoHandler(w http.ResponseWriter, r *http.Request) {
	var alumno Alumno
	err := json.NewDecoder(r.Body).Decode(&alumno)
	if err != nil {
		panic(err)
	}

	alumno.CreatedAt = time.Now()
	id++
	k := strconv.Itoa(id)
	alumnoStore[k] = alumno

	w.Header().Set("Content-Type", "appliacation/json")
	j, err := json.Marshal(alumno)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

//PutAlumnoHandler - PUT - /api/alumnos
func PutAlumnoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	k := vars["id"]
	var alumno_update Alumno
	err := json.NewDecoder(r.Body).Decode(&alumno_update)
	if err != nil {
		panic(err)
	}
	if alumno, ok := alumnoStore[k]; ok {
		alumno_update.CreatedAt = alumno.CreatedAt
		delete(alumnoStore, k)
		alumnoStore[k] = alumno_update

	} else {
		log.Printf("No encontro el id %s", k)
	}
	w.WriteHeader(http.StatusNoContent)
}

//DeleteAlumnoHandler - DELETE - /api/alumnos
func DeleteAlumnoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	k := vars["id"]

	if _, ok := alumnoStore[k]; ok {

		delete(alumnoStore, k)

	} else {
		log.Printf("No encontro el id %s", k)
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/api/alumnos", GetAlumnoHandler).Methods("GET")
	r.HandleFunc("/api/alumnos", PostAlumnoHandler).Methods("POST")
	r.HandleFunc("/api/alumnos/{id}", PutAlumnoHandler).Methods("PUT")
	r.HandleFunc("/api/alumnos/{id}", DeleteAlumnoHandler).Methods("DELETE")

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
