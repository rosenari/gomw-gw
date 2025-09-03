# gomw-gw

High Performance Go Mini WebSocket Gateway

![Image](https://github.com/user-attachments/assets/3c181172-1e98-4eca-b43b-93082febd3f8)

## Environment Variable

| name | description | default | required |
|--------|------|--------|------|
| `LISTEN_ADDR` | Server Listen Port | `:8080` | ❌ |
| `ONCONNECT_URL` | Onconnect Webhook URL | - | ❌ |
| `DISCONNECT_URL` | Disconnect Webhook URL | - | ❌ |

## Build

### Local Build
```bash
go mod tidy

go build -o gomw-gw ./app/cmd/server

./gomw-gw
```

### Docker Build
```bash
docker build -t gomw-gw .
```

### Docker Run
```bash
docker run --name gomw-gw \
  -p 8080:8080 \
  -e ONCONNECT_URL=https://your-webhook.com/connect \
  -e DISCONNECT_URL=https://your-webhook.com/disconnect \
  gomw-gw
```

### Docker Compose
```yaml
docker-compose up
```

## API Endpoint

### WebSocket Connection
- **URL**: `/ws`
- **Protocol**: WebSocket
- **Description**: Websocket Connection Endpoint

### Send Message
- **URL**: `/send`
- **Method**: `POST`
- **Content-Type**: `application/json`

**Request Body:**
```json
{
  "connection_id": "uuid-string",
  "message": "Message Contents"
}
```

**Response:**
```json
{
  "success": true,
  "connection_id": "uuid-string",
  "message": "Message sent successfully"
}
```

### Environment Info
- **URL**: `/env`
- **Method**: `GET`
- **Description**: Get current environment configuration

**Response:**
```json
{
  "listen_address": ":8080",
  "on_connect_url": "https://your-webhook.com/connect",
  "on_disconnect_url": "https://your-webhook.com/disconnect"
}
```

### Health Check
- **URL**: `/health`
- **Method**: `GET`
- **Description**: Check service health status

**Response:**
```json
{
  "status": "healthy",
  "active_connections": 5,
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### Connection Status
- **URL**: `/status`
- **Method**: `GET`
- **Description**: Get active connection list

**Response:**
```json
{
  "total_connections": 2,
  "connections": [
    {
      "connection_id": "uuid-1",
      "client_ip": "192.168.1.100",
      "connected_at": "2024-01-01T10:00:00Z",
      "query_params": {"token": ["abc123"]}
    }
  ]
}
```

## Webhook Payload

### On Connect (ONCONNECT_URL)
```json
{
  "connection_id": "uuid-string",
  "client_ip": "192.168.1.100",
  "query_params": {"param1": ["value1"]},
  "timestamp": "2024-01-01T10:00:00Z",
  "server_ip": "10.0.1.100",
  "server_port": "8080"
}
```

### On Disconnect (DISCONNECT_URL)
```json
{
  "connection_id": "uuid-string",
  "client_ip": "192.168.1.100",
  "timestamp": "2024-01-01T10:05:00Z",
  "server_ip": "10.0.1.100",
  "server_port": "8080"
}
```
