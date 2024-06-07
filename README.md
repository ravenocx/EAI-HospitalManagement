# Final Project EAI - Hospital Management

### Prerequisites
- Docker
- A Linux environment (WSL / MacOS should be fine)


### How to run
#### Build the container for database
```bash
docker-compose up
```
#### Build the image for each service
Navigate to the directory for each service:
```bash
cd [SERVICE_NAME]
```

Build the image:
```bash
docker build -t [SERVICE_NAME] .
```

Repeat until all service image is built

#### How to run the service:
```bash
docker run -d [SERVICE_NAME]
```

Note: Use the [SERVICE_NAME] for each service name
