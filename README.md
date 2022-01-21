# RaspberryPi automated feeder
A RaspberryPi automated feeder for pets. What this module does essentially is make sure it rotates a servo motor, connected to a specified pin on the RPI. The application uses MQTT 5 to handle requests and to send responses back.

The feeder can feed on request or by executing a pre-defined schedule.

## Configuration
There are several configuration sections for the feeder.

### General
General settings for the feeder.

| Key       | Description                                                                                                                              |
|-----------|------------------------------------------------------------------------------------------------------------------------------------------|
| dbPath    | The location in which to store the BoltDB database.                                                                                      |
| servoPin  | The control pin to which the servo motor is connected.                                                                                   |
| portionMs | The milliseconds the servo should rotate in order to drop 1 portion of food. That would be dependent on the food dispenser that is used. |

### MQTT
MQTT specific settings.

| Key               | Description                                                                                          |
|-------------------|------------------------------------------------------------------------------------------------------|
| server            | The URL to the MQTT broker.                                                                          |
| username          | The username to use to authenticate with the broker.                                                 |
| password          | The password to use to authenticate with the broker.                                                 |
| clientId          | The clientId to use to authenticate with the broker.                                                 |
| keepAlive         | The keep-alive timeout for the MQTT connection.                                                      |
| connectRetryDelay | The retry delay in seconds for attempting to reconnect to the MQTT broken if the connection is lost. |

An example configuration exists in `example_config.json`.


## MQTT topics used
The following MQTT topics are used by the feeder.

| Topic                        | Description                                                                                                                                                                                                                                                                                         |
|----------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| feeder/{clientId}/status   | The status of the feeder is available on this topic. The status message is persisted and states the current version of the feeder software and whether it is online or offline. The feeder implements LWT message such that when connection is lost the status is automatically updated to offline. |
| feeder/{clientId}/feed     | The feeder listens for messages on this topic for performing a manual feed. The message should contain the amount of portions to be dropped.                                                                                                                                                        |
| feeder/{clientId}/feed_log | The feed log is available on this topic. Every time the feeder drops food, sends a message on this topic stating the time and the portions that were dropped. If the feeder has lost connection with the broker, it will re-send the current feed log history on its next restart.                  |


