package handlers

import (
    "net/http"
    "encoding/json"
    "fmt"
    "log"
    "github.com/skrzepto/IPRO-VAULT-SERVER/models"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
    "github.com/jinzhu/gorm"
)

type Impl struct {
    DB *gorm.DB
}

func (i *Impl) InitDB() {
    var err error
    i.DB, err = gorm.Open("sqlite3", "ipro-vault.db")
    if err != nil {
        log.Fatalf("Got error when connect database, the error is '%v'", err)
    }
    i.DB.LogMode(true)
}

func (i *Impl) InitSchema() {
    i.DB.AutoMigrate(&models.SensorData{})
}

// POST /api/sensor_data
func (i *Impl) InsertNewSensorData(rw http.ResponseWriter, req *http.Request) {
    if req.Method == "POST" {
        decoder := json.NewDecoder(req.Body)
        var t models.SensorData
        err := decoder.Decode(&t)
        if err != nil {
            panic(err)
            //log.Fatal(err)
        }
        fmt.Printf("decoded to %#v\n", t)
        i.DB.Create(&t)
    }
}
