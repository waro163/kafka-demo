# 概述
kafka生产消费的demo，使用的go库为 gopkg.in/Shopify/sarama.v1

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
