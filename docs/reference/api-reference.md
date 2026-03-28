# API Reference

Complete reference for all HTTP endpoints in Laju Go.

## Base URL

```
Development: http://localhost:8080
Production: https://yourdomain.com
```

## Authentication

Most endpoints require authentication via session cookies. Sessions are automatically managed by the server.

### Session Cookie

```
Name: laju_session
HttpOnly: true
SameSite: Lax
Secure: false (development), true (production)
```

## Public Endpoints

### Home Page

```
GET /
```

**Response**: `200 OK`

Renders the landing page template.

---

### About Page

```
GET /about
```

**Response**: `200 OK`

Renders the about page template.

---

### Health Check

```
GET /health
```

**Response**: `200 OK`

```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

---

## Authentication Endpoints

### Show Login Form

```
GET /login
```

**Response**: `200 OK`

Renders the login page (guests only).

**Middleware**: `Guest` - Redirects authenticated users to `/app`

---

### Login

```
POST /login/login
```

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Responses**:

| Status | Description |
|--------|-------------|
| `302` | Successful login, redirect to `/app` |
| `400` | Invalid request body |
| `401` | Invalid credentials |

**Rate Limit**: 5 requests per 15 minutes

---

### Show Registration Form

```
GET /register
```

**Response**: `200 OK`

Renders the registration page (guests only).

**Middleware**: `Guest` - Redirects authenticated users to `/app`

---

### Register

```
POST /register/register
```

**Request Body**:
```json
{
  "name": "John Doe",
  "email": "user@example.com",
  "password": "password123"
}
```

**Responses**:

| Status | Description |
|--------|-------------|
| `302` | Successful registration, redirect to `/app` |
| `400` | Invalid request or user already exists |

**Rate Limit**: 3 requests per 15 minutes

---

### Logout

```
POST /logout
```

**Authentication**: Required

**Response**: `302 Found`

Redirects to `/login` and destroys session.

---

### Get Current User

```
GET /api/me
```

**Authentication**: Required

**Response**: `200 OK`

```json
{
  "id": 1,
  "email": "user@example.com",
  "name": "John Doe",
  "role": "user",
  "avatar": "/storage/avatars/avatar.png"
}
```

---

## Google OAuth Endpoints

### Start OAuth Flow

```
GET /auth/google
```

**Response**: `302 Found`

Redirects to Google OAuth consent screen.

---

### OAuth Callback

```
GET /auth/google/callback
```

**Query Parameters**:
- `code` - Authorization code from Google
- `state` - CSRF protection token

**Response**: `302 Found`

Redirects to `/app` on success or `/login` on failure.

---

## Password Reset Endpoints

### Show Forgot Password Form

```
GET /forgot-password
```

**Response**: `200 OK`

Renders the password reset request form.

---

### Send Reset Link

```
POST /forgot-password
```

**Request Body**:
```json
{
  "email": "user@example.com"
}
```

**Response**: `200 OK`

```json
{
  "message": "If the email exists, a reset link has been sent"
}
```

**Rate Limit**: 3 requests per hour

**Note**: Always returns success message to prevent email enumeration.

---

### Show Reset Password Form

```
GET /reset-password/:token
```

**URL Parameters**:
- `token` - Password reset token

**Response**: `200 OK`

Renders the password reset form.

---

### Reset Password

```
POST /reset-password/:token
```

**URL Parameters**:
- `token` - Password reset token

**Request Body**:
```json
{
  "password": "newpassword123"
}
```

**Responses**:

| Status | Description |
|--------|-------------|
| `302` | Successful reset, redirect to `/login` |
| `400` | Invalid or expired token |

---

## Protected App Endpoints

### Dashboard

```
GET /app
```

**Authentication**: Required

**Response**: `200 OK`

Renders the dashboard page via Inertia.js.

**Middleware**: `AuthRequired`, `CSRF`

---

### Profile Page

```
GET /app/profile
```

**Authentication**: Required

**Response**: `200 OK`

Renders the profile page via Inertia.js.

**Middleware**: `AuthRequired`, `CSRF`

---

### Update Profile

```
PUT /app/profile
```

**Authentication**: Required

**Request Body**:
```json
{
  "name": "John Doe",
  "email": "newemail@example.com"
}
```

**Responses**:

| Status | Description |
|--------|-------------|
| `302` | Successful update, redirect to profile |
| `400` | Invalid request or email already exists |

**Middleware**: `AuthRequired`, `CSRF`

---

### Update Password

```
PUT /app/profile/password
```

**Authentication**: Required

**Request Body**:
```json
{
  "current_password": "oldpassword123",
  "password": "newpassword123"
}
```

**Responses**:

| Status | Description |
|--------|-------------|
| `302` | Successful update |
| `400` | Invalid current password or weak new password |

**Middleware**: `AuthRequired`, `CSRF`

---

### File Upload

```
POST /upload
```

**Authentication**: Required

**Request**: `multipart/form-data`

**Form Data**:
- `avatar` - Image file (JPEG, PNG, GIF)

**Responses**:

| Status | Description |
|--------|-------------|
| `200` | Successful upload |
| `400` | Invalid file type or size |
| `500` | Failed to save file |

**File Constraints**:
- Max size: 5MB
- Allowed types: `image/jpeg`, `image/png`, `image/gif`

**Middleware**: `AuthRequired`, `CSRF`

---

## Admin Endpoints

### Admin Dashboard

```
GET /admin
```

**Authentication**: Required

**Authorization**: Admin role only

**Response**: `200 OK`

Renders the admin dashboard page.

**Middleware**: `AuthRequired`, `AdminRequired`

---

## Error Responses

### Common Error Format

```json
{
  "error": "Error message here"
}
```

### Validation Errors

```json
{
  "error": "Validation failed",
  "details": {
    "email": "Email is required",
    "password": "Password must be at least 8 characters"
  }
}
```

### Status Codes

| Code | Description |
|------|-------------|
| `200` | Success |
| `302` | Redirect |
| `400` | Bad Request - Invalid input |
| `401` | Unauthorized - Not logged in |
| `403` | Forbidden - Insufficient permissions |
| `404` | Not Found |
| `422` | Unprocessable Entity - Validation error |
| `429` | Too Many Requests - Rate limit exceeded |
| `500` | Internal Server Error |

---

## Rate Limiting

### Rate Limit Headers

```
X-RateLimit-Limit: 5
X-RateLimit-Remaining: 4
X-RateLimit-Reset: 1704110400
```

### Rate Limits by Endpoint

| Endpoint | Limit | Window |
|----------|-------|--------|
| `POST /login/login` | 5 | 15 minutes |
| `POST /register/register` | 3 | 15 minutes |
| `POST /forgot-password` | 3 | 1 hour |
| `POST /upload` | 50 | 1 hour |
| API endpoints | 100 | 15 minutes |

---

## CSRF Protection

### Getting CSRF Token

CSRF token is available in:

1. **Cookie**: `csrf_token` (for JavaScript access)
2. **Inertia Props**: `csrf` (for Svelte components)

### Sending CSRF Token

Include CSRF token in state-changing requests:

**Header**:
```
X-CSRF-Token: <token>
```

**Form Field**:
```html
<input type="hidden" name="csrf_token" value="<token>">
```

### CSRF Exempt Endpoints

These endpoints don't require CSRF validation:
- `GET` requests
- `HEAD` requests
- `OPTIONS` requests
- `POST /login/login`
- `POST /register/register`
- `POST /forgot-password`

---

## Inertia.js Integration

### Inertia Requests

For SPA navigation, include the `X-Inertia` header:

```
GET /app
X-Inertia: true
```

**Response**: `200 OK`

```json
{
  "component": "Dashboard",
  "props": {
    "user": { "id": 1, "name": "John" },
    "stats": { "total": 100 }
  },
  "url": "/app",
  "version": "1.0.0"
}
```

### Inertia Errors

```json
{
  "component": "Profile",
  "props": {
    "errors": {
      "email": "Email already taken"
    }
  }
}
```

---

## Examples

### Login Flow

```bash
# 1. Get login page
curl -i http://localhost:8080/login

