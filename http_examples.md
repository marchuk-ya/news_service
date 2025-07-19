# HTTP Examples for News Service API

## Base URL
```
http://localhost:8080
```

## Authentication
Currently, the API does not require authentication, but it's prepared for future use.

## Endpoints

### 1. Health Checks

#### GET /health
Basic health check.

```bash
curl -X GET http://localhost:8080/health
```

**Response:**
```json
{
  "status": "healthy",
  "service": "news_service"
}
```

#### GET /health/readiness
Readiness check with database connectivity verification.

```bash
curl -X GET http://localhost:8080/health/readiness
```

**Response:**
```json
{
  "status": "ready",
  "service": "news_service",
  "database": "connected"
}
```

#### GET /api/v1/health
API health check.

```bash
curl -X GET http://localhost:8080/api/v1/health
```

**Response:**
```json
{
  "status": "healthy",
  "service": "news_service",
  "version": "1.0.0"
}
```

#### GET /api/v1/health/readiness
API readiness check.

```bash
curl -X GET http://localhost:8080/api/v1/health/readiness
```

**Response:**
```json
{
  "status": "ready",
  "service": "news_service",
  "database": "connected",
  "version": "1.0.0"
}
```

### 2. News Management

#### GET /api/v1/news
Get list of all news articles with pagination.

```bash
# Basic request
curl -X GET http://localhost:8080/api/v1/news

# With pagination
curl -X GET "http://localhost:8080/api/v1/news?page=1&limit=10"

# With custom parameters
curl -X GET "http://localhost:8080/api/v1/news?page=2&limit=5"
```

**Response:**
```json
{
  "news": [
    {
      "id": "507f1f77bcf86cd799439011",
      "title": "Breaking News",
      "content": "This is a breaking news article content.",
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z"
    }
  ],
  "total": 100,
  "page": 1,
  "limit": 10,
  "total_pages": 10
}
```

#### POST /api/v1/news
Create a new news article.

```bash
# Basic request
curl -X POST http://localhost:8080/api/v1/news \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Breaking News",
    "content": "This is a breaking news article content."
  }'

# With more content
curl -X POST http://localhost:8080/api/v1/news \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Important Update",
    "content": "This is a very important update that contains detailed information about the latest developments in our system."
  }'
```

**Response:**
```json
{
  "id": "507f1f77bcf86cd799439011",
  "title": "Breaking News",
  "content": "This is a breaking news article content.",
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

#### GET /api/v1/news/{id}
Get a news article by ID.

```bash
# Replace {id} with actual news ID
curl -X GET http://localhost:8080/api/v1/news/507f1f77bcf86cd799439011

