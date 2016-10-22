package handlers

import (
	"bytes"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	//"fmt"
)

func Test_checkForFaults(t *testing.T) {
	utc, err := time.LoadLocation("UTC")
	if err != nil {
		t.Fatalf("Could not create time location: %v", err)
	}
	dt := time.Date(2016, 12, 20, 10, 28, 50, 0, utc) //year int, month Month, day, hour, min, sec, nsec int, loc *Location

	sd := SensorData{Location: "IIT Tower Vault 1 SW Corner", Serial_Number: 1,
		Temperature: 120.1, Pressure: 760.2, Humidity: 80.2,
		Date_Time: dt, Water_Level: 20}

	ml := MetricLimits{100.0, -20.0, 0.7, 0.0, 160, 10}

	f := checkForFaults(sd, ml)

	available_sensors := [4]string{"Temperature", "Humidity", "Pressure", "Water level"}

	for idx, val := range f.FaultMessages {
		if val != available_sensors[idx]+" exceeded threshold!" {
			t.Fatalf("First message is not abbout %s its: %v", available_sensors[idx],
				f.FaultMessages[0])
		}
	}
}

func Test_checkForFaults_Nil(t *testing.T) {
	utc, err := time.LoadLocation("UTC")
	if err != nil {
		t.Fatalf("Could not create time location: %v", err)
	}
	dt := time.Date(2016, 12, 20, 10, 28, 50, 0, utc) //year int, month Month, day, hour, min, sec, nsec int, loc *Location

	sd := SensorData{Location: "IIT Tower Vault 1 SW Corner", Serial_Number: 1,
		Temperature: 20.1, Pressure: 120.2, Humidity: .2,
		Date_Time: dt, Water_Level: 2}

	ml := MetricLimits{100.0, -20.0, 0.7, 0.0, 160, 10}

	f := checkForFaults(sd, ml)
	if f != nil {
		t.Fatalf("All sensor data were within range, still returned exceeds threshold: %v", err)
	}
}

func Test_GET_Sensor_Data_Empty(t *testing.T) {
	i := InitGlobal()
	req, err := http.NewRequest(
		http.MethodGet,
		"http://localhost:8082/api/sensor_data",
		nil,
	)

	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rec := httptest.NewRecorder()
	i.GET_SensorData(rec, req, nil)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200 got %d", rec.Code)
	}
	if rec.Header()["Content-Type"][0] != "application/json" {
		t.Errorf("expected header of application/json but got %v", rec.Header()["Content-Type"])
	}
	if rec.Body.String() != "{}" {
		t.Errorf("expected empty json body but got %v", rec.Body)
	}
}

func Test_GET_Sensor_Data(t *testing.T) {
	i := InitGlobal()
	utc, err := time.LoadLocation("UTC")
	if err != nil {
		t.Fatalf("Could not create time location: %v", err)
	}
	dt := time.Date(2016, 12, 20, 10, 28, 50, 0, utc) //year int, month Month, day, hour, min, sec, nsec int, loc *Location

	sd1 := SensorData{Location: "IIT Tower Vault 1 SW Corner", Serial_Number: 1,
		Temperature: 120.1, Pressure: 760.2, Humidity: 80.2,
		Date_Time: dt, Water_Level: 20}
	i.sd_table["1"] = sd1

	sd2 := SensorData{Location: "IIT Tower Vault 2 NW Corner", Serial_Number: 2,
		Temperature: 40.1, Pressure: 760.2, Humidity: 40.2,
		Date_Time: dt, Water_Level: 2}
	i.sd_table["2"] = sd2

	req, err := http.NewRequest(
		http.MethodGet,
		"http://localhost:8082/api/sensor_data",
		nil,
	)

	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rec := httptest.NewRecorder()
	i.GET_SensorData(rec, req, nil)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200 got %d", rec.Code)
	}
	if rec.Header()["Content-Type"][0] != "application/json" {
		t.Errorf("expected header of application/json but got %v", rec.Header()["Content-Type"])
	}
	expected := "{\"1\":{\"location\":\"IIT Tower Vault 1 SW Corner\",\"serial_number\":1,\"temperature\":120.1,\"pressure\":760.2,\"humidity\":80.2,\"datetime\":\"2016-12-20T10:28:50Z\",\"water_level\":20},\"2\":{\"location\":\"IIT Tower Vault 2 NW Corner\",\"serial_number\":2,\"temperature\":40.1,\"pressure\":760.2,\"humidity\":40.2,\"datetime\":\"2016-12-20T10:28:50Z\",\"water_level\":2}}"
	if rec.Body.String() != expected {
		t.Errorf("expected empty json body but got %v", rec.Body)
	}
	req, err = http.NewRequest(
		http.MethodGet,
		"http://localhost:8082/api/sensor_data",
		nil,
	)

	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rec = httptest.NewRecorder()
	i.GET_SensorData(rec, req, nil)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200 got %d", rec.Code)
	}
	if rec.Header()["Content-Type"][0] != "application/json" {
		t.Errorf("expected header of application/json but got %v", rec.Header()["Content-Type"])
	}
	expected = "{\"1\":{\"location\":\"IIT Tower Vault 1 SW Corner\",\"serial_number\":1,\"temperature\":120.1,\"pressure\":760.2,\"humidity\":80.2,\"datetime\":\"2016-12-20T10:28:50Z\",\"water_level\":20},\"2\":{\"location\":\"IIT Tower Vault 2 NW Corner\",\"serial_number\":2,\"temperature\":40.1,\"pressure\":760.2,\"humidity\":40.2,\"datetime\":\"2016-12-20T10:28:50Z\",\"water_level\":2}}"
	if rec.Body.String() != expected {
		t.Errorf("expected empty json body but got %v", rec.Body)
	}

	req, err = http.NewRequest(
		http.MethodGet,
		"http://localhost:8082/api/sensor_data/999",
		nil,
	)

	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rec = httptest.NewRecorder()
	ps := httprouter.Params{httprouter.Param{"sensor_id", "999"}}
	i.GET_SensorData_ID(rec, req, ps)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404 got %d", rec.Code)
	}

	req, err = http.NewRequest(
		http.MethodGet,
		"http://localhost:8082/api/sensor_data/1",
		nil,
	)

	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rec = httptest.NewRecorder()
	ps = httprouter.Params{httprouter.Param{"sensor_id", "1"}}
	i.GET_SensorData_ID(rec, req, ps)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200 got %d", rec.Code)
	}

	if rec.Header()["Content-Type"][0] != "application/json" {
		t.Errorf("expected header of application/json but got %v", rec.Header()["Content-Type"])
	}
	expected = "{\"location\":\"IIT Tower Vault 1 SW Corner\",\"serial_number\":1,\"temperature\":120.1,\"pressure\":760.2,\"humidity\":80.2,\"datetime\":\"2016-12-20T10:28:50Z\",\"water_level\":20}"
	if rec.Body.String() != expected {
		t.Errorf("expected sd1 json body but got %v", rec.Body)
	}
}

