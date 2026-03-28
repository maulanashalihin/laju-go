# Troubleshooting

Common issues and solutions for Laju Go development and deployment.

## Quick Reference

| Problem | Quick Solution |
|---------|----------------|
| Port already in use | `lsof -ti:8080 \| xargs kill -9` |
| Database locked | Remove WAL files, restart |
| Vite port detection fails | Delete `.vite-port`, restart |
| Session not persisting | Check `SESSION_SECRET` |
| OAuth redirect mismatch | Verify Google Cloud Console URIs |
| Email not sending | Use Gmail App Password |
| Migration failed | Run manually with goose |
| Binary won't start | Check `.env` file exists |

---

## Development Issues

### Port Already in Use

**Error**: `listen tcp :8080: bind: address already in use`

**Cause**: Another process is using port 8080

**Solutions**:

```bash
# macOS/Linux - Kill process on port 8080
lsof -ti:8080 | xargs kill -9

# Alternative: Find and kill manually
lsof -i :8080
kill -9 <PID>

# Or change port in .env
APP_PORT=8081
```

**Prevention**: Always stop servers properly with Ctrl+C

---

### Database Locked

**Error**: `database is locked`

**Cause**: SQLite database file is locked by another process or unclosed transaction

**Solutions**:

```bash
# 1. Stop the application
# Ctrl+C or kill process

# 2. Remove WAL files
rm data/app.db-shm data/app.db-wal

# 3. Restart application
go run main.go
```

**If problem persists**:

```bash
# Check for zombie processes
ps aux | grep laju-go
kill -9 <PID>

# Verify no other process has the file open
lsof data/app.db
```

**Prevention**:
- Ensure WAL mode is enabled: `PRAGMA journal_mode=WAL`
- Set busy timeout: `PRAGMA busy_timeout=5000`
- Close all database connections properly

---

### Vite Port Detection Fails

**Error**: Go server can't connect to Vite dev server

**Cause**: `.vite-port` file is stale or missing

**Solutions**:

```bash
# Delete port cache file
rm .vite-port

# Restart Vite
npm run dev

# Verify Vite is running
curl http://localhost:5173
```

**Prevention**: Always start Vite before Go server

---

### Air Not Rebuilding

**Problem**: Changes to `.go` files don't trigger rebuild

**Solutions**:

1. **Check `.air.toml` configuration**:
   ```toml
   [build]
   include_ext = ["go", "tpl", "tmpl", "html"]
   exclude_dir = ["assets", "tmp", "vendor", "node_modules"]
   ```

2. **Verify file is not excluded**:
   ```bash
   # Check if file matches exclude patterns
   ls -la app/handlers/*.go
   ```

3. **Restart Air**:
   ```bash
   # Stop Air (Ctrl+C)
   # Clear tmp
   rm -rf tmp/
   
   # Restart
   air
   ```

4. **Check Air logs**:
   ```bash
   cat air.log
   ```

---

### HMR Not Working

**Problem**: Frontend changes don't appear in browser

**Solutions**:

1. **Verify Vite is running**:
   ```bash
   curl http://localhost:5173
   ```

2. **Clear browser cache**:
   - Chrome: Ctrl+Shift+Delete
   - Or use Incognito mode

3. **Check browser console**:
   ```
   F12 → Console → Look for errors
   ```

4. **Restart Vite**:
   ```bash
   # Stop Vite
   # Clear node_modules/.vite
   rm -rf node_modules/.vite
   
   # Restart
   npm run dev
   ```

---

### Go Module Errors

**Error**: `missing go.sum entry for module`

**Solutions**:

```bash
# Clean module cache
go clean -modcache

# Re-download dependencies
go mod download

# Tidy dependencies
go mod tidy

# Verify go.sum
cat go.sum | head
```

**If problem persists**:

```bash
# Remove go.sum and regenerate
rm go.sum
go mod tidy
```

---

### Node Module Issues

**Error**: `Cannot find module` or version conflicts

**Solutions**:

```bash
# Clean reinstall
rm -rf node_modules package-lock.json

# Install dependencies
npm install

# Verify installation
npm list --depth=0
```

**If problem persists**:

```bash
# Clear npm cache
npm cache clean --force

# Reinstall
npm install
```

---

### TypeScript Errors

**Error**: Various TypeScript compilation errors

**Solutions**:

```bash
# Check TypeScript config
cat tsconfig.json

# Verify types are installed
npm list --depth=0 | grep types

# Reinstall types
npm install -D @types/node @types/express
```

---

## Authentication Issues

### Session Not Persisting

