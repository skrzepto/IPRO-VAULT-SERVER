package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/skrzepto/IPRO-VAULT-SERVER/handlers"
	"io"
	"log"
	"net/http"
)

func hello(rw http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	io.WriteString(rw, "Hello world!")
}

func main() {
	i := handlers.InitGlobal()
	router := httprouter.New()

	router.GET("/", hello)
	router.GET("/api/sensor_data", i.GET_SensorData)
	router.GET("/api/sensor_data/:sensor_id", i.GET_SensorData_ID)
	router.POST("/api/sensor_data/:sensor_id", i.POST_SensorData_ID)
	router.GET("/api/faults", i.GET_Faults)
	router.GET("/api/faults/:sensor_id", i.GET_Faults_ID)
	router.DELETE("/api/faults/:sensor_id", i.DELETE_Faults_ID)

	log.Fatal(http.ListenAndServe(":8082", router))
}
