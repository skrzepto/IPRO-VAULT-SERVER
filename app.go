package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func main() {
	i := InitGlobal()
	router := httprouter.New()

	router.GET("/", i.Index)
	router.GET("/node/:sensor_id", i.NodeStatus)
	router.GET("/api/sensor_data", i.GET_SensorData)
	router.GET("/api/sensor_data/:sensor_id", i.GET_SensorData_ID)
	router.POST("/api/sensor_data/:sensor_id", i.POST_SensorData_ID)
	router.GET("/api/faults", i.GET_Faults)
	router.GET("/api/faults/:sensor_id", i.GET_Faults_ID)
	router.DELETE("/api/faults/:sensor_id", i.DELETE_Faults_ID)
	router.GET("/api/faults/:sensor_id/delete", i.DELETE_Faults_ID)

	log.Fatal(http.ListenAndServe(":8082", router))
}
