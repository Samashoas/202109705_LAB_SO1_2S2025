package main

import (
	"net/http"
	"encoding/json"
	"time"
	"io"
)

const (
	ApiName = "API1"
	VmName = "VM1"
	Estudiante = "Juan Pablo Samayoa Ruiz"
	Carnet = "202109705"
	port = "8081"

	Api2IP = "http://127.0.0.1:8082"
	Api3IP = "http://192.168.122.114:8083"
)


type Response struct {
	Message string `json:"message"`
}

var httpClient = &http.Client{Timeout: 3 * time.Second}

func main() {
	// ENDPOINTS DE RESPUESTA DE LA API1
	http.HandleFunc("/api1/"+Carnet+"/respuesta-api2", func(w http.ResponseWriter, r *http.Request) {
		writeJson(w, Response{Message: "Hola, responde la API: " + ApiName + " en la " + VmName + ", desarrollada por el estudiante: " + Estudiante + " con carnet: " + Carnet})
	})

	http.HandleFunc("/api1/"+Carnet+"/respuesta-api3", func(w http.ResponseWriter, r *http.Request) {
		writeJson(w, Response{Message: "Hola, responde la API: " + ApiName + " en la " + VmName + ", desarrollada por el estudiante: " + Estudiante + " con carnet: " + Carnet})
	})

	// ENDPOINTS DE LLAMADA A OTRAS API
	http.HandleFunc("/api1/"+Carnet+"/llamar-api2", func(w http.ResponseWriter, r *http.Request) {
		forward(w, Api2IP+"/api2/"+Carnet+"/respuesta-api1")
	})
	http.HandleFunc("/api1/"+Carnet+"/llamar-api3", func(w http.ResponseWriter, r *http.Request) {
		forward(w, Api3IP+"/api3/"+Carnet+"/respuesta-api1")
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	http.ListenAndServe(":"+port, nil)
}

func forward(w http.ResponseWriter, url string) {
	resp, err := httpClient.Get(url)
	if err != nil {
		http.Error(w, "Error forwarding request: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
}

func writeJson(w http.ResponseWriter, response Response) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}