package main

import (
	"strings"

	sarama "github.com/IBM/sarama"
	"go.opentelemetry.io/otel/propagation"
)

type MessageCarrier struct {
	msg *sarama.ConsumerMessage
}

var _ propagation.TextMapCarrier = (*MessageCarrier)(nil)

func (mc *MessageCarrier) Get(key string) string {
	for _, h := range mc.msg.Headers {
		if strings.EqualFold(key, string(h.Key)) {
			return string(h.Value)
		}
	}
	return ""
}

func (mc *MessageCarrier) Set(key string, value string) {
	mc.msg.Headers = append(mc.msg.Headers, &sarama.RecordHeader{
		Key:   []byte(key),
		Value: []byte(value),
	})
}

func (mc *MessageCarrier) Keys() []string {
	var keys []string
	for _, h := range mc.msg.Headers {
		keys = append(keys, string(h.Key))
	}
	return keys
}
