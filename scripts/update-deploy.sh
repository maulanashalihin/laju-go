#!/bin/bash

# Laju Go - Update Deploy Script
# Updates an existing installation via git pull

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

echo -e "${BLUE}═══ UPDATE DEPLOY ═══${NC}"
echo ""

# Set default service name if not set
SERVICE_NAME=${SERVICE_NAME:-crm-maulanabuilds}

# Step 1: Stop service
echo -e "${YELLOW}[1/4] Stopping service...${NC}"
ssh "$SERVER_USER@$SERVER_HOST" "systemctl stop $SERVICE_NAME"
echo -e "${GREEN}      ✓ Service stopped${NC}"

# Step 2: Upload binary and dist folder
echo -e "${YELLOW}[2/4] Uploading binary and assets...${NC}"
scp laju-go "$SERVER_USER@$SERVER_HOST:$SERVER_PATH/laju-go"
scp -r dist "$SERVER_USER@$SERVER_HOST:$SERVER_PATH/dist"
ssh "$SERVER_USER@$SERVER_HOST" "chmod +x $SERVER_PATH/laju-go"
echo -e "${GREEN}      ✓ Assets uploaded${NC}"

# Step 3: Pull latest source code (for code changes, migrations, etc.)
echo -e "${YELLOW}[3/4] Pulling latest source code...${NC}"
ssh "$SERVER_USER@$SERVER_HOST" "
    cd $SERVER_PATH
    git fetch origin main
    git reset --hard origin/main  # Reset all source files to latest
"
echo -e "${GREEN}      ✓ Source code updated${NC}"

# Step 4: Restart service
echo -e "${YELLOW}[4/4] Restarting service...${NC}"
ssh "$SERVER_USER@$SERVER_HOST" "
    systemctl restart $SERVICE_NAME
"
sleep 2
echo -e "${GREEN}      ✓ Service restarted${NC}"

# Verify
echo ""
echo -e "${BLUE}Verifying service...${NC}"
ssh "$SERVER_USER@$SERVER_HOST" "systemctl is-active $SERVICE_NAME" || {
    echo -e "${RED}Service failed to start. Check logs:${NC}"
    ssh "$SERVER_USER@$SERVER_HOST" "journalctl -u $SERVICE_NAME -n 30 --no-pager"
    exit 1
}

# Show recent logs
echo ""
echo -e "${BLUE}Recent logs:${NC}"
ssh "$SERVER_USER@$SERVER_HOST" "journalctl -u $SERVICE_NAME -n 5 --no-pager"

echo ""
echo -e "${GREEN}═══ UPDATE DEPLOY COMPLETE ═══${NC}"
