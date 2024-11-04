Context and scope

This document describes the design of a Workout Tracker application built using Golang and PostgreSQL. The application allows users to create, update, and delete workouts, as well as generate reports for their workouts. It includes JWT authentication for user sessions and provides a RESTful API for interacting with the system.

The application is structured as a backend service, with no frontend implementation included in the current scope. It is designed to be run as a standalone server, exposing HTTP endpoints for client applications to consume.

Goals and non-goals

Goals:

Provide user authentication and authorization

Allow CRUD operations for workouts and exercises

Generate reports on user workouts

Ensure data persistence using PostgreSQL

Implement JWT-based session management

Non-goals:

Implement a frontend user interface

Provide real-time synchronization of workout data

Integrate with external fitness tracking devices or APIs

Implement advanced analytics or machine learning features

Design

System-context-diagram

sequenceDiagram
    participant Client
    participant API
    participant AuthService
    participant WorkoutService
    participant Database

    Client->>API: HTTP Request
    API->>AuthService: Authenticate User
    AuthService->>Database: Verify Credentials
    Database-->>AuthService: User Data
    AuthService-->>API: JWT Token
    API->>WorkoutService: Process Workout Data
    WorkoutService->>Database: CRUD Operations
    Database-->>WorkoutService: Workout Data
    WorkoutService-->>API: Response
    API-->>Client: HTTP Response

APIs

The application exposes the following main API endpoints:

Authentication:

POST /signup: Register a new user

POST /login: Authenticate a user and receive a JWT token

Workouts:

POST /workouts/create: Create a new workout

GET /workouts: Retrieve user's workouts

PUT /workouts/update: Update an existing workout

DELETE /workouts/delete: Delete a workout

GET /workouts/report: Generate a workout report

Exercises:

GET /exercises: Retrieve all exercises

GET /exercises/{id}: Retrieve a specific exercise

Data storage

The application uses PostgreSQL as its primary data store. The main tables are:

users: Stores user information

workouts: Stores workout details

exercises: Stores exercise information

workout_exercises: Junction table linking workouts and exercises

The database schema is managed using migrations, allowing for easy schema updates and version control.

Code structure

The application follows a typical Go project structure:

cmd/main.go: Entry point of the application

config/: Configuration management

handler/: HTTP request handlers

model/: Data models and business logic

repository/: Database interaction layer

service/: Business logic layer

util/: Utility functions

router/: API route definitions

This structure promotes separation of concerns and makes the codebase more maintainable and testable.

Alternatives considered

Using a NoSQL database:

Pros: Flexibility in data modeling, potential for better scalability

Cons: Less structured data, potential for data inconsistency, less familiar to many developers

Decision: Chose PostgreSQL for its ACID compliance and structured data model, which fits well with the workout tracking use case

Implementing GraphQL instead of REST:

Pros: More flexible querying, reduced over-fetching and under-fetching of data

Cons: Steeper learning curve, potential overkill for a relatively simple API

Decision: Chose REST for its simplicity and wide adoption, making it easier for clients to integrate

Using a third-party authentication service:

Pros: Reduced development time, potentially more secure

Cons: Dependency on external service, potential cost implications

Decision: Implemented custom JWT-based authentication for better control and to avoid external dependencies

By choosing PostgreSQL, REST API, and custom JWT authentication, the design prioritizes simplicity, familiarity, and control over the application's core functionality. This approach allows for easier development, maintenance, and potential future expansion of the Workout Tracker application.
