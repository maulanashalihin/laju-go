# Database

This guide covers database setup, migrations, and query building in Laju Go.

## Overview

Laju Go uses **SQLite** as the database with **Squirrel** as a query builder and **Goose** for migrations. This combination provides:

- **Zero configuration** - No database server to manage
- **Type-safe queries** - Squirrel builds parameterized SQL
- **Version control** - Goose manages schema migrations
- **Production-ready** - SQLite with WAL mode and optimizations

## Database Setup

### Connection Initialization

```go
// main.go
import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

func initDatabase() *sql.DB {
    dbPath := os.Getenv("DB_PATH")
    if dbPath == "" {
        dbPath = "data/app.db"
    }
    
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        log.Fatal("Failed to open database:", err)
    }
    
    // Apply production optimizations
    applySQLiteOptimizations(db)
    
    // Configure connection pool
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)
    
    return db
}
```

### SQLite Optimizations

```go
func applySQLiteOptimizations(db *sql.DB) {
    optimizations := []string{
        "PRAGMA journal_mode=WAL",           // Write-Ahead Logging
        "PRAGMA synchronous=NORMAL",         // Balance speed/durability
        "PRAGMA cache_size=-64000",          // 64MB cache
        "PRAGMA temp_store=MEMORY",          // Memory temp tables
        "PRAGMA busy_timeout=5000",          // 5 second lock wait
        "PRAGMA foreign_keys=ON",            // Enable foreign keys
    }
    
    for _, pragma := range optimizations {
        _, err := db.Exec(pragma)
        if err != nil {
            log.Printf("Warning: Failed to set %s: %v", pragma, err)
        }
    }
}
```

### Why These Settings?

| Setting | Value | Benefit |
|---------|-------|---------|
| `journal_mode=WAL` | WAL | Better write concurrency, readers don't block writers |
| `synchronous=NORMAL` | NORMAL | Safe for WAL mode, faster than FULL |
| `cache_size=-64000` | 64MB | Reduces disk I/O for frequent queries |
| `temp_store=MEMORY` | MEMORY | Faster temporary table operations |
| `busy_timeout=5000` | 5000ms | Automatic retry on database locks |
| `foreign_keys=ON` | ON | Enforce referential integrity |

## Database Migrations

### Installing Goose

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### Creating Migrations

```bash
# Create a new migration
goose -dir migrations create add_users_table

# Output: migrations/20240101120000_add_users_table.sql
```

### Migration File Structure

```sql
-- migrations/0001_create_users_table.sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    password TEXT,
    avatar TEXT DEFAULT '',
    role TEXT NOT NULL DEFAULT 'user',
    google_id TEXT UNIQUE,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_google_id ON users(google_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
```

### Running Migrations

```bash
# Run all pending migrations
goose -dir migrations sqlite3 data/app.db up

# Check migration status
goose -dir migrations sqlite3 data/app.db status

# Rollback last migration
goose -dir migrations sqlite3 data/app.db down

# Reset all migrations
goose -dir migrations sqlite3 data/app.db reset

# Run specific migration
goose -dir migrations sqlite3 data/app.db up-by-one
```

### Auto-Run Migrations on Startup

```go
// main.go
func runMigrations(db *sql.DB) {
    dbPath := os.Getenv("DB_PATH")
    if dbPath == "" {
        dbPath = "data/app.db"
    }
    
    err := goose.Up(db, "migrations")
    if err != nil {
        log.Fatal("Migration failed:", err)
    }
    
    log.Println("Migrations completed successfully")
}
```

## Query Building with Squirrel

### Installing Squirrel

```bash
go get github.com/Masterminds/squirrel
```

### Basic Queries

```go
import (
    "github.com/Masterminds/squirrel"
)

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)
```

#### Select

```go
// SELECT id, email, name FROM users WHERE id = ?
query := psql.
    Select("id", "email", "name").
    From("users").
    Where(squirrel.Eq{"id": userID})

sql, args, err := query.ToSql()
if err != nil {
    return nil, err
}

var user User
err = db.QueryRow(sql, args...).Scan(&user.ID, &user.Email, &user.Name)
```