# Example with real ID
curl -X GET http://localhost:8080/api/v1/news/507f1f77bcf86cd799439011
```

**Response:**
```json
{
  "id": "507f1f77bcf86cd799439011",
  "title": "Breaking News",
  "content": "This is a breaking news article content.",
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

#### PUT /api/v1/news/{id}
Update an existing news article.

```bash
# Replace {id} with actual news ID
curl -X PUT http://localhost:8080/api/v1/news/507f1f77bcf86cd799439011 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Breaking News",
    "content": "This is an updated breaking news article content."
  }'

# Example with real ID
curl -X PUT http://localhost:8080/api/v1/news/507f1f77bcf86cd799439011 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Breaking News",
    "content": "This is an updated breaking news article content with more details."
  }'
```

**Response:**
```json
{
  "id": "507f1f77bcf86cd799439011",
  "title": "Updated Breaking News",
  "content": "This is an updated breaking news article content.",
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T13:00:00Z"
}
```

#### DELETE /api/v1/news/{id}
Delete a news article by ID.

```bash
# Replace {id} with actual news ID
curl -X DELETE http://localhost:8080/api/v1/news/507f1f77bcf86cd799439011

# Example with real ID
curl -X DELETE http://localhost:8080/api/v1/news/507f1f77bcf86cd799439011
```

**Response:** 204 No Content

#### GET /api/v1/news/search
Search news articles by title and content.

```bash
# Basic search
curl -X GET "http://localhost:8080/api/v1/news/search?q=breaking"

# Search with pagination
curl -X GET "http://localhost:8080/api/v1/news/search?q=news&page=1&limit=5"

# Search with more complex query
curl -X GET "http://localhost:8080/api/v1/news/search?q=important%20update&page=1&limit=10"
```

**Response:**
```json
{
  "news": [
    {
      "id": "507f1f77bcf86cd799439011",
      "title": "Breaking News",
      "content": "This is a breaking news article content.",
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z"
    }
  ],
  "total": 5,
  "query": "breaking",
  "page": 1,
  "limit": 10,
  "total_pages": 1
}
```

### 3. Web Interface (HTML)

#### GET /
List all news articles (HTML).

```bash
curl -X GET http://localhost:8080/
```

#### GET /news/create
Create news form (HTML).

```bash
curl -X GET http://localhost:8080/news/create
```

#### GET /news/{id}
View news article (HTML).

```bash
curl -X GET http://localhost:8080/news/507f1f77bcf86cd799439011
```

#### GET /news/{id}/edit
Edit news form (HTML).

```bash
curl -X GET http://localhost:8080/news/507f1f77bcf86cd799439011/edit
```

#### GET /news/search
Search news articles (HTML).

```bash
curl -X GET "http://localhost:8080/news/search?q=breaking"
```

### 4. Documentation

#### GET /swagger/index.html
Swagger UI documentation.

```bash
curl -X GET http://localhost:8080/swagger/index.html
```

#### GET /docs
Redirect to Swagger UI.

```bash
curl -X GET http://localhost:8080/docs
```

## Error Responses

### Validation Error (400)
```json
{
  "error": "Invalid request body",
  "code": "VALIDATION_ERROR",
  "message": "Key: 'CreateNewsRequest.Title' Error:Field validation for 'Title' failed on the 'required' tag"
}
```

### Not Found Error (404)
```json
{
  "error": "Failed to get news",
  "code": "NOT_FOUND",
  "message": "news not found"
}
```

### Internal Server Error (500)
```json
{
  "error": "Failed to fetch news",
  "code": "INTERNAL_ERROR",
  "message": "database connection failed"
}
```

### Rate Limit Error (429)
```json
{
  "error": "Rate limit exceeded",
  "code": "RATE_LIMIT_EXCEEDED"
}
```

## Headers

### Request Headers
```bash
# For JSON API
-H "Content-Type: application/json"

# For authorization (future)
-H "Authorization: Bearer <token>"
```

### Response Headers
```
Content-Type: application/json
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Referrer-Policy: strict-origin-when-cross-origin
Access-Control-Allow-Origin: *
```

## Rate Limiting

API has a limit of **100 requests per minute** per IP address.

## Pagination

All endpoints that return lists support pagination:

- `page` (optional): Page number (default: 1)
- `limit` (optional): Number of items per page (default: 10, maximum: 100)

## Validation Rules

### News Validation
- **title**: required, minimum 3 characters, maximum 200 characters
- **content**: required, minimum 10 characters

### Pagination Validation
- **page**: minimum 1
- **limit**: minimum 1, maximum 100

### Search Validation
- **query**: required, cannot be empty

## Complete Workflow Example

```bash
# 1. Check health
curl -X GET http://localhost:8080/api/v1/health

# 2. Create news article
curl -X POST http://localhost:8080/api/v1/news \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test News",
    "content": "This is a test news article for demonstration."
  }'

# 3. Get list of news articles
curl -X GET http://localhost:8080/api/v1/news

# 4. Search news articles
curl -X GET "http://localhost:8080/api/v1/news/search?q=test"

# 5. Get specific news article (use ID from step 2)
curl -X GET http://localhost:8080/api/v1/news/507f1f77bcf86cd799439011

# 6. Update news article
curl -X PUT http://localhost:8080/api/v1/news/507f1f77bcf86cd799439011 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Test News",
    "content": "This is an updated test news article."
  }'

# 7. Delete news article
curl -X DELETE http://localhost:8080/api/v1/news/507f1f77bcf86cd799439011
```

## Testing with Different Tools

### Using wget
```bash
wget -qO- http://localhost:8080/api/v1/health
```

### Using httpie
```bash
# GET request
http GET localhost:8080/api/v1/news

# POST request
http POST localhost:8080/api/v1/news title="Test News" content="Test content"
```

### Using Postman
1. Import the collection
2. Set base URL: `http://localhost:8080`
3. Use the provided examples for each endpoint

### Using JavaScript (fetch)
```javascript
// GET request
fetch('http://localhost:8080/api/v1/news')
  .then(response => response.json())
  .then(data => console.log(data));

// POST request
fetch('http://localhost:8080/api/v1/news', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    title: 'Test News',
    content: 'Test content'
  })
})
.then(response => response.json())
.then(data => console.log(data));
``` 