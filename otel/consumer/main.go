package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"

	"kafka-demo/common"

	sarama "github.com/IBM/sarama"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

var (
	groupName  string
	topicNames string
	url        string = "localhost:9092"
)

func init() {
	flag.StringVar(&groupName, "groupname", "myGroup", "provide your group name to consume")
	flag.StringVar(&topicNames, "topics", "quickstart-topic", "provide your topic name, split by ','")
}

type customConsumerHandler struct{}

func (c *customConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *customConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *customConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, cliam sarama.ConsumerGroupClaim) error {
	for msg := range cliam.Messages() {
		log.Printf("topic: %s, key %s, value %s\n", msg.Topic, msg.Key, msg.Value)
		sess.MarkMessage(msg, "done")

		prop := otel.GetTextMapPropagator()
		parentCtx := context.Background()
		ctx := prop.Extract(parentCtx, &MessageCarrier{
			msg: msg,
		})
		trace := otel.Tracer("trace-consumer")
		spanName := "consumer-span-name"
		_, span := trace.Start(ctx, spanName)
		span.SetStatus(codes.Ok, "producer msg succes")
		span.End()
	}
	return nil
}

func main() {
	flag.Parse()
	// init otel provider
	tp, err := common.InitTracer("saram-consumer", common.TracingConfig{
		CollectorHost: "localhost:4317",
		SamplingRate:  1,
	})
	if err != nil {
		log.Println("init otel trace_provider error", err)
		return
	}
	defer tp.Shutdown(context.Background())

	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0
	config.Consumer.Return.Errors = true

	group, err := sarama.NewConsumerGroup([]string{url}, groupName, config)
	if err != nil {
		return
	}

	handle := new(customConsumerHandler)
	go func() {
		for {
			err := group.Consume(context.Background(), convertStr2List(topicNames), handle)
			if err != nil {
				log.Println("consumer occur error", err)
			}
		}
	}()

	go func() {
		for err := range group.Errors() {
			log.Println("group consumer error: ", err)
		}
	}()

	exitChan := make(chan os.Signal)
	signal.Notify(exitChan, os.Interrupt)
	<-exitChan
	group.Close()
	log.Println("close group comsumer, good bye:)")
}

func convertStr2List(name string) []string {
	return strings.Split(name, ",")
}