#### Insert

```go
// INSERT INTO users (email, name, password) VALUES (?, ?, ?)
query := psql.
    Insert("users").
    Columns("email", "name", "password").
    Values(email, name, password)

sql, args, err := query.ToSql()
if err != nil {
    return err
}

result, err := db.Exec(sql, args...)
if err != nil {
    return err
}

// Get inserted ID
id, _ := result.LastInsertId()
```

#### Update

```go
// UPDATE users SET name = ?, email = ? WHERE id = ?
query := psql.
    Update("users").
    Set("name", name).
    Set("email", email).
    Where(squirrel.Eq{"id": userID})

sql, args, err := query.ToSql()
if err != nil {
    return err
}

_, err = db.Exec(sql, args...)
```

#### Delete

```go
// DELETE FROM users WHERE id = ?
query := psql.
    Delete("users").
    Where(squirrel.Eq{"id": userID})

sql, args, err := query.ToSql()
if err != nil {
    return err
}

_, err = db.Exec(sql, args...)
```

### Advanced Queries

#### Joins

```go
// SELECT u.*, s.data as session_data
// FROM users u
// LEFT JOIN sessions s ON u.id = s.user_id
// WHERE u.id = ?
query := psql.
    Select("u.*", "s.data as session_data").
    From("users u").
    LeftJoin("sessions s ON u.id = s.user_id").
    Where(squirrel.Eq{"u.id": userID})

sql, args, err := query.ToSql()
```

#### Multiple Conditions

```go
// SELECT * FROM users
// WHERE (email = ? OR google_id = ?) AND role = ?
query := psql.
    Select("*").
    From("users").
    Where(squirrel.Or{
        squirrel.Eq{"email": email},
        squirrel.Eq{"google_id": googleID},
    }).
    Where(squirrel.Eq{"role": "user"})

sql, args, err := query.ToSql()
```

#### Like Operator

```go
// SELECT * FROM users WHERE name LIKE ?
query := psql.
    Select("*").
    From("users").
    Where(squirrel.Like{"name": "%john%"})

sql, args, err := query.ToSql()
```

#### IN Clause

```go
// SELECT * FROM users WHERE id IN (?, ?, ?)
query := psql.
    Select("*").
    From("users").
    Where(squirrel.Eq{"id": []int{1, 2, 3}})

sql, args, err := query.ToSql()
// args will be [1, 2, 3]
```

#### Ordering and Limiting

```go
// SELECT * FROM users ORDER BY created_at DESC LIMIT 10 OFFSET 0
query := psql.
    Select("*").
    From("users").
    OrderBy("created_at DESC").
    Limit(10).
    Offset(0)

sql, args, err := query.ToSql()
```

#### Count

```go
// SELECT COUNT(*) FROM users WHERE role = ?
query := psql.
    Select("COUNT(*) as count").
    From("users").
    Where(squirrel.Eq{"role": "user"})

sql, args, err := query.ToSql()

var count int
err = db.QueryRow(sql, args...).Scan(&count)
```

## Repository Pattern

### Repository Structure

```go
// app/repositories/user.repository.go
package repositories

import (
    "database/sql"
    "github.com/Masterminds/squirrel"
    "laju-go/app/models"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)
```

### CRUD Operations

