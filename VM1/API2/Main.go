package main

import (
	"net/http"
	"encoding/json"
)

const (
	ApiName = "API2"
	VmName = "VM1"
	Estudiante = "Juan Pablo Samayoa Ruiz"
	Carnet = "202109705"
	port = "8082"
)

type Response struct {
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/api2/"+Carnet+"/respuesta-api1", func(w http.ResponseWriter, r *http.Request) {
		writeJson(w, Response{Message: "Hola, responde la API: " + ApiName + " en la " + VmName + ", desarrollada por el estudiante: " + Estudiante + " con carnet: " + Carnet})
	})

	http.HandleFunc("/api2/"+Carnet+"/respuesta-api3", func(w http.ResponseWriter, r *http.Request) {
		writeJson(w, Response{Message: "Hola, responde la API: " + ApiName + " en la " + VmName + ", desarrollada por el estudiante: " + Estudiante + " con carnet: " + Carnet})
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