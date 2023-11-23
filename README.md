
# Ports-Service

## Introduction
This project is a Go-based microservices application. It demonstrates proficiency in hexagonal architecture, gRPC, Domain-Driven Design (DDD), and various streaming and data persistence techniques. The application provides two key adapter implementations: a gRPC streaming API and filesystem-based data streaming.

## Features
- **gRPC Streaming API**: Implements a gRPC-based streaming service, allowing efficient data transfer over the network.
- **Filesystem Data Streaming**: Supports streaming data from the filesystem, showcasing adaptability to different data sources.
- **In-Memory Database**: Utilizes an in-memory database for temporary data storage, ensuring fast data access.
- **Hexagonal Architecture**: Adheres to hexagonal architecture principles, promoting loose coupling and high modularity.
- **Domain-Driven Design (DDD)**: Implements DDD principles, aligning the solution with business requirements.
- **Debugging Capabilities**: Includes a feature to test the presence of specific keys by setting the `debugkey` flag or environment variable, which periodically attempts to retrieve the key from storage.
- **Persistent Data Storage**: Currently supports only persistent data storage without retrieval capabilities.
- **Docker Support**: Includes a Dockerfile for easy containerization and deployment.

## Architecture
The project is structured following the hexagonal architecture model, divided into:

- **Ports**: Defined in `port-for-ships.go` and `port-for-ships_repository.go`, representing the primary interfaces for our application.
- **Adapters**: Including `in-memory-db.go` for in-memory storage, `grpc-streaming.go` and `filesystem-streamer.go` for data streaming, and `store.go` for data persistence.
- **Domain**: Core business logic encapsulated within the service implementation.
