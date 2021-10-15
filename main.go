package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type Response struct {
	Output string `json:"output"`
}

type Request struct {
	Input string `json:"input"`
}

func scriptHandler(w http.ResponseWriter, r *http.Request) {
	// Parse JSON HTTP Request Body
	decoder := json.NewDecoder(r.Body)
	var req Request
	// Deserialize to Request Struct
	errDec := decoder.Decode(&req)
	if errDec != nil {
		panic(errDec)
	}
	// Base64 Decode Input
	inputDec, inputErr := base64.StdEncoding.DecodeString(req.Input)
	if inputErr != nil {
		panic(inputErr)
	}
	// Exec Command
	cmd := exec.CommandContext(r.Context(), "bash", "-c", string(inputDec))
	cmdOut, cmdErr := cmd.Output()
	if cmdErr != nil {
		panic(cmdErr)
	}
	// Base64 Encode Output
	outputEnc := base64.StdEncoding.EncodeToString([]byte(cmdOut))
	// Serialize to Response Struct
	res := &Response{string(outputEnc)}
	body, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	// HTTP Response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if _, err := w.Write(body); err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/", scriptHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
