package models

import (
    "time"
)
type SensorData struct {
    Name                string      `json:"name"`
    Location            string      `json:"location"`
    Serial_Number       int         `json:"serial_number"`
    Temperate           float64     `json:"temperature"`
    Pressure            float64     `json:"pressure"`
    Humidity            float64     `json:"humidity"`
    Water_Level         float64     `json:"water_level"`
    CreatedAt           time.Time   `json:"createdAt"`
    UpdatedAt           time.Time   `json:"updatedAt"`
}
