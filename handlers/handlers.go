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
	Temperature   float64   `json:"temperature"`
	Pressure      float64   `json:"pressure"`
	Humidity      float64   `json:"humidity"`
	Water_Level   float64   `json:"water_level"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type MetricLimits struct {
	TempMax   float64
	TempMin   float64
	HumMax    float64
	HumMin    float64
	PresMax   float64
	WLevelMax float64
}

type FaultEntry struct {
	DataEntry     *SensorData
	FaultMessages []string
}

type Impl struct {
	sd_table map[string]SensorData
	faults   map[string][]*FaultEntry
	limits   MetricLimits
	mu       sync.Mutex
}

func InitGlobal() *Impl {
	var i Impl
	i.sd_table = make(map[string]SensorData)
	i.faults = make(map[string][]*FaultEntry)
	i.limits = MetricLimits{100.0, -20.0, 0.7, 0.0, 160, 10}
	return &i
}

func checkForFaults(sd SensorData, lim MetricLimits) *FaultEntry {
	var msgs []string

	if sd.Temperature < lim.TempMin || lim.TempMax < sd.Temperature {
		msgs = append(msgs, "Temperature exceeded threshold!")
	}

	if sd.Humidity < lim.HumMin || lim.HumMax < sd.Humidity {
		msgs = append(msgs, "Humidity exceeded threshold!")
	}

	if sd.Pressure > lim.PresMax {
		msgs = append(msgs, "Pressure exceeded threshold!")
	}

	if sd.Water_Level > lim.WLevelMax {
		msgs = append(msgs, "Water level exceeded threshold!")
	}

	if len(msgs) == 0 {
		return nil
	}

	return &FaultEntry{&sd, msgs}
}

func (i *Impl) POST_SensorData_ID(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	var sd SensorData

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&sd)
	if err != nil {
		panic(err)
	}
	fmt.Printf("decoded to %#v\n", sd)

	faults := checkForFaults(sd, i.limits)
	sensor_id := ps.ByName("sensor_id")

	i.mu.Lock()
	i.sd_table[sensor_id] = sd
	if faults != nil {
		i.faults[sensor_id] = append(i.faults[sensor_id], faults)
	}
	i.mu.Unlock()
}

func (i *Impl) GET_SensorData_ID(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	sensor_id := ps.ByName("sensor_id")
	i.mu.Lock()
	sd, exist := i.sd_table[sensor_id]
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

func (i *Impl) GET_SensorData(rw http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	i.mu.Lock()
	js, err := json.Marshal(i.sd_table)
	i.mu.Unlock()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		panic(err)
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(js)
}

func (i *Impl) GET_Faults(rw http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	i.mu.Lock()
	js, err := json.Marshal(i.faults)
	i.mu.Unlock()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(js)
}

func (i *Impl) DELETE_Faults_ID(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	sensor_id := ps.ByName("sensor_id")

	i.mu.Lock()
	flts, exist := i.faults[sensor_id]
	if !exist {
		rw.WriteHeader(404)
		rw.Write([]byte("Can't find fault data for the id specified"))
		i.mu.Unlock()
		return
	}
	js, err := json.Marshal(flts)
	delete(i.faults, sensor_id)
	i.mu.Unlock()

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		panic(err)
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(js)
}
