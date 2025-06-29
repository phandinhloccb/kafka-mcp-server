package main

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

// listKafkaTopics get all topics from Kafka broker
func listKafkaTopics(broker string) ([]string, error) {
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to broker %s: %v", broker, err)
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		return nil, fmt.Errorf("cannot read partitions: %v", err)
	}

	// Create map to store unique topics
	topicMap := make(map[string]bool)
	for _, partition := range partitions {
		topicMap[partition.Topic] = true
	}

	// Convert map to slice
	var topics []string
	for topic := range topicMap {
		topics = append(topics, topic)
	}

	return topics, nil
}

// createKafkaTopic create a new topic with specified number of partitions
func createKafkaTopic(broker, topic string, partitions int) error {
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		return fmt.Errorf("cannot connect to broker %s: %v", broker, err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("cannot get controller information: %v", err)
	}

	controllerConn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	if err != nil {
		return fmt.Errorf("cannot connect to controller: %v", err)
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     partitions,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		return fmt.Errorf("cannot create topic: %v", err)
	}

	return nil
}

// produceKafkaMessage send message to topic
func produceKafkaMessage(broker, topic, key, message string) error {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	defer writer.Close()

	msg := kafka.Message{
		Value: []byte(message),
		Time:  time.Now(),
	}

	// Add key if provided
	if key != "" {
		msg.Key = []byte(key)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := writer.WriteMessages(ctx, msg)
	if err != nil {
		return fmt.Errorf("không thể gửi message: %v", err)
	}

	return nil
}

// consumeKafkaMessages read messages from topic
func consumeKafkaMessages(broker, topic string, count int) (string, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{broker},
		Topic:     topic,
		Partition: 0,
		MinBytes:  1,
		MaxBytes:  10e6,
		MaxWait:   1 * time.Second,
	})

	defer reader.Close()

	var messages []string
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := 0; i < count; i++ {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			if err == context.DeadlineExceeded || err == context.Canceled {
				if i == 0 {
					return "Topic is empty or no new messages", nil
				}
				break
			}
			return "", fmt.Errorf("cannot read message: %v", err)
		}

		keyStr := ""
		if len(msg.Key) > 0 {
			keyStr = fmt.Sprintf("Key: %s, ", string(msg.Key))
		}

		timestamp := msg.Time.Format("2006-01-02 15:04:05")
		messageStr := fmt.Sprintf("[%s] %sValue: %s (Partition: %d, Offset: %d)",
			timestamp, keyStr, string(msg.Value), msg.Partition, msg.Offset)

		messages = append(messages, messageStr)
	}

	if len(messages) == 0 {
		return "Topic is empty or no new messages", nil
	}

	result := ""
	for i, msg := range messages {
		result += fmt.Sprintf("%d. %s\n", i+1, msg)
	}

	return result, nil
}
