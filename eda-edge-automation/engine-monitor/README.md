# Engine monitor simulator

Engine monitor simulator for Event-Driven Ansible demo at Red Hat Summit Connect Italy 2023.

The application simulates the collection of sensor data from a DC engine and introduces a skew in offset vibrations after a predefined time.
All data are written to a Kafka topic and observed by an EDA rulebook.
Since an increase of the offset vibration could lead to an immineant breakage, the rulebook triggers the engine shutdown using a POST method exposed by the simulator REST API.

## Build

To build the application:

```bash
podman build -t engine-monitor .
```

## Run

Execute the `engine-monitor` with the following arguments:

```
engine-monitor -config <CONFIG_FILE_PATH> -port <PORT> -ttf <TIME_TO_FAIL>
```

The `-ttf` flag is the amount of time before introducing the offset vibration issue.

## Customize configuration

Apply the desired customizations to the `config.yaml` file. 

Update the `bootstrap-servers` field with the Kafka(s) list of endpoints in the pattern <HOST>:<PORT>.  
Update the `topic` field with the name of the topic receiveng the events.

For the sake of the demo no authentication protocol is enabled by default.

## Deploy to OpenShift/Microshift

Simply apply the files in the `manifest` directory to the namespace of choice:

```bash
kubectl apply -f manifests/ -n <NAMESPACE>
```

