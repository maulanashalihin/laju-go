# Performance Optimization

This guide covers SQLite optimization, connection pooling, and performance tuning for Laju Go in production.

## SQLite Optimization

### Applied Optimizations

Laju Go applies these optimizations automatically on startup:

```go
// cmd/laju-go/main.go
func initDatabase(dbPath string) (*sql.DB, error) {
    db, err := sql.Open("sqlite", dbPath)
    if err != nil {
        return nil, err
    }

    db.SetMaxOpenConns(15)
    db.SetMaxIdleConns(10)
    db.SetConnMaxLifetime(5 * time.Minute)
    db.SetConnMaxIdleTime(30 * time.Second)

    // Enable foreign keys
    db.Exec("PRAGMA foreign_keys = ON")

    // WAL mode for better concurrency
    db.Exec("PRAGMA journal_mode = WAL")
    db.Exec("PRAGMA synchronous = NORMAL")

    // 16MB page cache (negative value = KB)
    db.Exec("PRAGMA cache_size = -16000")

    // 256MB memory-mapped I/O for NVMe performance
    db.Exec("PRAGMA mmap_size = 268435456")

    // Store temp tables in memory
    db.Exec("PRAGMA temp_store = MEMORY")

    // 5 second busy timeout
    db.Exec("PRAGMA busy_timeout = 5000")

    return db, nil
}
```

### Optimization Details

| Setting | Value | Benefit |
|---------|-------|---------|
| `journal_mode=WAL` | WAL | Better write concurrency, readers don't block writers |
| `synchronous=NORMAL` | NORMAL | Safe for WAL mode, faster than FULL |
| `cache_size=-16000` | **16MB** | Reduces disk I/O for frequent queries |
| `mmap_size=268435456` | **256MB** | NVMe memory-mapped I/O for faster reads |
| `temp_store=MEMORY` | MEMORY | Faster temporary table operations |
| `busy_timeout=5000` | 5000ms | Automatic retry on database locks |
| `foreign_keys=ON` | ON | Enforce referential integrity |

### Verify Settings

```bash
sqlite3 data/app.db "PRAGMA journal_mode;"
# Output: wal

sqlite3 data/app.db "PRAGMA synchronous;"
# Output: 1 (NORMAL)

sqlite3 data/app.db "PRAGMA cache_size;"
# Output: -16000

sqlite3 data/app.db "PRAGMA mmap_size;"
# Output: 268435456
```

### Manual Optimization

If settings aren't applied, run manually:

```bash
sqlite3 data/app.db <<EOF
PRAGMA journal_mode=WAL;
PRAGMA synchronous=NORMAL;
PRAGMA cache_size=-16000;
PRAGMA mmap_size=268435456;
PRAGMA temp_store=MEMORY;
PRAGMA busy_timeout=5000;
PRAGMA foreign_keys=ON;
EOF
```

## Connection Pooling

### Configuration

```go
// cmd/laju-go/main.go
db.SetMaxOpenConns(15)                  // Maximum open connections
db.SetMaxIdleConns(10)                  // Keep idle connections ready
db.SetConnMaxLifetime(5 * time.Minute)  // Recycle connections
db.SetConnMaxIdleTime(30 * time.Second) // Free stale idle connections
```

### Tuning Connection Pool

| Setting | Default | Description |
|---------|---------|-------------|
| `MaxOpenConns` | 15 | Maximum concurrent connections |
| `MaxIdleConns` | 10 | Idle connections to keep warm |
| `ConnMaxLifetime` | 5m | Max connection reuse duration |
| `ConnMaxIdleTime` | 30s | Recycle stale idle connections |

### When to Adjust

**Increase `MaxOpenConns`** if:
- High concurrent traffic (>100 requests/sec)
- Database queries are slow
- CPU usage is low

**Decrease `MaxOpenConns`** if:
- Memory usage is high
- Database locks are frequent
- Server has limited resources (512MB RAM)

**Tuning by Server RAM**:

| RAM | MaxOpenConns | cache_size | mmap_size |
|-----|-------------|------------|-----------|
| 512MB ⚠️ | 10 | 8MB | 128MB |
| **1-2GB ✅** | **15** | **16MB** | **256MB** |
| 4GB | 25 | 32MB | 512MB |
| 8GB | 50 | 256MB | 1GB |
| 16GB+ | 100 | 500MB+ | 2GB |

## Index Optimization

### Existing Indexes

```sql
-- Users table
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_google_id ON users(google_id);

-- Sessions table
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
```

### Add Custom Indexes

