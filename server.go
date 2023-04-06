package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"ociswrapper/ocis"
)

func environmentHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	reqBody, _ := io.ReadAll(req.Body)

	if len(reqBody) == 0 || !json.Valid(reqBody) {
		http.Error(res, "Bad request", http.StatusBadRequest)
		return
	}

	var environments map[string]any
	json.Unmarshal(reqBody, &environments)

	res.WriteHeader(http.StatusOK)
	res.Header().Set("Content-Type", "application/json")

	resBody := make(map[string]string)
	resBody["status"] = "OK"
	resBody["message"] = "oCIS environment successfully set"

	jsonResponse, _ := json.Marshal(resBody)
	res.Write(jsonResponse)
}

func startServer() {
	var mux = http.NewServeMux()
	mux.HandleFunc("/", http.NotFound)
	mux.HandleFunc("/environment", environmentHandler)

	httpServer := &http.Server{
		Addr:    ":5000",
		Handler: mux,
	}

	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
		// os.Exit(1)
	}
}

func main() {
	ocis.InitOcis()
	// startServer()
}
