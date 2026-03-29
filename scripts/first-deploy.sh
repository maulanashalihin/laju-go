#!/bin/bash

# Laju Go - First Deploy Script
# Sets up the application and systemd service from scratch

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

source "$PROJECT_ROOT/.deploy"

echo -e "${BLUE}═══ FIRST DEPLOY ═══${NC}"
echo ""

# Interactive prompts for environment configuration
echo -e "${YELLOW}Application Port (default: 8080):${NC}"
read APP_PORT_INPUT
APP_PORT=${APP_PORT_INPUT:-8080}

echo -e "${YELLOW}Application URL (e.g., https://yourdomain.com):${NC}"
read APP_URL

# Auto-generate SESSION_SECRET
SESSION_SECRET=$(openssl rand -hex 32)

# Step 1: Create directory
echo -e "${YELLOW}[1/4] Creating remote directory...${NC}"
ssh "$SERVER_USER@$SERVER_HOST" "mkdir -p $SERVER_PATH"
echo -e "${GREEN}      ✓ Directory created${NC}"

# Step 2: Clone or setup repository
echo -e "${YELLOW}[2/4] Setting up repository...${NC}"
ssh "$SERVER_USER@$SERVER_HOST" "
    if [ -d '$SERVER_PATH/.git' ]; then
        echo '      Repository already exists, pulling changes...'
        cd $SERVER_PATH && git pull origin main
    else
        echo '      Cloning repository...'
        rm -rf $SERVER_PATH
        git clone $REPO_URL $SERVER_PATH
    fi
"
echo -e "${GREEN}      ✓ Repository ready${NC}"

# Step 3: Setup .env file
echo -e "${YELLOW}[3/4] Setting up environment file...${NC}"
ssh "$SERVER_USER@$SERVER_HOST" "
    if [ ! -f $SERVER_PATH/.env ]; then
        if [ -f $SERVER_PATH/.env.example ]; then
            cp $SERVER_PATH/.env.example $SERVER_PATH/.env
            # Auto-update production values
            sed -i 's/APP_PORT=8080/APP_PORT=$APP_PORT/g' $SERVER_PATH/.env
            sed -i \"s|APP_URL=http://localhost:8080|APP_URL=$APP_URL|g\" $SERVER_PATH/.env
            sed -i 's/APP_ENV=development/APP_ENV=production/g' $SERVER_PATH/.env
            sed -i 's/SESSION_SECRET=your-secret-key-change-this-in-production/SESSION_SECRET=$SESSION_SECRET/g' $SERVER_PATH/.env
            echo '      Created .env from .env.example (production-ready)'
        else
            echo '      ⚠️  .env.example not found - you must create .env manually'
        fi
    else
        echo '      .env already exists'
    fi
"
echo -e "${GREEN}      ✓ Environment configured${NC}"

# Step 4: Create and start systemd service
echo -e "${YELLOW}[4/4] Setting up systemd service...${NC}"

# Set default service name if not set
SERVICE_NAME=${SERVICE_NAME:-crm-maulanabuilds}

# Copy service file
scp "$PROJECT_ROOT/systemd/laju-go.service" "$SERVER_USER@$SERVER_HOST:/etc/systemd/system/$SERVICE_NAME.service"

# Update service file with correct path and service name
ssh "$SERVER_USER@$SERVER_HOST" "
    sed -i 's|/opt/crm-maulanabuilds|'$SERVER_PATH'|g' /etc/systemd/system/$SERVICE_NAME.service
    sed -i 's|SyslogIdentifier=crm-maulanabuilds|SyslogIdentifier='$SERVICE_NAME'|g' /etc/systemd/system/$SERVICE_NAME.service
"

# Enable and start
ssh "$SERVER_USER@$SERVER_HOST" "
    systemctl daemon-reload
    systemctl enable $SERVICE_NAME
    systemctl start $SERVICE_NAME
"

sleep 2
echo -e "${GREEN}      ✓ Service started${NC}"

# Verify
echo ""
echo -e "${BLUE}Verifying service...${NC}"
ssh "$SERVER_USER@$SERVER_HOST" "systemctl is-active $SERVICE_NAME" || {
    echo -e "${RED}Service failed to start. Check logs:${NC}"
    ssh "$SERVER_USER@$SERVER_HOST" "journalctl -u $SERVICE_NAME -n 30 --no-pager"
    exit 1
}

echo ""
echo -e "${GREEN}═══ FIRST DEPLOY COMPLETE ═══${NC}"
echo ""
echo -e "${YELLOW}Required: Configure OAuth & SMTP for full functionality:${NC}"
echo "  Edit .env: ssh $SERVER_USER@$SERVER_HOST 'nano $SERVER_PATH/.env'"
echo ""
echo "  # Google OAuth (get from console.cloud.google.com)"
echo "  GOOGLE_CLIENT_ID=your-client-id"
echo "  GOOGLE_CLIENT_SECRET=your-secret"
echo "  GOOGLE_REDIRECT_URL=$APP_URL/auth/google/callback"
echo ""
echo "  # SMTP (for password reset)"
echo "  SMTP_HOST=smtp.gmail.com"
echo "  SMTP_USER=your-email@gmail.com"
echo "  SMTP_PASS=your-app-password"
echo "  FROM_EMAIL=noreply@yourdomain.com"
echo ""
