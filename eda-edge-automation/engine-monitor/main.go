package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// loadConf loads configurations from YAML file
func loadConf(path string) (map[interface{}]interface{}, error) {

	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	data := make(map[interface{}]interface{})

	err = yaml.Unmarshal(yamlFile, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// fakeSensorData simulates the output of sensors attached to an engine
func fakeSensorData(t time.Time) string {

	vInRangeMin := 180
	vInRangeMax := 200
	vOutRangeMin := 201
	vOutRangeMax := 350

	// As of Go 1.20 there is no reason to call Seed with a random value
	v := rand.Intn(vInRangeMax-vInRangeMin) + vInRangeMin

	u, _ := time.ParseDuration("1m")
	if time.Since(t) > u {
		v = rand.Intn(vOutRangeMax-vOutRangeMin) + vOutRangeMin
	}
	return fmt.Sprintf(`{"rpms": 5000, "volts": 5, "offset-vibration (mm/s)": %d}`, v)
}

var engineOn bool

type Message struct {
	Status string `json:"status"`
	Fail   bool   `json:"fail"`
}

// engineShutDown is a dummy function to simulate a shutdown action
func engineShutdown(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		var m Message

		json.NewDecoder(r.Body).Decode(&m)

		if m.Status == "shutdown" {
			engineOn = false
			log.Printf("Engine shutdown request from %v, User-Agent: %v\n", r.Host, r.UserAgent())
		}
	} else {
		log.Printf("Invalid request from %v, User-Agent: %v\n", r.Host, r.UserAgent())
	}
}

// startServer starts the http server
func startServer(port string) {
	http.HandleFunc("/shutdown", engineShutdown)
	http.ListenAndServe(":"+port, nil)
}

func main() {

	fmt.Println("Starting engine monitor simulator.")

	engineOn = true

	yamlConfigMap := flag.String("config", "config.yaml", "Client YAML config file")
	listenPort := flag.String("port", "8080", "Default listen port")
	flag.Parse()

	// Load configuration from YAML
	cfg, err := loadConf(*yamlConfigMap)
	if err != nil {
		log.Printf("Error loading conf: %v\n", err)
		os.Exit(1)
	}

	// Populate kafkaConfigMap
	// To manage authetication mechanisms, uncomment the security.protocol and sasl fields
	var engineMonConfig = kafka.ConfigMap{
		"bootstrap.servers": cfg["bootstrap-servers"],
		//"security.protocol": cfg["security-protocol"],
		//"sasl.mechanisms":   cfg["sasl-mechanisms"],
		//"sasl.username":     cfg["sasl-username"],
		//"sasl.password":     cfg["sasl-password"],
	}

	// Define topic from config
	var topic = cfg["topic"].(string)

	// startTime is used to begin the
	var startTime = time.Now()

	// Start web service to handle control actions as asynchronous goroutine
	go startServer(*listenPort)
	log.Printf("Starting web service to handle control action")

	// Create new Kafka producer
	p, err := kafka.NewProducer(&engineMonConfig)
	if err != nil {
		log.Printf("Failed to create Kafka producer. %v", err)
		os.Exit(1)
	}

	// Defer Closing of producer channel
	defer p.Close()

	// Listen to all the events in the default event channel
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					log.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
				} else {
					log.Printf("Udated engine status to topic %s [%d] at offset %v\n", *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
				}
			case kafka.Error:
				log.Printf("Error: %v\n", ev)
			default:
				log.Printf("Ignored event: %s\n", ev)
			}
		}
	}()

	// Produce events until engineOn is true
	for engineOn == true {
		value := fmt.Sprintf(fakeSensorData(startTime))

		err := p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(value),
			Headers:        []kafka.Header{{Key: "test header", Value: []byte("binary test values")}},
		}, nil)

		if err != nil {
			if err.(kafka.Error).Code() == kafka.ErrQueueFull {
				time.Sleep(time.Second)
				continue
			}
			log.Printf("Failed to produce message: %v\n", err)
		}
		time.Sleep(2 * time.Second)
	}

}
