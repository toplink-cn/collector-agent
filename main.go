package main

import (
	"collector-agent/lib"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	run()
}

func run() {
	// 连接到RabbitMQ服务器
	conn, err := amqp.Dial("amqp://root:password@192.168.88.112:5672/")
	lib.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	var collector lib.Collector

	mainCh, retryCh := collector.Init(conn)
	defer mainCh.Close()
	defer retryCh.Close()

	// 处理接收到的消息
	forever := make(chan bool)

	go collector.ListenQ(collector.MainCh, collector.MainQ)
	go collector.ListenQ(collector.RetryCh, collector.RetryQ)

	log.Printf(" [*] Waiting for messages. To exit, press CTRL+C")
	<-forever
}
