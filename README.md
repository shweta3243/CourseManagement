# 📚 Course Management System

A RESTful API for managing courses, subjects, enrollments, and learners — built with **Go**, **Gin**, and **MySQL**.

## ✨ Features

- **JWT Authentication** — Secure login/register with role-based access (Admin/Learner)
- **Course Management** — CRUD operations for courses (Admin only)
- **Subject Management** — Create subjects and assign them to courses
- **Enrollment System** — Learners can enroll in courses
- **Admin Dashboard** — Concurrent stats fetching using goroutines, channels & WaitGroup
- **CORS Support** — Configured for Vue.js frontend integration
- **Clean Architecture** — Repository → Service → Controller layered design with Dependency Injection

## 🏗️ Project Structure

```
course-management/
├── main.go                  # Application entry point
├── config/
│   └── database.go          # Database connection setup
├── controllers/             # HTTP request handlers
│   ├── auth_controller.go
│   ├── course_controller.go
│   ├── dashboard_controller.go
│   ├── enrollment_controller.go
│   └── subject_controller.go
├── middleware/
│   ├── jwt.go               # JWT authentication middleware
│   └── logging.go           # Request logging middleware
├── models/                  # Data structures
│   ├── course.go
│   ├── dashboard.go
│   ├── enrollment.go
│   ├── subject.go
│   └── user.go
├── repository/              # Database access layer
│   ├── course_repository.go
│   ├── enrollment_repository.go
│   ├── subject_repository.go
│   └── user_repository.go
├── routes/
│   └── routes.go            # Route definitions
├── services/                # Business logic layer
│   ├── auth_service.go
│   ├── course_service.go
│   ├── dashboard_service.go
│   ├── enrollment_service.go
│   └── subject_service.go
└── utils/
    ├── password.go          # Password hashing utilities
    └── token.go             # JWT token utilities
```

## 🚀 Getting Started

### Prerequisites

- Go 1.24+
- MySQL

### Run

```bash
# Set environment variables
export DB_DSN="user:password@tcp(localhost:3306)/course_db?parseTime=true"
export JWT_SECRET="your-secret-key"
export PORT="8080"

# Run the application
go run main.go
```

## 📡 API Endpoints

### Public
| Method | Endpoint         | Description       |
|--------|-----------------|-------------------|
| POST   | `/api/register` | Register new user |
| POST   | `/api/login`    | Login & get JWT   |

### Protected (JWT Required)
| Method | Endpoint                    | Description              | Role    |
|--------|----------------------------|--------------------------|---------|
| GET    | `/api/courses`             | List all courses         | Any     |
| GET    | `/api/courses/:id`         | Get course details       | Any     |
| GET    | `/api/courses/:id/subjects`| Get subjects for course  | Any     |
| POST   | `/api/enrollment`          | Enroll in a course       | Learner |
| GET    | `/api/my-enrollments`      | My enrollments           | Learner |

### Admin Only
| Method | Endpoint                                    | Description               |
|--------|---------------------------------------------|---------------------------|
| POST   | `/api/courses`                              | Create course             |
| PUT    | `/api/courses/:id`                          | Update course             |
| DELETE | `/api/courses/:id`                          | Delete course             |
| POST   | `/api/subjects`                             | Create subject            |
| GET    | `/api/subjects`                             | List all subjects         |
| DELETE | `/api/subjects/:id`                         | Delete subject            |
| POST   | `/api/course-subject`                       | Assign subject to course  |
| DELETE | `/api/course-subject/:courseId/:subjectId`  | Remove subject from course|
| GET    | `/api/learners`                             | List all learners         |
| GET    | `/api/enrollments`                          | List all enrollments      |
| GET    | `/api/dashboard`                            | Dashboard statistics      |

## 🛠️ Tech Stack

- **Go** (Gin framework)
- **MySQL** (via `sqlx`)
- **JWT** for authentication
- **bcrypt** for password hashing
- **Goroutines & Channels** for concurrent dashboard queries

## 📄 License

MIT
