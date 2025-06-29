# Kafka MCP Server

A Model Context Protocol (MCP) Server written in Go to integrate Apache Kafka with Cursor IDE and other AI assistants.

## 🚀 Features

This MCP Server provides 4 main tools to interact with Kafka:

- **📋 list_topics** - List all topics in Kafka broker
- **➕ create_topic** - Create new topic in Kafka  
- **📤 produce_message** - Send message to Kafka topic
- **📥 consume_messages** - Read messages from Kafka topic

## 📋 System Requirements

- Go 1.21+
- Apache Kafka (can run via Docker)
- Cursor IDE (for MCP integration)

## 🛠️ Installation

### 1. Clone repository

```bash
git clone <repository-url>
cd kafka-mcp
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Build MCP server

```bash
go build -o mcp-kafka .
```

### 4. Start Kafka (if not already running)

Using Docker Compose:

```bash
# Create docker-compose.yml file
version: '3'
services:
  zookeeper:
    image: confluentinc/cp-zookeeper:7.5.0
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000

  broker:
    image: confluentinc/cp-kafka:7.5.0
    ports:
      - "9092:9092"
    depends_on:
      - zookeeper
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: 'zookeeper:2181'
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_INTERNAL:PLAINTEXT
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://localhost:9092,PLAINTEXT_INTERNAL://broker:29092
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1

# Start services
docker-compose up -d
```

## ⚙️ MCP Configuration in Cursor

### 1. Create MCP configuration file

Create file `~/.cursor/mcp.json`:

```json
{
  "mcpServers": {
    "kafka": {
      "command": "/path/to/your/kafka-mcp/mcp-kafka",
      "args": []
    }
  }
}
```

### 2. Restart Cursor

After configuration, restart Cursor IDE to load the MCP server.

## 🎯 Usage

### In Cursor IDE

After successful configuration, you can use these commands in chat:

```
// List topics
list topics in kafka broker localhost:9092

// Create new topic  
create topic "my-new-topic" in kafka broker localhost:9092

// Send message
send message "Hello World" to kafka topic "my-topic" on localhost:9092

// Read messages
read 5 messages from kafka topic "my-topic" on localhost:9092
```

### Manual Testing

You can test the MCP server directly:

```bash
# Start server
./mcp-kafka

# In another terminal, send JSON-RPC request
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}' | ./mcp-kafka
```

## 📚 API Reference

### Tools

#### 1. list_topics
```json
{
  "name": "list_topics",
  "description": "List all topics in Kafka broker",
  "inputSchema": {
    "type": "object",
    "properties": {
      "broker": {
        "type": "string",
        "description": "Kafka broker address (e.g. localhost:9092)"
      }
    },
    "required": ["broker"]
  }
}
```

#### 2. create_topic
```json
{
  "name": "create_topic", 
  "description": "Create a new topic in Kafka",
  "inputSchema": {
    "type": "object",
    "properties": {
      "broker": {
        "type": "string",
        "description": "Kafka broker address"
      },
      "topic": {
        "type": "string", 
        "description": "Name of the topic to create"
      },
      "partitions": {
        "type": "number",
        "description": "Number of partitions (default: 1)"
      }
    },
    "required": ["broker", "topic"]
  }
}
```

#### 3. produce_message
```json
{
  "name": "produce_message",
  "description": "Send message to Kafka topic", 
  "inputSchema": {
    "type": "object",
    "properties": {
      "broker": {
        "type": "string",
        "description": "Kafka broker address"
      },
      "topic": {
        "type": "string",
        "description": "Topic name" 
      },
      "message": {
        "type": "string",
        "description": "Message content"
      },
      "key": {
        "type": "string",
        "description": "Key for message (optional)"
      }
    },
    "required": ["broker", "topic", "message"]
  }
}
```

#### 4. consume_messages
```json
{
  "name": "consume_messages",
  "description": "Read messages from Kafka topic",
  "inputSchema": {
    "type": "object", 
    "properties": {
      "broker": {
        "type": "string",
        "description": "Kafka broker address"
      },
      "topic": {
        "type": "string",
        "description": "Topic name"
      },
      "count": {
        "type": "number",
        "description": "Number of messages to read (default: 10)"
      }
    },
    "required": ["broker", "topic"]
  }
}
```

## 🔧 Development

### Project Structure

```
kafka-mcp/
├── main.go           # Main MCP server
├── kafka_client.go   # Kafka client functions
├── go.mod           # Go dependencies
├── go.sum           # Go dependencies checksum
└── README.md        # Documentation
```

### Dependencies

- `github.com/mark3labs/mcp-go` - MCP protocol implementation
- `github.com/segmentio/kafka-go` - Kafka client for Go

### Build and Test

```bash
# Build
go build -o mcp-kafka .

# Test connection
echo '{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}' | ./mcp-kafka
```

## 🐛 Troubleshooting

### MCP Server Connection Issues

1. Check binary path in `mcp.json`
2. Ensure binary has execute permissions: `chmod +x mcp-kafka`
3. Test server manually before configuring Cursor

### Kafka Connection Errors

1. Check if Kafka is running: `docker ps | grep kafka`
2. Test connection: `telnet localhost 9092`
3. Check logs: `docker logs <kafka-container-id>`

### Tools Not Showing in Cursor

1. Check logs in Cursor Developer Tools
2. Restart Cursor after changing MCP configuration
3. Ensure JSON format in `mcp.json` is correct

## 📝 License

MIT License

## 🤝 Contributing

1. Fork the repository
2. Create feature branch: `git checkout -b feature/amazing-feature`
3. Commit changes: `git commit -m 'Add amazing feature'`
4. Push to branch: `git push origin feature/amazing-feature`
5. Create Pull Request

## 📞 Support

If you encounter any issues, please create an issue in this repository with detailed information about the error and your environment.
