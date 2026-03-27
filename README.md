# VeloStack Go

High-performance SaaS boilerplate built with **Go Fiber** + **Svelte 5** + **SQLite**.

## 🚀 Features

- **Backend**: Go Fiber (fasthttp) for blazing-fast performance
- **Frontend**: Svelte 5 with Inertia.js for reactive SPA experience
- **Database**: SQLite with Squirrel query builder
- **Authentication**: Email/Password + Google OAuth
- **Session Management**: Secure session-based authentication
- **Role-Based Access**: Admin/User roles out of the box
- **File Upload**: Avatar upload with validation
- **Database Migrations**: Using Goose for schema management
- **Docker Ready**: Multi-stage build for production deployment

## 📁 Project Structure

```
velostack-go/
├── main.go                  # Entry point
├── go.mod                   # Go dependencies
├── package.json             # Node.js dependencies
├── vite.config.js           # Vite configuration
│
├── app/                     # Go backend code
│   ├── config/              # Environment configuration
│   ├── models/              # Data structures (user.model.go, dto.model.go)
│   ├── repositories/        # Database queries (user.repo.go)
│   ├── services/            # Business logic (auth.service.go, inertia.service.go)
│   ├── handlers/            # HTTP controllers (auth.handler.go, etc.)
│   ├── middleware/          # Auth guards, CORS (auth.middleware.go)
│   └── session/             # Session management
│
├── routes/                  # Route definitions
│   └── web.go               # Web routes setup
│
├── frontend/                # Svelte 5 frontend
│   ├── src/
│   │   ├── components/      # Reusable UI components
│   │   ├── layouts/         # Page layouts
│   │   ├── pages/           # Page components (Auth/, App/, Admin/)
│   │   ├── main.ts          # Entry point
│   │   └── app.css          # Global styles
│   ├── package.json
│   └── vite.config.js
│
├── templates/               # HTML templates for backend rendering
│   └── inertia.html         # Inertia.js base template
│
├── migrations/              # Database migrations (Goose)
├── data/                    # SQLite database (gitignored)
├── storage/                 # User uploads (gitignored)
├── dist/                    # Built frontend assets (production)
└── public/                  # Static assets
```

## 🛠️ Tech Stack

| Layer | Technology |
|-------|------------|
| Web Framework | Go Fiber v2 |
| Template Engine | Fiber HTML Template |
| Database | SQLite3 |
| Query Builder | Squirrel |
| Session Store | Fiber Session (Memory) |
| OAuth | golang.org/x/oauth2 |
| Frontend | Svelte 5 |
| Build Tool | Vite |
| Styling | Tailwind CSS |
| SPA Router | Inertia.js |
| Migrations | Goose |

## 📦 Getting Started

### Prerequisites

- Go 1.22+
- Node.js 18+
- SQLite3

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/velostack/velostack-go.git
   cd velostack-go
   ```

2. **Install Go dependencies**
   ```bash
   go mod download
   ```

3. **Install Node.js dependencies**
   ```bash
   npm install
   ```

4. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env with your settings
   ```

