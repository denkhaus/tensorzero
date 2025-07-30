#!/bin/bash

# Setup script for TensorZero Go Client Integration Tests
# This script ensures the environment is properly configured for integration testing

set -e

echo "🚀 Setting up TensorZero Integration Tests..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if we're in the right directory
if [ ! -f "go.mod" ] || [ ! -d "docker" ]; then
    print_error "Please run this script from the root of the tensorzero-go repository"
    exit 1
fi

print_status "Checking prerequisites..."

# Check if Docker is installed and running
if ! command -v docker &> /dev/null; then
    print_error "Docker is not installed. Please install Docker first."
    exit 1
fi

if ! docker info &> /dev/null; then
    print_error "Docker is not running. Please start Docker first."
    exit 1
fi

print_success "Docker is available and running"

# Check if Docker Compose is available
if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
    print_error "Docker Compose is not available. Please install Docker Compose."
    exit 1
fi

print_success "Docker Compose is available"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed. Please install Go 1.21 or later."
    exit 1
fi

GO_VERSION=$(go version | grep -oE 'go[0-9]+\.[0-9]+' | cut -c3-)
if [ "$(printf '%s\n' "1.21" "$GO_VERSION" | sort -V | head -n1)" != "1.21" ]; then
    print_error "Go version $GO_VERSION is too old. Please install Go 1.21 or later."
    exit 1
fi

print_success "Go version $GO_VERSION is compatible"

# Check if .env file exists in docker directory
if [ ! -f "docker/.env" ]; then
    print_warning ".env file not found in docker/ directory"
    
    if [ -f "docker/example.env" ]; then
        print_status "Copying example.env to .env..."
        cp docker/example.env docker/.env
        print_success "Created docker/.env from example.env"
    else
        print_status "Creating docker/.env file..."
        cat > docker/.env << EOF
TENSORZERO_CLICKHOUSE_URL=http://chuser:chpassword@clickhouse:8123/tensorzero
TENSORZERO_GATEWAY_URL=http://gateway:3000
OPENROUTER_API_KEY=sk-or-v1-dummy-key-for-testing
EOF
        print_success "Created docker/.env with default values"
    fi
    
    print_warning "Please update OPENROUTER_API_KEY in docker/.env with a real API key for full functionality"
else
    print_success "docker/.env file exists"
fi

# Check OPENROUTER_API_KEY
if grep -q "dummy\|test\|example" docker/.env; then
    print_warning "OPENROUTER_API_KEY appears to be a dummy key. Some tests may fail without a real API key."
fi

# Download Go dependencies
print_status "Downloading Go dependencies..."
go mod download
go mod tidy
print_success "Go dependencies are up to date"

# Build the project to check for compilation errors
print_status "Building the project..."
if go build ./...; then
    print_success "Project builds successfully"
else
    print_error "Project build failed. Please fix compilation errors first."
    exit 1
fi

# Run unit tests to ensure they pass
print_status "Running unit tests..."
if make test-unit; then
    print_success "Unit tests pass"
else
    print_error "Unit tests failed. Please fix unit test issues first."
    exit 1
fi

# Check if TensorZero services are already running
print_status "Checking if TensorZero services are running..."
if curl -s http://localhost:3000/health > /dev/null 2>&1; then
    print_success "TensorZero gateway is already running"
    SERVICES_RUNNING=true
else
    print_status "TensorZero services are not running"
    SERVICES_RUNNING=false
fi

# Start TensorZero services if not running
if [ "$SERVICES_RUNNING" = false ]; then
    print_status "Starting TensorZero services..."
    cd docker
    
    # Stop any existing services first
    docker-compose down > /dev/null 2>&1 || true
    
    # Start services
    if docker-compose up -d; then
        print_success "TensorZero services started"
    else
        print_error "Failed to start TensorZero services"
        exit 1
    fi
    
    cd ..
    
    # Wait for services to be healthy
    print_status "Waiting for services to be healthy..."
    TIMEOUT=120
    ELAPSED=0
    
    while [ $ELAPSED -lt $TIMEOUT ]; do
        if curl -s http://localhost:3000/health > /dev/null 2>&1; then
            print_success "TensorZero gateway is healthy"
            break
        fi
        
        echo -n "."
        sleep 2
        ELAPSED=$((ELAPSED + 2))
    done
    
    if [ $ELAPSED -ge $TIMEOUT ]; then
        print_error "Timeout waiting for TensorZero services to become healthy"
        print_status "Checking service logs..."
        cd docker && docker-compose logs
        exit 1
    fi
fi

# Test basic connectivity
print_status "Testing basic connectivity..."
if curl -s http://localhost:3000/health | grep -q "ok\|healthy\|ready"; then
    print_success "TensorZero gateway is responding correctly"
else
    print_warning "TensorZero gateway health check returned unexpected response"
fi

# Run a quick integration test
print_status "Running a quick integration test..."
if timeout 30 go test -tags=integration ./tests -run TestIntegration_BasicInference -v; then
    print_success "Basic integration test passed"
else
    print_warning "Basic integration test failed or timed out. This might be due to API key issues."
fi

echo ""
print_success "🎉 Integration test environment is ready!"
echo ""
echo "Available commands:"
echo "  make test-integration  - Run all integration tests"
echo "  make test-all         - Run unit + integration tests"
echo "  make docker-logs      - View TensorZero service logs"
echo "  make docker-down      - Stop TensorZero services"
echo ""
echo "TensorZero services:"
echo "  Gateway: http://localhost:3000"
echo "  UI:      http://localhost:4000"
echo ""
echo "To run integration tests:"
echo "  go test -tags=integration ./tests -v"
echo ""

if grep -q "dummy\|test\|example" docker/.env; then
    print_warning "Remember to update OPENROUTER_API_KEY in docker/.env for full test functionality"
fi