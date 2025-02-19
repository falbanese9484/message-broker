# Message Broker Service

https://github.com/user-attachments/assets/309f9592-0008-49ac-901a-c944c9d70510

## Overview

The Message Broker Service is a centralized service designed to operate within a Kubernetes environment. It acts as a message broker for event-based and queued processes, facilitating communication between various components of a distributed system.

## Objectives

- **Lightweight JSON Brokering**: Supports JSON payloads for efficient message brokering.
- **Flexible TCP Protocol**: Provides a responsive TCP protocol adaptable to multiple languages, including PHP, Node.js, and Python.

## Architecture

The service is designed to be lightweight and scalable, leveraging Kubernetes for orchestration and management. It supports both event-based and queued message processing, providing a robust solution for distributed systems.

## Wish List on the Todo List
- **Built-in Frontend**: Includes a frontend interface for visibility and monitoring of message flows.
- **Fault Tolerance and Self-Healing**: Ensures high availability and reliability through fault tolerance and self-healing mechanisms.
- **HTTP Task Acceptance**: Capable of accepting HTTP requests for task submission from external sources.

## Worker Flow

### External Task Process Listener

1. **Worker Registration**: External systems call the service to register a worker, enabling it to start receiving tasks.
2. **Worker Initialization**: The system registers the worker and initiates a goroutine to send tasks to the client.
3. **Task Reception**: Each client acts as a TCP server, receiving tasks and performing designated actions.

### Register Worker

- **Queue ID**: Identifies the queue to which the worker is assigned.
- **TCP Connection**: Specifies the external TCP address for task delivery.

## Usage

1. **Deploy the Service**: Deploy the message broker service within your Kubernetes cluster.
2. **Register Workers**: Use the provided API to register workers and specify their queue assignments and TCP connections.
3. **Submit Tasks**: Send HTTP requests to the service to submit new tasks for processing.
4. **Monitor**: Use the built-in frontend to monitor task processing and system health.

## Contributing

We welcome contributions to improve the Message Broker Service. Please follow the contribution guidelines outlined in the repository.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.
