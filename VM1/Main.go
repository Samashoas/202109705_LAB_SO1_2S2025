package main

import(
	"net/http"
	"encoding/json"
)

type Response struct{
	Message string `json:"message"`
}

func hola(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Hola Mundo!"})
}

func main(){
	http.HandleFunc("/hola", hola)
	http.ListenAndServe(":8080", nil)
}