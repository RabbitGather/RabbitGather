# RabbitGather Project Index

- RabbitGather-homepage
  - Host on CloudStorage

## Interest News

- RabbitGather-interest-reader
  - Provide API for getting news
  - Read needed news from Database

- RabbitGather-interest-manager
  - Provide API to control everything about interest news

- RabbitGather-interest-crawler
    - Crawl news from web
    - Save news html to CloudStorage

- RabbitGather-interest-parser
  - Parse news html from CloudStorage to structured data and save to CloudSQL


## Library

- RabbitGather-proto
  - Define proto file for all services

- go-meowalien-lib
  - golang tools 

## Infrastructure

- ArtifactRegistry
- CloudStorage
- CloudSQL

- GCP Pub/Sub
  - for microservice communication
- Kafka
  - for data streaming

- RabbitMQ