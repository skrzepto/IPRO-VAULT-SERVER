package main

import (
    "net/http"
    "io"
    "log"
    "github.com/skrzepto/IPRO-VAULT-SERVER/handlers"
)

func hello(rw http.ResponseWriter, req *http.Request) {
    io.WriteString(rw, "Hello world!")
}

func main() {
    http.HandleFunc("/", hello)
    http.HandleFunc("/api/sensor_data", handlers.InsertNewSensorData)
    log.Fatal(http.ListenAndServe(":8082", nil))
}
