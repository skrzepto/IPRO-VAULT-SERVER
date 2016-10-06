package handlers

import (
    "net/http"
    "encoding/json"
    "fmt"
    // "log"
    "github.com/skrzepto/IPRO-VAULT-SERVER/models"
)

// POST /api/sensor_data
func InsertNewSensorData(rw http.ResponseWriter, req *http.Request) {
    if req.Method == "POST" {
        decoder := json.NewDecoder(req.Body)
        var t models.SensorData
        err := decoder.Decode(&t)
        if err != nil {
            panic(err)
            //log.Fatal(err)
        }
        fmt.Printf("decoded to %#v\n", t)
    }
}
