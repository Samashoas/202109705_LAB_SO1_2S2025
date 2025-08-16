package main

import (
	"net/http"
	"encoding/json"
)

const (
	ApiName = "API1"
	VmName = "VM1"
	Estudiante = "Juan Pablo Samayoa Ruiz"
	Carnet = "202109705"
	port = "8080"
)

type Response struct {
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/api1/"+Carnet+"/llamar-api2", func(w http.ResponseWriter, r *http.Request) {
		writeJson(w, Response{Message: "Hola, responde la API2: " + ApiName + " en la " + VmName + ", desarrollada por el estudiante: " + Estudiante + " con carnet: " + Carnet})
	})

	http.HandleFunc("/api1/"+Carnet+"/llamar-api3", func(w http.ResponseWriter, r *http.Request) {
		writeJson(w, Response{Message: "Hola, responde la API3: " + ApiName + " en la " + VmName + ", desarrollada por el estudiante: " + Estudiante + " con carnet: " + Carnet})
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