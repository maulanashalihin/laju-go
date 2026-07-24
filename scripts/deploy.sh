#!/bin/bash

# Laju Go - One-Click Deploy Script
# Rsyncs source code to server, builds on server, restarts service.
# Uses systemd --user service (no root/sudo needed for the service itself).
#
# Requirements:
#   - SSH key access to server
#   - .deploy file configured (cp .deploy.example .deploy)
#   - Server has: Go, Node/npm, sqlite3, rsync
#   - Server user has passwordless sudo (only for `loginctl enable-linger`,
#     run once during first-deploy)

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║     Laju Go - One-Click Deploy        ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""

# Load configuration
if [ ! -f "$PROJECT_ROOT/.deploy" ]; then
    echo -e "${RED}Error: .deploy file not found!${NC}"
    echo -e "${YELLOW}Please copy .deploy.example to .deploy and configure it first.${NC}"
    echo ""
    echo "  cp .deploy.example .deploy"
    echo ""
    exit 1
fi

source "$PROJECT_ROOT/.deploy"

# Set defaults
APP_NAME=${APP_NAME:-laju-go}
SERVICE_NAME=${SERVICE_NAME:-$APP_NAME}
APP_PORT=${APP_PORT:-8080}

# Validate required variables
if [ -z "$SERVER_USER" ] || [ -z "$SERVER_HOST" ] || [ -z "$SERVER_PATH" ]; then
    echo -e "${RED}Error: Missing required variables in .deploy file${NC}"
    echo "Please ensure SERVER_USER, SERVER_HOST, and SERVER_PATH are set."
    exit 1
fi

echo -e "${GREEN}App:      ${YELLOW}$APP_NAME${NC}"
echo -e "${GREEN}Server:   ${YELLOW}$SERVER_USER@$SERVER_HOST${NC}"
echo -e "${GREEN}Path:     ${YELLOW}$SERVER_PATH${NC}"
echo -e "${GREEN}Port:     ${YELLOW}$APP_PORT${NC}"
echo ""

# Test SSH connection
echo -e "${BLUE}Testing SSH connection...${NC}"
if ! ssh -o ConnectTimeout=10 -o BatchMode=yes "$SERVER_USER@$SERVER_HOST" "echo OK" > /dev/null 2>&1; then
    echo -e "${RED}Error: Cannot connect to server via SSH${NC}"
    echo "Please check your SSH credentials and server accessibility."
    exit 1
fi
echo -e "${GREEN}✓ SSH connection successful${NC}"
echo ""

# Detect first vs update deploy by checking if user service exists
echo -e "${BLUE}Checking deployment status...${NC}"
IS_FIRST=false
if ssh "$SERVER_USER@$SERVER_HOST" "systemctl --user is-active $SERVICE_NAME" > /dev/null 2>&1; then
    echo -e "${GREEN}→ Existing deployment detected${NC}"
else
    echo -e "${YELLOW}→ No existing deployment found${NC}"
    IS_FIRST=true
fi
echo ""

# Rsync source code to server (exclude build artifacts, data, env, etc.)
echo -e "${BLUE}Syncing source code to server...${NC}"
rsync -az --delete \
    --exclude='.git' \
    --exclude='node_modules' \
    --exclude='data' \
    --exclude='storage' \
    --exclude='dist' \
    --exclude='.env' \
    --exclude='.deploy' \
    --exclude='.vite-port' \
    --exclude='*.db' \
    --exclude='*.db-journal' \
    --exclude='*.db-wal' \
    --exclude='*.db-shm' \
    --exclude="$APP_NAME" \
    --exclude="$APP_NAME.exe" \
    --exclude='.DS_Store' \
    --exclude='.air.toml' \
    --exclude='.zero' \
    --exclude='.pi' \
    "$PROJECT_ROOT/" "$SERVER_USER@$SERVER_HOST:$SERVER_PATH/"
echo -e "${GREEN}✓ Source code synced${NC}"
echo ""

# Build on server
echo -e "${BLUE}Building on server...${NC}"
ssh "$SERVER_USER@$SERVER_HOST" "cd $SERVER_PATH && npm run build 2>&1 && go build -o $APP_NAME ./cmd/laju-go 2>&1 && chmod +x $APP_NAME"
echo -e "${GREEN}✓ Built on server${NC}"
echo ""

# Run first-deploy or update-deploy
if [ "$IS_FIRST" = true ]; then
    "$SCRIPT_DIR/first-deploy.sh"
else
    "$SCRIPT_DIR/update-deploy.sh"
fi

# Final status
echo ""
echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║        Deployment Status             ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""

if ssh "$SERVER_USER@$SERVER_HOST" "systemctl --user is-active $SERVICE_NAME" &>/dev/null; then
    echo -e "${GREEN}✓ Service $SERVICE_NAME is running${NC}"
else
    echo -e "${RED}✗ Service $SERVICE_NAME is not running${NC}"
fi
if ssh "$SERVER_USER@$SERVER_HOST" "systemctl --user is-enabled $SERVICE_NAME" &>/dev/null; then
    echo -e "${GREEN}✓ Service enabled (auto-start on boot)${NC}"
else
    echo -e "${YELLOW}! Service not enabled${NC}"
fi

echo ""
echo -e "${CYAN}Useful commands:${NC}"
echo "  View logs:     ssh $SERVER_USER@$SERVER_HOST 'journalctl --user -u $SERVICE_NAME -f'"
echo "  Check status:  ssh $SERVER_USER@$SERVER_HOST 'systemctl --user status $SERVICE_NAME'"
echo "  Restart:       ssh $SERVER_USER@$SERVER_HOST 'systemctl --user restart $SERVICE_NAME'"
echo ""
echo -e "${GREEN}Deployment complete!${NC}"
