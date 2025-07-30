#!/bin/bash

# Comprehensive Integration Test Runner for TensorZero Go Client
# This script runs all integration tests with proper reporting and cleanup

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Function to print colored output
print_header() {
    echo -e "${PURPLE}================================${NC}"
    echo -e "${PURPLE}$1${NC}"
    echo -e "${PURPLE}================================${NC}"
}

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

print_test() {
    echo -e "${CYAN}[TEST]${NC} $1"
}

# Configuration
VERBOSE=${VERBOSE:-false}
BENCHMARK=${BENCHMARK:-false}
CLEANUP=${CLEANUP:-true}
TIMEOUT=${TIMEOUT:-300}
PARALLEL=${PARALLEL:-1}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -b|--benchmark)
            BENCHMARK=true
            shift
            ;;
        --no-cleanup)
            CLEANUP=false
            shift
            ;;
        -t|--timeout)
            TIMEOUT="$2"
            shift 2
            ;;
        -p|--parallel)
            PARALLEL="$2"
            shift 2
            ;;
        -h|--help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  -v, --verbose     Enable verbose output"
            echo "  -b, --benchmark   Run benchmark tests"
            echo "  --no-cleanup      Skip cleanup after tests"
            echo "  -t, --timeout     Test timeout in seconds (default: 300)"
            echo "  -p, --parallel    Number of parallel test processes (default: 1)"
            echo "  -h, --help        Show this help message"
            echo ""
            echo "Environment variables:"
            echo "  VERBOSE=true      Enable verbose output"
            echo "  BENCHMARK=true    Run benchmark tests"
            echo "  CLEANUP=false     Skip cleanup"
            echo "  TIMEOUT=300       Test timeout in seconds"
            echo "  PARALLEL=1        Number of parallel processes"
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            exit 1
            ;;
    esac
done

print_header "TensorZero Go Client Integration Tests"

# Check if we're in the right directory
if [ ! -f "go.mod" ] || [ ! -d "tests" ]; then
    print_error "Please run this script from the root of the tensorzero-go repository"
    exit 1
fi

# Setup test environment
print_status "Setting up test environment..."
if [ -f "scripts/setup-integration-tests.sh" ]; then
    if ! bash scripts/setup-integration-tests.sh; then
        print_error "Failed to setup integration test environment"
        exit 1
    fi
else
    print_warning "Setup script not found, assuming environment is ready"
fi

# Create test results directory
RESULTS_DIR="test-results/$(date +%Y%m%d_%H%M%S)"
mkdir -p "$RESULTS_DIR"

print_status "Test results will be saved to: $RESULTS_DIR"

# Function to run tests with proper logging
run_test_suite() {
    local test_name="$1"
    local test_pattern="$2"
    local extra_args="$3"
    
    print_test "Running $test_name..."
    
    local cmd_args="-tags=integration ./tests"
    
    if [ "$VERBOSE" = true ]; then
        cmd_args="$cmd_args -v"
    fi
    
    if [ -n "$test_pattern" ]; then
        cmd_args="$cmd_args -run $test_pattern"
    fi
    
    if [ -n "$extra_args" ]; then
        cmd_args="$cmd_args $extra_args"
    fi
    
    cmd_args="$cmd_args -timeout ${TIMEOUT}s"
    cmd_args="$cmd_args -parallel $PARALLEL"
    
    local output_file="$RESULTS_DIR/${test_name}.log"
    local json_file="$RESULTS_DIR/${test_name}.json"
    
    print_status "Command: go test $cmd_args"
    
    # Run the test and capture output
    if go test $cmd_args -json > "$json_file" 2>&1; then
        print_success "$test_name passed"
        
        # Extract summary from JSON output
        if command -v jq &> /dev/null; then
            local passed=$(jq -r 'select(.Action=="pass" and .Test==null) | .Package' "$json_file" | wc -l)
            local failed=$(jq -r 'select(.Action=="fail" and .Test==null) | .Package' "$json_file" | wc -l)
            print_status "$test_name summary: $passed passed, $failed failed"
        fi
        
        return 0
    else
        print_error "$test_name failed"
        
        # Show recent output for debugging
        if [ "$VERBOSE" = true ]; then
            echo "Recent output:"
            tail -20 "$json_file" | jq -r 'select(.Output) | .Output' 2>/dev/null || tail -20 "$json_file"
        fi
        
        return 1
    fi
}