# 2. Submit login
curl -i -X POST http://localhost:8080/login/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}' \
  -c cookies.txt

# 3. Access protected route
curl -i http://localhost:8080/app \
  -b cookies.txt
```

### Registration Flow

```bash
# 1. Register
curl -i -X POST http://localhost:8080/register/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"user@example.com","password":"password123"}' \
  -c cookies.txt

# 2. Access dashboard
curl -i http://localhost:8080/app -b cookies.txt
```

### Password Reset Flow

```bash
# 1. Request reset
curl -i -X POST http://localhost:8080/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com"}'

# 2. Reset password (with token from email)
curl -i -X POST http://localhost:8080/reset-password/abc123 \
  -H "Content-Type: application/json" \
  -d '{"password":"newpassword123"}'
```

### File Upload

```bash
curl -i -X POST http://localhost:8080/upload \
  -b cookies.txt \
  -F "avatar=@/path/to/avatar.jpg"
```

---

## Best Practices

### 1. Handle Errors Gracefully

```javascript
// Svelte example
import { router } from '@inertiajs/svelte';

router.post('/login', formData, {
  onError: (errors) => {
    console.log('Login failed:', errors);
  },
  onSuccess: () => {
    console.log('Login successful!');
  },
});
```

### 2. Use CSRF Tokens

```javascript
// Get CSRF token from cookie
function getCsrfToken() {
  const match = document.cookie.match(/csrf_token=([^;]+)/);
  return match ? match[1] : null;
}

// Include in fetch
fetch('/app/profile', {
  method: 'PUT',
  headers: {
    'X-CSRF-Token': getCsrfToken(),
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({ name: 'John' }),
});
```

### 3. Handle Rate Limits

```javascript
// Check for rate limit headers
if (response.status === 429) {
  const resetTime = response.headers.get('X-RateLimit-Reset');
  const waitTime = resetTime - Date.now() / 1000;
  console.log(`Rate limited. Try again in ${waitTime} seconds`);
}
```

---

## Next Steps

- [Routing Guide](../guide/routing.md) - Route definitions
- [Handlers Guide](../guide/handlers.md) - Building handlers
- [Authentication Guide](../guide/authentication.md) - Auth implementation
