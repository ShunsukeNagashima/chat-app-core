# Real-time Chat-API

## Overview
This is the backend REST API for a real-time chat application.


## Technologies Used
- Go (latest stable version), Gin Framework
- Amazon DynamoDB
- ElasticSearch
- Docker

## Directory Structure
Below represents the directory structure of the project:

```
.
├── cmd                  # contains the entry point to the project
│   └── main.go
├── pkg                  # houses the main code of the application.
│   ├── apperror         # manage specific errors of application
│   ├── clock            # provides functionality related to time
│   ├── domain           # contains the domain models, repository interfaces, and use cases
│   │   ├── model
│   │   ├── repository
│   │   └── usecase
│   ├── infra            # implements concrete details like persistence.
│   │   ├── auth
│   │   └── repository
│   ├── interface
│   │   ├── controller
│   │   └── route
│   └── usecase
