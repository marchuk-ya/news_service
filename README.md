# News Service

A modern news management service built with Go, HTMx, and MongoDB.

## Features

- Create, read, update, and delete news articles
- Server-side rendered views with HTMx for smooth interactions
- MongoDB for data storage
- Responsive UI with Tailwind CSS
- Pagination and search functionality
- Docker support for easy deployment

## Prerequisites

- Go 1.21 or later
- MongoDB
- Docker and Docker Compose (optional)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/news_service.git
cd news_service
```

2. Install dependencies:
```bash
make deps
```

3. Set up environment variables (optional):
```bash
export MONGODB_URI=mongodb://localhost:27017
export MONGODB_DATABASE=news_service
export PORT=8080
```

## Running the Application

### Using Make

```bash
# Build the application
make build

# Run the application
make run

# Run tests
make test
```

### Using Docker

```bash
# Build and start containers
make docker-build
make docker-up

# Stop containers
make docker-down
```

## API Endpoints

- `GET /` - List all news articles
- `GET /news/create` - Show create form
- `POST /news` - Create new article
- `GET /news/:id` - View article
- `GET /news/:id/edit` - Show edit form
- `PUT /news/:id` - Update article
- `DELETE /news/:id` - Delete article
- `GET /news/search` - Search articles

## Project Structure

```
news_service/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── domain/
│   │   └── news.go
│   ├── repository/
│   │   └── mongodb/
│   │       └── news.go
│   ├── service/
│   │   └── news.go
│   └── handler/
│       └── news.go
├── web/
│   ├── templates/
│   │   ├── layout.html
│   │   └── news/
│   │       ├── list.html
│   │       ├── create.html
│   │       └── edit.html
│   └── static/
│       └── css/
│           └── main.css
├── docker/
│   ├── Dockerfile
│   └── docker-compose.yml
├── Makefile
├── go.mod
└── README.md
```

## Testing

The project includes both unit tests and integration tests. To run the tests:

```bash
make test
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 