# IPRO-VAULT-SERVER

## Purpose
----
To aggregate sensor data from various raspberry pi and alert if any of the pi
read sensor data that is outside the specified safe range.

## API

### GET /api/
Display general information of the nodes (eg. which sensors are active, the
range for safe conditions, location of the node, etc…)

### GET /api/sensor_data
Display the current state of all nodes with their sensor values, timestamp,
and current status (green, red)

### GET /api/sensor_data/:sensor_id
Display the sensor data for a specific ResponseWriter

### POST /api/sensor_data/:sensor_id
Send data from rpi of its current sensor values

### GET /api/faults
Display current nodes which are faulting and note which sensor, with
timestamp of fault

### DELETE /api/faults/:sensor_id
Remove sensor from fault table