**Problem**: Users are logged out after server restart

**Solutions**:

1. **Check `SESSION_SECRET` is set**:
   ```bash
   grep SESSION_SECRET .env
   ```

2. **Verify session store is database-backed**:
   ```go
   // app/session/session.go
   // Should save to database, not just cookies
   ```

3. **Check database has sessions table**:
   ```bash
   sqlite3 data/app.db ".tables"
   # Should show: sessions  users  goose_db_version
   ```

---

### OAuth State Mismatch

**Error**: "Invalid state" on OAuth callback

**Cause**: State token doesn't match between request and callback

**Solutions**:

1. **Verify state is stored in session**:
   ```go
   func (h *AuthHandler) storeState(c *fiber.Ctx, state string) {
       sess, _ := h.store.Get(c)
       sess.Set("oauth_state", state)
       sess.Save()
   }
   ```

2. **Check session secret is consistent**:
   ```bash
   # Should be same across restarts
   grep SESSION_SECRET .env
   ```

3. **Verify cookies are being sent**:
   - Check browser DevTools → Application → Cookies
   - Ensure `laju_session` cookie exists

---

### OAuth Redirect URI Mismatch

**Error**: `redirect_uri_mismatch` from Google

**Cause**: Redirect URI in request doesn't match Google Cloud Console

**Solutions**:

1. **Check `.env` redirect URL**:
   ```bash
   grep GOOGLE_REDIRECT_URL .env
   # Should be: http://localhost:8080/auth/google/callback
   ```

