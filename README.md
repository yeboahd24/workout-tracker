# Workout Tracker

PROJECT_URL: `https://roadmap.sh/projects/fitness-workout-tracker`

## Introduction

This is a simple workout tracker application built using Golang and PostgreSQL. It allows users to create, update, and delete workouts, as well as generate reports for their workouts. The application also includes JWT authentication for user sessions.

## Features

- User registration and login
- CRUD operations for workouts and exercises
- JWT authentication for user sessions
- Report generation for workouts

## Installation

To install the application, follow these steps:

1. Clone the repository:

```bash
git clone https://github.com/yeboahd24/workout-tracker.git
```

2. Change into the project directory:

```bash
cd workout-tracker
```

3. Build the application:

```bash
go build -o workout-tracker cmd/main.go
```

4. Run the application:

```bash
./workout-tracker
```

## Usage

### User Registration and Login

To register a new user, send a POST request to the `/signup` endpoint with the following JSON payload:

```json
{
  "username": "johndoe",
  "email": "johndoe@example.com",
  "password": "password123"
}
```

To login to the application, send a POST request to the `/login` endpoint with the following JSON payload:

```json
{
  "username": "johndoe",
  "password": "password123"
}
```

### CRUD Operations for Workouts and Exercises

#### Create a Workout

To create a new workout, send a POST request to the `/workouts/create` endpoint with the following JSON payload:

```json
{
  "name": "My Workout",
  "description": "A sample workout",
  "scheduled_for": "2023-01-01T00:00:00Z",
  "exercises": [
    {
      "exercise_id": 1,
      "sets": 3,
      "reps": 10,
      "weight": 70
    },
    {
      "exercise_id": 2,
      "sets": 3,
      "reps": 10,
      "weight": 70
    }
  ]
}
```

#### Get Workouts by User

To get all workouts for a user, send a GET request to the `/workouts` endpoint with the following query parameters:

- `start_date`: The start date of the workouts to fetch (in the format "YYYY-MM-DD").
- `end_date`: The end date of the workouts to fetch (in the format "YYYY-MM-DD"). If not provided, the current date will be used.

Example:

```bash
curl -X GET "http://localhost:8080/workouts?start_date=2023-01-01&end_date=2023-01-31"
```

#### Update a Workout

To update an existing workout, send a PUT request to the `/workouts/update` endpoint with the following JSON payload:

```json
{
  "id": 1,
  "name": "My Updated Workout",
  "description": "An updated workout",
  "scheduled_for": "2023-01-01T00:00:00Z",
  "exercises": [
    {
      "exercise_id": 1,
      "sets": 3,
      "reps": 10,
      "weight": 70
    },
    {
      "exercise_id": 2,
      "sets": 3,
      "reps": 10,
      "weight": 70
    }
  ]
}
```

#### Delete a Workout

To delete a workout, send a DELETE request to the `/workouts/delete` endpoint with the following query parameter:

- `id`: The ID of the workout to delete.

Example:

```bash
curl -X DELETE "http://localhost:8080/workouts/delete?id=1"
```

#### Generate a Workout Report

To generate a workout report, send a GET request to the `/workouts/report` endpoint with the following query parameters:

- `start_date`: The start date of the workouts to fetch (in the format "YYYY-MM-DD").
- `end_date`: The end date of the workouts to fetch (in the format "YYYY-MM-DD"). If not provided, the current date will be used.

Example:

```bash
curl -X GET "http://localhost:8080/workouts/report?start_date=2023-01-01&end_date=2023-01-31"
```

There is a postman collection available in the `postman` directory for testing the application.

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request on the GitHub repository.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.