# Function to run benchmarks
run_benchmarks() {
    print_test "Running benchmark tests..."
    
    local cmd_args="-tags=integration ./tests -bench=. -benchmem"
    
    if [ "$VERBOSE" = true ]; then
        cmd_args="$cmd_args -v"
    fi
    
    cmd_args="$cmd_args -timeout ${TIMEOUT}s"
    
    local output_file="$RESULTS_DIR/benchmarks.log"
    
    if go test $cmd_args > "$output_file" 2>&1; then
        print_success "Benchmarks completed"
        
        # Show benchmark summary
        echo "Benchmark results:"
        grep "Benchmark" "$output_file" | head -10
        
        return 0
    else
        print_error "Benchmarks failed"
        
        if [ "$VERBOSE" = true ]; then
            echo "Recent output:"
            tail -20 "$output_file"
        fi
        
        return 1
    fi
}

# Track test results
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Test suites to run
declare -a TEST_SUITES=(
    "basic:TestIntegration_BasicInference|TestIntegration_JSONInference|TestIntegration_StreamingInference:"
    "feedback:TestIntegration_Feedback:"
    "evaluation:TestIntegration_DynamicEvaluation:"
    "datapoints:TestIntegration_DatapointOperations:"
    "filtering:TestIntegration_ListInferences|TestIntegration_ComplexFiltering|TestIntegration_TimeFiltering:"
    "advanced:TestIntegration_MetricFiltering|TestIntegration_OrderingAndPagination|TestIntegration_NotFilter:"
    "reliability:TestIntegration_ConcurrentRequests|TestIntegration_LongRunningStream|TestIntegration_RetryLogic:"
    "cleanup:TestIntegration_ResourceCleanup|TestIntegration_MemoryUsage:"
    "errors:TestIntegration_ErrorHandling|TestIntegration_ContextCancellation:"
)

print_header "Running Integration Test Suites"

# Run each test suite
for suite in "${TEST_SUITES[@]}"; do
    IFS=':' read -r suite_name test_pattern extra_args <<< "$suite"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    if run_test_suite "$suite_name" "$test_pattern" "$extra_args"; then
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    
    echo "" # Add spacing between test suites
done

# Run benchmarks if requested
if [ "$BENCHMARK" = true ]; then
    print_header "Running Benchmark Tests"
    
    if run_benchmarks; then
        print_success "All benchmarks completed"
    else
        print_warning "Some benchmarks failed"
    fi
fi

# Generate test report
print_header "Test Results Summary"

echo "Test Suites: $TOTAL_TESTS"
echo "Passed: $PASSED_TESTS"
echo "Failed: $FAILED_TESTS"
echo "Success Rate: $(( PASSED_TESTS * 100 / TOTAL_TESTS ))%"
echo ""
echo "Results saved to: $RESULTS_DIR"

# Generate HTML report if possible
if command -v go-junit-report &> /dev/null; then
    print_status "Generating JUnit report..."
    
    for json_file in "$RESULTS_DIR"/*.json; do
        if [ -f "$json_file" ]; then
            local base_name=$(basename "$json_file" .json)
            go-junit-report < "$json_file" > "$RESULTS_DIR/${base_name}.xml" 2>/dev/null || true
        fi
    done
fi

# Cleanup if requested
if [ "$CLEANUP" = true ]; then
    print_status "Cleaning up test environment..."
    make docker-down > /dev/null 2>&1 || true
    print_success "Cleanup completed"
fi

# Final status
if [ $FAILED_TESTS -eq 0 ]; then
    print_success "🎉 All integration tests passed!"
    exit 0
else
    print_error "❌ $FAILED_TESTS test suite(s) failed"
    exit 1
fi