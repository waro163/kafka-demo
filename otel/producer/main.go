package main

import (
	"context"
	"fmt"
	"kafka-demo/common"
	"log"
	"os"
	"os/signal"

	sarama "github.com/IBM/sarama"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

var (
	url string = "localhost:9092"
)

func main() {
	// init otel provider
	tp, err := common.InitTracer("saram-producer", common.TracingConfig{
		CollectorHost: "localhost:4317",
		SamplingRate:  1,
	})
	if err != nil {
		log.Println("init otel trace_provider error", err)
		return
	}
	defer tp.Shutdown(context.Background())

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

			msg := generateMsg(topicName, message)
			producer.Input() <- msg
		}
	}()

	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, os.Interrupt)
	<-exitSignal
	producer.Close()
	log.Println("producer exit")
}

func generateMsg(topicName, message string) *sarama.ProducerMessage {
	prop := otel.GetTextMapPropagator()
	parentCtx := context.Background()
	msg := &sarama.ProducerMessage{
		Topic: topicName,
		Key:   sarama.StringEncoder(message),
		Value: sarama.StringEncoder(message),
	}
	// usually will extract request header into parentCtx
	// prop.Extract(parentCtx, propagation.HeaderCarrier(c.Request.Header))
	ctx := prop.Extract(parentCtx, &MessageCarrier{
		msg: new(sarama.ProducerMessage),
	})
	trace := otel.Tracer("trace-producer")
	spanName := "producer-span-name"
	sonCtx, span := trace.Start(ctx, spanName)
	defer span.End()
	// send msg to kafka
	prop.Inject(sonCtx, &MessageCarrier{msg: msg})
	span.SetStatus(codes.Ok, "producer msg succes")
	return msg
}
