# Async CSV Importer with Worker Pool

**Project**: `csv-importer` — a Go-based project demonstrating asynchronous processing of uploaded CSV files using Redis as a queue and PostgreSQL as a database.

📌 **Primary Goal**: Accept CSV uploads via an HTTP API, enqueue processing jobs in Redis, and process these files in the background (worker pool), loading their contents into PostgreSQL.

---

## 📋 Repository Structure

```
async_csv_importer/
├── cmd/
│   └── api/
│       └── main.go           # HTTP server: accepts CSV and enqueues job in Redis
├── internal/
│   └── handlers/
│       └── upload.go         # Logic for receiving file and publishing job to Redis
├── Dockerfile                # Builds Go application and packages it into a minimal image
├── docker-compose.yml        # Configures Docker Compose for API, Redis, PostgreSQL
├── go.mod                    # Go modules file
└── README.md                 # This is what you are reading now
```

> **Note**: This version only includes the API that publishes jobs to Redis. The worker service is not implemented yet but can be added easily.

---

## 🚀 Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/nikitam/async_csv_importer.git
cd async_csv_importer
```

### 2. Configure Environment Variables for PostgreSQL

Create a `.env` file in the project root containing:

```dotenv
POSTGRES_USER=postgres
POSTGRES_PASSWORD=secret
POSTGRES_DB=default
```

> These variables are used in `docker-compose.yml` for the PostgreSQL container.

### 3. Start Services with Docker Compose

```bash
docker compose up -d
```

* **api**: Runs the HTTP server on port 8080
* **redis**: Redis 8, port 6379
* **postgres**: PostgreSQL 17, port 5432, database `importer`

### 4. Verify Running Containers

```bash
docker compose ps
```

Expected output (example):

```
    Name                      Command               State           Ports
----------------------------------------------------------------------------
csv_importer_api_1      /server               Up      0.0.0.0:8080->8080/tcp
csv_importer_redis_1    docker-entrypoint ...  Up      0.0.0.0:6379->6379/tcp
csv_importer_postgres_1 docker-entrypoint ...  Up      0.0.0.0:5432->5432/tcp
```

### 5. Test CSV Upload

1. Create a sample CSV in the project root, e.g.:

   ```bash
   echo "id,name,email" > test.csv
   echo "1,Nikita,nikita@example.com" >> test.csv
   ```
2. Send an HTTP request to upload it:

   ```bash
   curl -i -X POST \
        -F "file=@test.csv" \
        http://localhost:8080/upload
   ```
3. You should receive:

   ```txt
   HTTP/1.1 202 Accepted
   File test.csv uploaded successfully
   ```

### 6. Verify Job Queue in Redis

```bash
# Find the Redis container ID (for example, async_csv_importer_redis_1)
redis_container=$(docker ps -qf "ancestor=redis:8")

docker exec -it $redis_container redis-cli LRANGE import_jobs 0 -1
```

You should see a JSON object representing the job:

```json
1) "{\"file_path\":\"/tmp/1678901234567890000_test.csv\",\"uploaded\":1678901234567890000}"
```

If a job appears, the setup is correct.

---

## 📝 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

> **Contact**: Nikita Makovei ([hello.nikita.makovei@gmail.com](mailto:hello.nikita.makovei@gmail.com))