5. **Set up Google OAuth (optional)**
   - Go to [Google Cloud Console](https://console.cloud.google.com/)
   - Create a new project
   - Enable Google+ API
   - Create OAuth 2.0 credentials
   - Add `http://localhost:8080/auth/google/callback` to authorized redirect URIs
   - Copy Client ID and Secret to `.env`

### Development

1. **Start the Go server**
   ```bash
   go run .
   ```

2. **Start the Vite dev server (in another terminal)**
   ```bash
   cd frontend && npm run dev
   ```

3. **Open your browser**
   - App: http://localhost:8000
   - Vite will auto-reload on frontend changes
   - Go server will auto-reload on backend changes (Go 1.22+)

### Building for Production

1. **Build frontend assets**
   ```bash
   npm run build
   ```

2. **Build Go binary**
   ```bash
   go build -o velostack-go .
   ```

3. **Run the binary**
   ```bash
   ./velostack-go
   ```

### Using Docker

1. **Build the image**
   ```bash
   docker build -t velostack-go .
   ```

2. **Run the container**
   ```bash
   docker run -p 8080:8080 \
     -v $(pwd)/data:/root/data \
     -v $(pwd)/storage:/root/storage \
     velostack-go
   ```

## 🔐 Default Admin Setup

To create an admin user, you can manually set the role in the database:

```sql
UPDATE users SET role = 'admin' WHERE email = 'your@email.com';
```

## 📝 API Endpoints

### Public Routes
- `GET /` - Home page
- `GET /about` - About page
- `GET /login` - Login page
- `POST /login/login` - User login
- `POST /login/register` - User registration
- `GET /auth/google` - Google OAuth login
- `GET /auth/google/callback` - Google OAuth callback

### Protected Routes
- `GET /app` - Dashboard
- `GET /app/profile` - User profile
- `PUT /app/profile` - Update profile
- `POST /app/upload` - File upload
- `POST /logout` - Logout

### Admin Routes
- `GET /admin` - Admin dashboard

## 🗄️ Database Migrations

Migrations are automatically run on startup. To manually run migrations:

```bash
# Install goose
go install github.com/pressly/goose/v3/cmd/goose@latest

# Run migrations
goose -dir migrations sqlite3 ./data/app.db up

# Check status
goose -dir migrations sqlite3 ./data/app.db status
```

## 🚀 Production Deployment (Ubuntu Linux)

### Prerequisites

```bash
# Install Go
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Install build dependencies
sudo apt update
sudo apt install -y build-essential

# Install Node.js (for building frontend)
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt install -y nodejs
```

### Build & Deploy

```bash
# 1. Clone repository
git clone https://github.com/velostack/velostack-go.git /opt/velostack-go
cd /opt/velostack-go

# 2. Install dependencies
go mod download
npm install

# 3. Build frontend assets
npm run build

# 4. Build Go binary
go build -o velostack-go .

# 5. Configure environment
cp .env.example .env
nano .env  # Edit with your settings
```

### Create Systemd Service

```bash
sudo nano /etc/systemd/system/velostack-go.service
```

Add the following content:

```ini
[Unit]
Description=VeloStack Go Application
After=network.target

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/opt/velostack-go
ExecStart=/opt/velostack-go/velostack-go
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

# Security hardening
NoNewPrivileges=true
PrivateTmp=true

# Environment
Environment="PATH=/usr/local/go/bin:/usr/bin:/bin"

[Install]
WantedBy=multi-user.target
```

### Setup Permissions

```bash
# Create data and storage directories
sudo mkdir -p /opt/velostack-go/data /opt/velostack-go/storage/avatars

# Set ownership
sudo chown -R www-data:www-data /opt/velostack-go

# Set permissions (SQLite needs write access to data directory)
sudo chmod 755 /opt/velostack-go
sudo chmod 770 /opt/velostack-go/data
sudo chmod 770 /opt/velostack-go/storage
```

### Start & Enable Service

```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable on boot
sudo systemctl enable velostack-go

# Start service
sudo systemctl start velostack-go

# Check status
sudo systemctl status velostack-go

# View logs
sudo journalctl -u velostack-go -f
```

### SQLite Production Optimizations

The application includes the following SQLite optimizations:

| Setting | Value | Purpose |
|---------|-------|---------|
| `journal_mode` | WAL | Write-Ahead Logging for better concurrency |
| `synchronous` | NORMAL | Balance between durability and speed |
| `cache_size` | 64MB | Memory cache for faster queries |
| `temp_store` | MEMORY | Faster temporary table operations |
| `busy_timeout` | 5000ms | Wait for database locks instead of failing |
| Connection Pool | 25 max open, 5 idle | Efficient connection management |

### Database Maintenance

**Backup (online, no downtime):**

```bash
# Create backup using SQLite backup API
sqlite3 /opt/velostack-go/data/app.db ".backup '/opt/velostack-go/backups/app-$(date +%Y%m%d).db'"
```

**Checkpoint WAL (optional, for maintenance):**

```bash
# Checkpoint WAL to main database file
sqlite3 /opt/velostack-go/data/app.db "PRAGMA wal_checkpoint(TRUNCATE);"
```

**Vacuum (reclaim space, requires downtime):**

```bash
# Stop service first
sudo systemctl stop velostack-go

# Vacuum database
sqlite3 /opt/velostack-go/data/app.db "VACUUM;"

# Restart service
sudo systemctl start velostack-go
```

### Nginx Reverse Proxy (Optional)

```bash
sudo apt install -y nginx
sudo nano /etc/nginx/sites-available/velostack-go
```

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

```bash
# Enable site
sudo ln -s /etc/nginx/sites-available/velostack-go /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### SSL with Let's Encrypt

```bash
sudo apt install -y certbot python3-certbot-nginx
sudo certbot --nginx -d your-domain.com
```

## 🧪 Testing

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...
```

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 🙏 Acknowledgments

- [Go Fiber](https://gofiber.io/) - Fast web framework
- [Svelte](https://svelte.dev/) - Cybernetically enhanced web apps
- [Inertia.js](https://inertiajs.com/) - Server-driven SPA
