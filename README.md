# 概述

kafka生产消费的demo，使用的go库为 ~~gopkg.in/Shopify/sarama.v1~~ github.com/IBM/sarama

另一个库为<https://github.com/segmentio/kafka-go>，暂时没用使用对比

# before run

```
go mod tidy

go mod vendor
```

# simple use

## producer
```
go run -mod=vendor main/producer/*
```

## consumer
```
go run -mod=vendor main/consumer/* --topics=quickstart-topic,test --groupname=myGroup
```

# open telemetry use

## producer
```
go run -mod=vendor otel/producer/*
```

## consumer
```
go run -mod=vendor otel/consumer/* --topics=quickstart-topic,test --groupname=myGroup
```
