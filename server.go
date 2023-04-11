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
	err := json.Unmarshal(reqBody, &environments)
	if err != nil {
		fmt.Println(err)
	}

	ocisStatus := ocis.RestartOcisServer(&wg, environments)

	resBody := make(map[string]string)

	if ocisStatus {
		res.WriteHeader(http.StatusOK)
		resBody["status"] = "OK"
		resBody["message"] = "oCIS environment successfully set"
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
	mux.HandleFunc("/environment", environmentHandler)
	// Todo: create rollback handler that removes the provided envs
	mux.HandleFunc("/rollback", environmentHandler)

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
