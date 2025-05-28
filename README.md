# News Service

### Features

- Create, read, search, update, and delete news articles

### Run

#### Using Make

```bash
# Build the application
make build

# Run the application
make run
```

#### Using Docker

```bash
# Build and start containers
make up

# Stop containers
make down
```

### API Endpoints

- `GET /` - List all news articles
- `GET /news/create` - Show create form
- `POST /news` - Create new article
- `GET /news/:id` - View article
- `GET /news/:id/edit` - Show edit form
- `PUT /news/:id` - Update article
- `DELETE /news/:id` - Delete article
- `GET /news/search` - Search articles


## Testing

```bash
# run unit-tests
make test-unit

# run integration tests
make test-integration
```
