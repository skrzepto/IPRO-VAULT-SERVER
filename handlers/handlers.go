package handlers

import (
    "net/http"
    "encoding/json"
    "fmt"
    "sync"
    "github.com/julienschmidt/httprouter"
    "github.com/skrzepto/IPRO-VAULT-SERVER/models"
)

type Impl struct {
  sd_map map[string]models.SensorData
  mu sync.Mutex
}

func InitGlobal() *Impl {
    var i Impl
    i.sd_map = make(map[string]models.SensorData)
    return &i
}


func (i *Impl) POST_SensorData_ID(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
    decoder := json.NewDecoder(req.Body)
    var sd models.SensorData
    err := decoder.Decode(&sd)
    if err != nil {
        panic(err)
    }
    fmt.Printf("decoded to %#v\n", sd)
    //i.DB.Create(&t) change this to array
    sensor_id := ps.ByName("sensor_id")
    i.mu.Lock()
    i.sd_map[sensor_id] = sd
    i.mu.Unlock()
}

func (i *Impl) GET_SensorData_ID(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
    sensor_id := ps.ByName("sensor_id")
    i.mu.Lock()
    sd := i.sd_map[sensor_id]
    i.mu.Unlock()

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
      js, err := json.Marshal(i.sd_map)
      if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
        panic(err)
      }
      rw.Header().Set("Content-Type", "application/json")
      rw.Write(js)
    }
}
