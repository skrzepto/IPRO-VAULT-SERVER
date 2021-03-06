package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"sync"
	"time"
	"html/template"
	"path"
)

type SensorData struct {
	Location      string    `json:"location"`
	Serial_Number int       `json:"serial_number"`
	Temperature   float64   `json:"temperature"`
	Pressure      float64   `json:"pressure"`
	Humidity      float64   `json:"humidity"`
	Date_Time     time.Time `json:"datetime"`
	Water_Level   float64   `json:"water_level"`
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
	i.limits = MetricLimits{30.0, 20.0, 20, 0.0, 1.1, 10}
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


/***
______               _     _____          _
|  ___|             | |   |  ___|        | |
| |_ _ __ ___  _ __ | |_  | |__ _ __   __| |
|  _| '__/ _ \| '_ \| __| |  __| '_ \ / _` |
| | | | | (_) | | | | |_  | |__| | | | (_| |
\_| |_|  \___/|_| |_|\__| \____|_| |_|\__,_|
***/
func (i *Impl) Index(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	i.mu.Lock()
	data := i.sd_table
	i.mu.Unlock()

	fp := path.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
	    http.Error(rw, err.Error(), http.StatusInternalServerError)
	    return
	}
	if err := tmpl.Execute(rw, data); err != nil {
	    http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

type nodedata struct {
    SD SensorData
		Faults	[]*FaultEntry
		IsFaulting bool
}

func (i *Impl) NodeStatus(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	sensor_id := ps.ByName("sensor_id")
	var data nodedata
	i.mu.Lock()
	if _, ok := i.sd_table[sensor_id]; ok {
    data.SD = i.sd_table[sensor_id]
	} else {
		i.mu.Unlock()
		rw.WriteHeader(404)
		rw.Write([]byte("Can't find sendor data for the id specified"))
		return
	}
	i.mu.Unlock()

	if _, ok := i.faults[sensor_id]; ok {
    data.Faults = i.faults[sensor_id]
		data.IsFaulting = true
	} else {
		data.IsFaulting = false
	}

	//fmt.Printf("decoded to %#v\n", data)

	fp := path.Join("templates", "rpi.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
	    http.Error(rw, err.Error(), http.StatusInternalServerError)
	    return
	}
	if err := tmpl.Execute(rw, data); err != nil {
	    http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}



/***
______ _____ _____ _____    ___ ______ _____
| ___ |  ___/  ___|_   _|  / _ \| ___ |_   _|
| |_/ | |__ \ `--.  | |   / /_\ | |_/ / | |
|    /|  __| `--. \ | |   |  _  |  __/  | |
| |\ \| |___/\__/ / | |   | | | | |    _| |_
\_| \_\____/\____/  \_/   \_| |_\_|    \___/
***/
func (i *Impl) POST_SensorData_ID(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	var sd SensorData

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&sd)
	if err != nil {
		panic(err)
	}
	fmt.Printf("decoded to %#v\n", sd)
    sd.Pressure = sd.Pressure/100000.0
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

func (i *Impl) GET_Faults_ID(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	sensor_id := ps.ByName("sensor_id")
	i.mu.Lock()
	fault, ok := i.faults[sensor_id]
	if !ok {
		rw.WriteHeader(404)
		rw.Write([]byte("Can't find fault data for the id specified"))
		i.mu.Unlock()
		return
	}
	js, err := json.Marshal(fault[len(fault)-1])
	i.mu.Unlock()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(js)
}

func (i *Impl) getLatestFaults() map[string]*FaultEntry {
	res := make(map[string]*FaultEntry)
	i.mu.Lock()
	if len(i.faults) > 0 {
		for key, val := range i.faults {
			res[key] = val[len(val)-1]
		}
	}
	i.mu.Unlock()
	return res
}

func (i *Impl) GET_Faults(rw http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	js, err := json.Marshal(i.getLatestFaults())

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