2. **Verify Google Cloud Console**:
   - Go to [Google Cloud Console](https://console.cloud.google.com/)
   - APIs & Services → Credentials
   - Edit OAuth client
   - Add exact redirect URI

3. **Match exactly** (including trailing slash):
   ```
   ✅ http://localhost:8080/auth/google/callback
   ❌ http://localhost:8080/auth/google/callback/
   ```

---

### Password Reset Email Not Received

**Problem**: User doesn't receive password reset email

**Solutions**:

1. **Check SMTP configuration**:
   ```bash
   grep SMTP_ .env
   ```

2. **For Gmail, use App Password**:
   - Go to Google Account → Security
   - Enable 2-Factor Authentication
   - Generate App Password
   - Use 16-character password in `.env`

3. **Check mailer logs**:
   ```bash
   # Look for email sending errors
   journalctl -u laju-go | grep -i mail
   ```

4. **Test SMTP connection**:
   ```bash
   # Using telnet
   telnet smtp.gmail.com 587
   ```

5. **Check spam folder**:
   - Email might be marked as spam

---

## Database Issues

### Migration Failed

**Error**: `Migration failed` or `duplicate column`

**Solutions**:

1. **Check migration status**:
   ```bash
   goose -dir migrations sqlite3 data/app.db status
   ```

2. **Run migrations manually**:
   ```bash
   goose -dir migrations sqlite3 data/app.db up
   ```

3. **If migration is broken**:
   ```bash
   # Rollback last migration
   goose -dir migrations sqlite3 data/app.db down
   
   # Fix migration file
   # Run again
   goose -dir migrations sqlite3 data/app.db up
   ```

4. **Reset all migrations** (destructive):
   ```bash
   # Delete database
   rm data/app.db
   
   # Restart server (migrations run automatically)
   go run main.go
   ```

---

### Database Not Found

**Error**: `unable to open database file`

**Solutions**:

```bash
# Create data directory
mkdir -p data

# Set permissions
chmod 755 data

# Restart server
go run main.go
```

**For production**:

```bash
# Create directory
sudo mkdir -p /var/lib/laju

# Set ownership
sudo chown www-data:www-data /var/lib/laju

# Set permissions
sudo chmod 755 /var/lib/laju
```

---

### Foreign Key Constraint Failed

**Error**: `FOREIGN KEY constraint failed`

**Cause**: Trying to insert/update data that violates foreign key constraints

**Solutions**:

1. **Check data integrity**:
   ```bash
   sqlite3 data/app.db "PRAGMA foreign_key_check;"
   ```

2. **Disable foreign keys temporarily** (not recommended):
   ```bash
   sqlite3 data/app.db "PRAGMA foreign_keys=OFF;"
   ```

3. **Fix data in correct order**:
   - Insert parent records first
   - Then insert child records

---

## Production Issues

### Service Won't Start

**Error**: `Failed to start Laju Go Application`

**Solutions**:

1. **Check logs**:
   ```bash
   sudo journalctl -u laju-go -n 50
   ```

2. **Common issues**:

   **Missing `.env`**:
   ```bash
   ls -la /opt/laju-go/.env
   ```

   **Wrong permissions**:
   ```bash
   sudo chown www-data:www-data /opt/laju-go
   sudo chmod 755 /opt/laju-go
   ```

   **Port in use**:
   ```bash
   sudo lsof -i :8080
   ```

3. **Test binary manually**:
   ```bash
   cd /opt/laju-go
   sudo -u www-data ./laju-go
   ```

---

### Nginx 502 Bad Gateway

**Error**: Nginx returns 502

**Cause**: Nginx can't connect to Go application

**Solutions**:

1. **Check if app is running**:
   ```bash
   curl http://localhost:8080/health
   ```

2. **Check Nginx logs**:
   ```bash
   sudo tail -f /var/log/nginx/error.log
   ```

3. **Verify Nginx config**:
   ```bash
   sudo nginx -t
   ```

4. **Check upstream is correct**:
   ```nginx
   upstream laju_app {
       server 127.0.0.1:8080;
   }
   ```

---

### SSL Certificate Issues

**Error**: Certificate expired or invalid

**Solutions**:

1. **Check certificate status**:
   ```bash
   sudo certbot certificates
   ```

2. **Renew certificate**:
   ```bash
   sudo certbot renew
   ```

3. **Force renewal**:
   ```bash
   sudo certbot renew --force-renewal
   ```

4. **Verify Nginx config**:
   ```bash
   sudo nginx -t
   sudo systemctl reload nginx
   ```

---

### High Memory Usage

**Problem**: Application using too much memory

**Solutions**:

1. **Reduce connection pool**:
   ```go
   db.SetMaxOpenConns(10)  // Reduce from 25
   ```

2. **Lower SQLite cache**:
   ```sql
   PRAGMA cache_size=-16000;  // Reduce to 16MB
   ```

3. **Check for memory leaks**:
   ```bash
   # Monitor memory over time
   watch -n 1 'ps aux | grep laju-go'
   ```

---

## Build Issues

### Frontend Build Fails

**Error**: Various build errors

**Solutions**:

```bash
# Clear cache
rm -rf node_modules/.vite
rm -rf dist

# Reinstall dependencies
npm install

# Rebuild
npm run build
```

**Check Node version**:
```bash
node --version  # Should be 18+
```

---

### Go Build Fails

**Error**: Compilation errors

**Solutions**:

```bash
# Check Go version
go version  # Should be 1.26+

# Clean build cache
go clean -cache

# Rebuild
go build -o laju-go .
```

**Check CGO for SQLite**:
```bash
# Install build tools
sudo apt install build-essential

# Rebuild
go build -o laju-go .
```

---

## Performance Issues

### Slow Queries

**Problem**: Database queries are slow

**Solutions**:

1. **Analyze query plan**:
   ```bash
   sqlite3 data/app.db "EXPLAIN QUERY PLAN SELECT * FROM users WHERE email = 'test@example.com';"
   ```

2. **Add index if full table scan**:
   ```sql
   CREATE INDEX idx_users_email ON users(email);
   ```

3. **Check database statistics**:
   ```bash
   sqlite3 data/app.db "ANALYZE;"
   ```

---

### Slow Page Loads

**Problem**: Pages take long to load

**Solutions**:

1. **Check frontend bundle size**:
   ```bash
   npm run build
   ls -lh dist/assets/
   ```

2. **Enable compression**:
   ```go
   import "github.com/gofiber/compression/v2"
   app.Use(compression.New())
   ```

3. **Check network tab**:
   - F12 → Network → Look for slow requests

---

## Getting Help

### Collect Debug Information

```bash
# System info
uname -a
go version
node --version

# Application logs
journalctl -u laju-go -n 100

# Database info
sqlite3 data/app.db ".schema"

# Network info
netstat -anp | grep 8080
```

### Report Issues

When reporting issues, include:

1. **Error message** (exact text)
2. **Steps to reproduce**
3. **Environment** (OS, Go version, Node version)
4. **Logs** (relevant excerpts)
5. **Configuration** (`.env` with secrets removed)

### Resources

- **GitHub Issues**: [maulanashalihin/laju-go/issues](https://github.com/maulanashalihin/laju-go/issues)
- **GitHub Discussions**: [maulanashalihin/laju-go/discussions](https://github.com/maulanashalihin/laju-go/discussions)
- **Documentation**: [docs/](./)

---

## Next Steps

- [API Reference](api-reference.md) - Complete endpoint documentation
- [Environment Guide](environment.md) - Configuration reference
- [Deployment Guide](../deployment/production.md) - Production setup
