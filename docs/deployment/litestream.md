# Litestream — SQLite Disaster Recovery

[Litestream](https://litestream.io/) is a standalone disaster recovery tool for SQLite that continuously replicates your database to remote storage (S3, S3-compatible, SFTP, or local filesystem) in real-time.

## Why Litestream?

SQLite's `sqlite3 .backup` (cron-based backup) has gaps:

| Aspect | Cron Backup | Litestream |
|--------|------------|------------|
| **Recovery Point** | Last cron run (up to 24h loss) | Continuous (seconds of data loss) |
| **Point-in-Time Recovery** | ❌ No | ✅ Yes — restore to any timestamp |
| **Server Loss** | ❌ Backup lost with server | ✅ Data safe in S3 |
| **Setup Complexity** | Simple (cron job) | Moderate (config file) |

**Bottom line**: Use Litestream for production. Use cron backup for staging/development.

## How It Works

```
┌──────────────────┐      ┌────────────────┐      ┌──────────────────────┐
│   Laju Go App    │      │   Litestream   │      │   S3 / S3-compatible │
│   (Go Fiber)     │──────▶  (background)  │──────▶   (continuous sync)  │
│   data/app.db    │ WAL  │  /usr/bin/     │ LTX  │   s3://bucket/app/   │
└──────────────────┘      │  litestream    │      └──────────────────────┘
                          └────────────────┘
```

Litestream reads the SQLite **WAL** (Write-Ahead Log) as changes happen and uploads them as LTX files to S3. It never writes to your database — zero risk of corruption.

## Installation

### Linux (amd64)

```bash
# Download latest binary
curl -fsSL https://github.com/benbjohnson/litestream/releases/latest/download/litestream-linux-amd64.tar.gz \
  -o litestream.tar.gz

# Extract and install
tar xzf litestream.tar.gz
sudo mv litestream /usr/local/bin/litestream

# Verify
litestream version
```

### macOS

```bash
brew install litestream
```

## Configuration

Create `/etc/litestream.yml`:

```yaml
# /etc/litestream.yml
dbs:
  - path: /opt/laju-go/data/app.db
    replicas:
      - type: s3
        bucket: your-laju-backups
        path: app/production
        region: ap-southeast-1
        access-key-id: ${AWS_ACCESS_KEY_ID}
        secret-access-key: ${AWS_SECRET_ACCESS_KEY}
      - type: s3
        bucket: your-laju-backups-dr
        path: app/production
        region: ap-southeast-3
        access-key-id: ${AWS_ACCESS_KEY_ID_DR}
        secret-access-key: ${AWS_SECRET_ACCESS_KEY_DR}
```

### Multi-Replica (Recommended)

Always configure **two replicas** in different regions. If one region goes down, you can restore from the other:

```yaml
replicas:
  - type: s3
    bucket: laju-backups-primary
    path: app/production
    region: ap-southeast-1  # Singapore
  - type: s3
    bucket: laju-backups-dr
    path: app/production
    region: ap-northeast-1  # Tokyo
```

### S3-Compatible (MinIO, DigitalOcean Spaces, Backblaze B2)

```yaml
replicas:
  - type: s3
    bucket: my-backups
    path: laju/production
    endpoint: https://sgp1.digitaloceanspaces.com
    region: sgp1
    access-key-id: ${SPACES_KEY}
    secret-access-key: ${SPACES_SECRET}
```

### Local File (Additional)

For an extra layer of safety, add a local replica:

```yaml
replicas:
  - type: s3
    bucket: laju-backups
    path: app/production
    region: ap-southeast-1
  - type: file
    path: /opt/laju-go/backups/litestream
```

## Running as Systemd Service

Create `/etc/systemd/system/litestream.service`:

```ini
[Unit]
Description=Litestream SQLite Replication
Documentation=https://litestream.io
After=network.target
Requires=laju-go.service

[Service]
Type=simple
User=root
EnvironmentFile=/opt/laju-go/.env
ExecStart=/usr/local/bin/litestream replicate --config /etc/litestream.yml
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

**Important**: Litestream must run **after** the app starts (it needs the WAL file to exist).

### Enable and Start

```bash
sudo systemctl daemon-reload
sudo systemctl enable litestream
sudo systemctl start litestream

# Check status
sudo systemctl status litestream

# View logs
sudo journalctl -u litestream -f
```

### Verify Replication

```bash
# List all snapshots and LTX files
litestream databases -config /etc/litestream.yml

# Check generation status
litestream generations -config /etc/litestream.yml /opt/laju-go/data/app.db
```

## Restoring from Backup

### Restore to Latest State

```bash
# Stop the app first
sudo systemctl stop laju-go

# Restore latest from S3
litestream restore -config /etc/litestream.yml \
  -o /opt/laju-go/data/app.db \
  s3://laju-backups/app/production

# Start the app
sudo systemctl start laju-go
```

### Point-in-Time Recovery

Restore to a specific moment (e.g., just before a disaster):

```bash
litestream restore -config /etc/litestream.yml \
  -o /opt/laju-go/data/app.db \
  -timestamp "2026-05-06T14:30:00Z" \
  s3://laju-backups/app/production
```

### Restore to a New Server

On a brand new server with no database:

```bash
# 1. Install laju-go binary + migrations + dist (see deployment guide)

# 2. Install Litestream
sudo apt install -y litestream  # or download binary

# 3. Restore database
litestream restore \
  -if-db-not-exists \
  -o /opt/laju-go/data/app.db \
  s3://laju-backups/app/production

# 4. Start litestream replication
sudo systemctl start litestream

# 5. Start laju-go
sudo systemctl start laju-go
```

### Integrity Check

Always verify after restore:

```bash
# Check database integrity
sqlite3 /opt/laju-go/data/app.db "PRAGMA integrity_check;"

# Verify Litestream replication health
litestream databases -config /etc/litestream.yml
```

## Full Disaster Recovery Playbook

When your server is completely gone:

```bash
# 1. Spin up a new server (any provider)
# 2. Install laju-go, litestream, systemd service
# 3. Restore database from S3
litestream restore \
  -if-db-not-exists \
  -o /opt/laju-go/data/app.db \
  s3://laju-backups/app/production

# 4. Run migrations (in case schema changed after backup)
/opt/laju-go/laju-go  # migrations run automatically on startup

# 5. Start replication + app
sudo systemctl start litestream
sudo systemctl start laju-go
```

**Total recovery time**: ~5-10 minutes (mostly DNS propagation).

## Monitoring

### Health Check Script

```bash
#!/bin/bash
# /opt/laju-go/scripts/check-litestream.sh

# Check if litestream is running
if ! pgrep -x litestream > /dev/null; then
    echo "Litestream is not running!"
    sudo systemctl restart litestream
    exit 1
fi

# Check last replication time (compare against WAL write time)
LAST_REPLICA=$(ls -lt /opt/laju-go/backups/litestream/ 2>/dev/null | head -2 | tail -1)
echo "Litestream is running. Last replica: $LAST_REPLICA"
```

### Prometheus / Grafana

Litestream exposes Prometheus metrics:

```yaml
# Add to /etc/litestream.yml
metrics:
  addr: ":9399"
```

Then scrape `http://localhost:9399/metrics`.

## Cost Estimation

Litestream uploads only WAL changes (typically **< 1% of DB size per day** for most apps):

| Database Size | Daily LTX | Monthly S3 Cost |
|---------------|-----------|-----------------|
| 100MB | ~1-2 MB | < $0.01 |
| 1GB | ~10-20 MB | ~$0.01 |
| 10GB | ~100-200 MB | ~$0.10 |

S3 costs are negligible. The main cost is the S3 bucket itself (~$0.023/GB/month).

## Comparison: Litestream vs Cron Backup

| Feature | Litestream | Cron + `.backup` |
|---------|-----------|------------------|
| **Recovery Point Objective** | Seconds | Hours |
| **Recovery Time Objective** | Minutes | Minutes |
| **Point-in-Time Recovery** | ✅ Yes | ❌ No |
| **Cross-Region DR** | ✅ Yes (multi-replica) | ❌ Manual |
| **Automated Restore** | ✅ One command | ❌ Manual SCP |
| **Complexity** | Moderate | Simple |
| **Cost** | ~$0.10/mo S3 | Free |

## Next Steps

- [Production Deployment](production.md) — Complete production setup
- [SQLite Configuration](sqlite-configuration.md) — Tuning SQLite for your server
