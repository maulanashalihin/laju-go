#!/bin/bash

# Laju Go - One-Click Deploy Script
# Auto-detects first deploy vs update and handles everything

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
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

# Validate required variables
if [ -z "$SERVER_USER" ] || [ -z "$SERVER_HOST" ] || [ -z "$SERVER_PATH" ]; then
    echo -e "${RED}Error: Missing required variables in .deploy file${NC}"
    echo "Please ensure SERVER_USER, SERVER_HOST, and SERVER_PATH are set."
    exit 1
fi

if [ -z "$REPO_URL" ]; then
    echo -e "${RED}Error: REPO_URL is not set in .deploy file${NC}"
    exit 1
fi

echo -e "${GREEN}Server:   ${YELLOW}$SERVER_USER@$SERVER_HOST${NC}"
echo -e "${GREEN}Path:     ${YELLOW}$SERVER_PATH${NC}"
echo -e "${GREEN}Repo:     ${YELLOW}$REPO_URL${NC}"
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

# Build assets locally before deploy
echo -e "${BLUE}Building assets locally...${NC}"

# Build frontend
echo -e "${YELLOW}Building frontend...${NC}"
npm run build
echo -e "${GREEN}✓ Frontend built${NC}"

# Build Go binary for Linux
echo -e "${YELLOW}Building Go binary (linux/amd64)...${NC}"
GOOS=linux GOARCH=amd64 go build -o laju-go .
echo -e "${GREEN}✓ Binary built${NC}"

# Commit and push all changes (source code + build artifacts)
echo -e "${YELLOW}Committing all changes...${NC}"
git add .

# Check if there are changes to commit
if git diff --staged --quiet; then
    echo -e "${CYAN}      No changes to commit${NC}"
else
    # Count changed files for commit message
    CHANGED_FILES=$(git diff --staged --name-only | wc -l | tr -d ' ')
    git commit -m "Deploy: ${CHANGED_FILES} files changed"
fi

git push origin main
echo -e "${GREEN}✓ Changes pushed${NC}"

echo ""

# Check if deployment exists (check both service and directory)
echo -e "${BLUE}Checking deployment status...${NC}"

# Check if directory and git repo exist
if ssh "$SERVER_USER@$SERVER_HOST" "[ -d '$SERVER_PATH/.git' ]" 2>/dev/null; then
    echo -e "${GREEN}→ Existing deployment detected${NC}"
    echo ""
    "$SCRIPT_DIR/update-deploy.sh"
else
    echo -e "${YELLOW}→ No existing deployment found${NC}"
    echo ""
    "$SCRIPT_DIR/first-deploy.sh"
fi

# Final status
echo ""
echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║        Deployment Status             ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""

# Set default service name if not set
SERVICE_NAME=${SERVICE_NAME:-crm-maulanabuilds}

ssh "$SERVER_USER@$SERVER_HOST" "systemctl is-active $SERVICE_NAME &>/dev/null && echo -e '${GREEN}✓ Service is running${NC}' || echo -e '${RED}✗ Service is not running${NC}'"
ssh "$SERVER_USER@$SERVER_HOST" "systemctl is-enabled $SERVICE_NAME &>/dev/null && echo -e '${GREEN}✓ Service enabled (auto-start on boot)${NC}' || echo -e '${YELLOW}! Service not enabled${NC}'"

echo ""
echo -e "${CYAN}Useful commands:${NC}"
echo "  View logs:     ssh $SERVER_USER@$SERVER_HOST 'journalctl -u $SERVICE_NAME -f'"
echo "  Check status:  ssh $SERVER_USER@$SERVER_HOST 'systemctl status $SERVICE_NAME'"
echo "  Restart:       ssh $SERVER_USER@$SERVER_HOST 'systemctl restart $SERVICE_NAME'"
echo ""
echo -e "${GREEN}Deployment complete!${NC}"
