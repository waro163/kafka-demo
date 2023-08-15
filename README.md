# 概述

kafka生产消费的demo，使用的go库为 ~~gopkg.in/Shopify/sarama.v1~~ github.com/IBM/sarama

另一个库为<https://github.com/segmentio/kafka-go>，暂时没用使用对比

# before run

```
go mod tidy

go mod vendor
```

# producer

```
go run -mod=vendor main/producer.go
```

# consumer

```
go run -mod=vendor main/worker.go --topics=quickstart-topic,test --groupname=myGroup
```
