package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	sarama "github.com/IBM/sarama"
)

var (
	url string = "localhost:9092"
)

func main() {
	conf := sarama.NewConfig()
	producer, err := sarama.NewAsyncProducer([]string{url}, conf)
	if err != nil {
		return
	}

	go func() {
		select {
		case <-producer.Errors():
			log.Println("produce error")
		case <-producer.Successes():
			log.Println("produce success")
		}
	}()

	var topicName, message string
	go func() {
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
	}()

	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, os.Interrupt)
	<-exitSignal
	producer.Close()
	log.Println("producer exit")
}