```go
// Create
func (r *UserRepository) Create(email, name, password string) (*models.User, error) {
    query := psql.
        Insert("users").
        Columns("email", "name", "password").
        Values(email, name, password)
    
    sql, args, err := query.ToSql()
    if err != nil {
        return nil, err
    }
    
    result, err := r.db.Exec(sql, args...)
    if err != nil {
        return nil, err
    }
    
    id, _ := result.LastInsertId()
    return r.GetByID(int(id))
}

// Read by ID
func (r *UserRepository) GetByID(id int) (*models.User, error) {
    query := psql.
        Select("id", "email", "name", "password", "avatar", "role", "google_id", "email_verified", "created_at", "updated_at").
        From("users").
        Where(squirrel.Eq{"id": id}).
        Limit(1)
    
    sql, args, err := query.ToSql()
    if err != nil {
        return nil, err
    }
    
    var user models.User
    err = r.db.QueryRow(sql, args...).Scan(
        &user.ID, &user.Email, &user.Name, &user.Password,
        &user.Avatar, &user.Role, &user.GoogleID, &user.EmailVerified,
        &user.CreatedAt, &user.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, ErrUserNotFound
        }
        return nil, err
    }
    
    return &user, nil
}

// Read by Email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
    query := psql.
        Select("id", "email", "name", "password", "avatar", "role", "google_id", "email_verified", "created_at", "updated_at").
        From("users").
        Where(squirrel.Eq{"email": email}).
        Limit(1)
    
    sql, args, err := query.ToSql()
    if err != nil {
        return nil, err
    }
    
    var user models.User
    err = r.db.QueryRow(sql, args...).Scan(
        &user.ID, &user.Email, &user.Name, &user.Password,
        &user.Avatar, &user.Role, &user.GoogleID, &user.EmailVerified,
        &user.CreatedAt, &user.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, ErrUserNotFound
        }
        return nil, err
    }
    
    return &user, nil
}

// Update
func (r *UserRepository) Update(id int, name, email string) error {
    query := psql.
        Update("users").
        Set("name", name).
        Set("email", email).
        Set("updated_at", squirrel.Expr("CURRENT_TIMESTAMP")).
        Where(squirrel.Eq{"id": id})
    
    sql, args, err := query.ToSql()
    if err != nil {
        return err
    }
    
    _, err = r.db.Exec(sql, args...)
    return err
}

// Update Password
func (r *UserRepository) UpdatePassword(id int, password string) error {
    query := psql.
        Update("users").
        Set("password", password).
        Set("updated_at", squirrel.Expr("CURRENT_TIMESTAMP")).
        Where(squirrel.Eq{"id": id})
    
    sql, args, err := query.ToSql()
    if err != nil {
        return err
    }
    
    _, err = r.db.Exec(sql, args...)
    return err
}

// Update Avatar
func (r *UserRepository) UpdateAvatar(id int, avatar string) error {
    query := psql.
        Update("users").
        Set("avatar", avatar).
        Set("updated_at", squirrel.Expr("CURRENT_TIMESTAMP")).
        Where(squirrel.Eq{"id": id})
    
    sql, args, err := query.ToSql()
    if err != nil {
        return err
    }
    
    _, err = r.db.Exec(sql, args...)
    return err
}

// Delete
func (r *UserRepository) Delete(id int) error {
    query := psql.
        Delete("users").
        Where(squirrel.Eq{"id": id})
    
    sql, args, err := query.ToSql()
    if err != nil {
        return err
    }
    
    _, err = r.db.Exec(sql, args...)
    return err
}

// Get by Google ID
func (r *UserRepository) GetByGoogleID(googleID string) (*models.User, error) {
    query := psql.
        Select("id", "email", "name", "password", "avatar", "role", "google_id", "email_verified", "created_at", "updated_at").
        From("users").
        Where(squirrel.Eq{"google_id": googleID}).
        Limit(1)
    
    sql, args, err := query.ToSql()
    if err != nil {
        return nil, err
    }
    
    var user models.User
    err = r.db.QueryRow(sql, args...).Scan(
        &user.ID, &user.Email, &user.Name, &user.Password,
        &user.Avatar, &user.Role, &user.GoogleID, &user.EmailVerified,
        &user.CreatedAt, &user.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, ErrUserNotFound
        }
        return nil, err
    }
    
    return &user, nil
}
```

## Database Schema

### Users Table

```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    password TEXT,
    avatar TEXT DEFAULT '',
    role TEXT NOT NULL DEFAULT 'user',
    google_id TEXT UNIQUE,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_google_id ON users(google_id);
```

