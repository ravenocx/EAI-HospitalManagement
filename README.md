# Final Project EAI - Hospital Management

### Prerequisites
- Docker
- A Linux environment (WSL / MacOS should be fine)


### Setting up database
### Install golang migrate
```bash
$ wget -qO- https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz
$ sudo mv migrate.linux-amd64 /usr/local/bin/migrate
$ sudo chmod +x /usr/local/bin/migrate
```
### Verify the installation
```bash
migrate -version
```

### Run the migration 
Note : run the migration for AuthService->Patient->MedicalRecord
```bash
migrate -database "postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable" -path [SERVICE_NAME]//db/migrations -verbose up
```
Repeat until all service table is created


### How to run
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