```sql
-- For frequent queries
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_role ON users(role);

-- Composite indexes for multi-column queries
CREATE INDEX idx_users_email_role ON users(email, role);
```

### Analyze Query Performance

```bash
# Enable query analysis
sqlite3 data/app.db "EXPLAIN QUERY PLAN SELECT * FROM users WHERE email = 'test@example.com';"

# Output example:
# 0|0|0|SEARCH TABLE users USING INDEX idx_users_email (email=?)
```

### Index Maintenance

```bash
# Analyze database statistics
sqlite3 data/app.db "ANALYZE;"

# Check index usage
sqlite3 data/app.db "SELECT * FROM sqlite_stat1;"
```

## Query Optimization with sqlc

Laju Go uses [sqlc](https://sqlc.dev/) for type-safe query generation. Optimizations happen at the SQL level — write efficient queries and sqlc generates efficient Go code.

### Use Parameterized Queries

sqlc queries are always parameterized — SQL injection is impossible by design:

```sql
-- queries/user.sql
-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ? LIMIT 1;
```

The `?` placeholder is type-checked at compile time by sqlc.

### Select Only Needed Columns

```sql
-- ✅ Good: Select specific columns
-- name: GetUserName :one
SELECT id, name FROM users WHERE id = ?;

-- ❌ Avoid SELECT * when you only need a few columns
-- name: GetUserByID :one  -- pulls all columns including JSON session data
SELECT * FROM users WHERE id = ?;
```

### Use LIMIT for Large Tables

```sql
-- name: ListRecentUsers :many
SELECT id, email, name FROM users ORDER BY created_at DESC LIMIT ?;
```

### Batch Operations

For batch operations, write the SQL directly:

```sql
-- name: CreateUsers :exec
INSERT INTO users (email, name, password) VALUES (?, ?, ?);
```

Insert multiple rows in a transaction:

```go
tx, _ := db.Begin()
for _, user := range users {
    tx.Exec("INSERT INTO users (email, name, password) VALUES (?, ?, ?)",
        user.Email, user.Name, user.Password)
}
tx.Commit()
```

## WAL Mode Management

### Check WAL Files

```bash
# List WAL files
ls -lh data/app.db*

# Output:
# app.db      - Main database
# app.db-shm  - Shared memory file
# app.db-wal  - Write-ahead log file
```

### Checkpoint WAL

```bash
# Manual checkpoint (copy WAL to main database)
sqlite3 data/app.db "PRAGMA wal_checkpoint(PASSIVE);"

# Full checkpoint (wait for all readers to finish)
sqlite3 data/app.db "PRAGMA wal_checkpoint(FULL);"

# Truncate checkpoint (checkpoint then truncate WAL)
sqlite3 data/app.db "PRAGMA wal_checkpoint(TRUNCATE);"
```

### Automate Checkpoint

```go
// Periodic checkpoint (run in goroutine)
func autoCheckpoint(db *sql.DB) {
    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()
    
    for range ticker.C {
        _, err := db.Exec("PRAGMA wal_checkpoint(PASSIVE)")
        if err != nil {
            log.Printf("Checkpoint error: %v", err)
        }
    }
}
```

## Database Maintenance

### Vacuum

Reclaim unused space (requires downtime):

```bash
# Stop application
sudo systemctl stop laju-go

# Vacuum database
sqlite3 data/app.db "VACUUM;"

# Start application
sudo systemctl start laju-go
```

### Integrity Check

```bash
# Check database integrity
sqlite3 data/app.db "PRAGMA integrity_check;"

# Quick check
sqlite3 data/app.db "PRAGMA quick_check;"
```

### Backup and Restore

```bash
# Online backup (no downtime)
sqlite3 data/app.db ".backup 'data/app-backup.db'"

# Restore from backup
cp data/app-backup.db data/app.db
```

## Application-Level Optimization

### User Profile Cache

Laju Go includes a built-in **in-memory TTL cache** for user profiles (`app/cache/user_cache.go`):

```go
// UserService uses the cache automatically
func (s *UserService) GetProfile(userID int64) (*models.UserResponse, error) {
    // Check cache first
    if user := s.cache.Get(userID); user != nil {
        response := user.ToResponse()
        return &response, nil
    }

    // Cache miss: query DB via sqlc
    user, err := s.querier.GetUserByID(context.Background(), userID)
    if err != nil {
        return nil, err
    }

    // Store in cache
    s.cache.Set(user)
    response := user.ToResponse()
    return &response, nil
}
```

**Cache features**:
- **TTL-based expiry** — configurable via `USER_CACHE_TTL` env var (default: 15m)
- **Auto-invalidation** — cache cleared automatically on profile updates, password changes, avatar uploads
- **Thread-safe** — uses `sync.RWMutex`
- **Configurable** — set to `0` to disable caching

### Avoid N+1 Queries

With sqlc, write a single query instead of looping:

```sql
-- ❌ Bad: N queries (one per user)
for _, id := range userIDs {
    user, _ := querier.GetUserByID(ctx, id)
}

-- ✅ Good: 1 query
-- name: GetUsersByIDs :many
SELECT * FROM users WHERE id IN (sqlc.slice('ids'));
```

### Pagination

```sql
-- queries/user.sql
-- name: ListUsersPaginated :many
SELECT id, email, name, created_at FROM users
ORDER BY created_at DESC LIMIT ? OFFSET ?;
```

```go
offset := (page - 1) * limit
users, err := s.querier.ListUsersPaginated(ctx, limit, offset)
```

## Frontend Optimization

### Build Optimization

Vite handles code splitting automatically. For custom chunks:

```javascript
// vite.config.js
export default defineConfig({
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          'vendor': ['@inertiajs/svelte'],  // or @inertiajs/react
          'utils': ['dayjs', 'axios'],
        },
      },
    },
  },
})
```

### Lazy Loading Components

```javascript
// Works with any framework via dynamic import + Inertia
const HeavyComponent = () => import('./HeavyComponent');
```

### Asset Optimization

```bash
# Compress images
npm install -g imagemin-cli
imagemin public/images/* --out-dir=public/images

# Use WebP format for smaller file sizes
```

## Monitoring Performance

### Query Logging

```go
// Enable query logging (development only)
func logQueries(db *sql.DB) {
    db.SetConnMaxLifetime(0)
    
    // Wrap with logging driver
    // Use: github.com/xo/dburl or similar
}
```

### Response Time Monitoring

```go
// Middleware to track response times
app.Use(func(c *fiber.Ctx) error {
    start := time.Now()
    
    err := c.Next()
    
    duration := time.Since(start)
    log.Printf("%s %s - %d - %v", c.Method(), c.Path(), c.Response().StatusCode(), duration)
    
    return err
})
```

### Resource Monitoring

```bash
# Memory usage
ps aux | grep laju-go

# CPU usage
top -p $(pgrep laju-go)

# Disk I/O
iotop -o

# Network connections
netstat -anp | grep laju-go
```

## Benchmarking

### Load Testing

```bash
# Install hey (HTTP load generator)
go install github.com/rakyll/hey@latest

# Benchmark homepage
hey -n 1000 -c 10 http://localhost:8080/

# Benchmark login
hey -n 100 -c 5 -m POST -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password"}' \
  http://localhost:8080/login/login
```

### Database Benchmarking

```bash
# SQLite benchmark
sqlite3 data/app.db <<EOF
.timer on
SELECT COUNT(*) FROM users;
SELECT * FROM users WHERE email = 'test@example.com';
EOF
```

## Production Checklist

- [ ] WAL mode enabled
- [ ] Connection pool configured
- [ ] Indexes created for frequent queries
- [ ] Query performance analyzed
- [ ] User profile cache configured (`USER_CACHE_TTL` env var)
- [ ] sqlc queries use LIMIT/OFFSET for pagination
- [ ] Frontend assets optimized (code splitting, compression)
- [ ] Monitoring in place
- [ ] Regular backup schedule
- [ ] Database maintenance scheduled

## Troubleshooting

### Slow Queries

**Solution**: Use EXPLAIN QUERY PLAN

```bash
sqlite3 data/app.db "EXPLAIN QUERY PLAN SELECT * FROM users WHERE email = 'test@example.com';"
```

Add index if full table scan:

```sql
CREATE INDEX idx_users_email ON users(email);
```

### Database Locked

**Solution**: 

1. Check WAL mode is enabled
2. Increase busy_timeout
3. Reduce concurrent writes
4. Check for long-running transactions

```sql
PRAGMA journal_mode=WAL;
PRAGMA busy_timeout=10000;
```

### High Memory Usage

**Solution**:

1. Reduce connection pool size
2. Lower cache_size pragma
3. Check for memory leaks

```go
db.SetMaxOpenConns(10)  // Reduce from 15
db.Exec("PRAGMA cache_size=-8000")  // Reduce to 8MB
```

## Next Steps

- [Production Deployment](production.md) - Complete deployment guide
- [SQLite Configuration](sqlite-configuration.md) - Tuning SQLite for your server
