# RabbitGather Project Index

## Service

### Frontend
- RabbitGather-homepage
  - Host on `CloudStorage`
  - Write in `Vue.js`


### Backend
- Provide API through gRPC
- Microservice architecture

#### Interest Service
Offers gRPC API.

Handles web crawling, monitoring, and scheduled tasks, maintaining tasks specified by users.
- Dependency
  - MongoDB


## Library
  - RabbitGather-proto
    - Define proto file for all services

  - go-meowalien-lib
    - golang tools 


## CI/CD
  - RabbitGather-cicd
    - Manage all CI/CD workflow


## Infrastructure
- Envoy Proxy ( GRPC-web proxy )
  - https://github.com/envoyproxy/envoy

- Centrifugo ( Websocket-grpc proxy )
  - https://centrifugal.dev

- RabbitMQ-Cluster on GCP
  - For decouple services
  - For task queue
  - For message broadcast

- ArtifactRegistry
  - Use for store docker images

- GitHubAction
  - Build and push docker image to `ArtifactRegistry`

- CloudStorage
  - Use for store and host web page and static files
  - Use for store user uploaded files






[//]: # ()
[//]: # (- MySQL&#40;CloudSQL&#41;)

[//]: # (  - Use for store transactional data)

[//]: # (  - Use to store key data that need to be consistent, reliable, and ACID)

[//]: # ()
[//]: # (- Firestore)

[//]: # (  - Use for store semi-structured data for fast query and search)

[//]: # (  - Use for store user generated data)

[//]: # ()
[//]: # (- Kafka/CloudPubSub/CloudDataFlow)

[//]: # (    - Use for analyzing big data)



