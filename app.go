package main

import (
    "net/http"
    "io"
    "log"
    "github.com/julienschmidt/httprouter"
    "github.com/skrzepto/IPRO-VAULT-SERVER/handlers"
)

func hello(rw http.ResponseWriter, req *http.Request, _ httprouter.Params) {
    io.WriteString(rw, "Hello world!")
}

func main() {
    i := handlers.InitGlobal()
    router := httprouter.New()

    router.GET("/", hello)
    router.GET("/api/sensor_data", i.SensorData)
    router.POST("/api/sensor_data", i.SensorData)
    router.GET("/api/sensor_data/:sensor_id", i.GET_SensorData_ID)
    router.POST("/api/sensor_data/:sensor_id", i.POST_SensorData_ID)

    log.Fatal(http.ListenAndServe(":8082", router))
}
