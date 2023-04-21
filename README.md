# ChatGPT Clone integration project

![image](https://user-images.githubusercontent.com/15023064/233719766-0c9fe7d8-3afe-44d2-a79c-092e3f3aebbb.png)

## Technologies used
1. Golang
2. Next.JS
3. OpenAI API
4. Prisma
5. SQLC
6. gRPC

## Getting started with protobuf

Protobuf configuration: https://grpc.io/docs/languages/go/quickstart/

# System Architecture

This is a simple system architecture for a frontend web application that communicates with BFF (backend for frontend) and this backend communicates with a microservice and the OpenAI API.

## Architecture Diagram

```mermaid
flowchart LR
    A[Next.JS] -->|HTTP| B(BFF)
    B --> C{Get all messages stored/Read realtime stream}
    C -->|Stored| D[(MySQL)]
    C -->|Realtime| E[Golang Microservice]
    E -->|OpenAPI Client| F[OpenAI API]
```


