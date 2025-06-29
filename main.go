package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"Kafka MCP Server",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	listTopicsTool := mcp.NewTool("list_topics",
		mcp.WithDescription("List all topics in Kafka broker"),
		mcp.WithString("broker",
			mcp.Required(),
			mcp.Description("Kafka broker address (e.g. localhost:9092)"),
		),
	)

	s.AddTool(listTopicsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		broker, err := request.RequireString("broker")
		if err != nil {
			return mcp.NewToolResultError("Missing broker information"), nil
		}

		topics, err := listKafkaTopics(broker)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error listing topics: %v", err)), nil
		}

		if len(topics) == 0 {
			return mcp.NewToolResultText("No topics found in broker"), nil
		}

		result := "List of topics:\n"
		for i, topic := range topics {
			result += fmt.Sprintf("%d. %s\n", i+1, topic)
		}

		return mcp.NewToolResultText(result), nil
	})

	createTopicTool := mcp.NewTool("create_topic",
		mcp.WithDescription("Create a new topic in Kafka"),
		mcp.WithString("broker",
			mcp.Required(),
			mcp.Description("Kafka broker address"),
		),
		mcp.WithString("topic",
			mcp.Required(),
			mcp.Description("Name of the topic to create"),
		),
		mcp.WithNumber("partitions",
			mcp.Description("Number of partitions (default: 1)"),
		),
	)

	s.AddTool(createTopicTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		broker, err := request.RequireString("broker")
		if err != nil {
			return mcp.NewToolResultError("Missing broker information"), nil
		}

		topicName, err := request.RequireString("topic")
		if err != nil {
			return mcp.NewToolResultError("Missing topic name"), nil
		}

		partitions := 1
		if p := request.GetFloat("partitions", 1.0); p > 0 {
			partitions = int(p)
		}

		err = createKafkaTopic(broker, topicName, partitions)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error creating topic: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("✅ Successfully created topic '%s' with %d partitions", topicName, partitions)), nil
	})

	// Tool 3: Produce Message
	produceTool := mcp.NewTool("produce_message",
		mcp.WithDescription("Send message to Kafka topic"),
		mcp.WithString("broker",
			mcp.Required(),
			mcp.Description("Kafka broker address"),
		),
		mcp.WithString("topic",
			mcp.Required(),
			mcp.Description("Topic name"),
		),
		mcp.WithString("message",
			mcp.Required(),
			mcp.Description("Message content"),
		),
		mcp.WithString("key",
			mcp.Description("Key for message (optional)"),
		),
	)

	s.AddTool(produceTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		broker, err := request.RequireString("broker")
		if err != nil {
			return mcp.NewToolResultError("Missing broker information"), nil
		}

		topic, err := request.RequireString("topic")
		if err != nil {
			return mcp.NewToolResultError("Missing topic name"), nil
		}

		message, err := request.RequireString("message")
		if err != nil {
			return mcp.NewToolResultError("Missing message content"), nil
		}

		key := request.GetString("key", "") // Key là optional

		err = produceKafkaMessage(broker, topic, key, message)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error sending message: %v", err)), nil
		}

		keyInfo := ""
		if key != "" {
			keyInfo = fmt.Sprintf(" with key '%s'", key)
		}

		return mcp.NewToolResultText(fmt.Sprintf("✅ Successfully sent message to topic '%s'%s", topic, keyInfo)), nil
	})

	// Tool 4: Consume Messages
	consumeTool := mcp.NewTool("consume_messages",
		mcp.WithDescription("Read messages from Kafka topic"),
		mcp.WithString("broker",
			mcp.Required(),
			mcp.Description("Kafka broker address"),
		),
		mcp.WithString("topic",
			mcp.Required(),
			mcp.Description("Topic name"),
		),
		mcp.WithNumber("count",
			mcp.Description("Number of messages to read (default: 10)"),
		),
	)

	s.AddTool(consumeTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		broker, err := request.RequireString("broker")
		if err != nil {
			return mcp.NewToolResultError("Missing broker information"), nil
		}

		topic, err := request.RequireString("topic")
		if err != nil {
			return mcp.NewToolResultError("Missing topic name"), nil
		}

		count := 10
		if c := request.GetFloat("count", 10.0); c > 0 {
			count = int(c)
		}

		messages, err := consumeKafkaMessages(broker, topic, count)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error reading messages: %v", err)), nil
		}

		if messages == "" {
			return mcp.NewToolResultText(fmt.Sprintf("No messages found in topic '%s'", topic)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("📨 Messages in topic '%s':\n\n%s", topic, messages)), nil
	})

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("❌ Error server: %v", err)
	}
}
