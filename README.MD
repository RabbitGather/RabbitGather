# RabbitGather Project Index

- RabbitGather-homepage
  - Host on `CloudStorage`
  - Write in `Vue.js`

## Interest News

- RabbitGather-interest-reader
  - Provide API for getting data from `Firestore`

- RabbitGather-interest-manager
  - Provide GRPC API for update, delete, create news or tasks
  - Create a task in `Firestore` and write down what to crawl
  - Push taskID(crawl) to `Pub/Sub`

- RabbitGather-interest-crawler
  - Consume task from `Pub/Sub`
  - Read task from `Firestore`
  - Execute task
  - Save the result to `CloudStorage`
  - Set result url to `Firestore` 
  - Push taskID(parse) to `Pub/Sub`

- RabbitGather-interest-parser
  - Consume task from `Pub/Sub`
  - Read task from `Firestore`
  - Parse news html from `CloudStorage` to structured data and save to `Firestore`


## Library

- RabbitGather-proto
  - Define proto file for all services

- go-meowalien-lib
  - golang tools 

## Infrastructure

- ArtifactRegistry
- CloudStorage
- CloudSQL
- Firestore
- Pub/Sub

- Kafka
- RabbitMQ

- RabbitMQ