func Test_POST_Sensor_Data(t *testing.T) {
	i := InitGlobal()
	payload := []byte(`{
      "serial_number": 1,
      "location": "IIT Tower Vault 1 SW Corner",
      "temperature": 20.1,
      "pressure": 760.2,
      "humidity": 80.2,
      "water_level": 2,
      "datetime": "2016-10-21T23:25:40.573559+00:00"
  }`)
	req, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:8082/api/sensor_data/0",
		bytes.NewReader(payload),
	)

	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rec := httptest.NewRecorder()
	ps := httprouter.Params{httprouter.Param{"sensor_id", "1"}}
	i.POST_SensorData_ID(rec, req, ps)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200 got %d", rec.Code)
	}

	// Make sure the sd_table has the data
	if i.sd_table["1"].Serial_Number != 1 {
		t.Errorf("expected sensor_number to be 1 but got %d", i.sd_table["1"].Serial_Number)
	}

	// The json sent is faulty verify if its inserted
	if *i.faults["1"][0].DataEntry != i.sd_table["1"] {
		t.Errorf("fault dataentry and sd_table should be equal, dataentry: %v  sd_table %v", *i.faults["1"][0].DataEntry, i.sd_table["1"])
	}
}

func Test_GET_DELETE_Faults_ID(t *testing.T) {
	i := InitGlobal()
	utc, err := time.LoadLocation("UTC")
	if err != nil {
		t.Fatalf("Could not create time location: %v", err)
	}
	dt := time.Date(2016, 12, 20, 10, 28, 50, 0, utc) //year int, month Month, day, hour, min, sec, nsec int, loc *Location

	sd1 := SensorData{Location: "IIT Tower Vault 1 SW Corner", Serial_Number: 1,
		Temperature: 120.1, Pressure: 760.2, Humidity: 80.2,
		Date_Time: dt, Water_Level: 20}
	i.sd_table["1"] = sd1
	fe := FaultEntry{&sd1, []string{"Temperature exceeded threshold!"}}
	i.faults["1"] = append(i.faults["1"], &fe)

	req, err := http.NewRequest(
		http.MethodGet,
		"http://localhost:8082/api/faults/1",
		nil,
	)

	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rec := httptest.NewRecorder()
	ps := httprouter.Params{httprouter.Param{"sensor_id", "1"}}
	i.GET_Faults_ID(rec, req, ps)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200 got %d", rec.Code)
	}
	expected := "{\"DataEntry\":{\"location\":\"IIT Tower Vault 1 SW Corner\",\"serial_number\":1,\"temperature\":120.1,\"pressure\":760.2,\"humidity\":80.2,\"datetime\":\"2016-12-20T10:28:50Z\",\"water_level\":20},\"FaultMessages\":[\"Temperature exceeded threshold!\"]}"
	if expected != rec.Body.String() {
		t.Errorf("expected string not equal to response body, got %s", rec.Body.String())
	}

	//Delete the faults for sensor id 1
	req, err = http.NewRequest(
		http.MethodDelete,
		"http://localhost:8082/api/faults/1",
		nil,
	)

	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rec = httptest.NewRecorder()
	ps = httprouter.Params{httprouter.Param{"sensor_id", "1"}}
	i.DELETE_Faults_ID(rec, req, ps)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200 got %d", rec.Code)
	}

	//Delete the faults for sensor id 999 which doesnt exist
	req, err = http.NewRequest(
		http.MethodDelete,
		"http://localhost:8082/api/faults/999",
		nil,
	)

	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rec = httptest.NewRecorder()
	ps = httprouter.Params{httprouter.Param{"sensor_id", "999"}}
	i.DELETE_Faults_ID(rec, req, ps)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404 got %d", rec.Code)
	}

	//Try to get faults for sensor id 1 which was just deleted
	req, err = http.NewRequest(
		http.MethodGet,
		"http://localhost:8082/api/faults/1",
		nil,
	)

	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rec = httptest.NewRecorder()
	ps = httprouter.Params{httprouter.Param{"sensor_id", "1"}}
	i.GET_Faults_ID(rec, req, ps)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404 got %d", rec.Code)
	}

	expected = "Can't find fault data for the id specified"
	if expected != rec.Body.String() {
		t.Errorf("expected string not equal to response body, got %s", rec.Body.String())
	}
}
