# 🍏 Fruit Store Challenge

The goal of this challenge is to create a complete CRUD for a fruit store, using Go, Docker Compose, JWT for authentication/authorization, dependency injection, and caching with Redis.

## 📦 Architecture

* API (`fruit-api`)
  * CRUD for fruits (directly stored in the database)
  * User registration → instead of writing directly to the database, it publishes an event `users.created` to Kafka

* User Consumer (`user-consumer`)
  * Kafka consumer that listens to the topic `users.created`
  * Persists users into the same PostgreSQL database used by the API (for practicality)

* Database (PostgreSQL)
  * Shared between API and Consumer
  * Stores both fruits and users

* Redis
  * Used as a cache for faster queries (e.g., fruits lookup)

* Kafka + Zookeeper
  * Messaging broker (topic `users.created`)
  * Ensures that user-created events are never lost, even if the consumer is down

* Kafka UI
  * Simple web interface to inspect topics and messages

## 🚀 Running the project
```
docker-compose up --build
```

Available services:

API → http://localhost:8080

Kafka UI → http://localhost:8082

Postgres → localhost:5432 (DB_USER, DB_PASSWORD, DB_NAME from .env)

## 🔑 Main Endpoints
Fruits (`/fruits`)
* GET /fruits → List all fruits
* GET /fruits/:id → Get a fruit by ID
* POST /fruits → Create a fruit
* PUT /fruits/:id → Update a fruit
* DELETE /fruits/:id → Delete a fruit

Users (`/users`)
* GET /users → List all users (stored by the consumer)
* POST /users → Create a user → publishes an event to Kafka → user-consumer listens and stores the user in the database

Authentication (`/auth`)

* POST /auth/login → Authenticate a user and return a JWT

  * Example payload:
    ```
    {
        "username": "admin",
        "password": "admin"
    }
    ```
 
  * Use this seed user for initial access:
    * Username: admin
    * Password: admin

## 🔐 Permissions and Error Responses

401 Unauthorized → request without valid authentication (missing/invalid token)

404 Not Found → resource does not exist

400 Bad Request → invalid input payload

## 📝 Postman Collection

A Postman collection is provided for testing all API endpoints easily.

* Importing the collection:
  1. Download the `fruit-store-challenge.postman_collection.json` file.
  2. Open Postman → click Import → select Upload Files → choose the JSON file.
  3. All endpoints (Fruits, Users, Auth) will be imported with proper paths and example payloads.

* Seed user:
  * Use the default admin credentials in Postman to authenticate and obtain a JWT token:
    * Username: admin
    * Password: admin

* Notes:
  * Use the token returned from `/auth/login` for protected endpoints.
  * The collection simplifies testing CRUD operations and Kafka-related user creation.

## 📂 Project Structure
```
.
├── cmd
│   ├── api
│   │   └── main.go                                # entrypoint for the API, sets up Gin router and server
│   └── user-consumer
│       └── main.go                                # entrypoint for the Kafka consumer, listens to users.created topic and stores in DB
├── internal
│   ├── app
│   │   └── app.go                                 # configures the Gin app, applies middleware and routes
│   ├── config
│   │   └── config.go                              # reads environment variables and app configuration
│   ├── database
│   │   └── database.go                            # Postgres connection via GORM
│   ├── middleware
│   │   └── auth.go                                # JWT authentication and authorization
│   ├── models
│   │   ├── fruit.go                               # Fruit entity and data model
│   │   └── user.go                                # User entity and data model
│   ├── routes
│   │   ├── auth.go                                # login and authentication endpoints
│   │   ├── fruits.go                              # CRUD endpoints for fruits
│   │   └── users.go                               # CRUD/listing endpoints for users (admin only)
│   ├── services
│   │   ├── hash.go                                # password hashing functions
│   │   ├── jwt.go                                 # JWT generation and validation
│   │   ├── kafka.go                               # Kafka producer/consumer helpers
│   │   └── redis.go                               # Redis caching helpers
│   └── util
│       └── respond.go                             # JSON response helpers (OK, Error, BadRequest, Unauthorized, etc.)
├── docker-compose.yml                             # orchestrates API, Kafka, Redis, Postgres, and consumer
├── Dockerfile                                     # Dockerfile to build the project image
├── go.mod                                         # Go module and dependencies
├── go.sum                                         # checksums for dependencies
├── README.md                                      # project documentation
└── fruit-store-challenge.postman_collection.json # Postman collection to test all API endpoints
```

## 📌 Design Decisions

1. Single database (Postgres):
   * Both API and Consumer share the same DB for practicality in this challenge.
   * In production, this could be split (e.g., isolated user service).
   * PostgreSQL was chosen because:
     * It is a robust relational database, easy to set up and widely used.
     * Supports transactions, foreign keys, and relational integrity, which are useful for this kind of CRUD project.
     * Integrates seamlessly with GORM, the Go ORM used in this project.
2. Event-driven with Kafka:
   * Chosen to ensure persistence and durability of events.
   * Prevents message loss when consumers are temporarily unavailable.

3. Response helpers (util/respond.go):
   * Provide consistency in error handling (401, 404, etc).