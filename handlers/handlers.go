package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"sync"
	"time"
)

type SensorData struct {
	Name          string    `json:"name"`
	Location      string    `json:"location"`
	Serial_Number int       `json:"serial_number"`
	Temperate     float64   `json:"temperature"`
	Pressure      float64   `json:"pressure"`
	Humidity      float64   `json:"humidity"`
	Water_Level   float64   `json:"water_level"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type Impl struct {
	sd_table map[string]SensorData
	mu       sync.Mutex
}

func InitGlobal() *Impl {
	var i Impl
	i.sd_table = make(map[string]SensorData)
	return &i
}

func (i *Impl) POST_SensorData_ID(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(req.Body)
	var sd SensorData
	err := decoder.Decode(&sd)
	if err != nil {
		panic(err)
	}
	fmt.Printf("decoded to %#v\n", sd)
	//i.DB.Create(&t) change this to array
	sensor_id := ps.ByName("sensor_id")
	i.mu.Lock()
	i.sd_table[sensor_id] = sd
	i.mu.Unlock()
}

func (i *Impl) GET_SensorData_ID(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	sensor_id := ps.ByName("sensor_id")
	i.mu.Lock()
	sd, exist:= i.sd_table[sensor_id]
	i.mu.Unlock()
  if !exist {
    rw.WriteHeader(404)
    rw.Write([]byte("Can't find sendor data for the id specified"))
    return
  }
	js, err := json.Marshal(sd)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		panic(err)
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(js)
}

// [GET] /api/sensor_data
func (i *Impl) SensorData(rw http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if req.Method == "GET" {
		js, err := json.Marshal(i.sd_table)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			panic(err)
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(js)
	}
}
