package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"ociswrapper/ocis"
	"sync"
)

var wg sync.WaitGroup

var httpServer = &http.Server{
	Addr: ":5000",
}

func configHandler(res http.ResponseWriter, req *http.Request) {
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
	err := json.Unmarshal(reqBody, &environments)
	if err != nil {
		fmt.Println(err)
	}

	ocisStatus := ocis.RestartOcisServer(&wg, environments)

	resBody := make(map[string]string)

	if ocisStatus {
		res.WriteHeader(http.StatusOK)
		resBody["status"] = "OK"
		resBody["message"] = "oCIS server is running with new configurations"
	} else {
		res.WriteHeader(http.StatusInternalServerError)
		resBody["status"] = "ERROR"
		resBody["message"] = "Internal server error"
	}
	res.Header().Set("Content-Type", "application/json")

	jsonResponse, _ := json.Marshal(resBody)
	res.Write(jsonResponse)
}

func rollbackHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ocisStatus := ocis.RestartOcisServer(&wg, nil)

	resBody := make(map[string]string)

	if ocisStatus {
		res.WriteHeader(http.StatusOK)
		resBody["status"] = "OK"
		resBody["message"] = "oCIS server is running"
	} else {
		res.WriteHeader(http.StatusInternalServerError)
		resBody["status"] = "ERROR"
		resBody["message"] = "Internal server error"
	}
	res.Header().Set("Content-Type", "application/json")

	jsonResponse, _ := json.Marshal(resBody)
	res.Write(jsonResponse)
}

func startServer(wg *sync.WaitGroup) {
	defer wg.Done()

	var mux = http.NewServeMux()
	mux.HandleFunc("/", http.NotFound)
	mux.HandleFunc("/config", configHandler)
	mux.HandleFunc("/rollback", rollbackHandler)

	httpServer.Handler = mux

	err := httpServer.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func serve() {
	wg.Add(1)
	go ocis.StartOcis(&wg, nil)
	wg.Add(1)
	go startServer(&wg)
	wg.Wait()
}

func main() {
	out, err := ocis.InitOcis()
	if err != "" {
		panic(err)
	}
	fmt.Println(out)

	serve()
}
