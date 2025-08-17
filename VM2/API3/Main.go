package main

import (
	"net/http"
	"encoding/json"
)

const (
	ApiName = "API3"
	VmName = "VM1"
	Estudiante = "Juan Pablo Samayoa Ruiz"
	Carnet = "202109705"
	port = "8083"
)

type Response struct {
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/api3/"+Carnet+"/llamar-api1", func(w http.ResponseWriter, r *http.Request) {
		writeJson(w, Response{Message: "Hola, responde la API3.1: " + ApiName + " en la " + VmName + ", desarrollada por el estudiante: " + Estudiante + " con carnet: " + Carnet})
	})

	http.HandleFunc("/api3/"+Carnet+"/llamar-api2", func(w http.ResponseWriter, r *http.Request) {
		writeJson(w, Response{Message: "Hola, responde la API3.2: " + ApiName + " en la " + VmName + ", desarrollada por el estudiante: " + Estudiante + " con carnet: " + Carnet})
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	http.ListenAndServe(":"+port, nil)
}

func writeJson(w http.ResponseWriter, response Response) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}