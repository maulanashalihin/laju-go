# Performance Optimization

This guide covers SQLite optimization, connection pooling, and performance tuning for Laju Go in production.

## SQLite Optimization

### Applied Optimizations

Laju Go applies these optimizations automatically on startup:

```go
// main.go
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

### Optimization Details

| Setting | Value | Benefit |
|---------|-------|---------|
| `journal_mode=WAL` | WAL | Better write concurrency, readers don't block writers |
| `synchronous=NORMAL` | NORMAL | Safe for WAL mode, faster than FULL |
| `cache_size=-64000` | 64MB | Reduces disk I/O for frequent queries |
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
# Output: -64000
```

### Manual Optimization

If settings aren't applied, run manually:

```bash
sqlite3 data/app.db <<EOF
PRAGMA journal_mode=WAL;
PRAGMA synchronous=NORMAL;
PRAGMA cache_size=-64000;
PRAGMA temp_store=MEMORY;
PRAGMA busy_timeout=5000;
EOF
```

## Connection Pooling

### Configuration

```go
// main.go
db.SetMaxOpenConns(25)           // Maximum open connections
db.SetMaxIdleConns(5)            // Idle connections to keep
db.SetConnMaxLifetime(5 * time.Minute)  // Max connection lifetime
```

### Tuning Connection Pool

| Setting | Recommended | Description |
|---------|-------------|-------------|
| `MaxOpenConns` | 25 | Maximum concurrent connections |
| `MaxIdleConns` | 5 | Minimum idle connections |
| `ConnMaxLifetime` | 5m | Connection reuse duration |

### When to Adjust

**Increase `MaxOpenConns`** if:
- High concurrent traffic (>100 requests/sec)
- Database queries are slow
- CPU usage is low

**Decrease `MaxOpenConns`** if:
- Memory usage is high
- Database locks are frequent
- Server has limited resources

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

## Query Optimization

### Use Parameterized Queries

```go
// ✅ Good: Parameterized query (fast, safe)
query := psql.
    Select("*").
    From("users").
    Where(squirrel.Eq{"email": email})

// ❌ Bad: String concatenation (slow, SQL injection risk)
query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email)
```

### Select Only Needed Columns

```go
// ✅ Good: Select specific columns
query := psql.
    Select("id", "email", "name").
    From("users")

// ❌ Bad: Select all columns
query := psql.
    Select("*").
    From("users")
```

### Use LIMIT for Large Tables

```go
// ✅ Good: Limit results
query := psql.
    Select("*").
    From("users").
    Limit(100)

// ❌ Bad: No limit on large table
query := psql.
    Select("*").
    From("users")
```

### Batch Operations

```go
// ✅ Good: Batch insert
func (r *UserRepository) BatchInsert(users []User) error {
    query := psql.Insert("users").Columns("email", "name", "password")
    
    for _, user := range users {
        query = query.Values(user.Email, user.Name, user.Password)
    }
    
    sql, args, err := query.ToSql()
    if err != nil {
        return err
    }
    
    _, err = r.db.Exec(sql, args...)
    return err
}

// ❌ Bad: Individual inserts
for _, user := range users {
    r.db.Exec("INSERT INTO users ...")  // Slow!
}
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

### Caching

Implement caching for frequently accessed data:

```go
// Simple in-memory cache
type Cache struct {
    data sync.Map
}

func (c *Cache) Get(key string) (interface{}, bool) {
    return c.data.Load(key)
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
    c.data.Store(key, value)
    
    // Auto-expire
    go func() {
        time.Sleep(ttl)
        c.data.Delete(key)
    }()
}

// Usage
func (s *UserService) GetByID(id int) (*User, error) {
    // Check cache first
    if cached, ok := s.cache.Get(fmt.Sprintf("user:%d", id)); ok {
        return cached.(*User), nil
    }
    
    // Fetch from database
    user, err := s.userRepo.GetByID(id)
    if err != nil {
        return nil, err
    }
    
    // Cache for 5 minutes
    s.cache.Set(fmt.Sprintf("user:%d", id), user, 5*time.Minute)
    
    return user, nil
}
```

### Lazy Loading

```go
// ✅ Good: Load related data only when needed
type User struct {
    ID       int
    Email    string
    Name     string
    Sessions []Session `json:"-"`  // Don't serialize by default
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
    user, _ := h.userRepo.GetByID(id)
    
    // Only load sessions if requested
    if c.Query("include") == "sessions" {
        user.Sessions, _ = h.sessionRepo.GetByUserID(user.ID)
    }
    
    return c.JSON(user)
}
```

### Pagination

```go
// ✅ Good: Paginate large result sets
func (r *UserRepository) GetUsers(page, limit int) ([]User, error) {
    offset := (page - 1) * limit
    
    query := psql.
        Select("id", "email", "name", "created_at").
        From("users").
        OrderBy("created_at DESC").
        Limit(limit).
        Offset(offset)
    
    // ...
}
```

## Frontend Optimization

### Build Optimization

```javascript
// vite.config.js
export default defineConfig({
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          // Split vendor chunks
          'vendor': ['@inertiajs/svelte'],
          'utils': ['dayjs', 'axios'],
        },
      },
    },
  },
})
```

### Lazy Loading Components

```svelte
<!-- ✅ Good: Lazy load heavy components -->
<script>
  import { lazy } from 'svelte';
  const HeavyComponent = lazy(() => import('./HeavyComponent.svelte'));
</script>

{#if showHeavyComponent}
  <HeavyComponent />
{/if}
```

### Asset Optimization

```bash
# Compress images
npm install -g imagemin-cli
imagemin public/images/* --out-dir=public/images

# Use WebP format
# Convert PNG/JPG to WebP for smaller file sizes
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
- [ ] Caching implemented for hot data
- [ ] Pagination for large result sets
- [ ] Frontend assets optimized
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
db.SetMaxOpenConns(10)  // Reduce from 25
db.Exec("PRAGMA cache_size=-16000")  // Reduce to 16MB
```

## Next Steps

- [Production Deployment](production.md) - Complete deployment guide
- [Monitoring Guide](monitoring.md) - Application monitoring
- [Scaling Guide](scaling.md) - Horizontal scaling strategies
