#!/bin/bash

# Laju Go - First Deploy Script
# Sets up the application and systemd user service from scratch.
# Run AFTER deploy.sh has uploaded artifacts to the server.
#
# Uses systemd --user (no root/sudo needed for the service itself).
# Linger is enabled once (needs sudo) so the service survives logout.
#
# Prerequisites on server:
#   - sqlite3 installed (for seed_public.sql)
#   - Go installed (for seed_admin.go via `go run`)
#   - User has passwordless sudo (for `loginctl enable-linger` only)

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

# Set defaults
APP_NAME=${APP_NAME:-laju-go}
SERVICE_NAME=${SERVICE_NAME:-$APP_NAME}
APP_PORT=${APP_PORT:-8080}

echo -e "${BLUE}═══ FIRST DEPLOY ═══${NC}"
echo ""

# Interactive prompts for environment configuration
echo -e "${YELLOW}Application Port (default: $APP_PORT):${NC}"
read -r APP_PORT_INPUT
APP_PORT=${APP_PORT_INPUT:-$APP_PORT}

echo -e "${YELLOW}Application URL (e.g., https://yourdomain.com):${NC}"
read -r APP_URL

# Auto-generate SESSION_SECRET
SESSION_SECRET=$(openssl rand -hex 32)

# Step 1: Create remote directories (in home dir — no sudo needed)
echo -e "${YELLOW}[1/6] Creating remote directories...${NC}"
ssh "$SERVER_USER@$SERVER_HOST" "mkdir -p $SERVER_PATH/data $SERVER_PATH/storage $SERVER_PATH/backups ~/.config/systemd/user"
echo -e "${GREEN}      ✓ Directories created${NC}"

# Step 2: Setup .env file
echo -e "${YELLOW}[2/6] Setting up environment file...${NC}"
scp "$PROJECT_ROOT/.env.example" "$SERVER_USER@$SERVER_HOST:$SERVER_PATH/"

ssh "$SERVER_USER@$SERVER_HOST" "
    if [ ! -f $SERVER_PATH/.env ]; then
        cp $SERVER_PATH/.env.example $SERVER_PATH/.env
        sed -i 's/APP_PORT=8080/APP_PORT=$APP_PORT/g' $SERVER_PATH/.env
        sed -i \"s|APP_URL=http://localhost:8080|APP_URL=$APP_URL|g\" $SERVER_PATH/.env
        sed -i 's/APP_ENV=development/APP_ENV=production/g' $SERVER_PATH/.env
        sed -i 's/SESSION_SECRET=change-this-in-production/SESSION_SECRET=$SESSION_SECRET/g' $SERVER_PATH/.env
        sed -i \"s|DB_PATH=./data/app.db|DB_PATH=$SERVER_PATH/data/app.db|g\" $SERVER_PATH/.env
        echo '      Created .env from .env.example (production-ready)'
    else
        echo '      .env already exists, skipping'
    fi
"
echo -e "${GREEN}      ✓ Environment configured${NC}"

# Step 3: Enable linger (needs sudo once — so user services survive logout)
echo -e "${YELLOW}[3/6] Enabling linger (sudo needed once)...${NC}"
ssh "$SERVER_USER@$SERVER_HOST" "sudo loginctl enable-linger $SERVER_USER"
echo -e "${GREEN}      ✓ Linger enabled${NC}"

# Step 4: Upload systemd user service file (no sudo — it's in ~/.config)
echo -e "${YELLOW}[4/6] Setting up systemd user service...${NC}"

SERVICE_FILE="$PROJECT_ROOT/systemd/$APP_NAME.service"

if [ -f "$SERVICE_FILE" ]; then
    scp "$SERVICE_FILE" "$SERVER_USER@$SERVER_HOST:~/.config/systemd/user/$SERVICE_NAME.service"
else
    echo -e "${RED}Error: Service file not found at $SERVICE_FILE${NC}"
    exit 1
fi

ssh "$SERVER_USER@$SERVER_HOST" "
    systemctl --user daemon-reload
    systemctl --user enable $SERVICE_NAME
    systemctl --user start $SERVICE_NAME
"

sleep 2
echo -e "${GREEN}      ✓ Service created and started${NC}"

# Step 5: Set permissions for data directories
echo -e "${YELLOW}[5/6] Setting up permissions...${NC}"
ssh "$SERVER_USER@$SERVER_HOST" "
    chmod 755 $SERVER_PATH/data
    chmod 770 $SERVER_PATH/storage
    chmod 770 $SERVER_PATH/backups
"
echo -e "${GREEN}      ✓ Permissions set${NC}"

# Step 6: Seed admin user (if seed_admin.go was uploaded)
echo -e "${YELLOW}[6/6] Seeding admin user...${NC}"
if ssh "$SERVER_USER@$SERVER_HOST" "test -f $SERVER_PATH/seed_admin.go" 2>/dev/null; then
    ssh "$SERVER_USER@$SERVER_HOST" "cd $SERVER_PATH && go run seed_admin.go -db ./data/app.db" 2>&1 || \
        echo -e "${YELLOW}      ! seed_admin.go failed (you can create admin manually later)${NC}"
else
    echo -e "${YELLOW}      ! seed_admin.go not found, skipping${NC}"
fi
echo -e "${GREEN}      ✓ Admin seed done${NC}"

# Verify
echo ""
echo -e "${BLUE}Verifying service...${NC}"
if ssh "$SERVER_USER@$SERVER_HOST" "systemctl --user is-active $SERVICE_NAME" > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Service is running${NC}"
else
    echo -e "${RED}Service failed to start. Check logs:${NC}"
    ssh "$SERVER_USER@$SERVER_HOST" "journalctl --user -u $SERVICE_NAME -n 30 --no-pager"
    exit 1
fi

echo ""
echo -e "${GREEN}═══ FIRST DEPLOY COMPLETE ═══${NC}"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "  1. Configure .env if needed:"
echo "     ssh $SERVER_USER@$SERVER_HOST 'nano $SERVER_PATH/.env'"
echo ""
echo "  2. Seed public data (if you have seed_public.sql):"
echo "     ssh $SERVER_USER@$SERVER_HOST 'sqlite3 $SERVER_PATH/data/app.db < $SERVER_PATH/seed_public.sql'"
echo ""
echo "  3. Default admin credentials (change password after first login):"
echo "     email:    admin@laju-go.local"
echo "     password: Admin123!"
echo ""
