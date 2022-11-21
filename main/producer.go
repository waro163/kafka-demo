package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	sarama "gopkg.in/Shopify/sarama.v1"
)

var (
	// topicName string
	url string = "localhost:9092"
)

// func init() {
// 	flag.StringVar(&topicName, "topic", "quickstart-topic", "provide which topic you will produce message")
// }
func main() {
	// flag.Parse()
	conf := sarama.NewConfig()
	producer, err := sarama.NewAsyncProducer([]string{url}, conf)
	if err != nil {
		return
	}
	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, os.Interrupt)

	go func() {
		select {
		case <-producer.Errors():
			log.Println("produce error")
		case <-producer.Successes():
			log.Println("produce success")
		case <-exitSignal:
			producer.Close()
			log.Println("producer exit")
			os.Exit(0)
		}
	}()
	var topicName, message string

	for {
		fmt.Println("input your topic:")
		fmt.Scanf("%s", &topicName)
		fmt.Println("input your message:")
		fmt.Scanf("%s", &message)

		msg := &sarama.ProducerMessage{
			Topic: topicName,
			Key:   sarama.StringEncoder(message),
			Value: sarama.StringEncoder(message),
		}
		producer.Input() <- msg
	}
}
