# Go SOA Example

This is a simple Service Oriented Architecture (SOA) example built with Go. The project demonstrates key SOA principles including:

- Service autonomy
- Service loose coupling
- Service contracts
- Service composability
- Service reusability

## Services

1. **Product Service** - Manages product information (port 8081)
2. **User Service** - Manages user accounts (port 8082)
3. **Order Service** - Handles order processing (port 8083)
4. **API Gateway** - Routes requests to appropriate services (port 8080)

## Running the Project

### Prerequisites
- Go 1.24+
- Docker and Docker Compose (optional)

### Running with Docker

```bash
docker-compose up --build