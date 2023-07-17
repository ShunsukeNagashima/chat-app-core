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
├── pkg                  # houses the main code of the application
│   ├── apperror         # manage specific errors of application
│   ├── clock            # provides functionality related to time
│   ├── domain           # contains the domain models, repository interfaces, and use cases
│   │   ├── model
│   │   ├── repository
│   │   └── usecase
│   ├── infra            # implements concrete details like persistence.
│   │   ├── auth
│   │   └── repository
│   ├── interface        # handles input and output of data
│   │   ├── controller
│   │   └── route
│   └── usecase          # execute the business logic
```

## Prerequisites

To run this project locally, you will need to install a few tools.

### Docker Desktop

This application uses Docker and docker-compose. To run these locally, you will need to install Docker Desktop.

Please refer to the official documentation for installation instructions: [Docker Desktop Installation Guide](https://docs.docker.com/desktop/)

### GolangCI-Lint

GolangCI-Lint is a tool that statically analyzes your Go code to look for issues. You will need to install this tool.

Please refer to the official documentation for installation instructions: [GolangCI-Lint Installation Guide](https://golangci-lint.run/usage/install/)

### Air

Air is used for hot reloading in this application. You will need to install this tool as well.

Please refer to the official repository for installation instructions: [Air GitHub Repository](https://github.com/cosmtrek/air)

## Set Env and Credentials
Before running the application, you will need to set up your environment variables and Firebase credentials.

### Environment Variables
```
cp .env.example .env.local
```

### Firebase Credentials
This application uses Firebase for user authentication. To set up your Firebase credentials:

1. Create a new file named `secrets/firebase-credentials.json`.
2. Follow the instructions in the [Firebase Documentation](https://firebase.google.com/docs/admin/setup) to get your Firebase credentials.
3. Paste your Firebase credentials into the `secrets/firebase-credentials.json` file.

## Running the Application

This application uses `make` commands for easy operation. Here are the basic commands you need to know:

- To start the application, run:
```
make up
```

- To stop the application, run:
```
make down
```

- To migrate the database, run:
```
make migrate
```

For other commands such as testing, please check the `Makefile`.
