# IPRO-VAULT and SERVER

## Purpose
To aggregate sensor data from various raspberry pi and alert if any of the pi read sensor data that is outside the specified safe range. This will alert IIT maintenance staff when a vault is needed of repair within the frequency of the data sent which is set to every ~15 minutes.

------

## Server API

### GET /api/
Display general information of the nodes (eg. which sensors are active, the range for safe conditions, location of the node, etc…)

### GET /api/sensor_data
Display the current state of all nodes with their sensor values, timestamp, and current status (green, red)

### GET /api/sensor_data/:sensor_id
Display the sensor data for a specific ResponseWriter

### POST /api/sensor_data/:sensor_id
Send data from rpi of its current sensor values

example json
```
{     
    "serial_number": 1,
    "name": "RPI 3",
    "location": "IIT Tower Vault 1 SW Corner",
    "temperature": 20.1,
    "pressure": 760.2,
    "humidity": 80.2,
    "water_level": 2
}
```

### GET /api/faults
Display current nodes which are faulting and note which sensor, with timestamp of fault

### GET /api/faults/:sensor_id
Display the specied nodes detailed history of faults from last deletion

### DELETE /api/faults/:sensor_id
Remove sensor from fault table

------

## Design Decisions / Concerns

Initially the thought was to use a DB to persist the data, but after careful consideration of the our initial scope it is not needed.

#### Initial Scope
Be able to display the most recent data sent and show which nodes are faulting

#### Future Scope
Log each of the nodes data into a db so in another application the possibility of analysis can be done to predict future vault failures and at which locations. (eg. does a particular time of year cause the vault to fault more than others, is there a location that faults the most and needs to be manually monitored more frequently, etc…)

Currently anyone can upload data pretending to be a node. If time permits add the ability to use authentication tokens with the messages so a single source is only allowed to update the  nodes data.


------

## Use Case / Story

### Node uploads new data

A node sends its current sensor data to the server. If the node or the server loose connection, ignore retry unless the data is outside the specified limits. If one of the sensors is outside the limits log the data and continue iteration and send the current and whatever is logged in the file.

### Node uploads data outside limit ranges

Update most recent data and add the node into the fault data table including why it’s faulting.

### User wants to see a glance of all sensor data

Display each sensor data and information about the sensor with a human readable status (green, red)

### User wants to see sensor data for only one node

Just display the data for that particular sensor

### User wants to see the fault table

Display all faults that have occurred with reasons why

### User wants to remove a node from the fault table

After checking out the node and conditions are good. The user must remove the node from the fault table.

-----

## JSON data example

```
{
    "serial_number": 1,
    "location": "IIT Tower Vault 1 SW Corner",
    "temperature": 20.1,
    "pressure": 760.2,
    "humidity": 80.2,
    "water_level": 2
}

```