### Sessions Table

```sql
CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    data TEXT NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
```

## Transactions

```go
func (r *UserRepository) TransferCredits(fromID, toID int, amount int64) error {
    // Begin transaction
    tx, err := r.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    // Deduct from source
    _, err = tx.Exec(
        "UPDATE users SET credits = credits - ? WHERE id = ?",
        amount, fromID,
    )
    if err != nil {
        return err
    }
    
    // Add to destination
    _, err = tx.Exec(
        "UPDATE users SET credits = credits + ? WHERE id = ?",
        amount, toID,
    )
    if err != nil {
        return err
    }
    
    // Commit transaction
    return tx.Commit()
}
```

## Database Utilities

### Health Check

```go
func CheckDatabaseHealth(db *sql.DB) error {
    return db.Ping()
}
```

### Get Database Statistics

```go
func GetDatabaseStats(db *sql.DB) (map[string]interface{}, error) {
    stats := make(map[string]interface{})
    
    // Get user count
    var userCount int
    err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
    if err != nil {
        return nil, err
    }
    stats["user_count"] = userCount
    
    // Get session count
    var sessionCount int
    err = db.QueryRow("SELECT COUNT(*) FROM sessions").Scan(&sessionCount)
    if err != nil {
        return nil, err
    }
    stats["session_count"] = sessionCount
    
    return stats, nil
}
```

## Best Practices

### 1. Use Parameterized Queries

Always use Squirrel or parameterized queries to prevent SQL injection:

```go
// ❌ Bad: SQL injection vulnerability
query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email)

// ✅ Good: Parameterized query
query := psql.
    Select("*").
    From("users").
    Where(squirrel.Eq{"email": email})
```

### 2. Handle sql.ErrNoRows

```go
// ✅ Good: Specific error handling
err = db.QueryRow(sql, args...).Scan(&user.ID)
if err != nil {
    if err == sql.ErrNoRows {
        return nil, ErrUserNotFound
    }
    return nil, err
}
```

### 3. Use Transactions for Multiple Writes

```go
// ✅ Good: Transaction for data integrity
tx, err := db.Begin()
if err != nil {
    return err
}
defer tx.Rollback()

// Multiple operations
_, err = tx.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, fromID)
_, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, toID)

return tx.Commit()
```

### 4. Index Frequently Queried Columns

```sql
-- Index for email lookups
CREATE INDEX idx_users_email ON users(email);

-- Index for foreign keys
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
```

### 5. Use Connection Pooling

```go
db.SetMaxOpenConns(25)    // Maximum open connections
db.SetMaxIdleConns(5)     // Idle connections to keep
db.SetConnMaxLifetime(5 * time.Minute)  // Max connection lifetime
```

### 6. Close Resources

```go
// ✅ Good: Close rows after use
rows, err := db.Query(sql, args...)
if err != nil {
    return err
}
defer rows.Close()

for rows.Next() {
    // Scan row
}
```

## Troubleshooting

### Database Locked

**Problem**: `database is locked`

**Solutions**:
1. Enable WAL mode: `PRAGMA journal_mode=WAL`
2. Set busy timeout: `PRAGMA busy_timeout=5000`
3. Reduce concurrent writes
4. Check for unclosed transactions

### Migration Failed

**Problem**: Migration fails on startup

**Solutions**:
1. Check migration syntax
2. Verify database path
3. Run migrations manually: `goose -dir migrations sqlite3 data/app.db up`
4. Check goose_db_version table

### Connection Issues

**Problem**: `unable to open database file`

**Solutions**:
1. Ensure directory exists: `mkdir -p data`
2. Check permissions: `chmod 755 data`
3. Verify DB_PATH in .env

## Next Steps

- [Authentication Guide](authentication.md) - User authentication and sessions
- [Architecture Guide](architecture.md) - Repository pattern in context
- [Deployment Guide](../deployment/production.md) - Production database setup
