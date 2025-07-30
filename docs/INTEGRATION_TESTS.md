# TensorZero Go Client - Integration Tests

This document provides comprehensive information about the integration test suite for the TensorZero Go client.

## 🎯 **Overview**

The integration test suite validates the TensorZero Go client against a real TensorZero instance, ensuring:
- **API Compatibility** - All endpoints work correctly
- **Performance** - Response times and throughput meet expectations  
- **Reliability** - Concurrent usage and error handling work properly
- **Feature Completeness** - All client features function as expected

## 📁 **Test Structure**

### **Test Files**
```
tests/
├── integration_test.go           # Core integration tests
├── integration_advanced_test.go  # Advanced filtering and operations
├── integration_benchmarks_test.go # Performance benchmarks
├── integration_reliability_test.go # Reliability and stress tests
└── openai_test.go               # OpenAI compatibility tests
```

### **Test Categories**

#### **1. Core Functionality Tests**
- **Basic Inference** - Simple chat completions
- **JSON Inference** - Structured output generation
- **Streaming** - Real-time response streaming
- **Feedback** - Rating and metric submission
- **Dynamic Evaluation** - A/B testing workflows

#### **2. Data Management Tests**
- **Datapoint Operations** - CRUD operations for training data
- **Bulk Operations** - Efficient batch processing
- **Dataset Management** - Multi-dataset workflows

#### **3. Advanced Filtering Tests**
- **List Inferences** - Query historical inferences
- **Complex Filtering** - AND/OR/NOT filter combinations
- **Time-based Filtering** - Date range queries
- **Metric Filtering** - Performance-based queries
- **Pagination** - Large result set handling

#### **4. Reliability Tests**
- **Concurrent Requests** - Multi-threaded usage
- **Long-running Streams** - Extended streaming sessions
- **Error Handling** - Graceful failure management
- **Resource Cleanup** - Memory and connection management
- **Context Cancellation** - Proper timeout handling

#### **5. Performance Benchmarks**
- **Inference Latency** - Response time measurements
- **Concurrent Throughput** - Multi-user performance
- **Streaming Performance** - Real-time data flow
- **Memory Usage** - Resource consumption analysis

## 🚀 **Running Integration Tests**

### **Quick Start**
```bash
# Setup environment and run all tests
./scripts/setup-integration-tests.sh
make test-integration
```

### **Advanced Usage**
```bash
# Run specific test categories
go test -tags=integration ./tests -run TestIntegration_BasicInference -v

# Run with custom timeout
go test -tags=integration ./tests -timeout 60s

# Run benchmarks
go test -tags=integration ./tests -bench=. -benchmem

# Run comprehensive test suite
./scripts/run-integration-tests.sh --verbose --benchmark
```

### **Test Runner Options**
```bash
./scripts/run-integration-tests.sh [OPTIONS]

Options:
  -v, --verbose     Enable detailed output
  -b, --benchmark   Include performance benchmarks
  --no-cleanup      Keep services running after tests
  -t, --timeout     Test timeout in seconds (default: 300)
  -p, --parallel    Parallel test processes (default: 1)
```

## ⚙️ **Environment Setup**

### **Prerequisites**
- **Go 1.21+** - Latest Go version
- **Docker & Docker Compose** - Container runtime
- **OpenRouter API Key** - For model access (optional for dry-run tests)

### **Configuration**
```bash
# Required environment file
docker/.env:
TENSORZERO_CLICKHOUSE_URL=http://chuser:chpassword@clickhouse:8123/tensorzero
TENSORZERO_GATEWAY_URL=http://gateway:3000
OPENROUTER_API_KEY=your_openrouter_api_key_here
```

### **Automatic Setup**
```bash
# Run the setup script for automated configuration
./scripts/setup-integration-tests.sh
```

## 📊 **Test Results & Reporting**

### **Output Formats**
- **Console Output** - Real-time test progress
- **JSON Reports** - Machine-readable results
- **JUnit XML** - CI/CD integration format
- **Benchmark Reports** - Performance metrics

### **Result Location**
```
test-results/
└── YYYYMMDD_HHMMSS/
    ├── basic.json          # Core functionality results
    ├── advanced.json       # Advanced feature results
    ├── reliability.json    # Stress test results
    ├── benchmarks.log      # Performance metrics
    └── *.xml              # JUnit reports (if available)
```

### **Success Criteria**
- **Unit Tests**: 100% pass rate required
- **Integration Tests**: 80%+ pass rate (allows for API key issues)
- **Concurrent Tests**: 80%+ success rate under load
- **Performance**: Response times within acceptable ranges

## 🔧 **Troubleshooting**

### **Common Issues**

#### **1. Services Not Starting**
```bash
# Check Docker status
docker ps
make docker-logs

# Restart services
make docker-down
make docker-up
```

#### **2. API Key Issues**
```bash
# Verify API key in environment
grep OPENROUTER_API_KEY docker/.env

# Test with dry-run mode
go test -tags=integration ./tests -run TestIntegration_BasicInference
```

#### **3. Network Connectivity**
```bash
# Test gateway health
curl http://localhost:3000/health

# Check service status
make docker-status
```

#### **4. Test Timeouts**
```bash
# Increase timeout for slow networks
go test -tags=integration ./tests -timeout 120s

# Run tests sequentially
go test -tags=integration ./tests -parallel 1
```

### **Debug Mode**
```bash
# Enable verbose logging
VERBOSE=true ./scripts/run-integration-tests.sh

# Run single test with full output
go test -tags=integration ./tests -run TestIntegration_BasicInference -v
```

## 📈 **Performance Expectations**

### **Latency Targets**
- **Basic Inference**: < 2s (with real API)
- **Streaming Start**: < 1s for first chunk
- **Feedback Submission**: < 500ms
- **List Operations**: < 1s for 100 items

### **Throughput Targets**
- **Concurrent Requests**: 10+ simultaneous users
- **Streaming**: Multiple concurrent streams
- **Bulk Operations**: 100+ datapoints/batch

### **Resource Usage**
- **Memory**: < 100MB for client
- **Connections**: Proper cleanup after tests
- **Goroutines**: No leaks after completion

## 🔄 **CI/CD Integration**

### **GitHub Actions Example**
```yaml
name: Integration Tests
on: [push, pull_request]

jobs:
  integration:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: '1.21'
      
      - name: Setup Integration Tests
        run: ./scripts/setup-integration-tests.sh
        env:
          OPENROUTER_API_KEY: ${{ secrets.OPENROUTER_API_KEY }}
      
      - name: Run Integration Tests
        run: ./scripts/run-integration-tests.sh --verbose
      
      - name: Upload Test Results
        uses: actions/upload-artifact@v3
        if: always()
        with:
          name: test-results
          path: test-results/
```

### **Makefile Integration**
```bash
# Available make targets
make test-unit          # Unit tests only
make test-integration   # Integration tests only  
make test-all          # Unit + Integration tests
make ci-test           # CI-optimized test run
```

## 📝 **Contributing**

### **Adding New Tests**
1. **Create test function** with `TestIntegration_` prefix
2. **Use build tag** `//go:build integration`
3. **Include proper setup** with `setupClient(t)`
4. **Add cleanup** with `defer client.Close()`
5. **Use appropriate timeouts** and error handling

### **Test Categories**
- **Basic Tests**: Core functionality validation
- **Advanced Tests**: Complex feature combinations
- **Reliability Tests**: Stress and error scenarios
- **Benchmark Tests**: Performance measurements

### **Best Practices**
- **Isolated Tests**: Each test should be independent
- **Proper Cleanup**: Always close clients and clean resources
- **Error Handling**: Distinguish between expected and unexpected errors
- **Timeouts**: Use appropriate timeouts for different operations
- **Logging**: Include helpful debug information

## 🎯 **Test Coverage Goals**

### **API Coverage**
- ✅ **Inference API** - All request types and options
- ✅ **Streaming API** - Real-time response handling  
- ✅ **Feedback API** - All metric types and targets
- ✅ **Evaluation API** - Dynamic A/B testing workflows
- ✅ **Datapoint API** - CRUD operations and bulk processing
- ✅ **List API** - Filtering, sorting, and pagination

### **Feature Coverage**
- ✅ **Error Handling** - All error types and scenarios
- ✅ **Authentication** - API key validation
- ✅ **Timeouts** - Context cancellation and deadlines
- ✅ **Concurrency** - Multi-threaded usage patterns
- ✅ **Resource Management** - Connection pooling and cleanup

### **Quality Metrics**
- **Test Coverage**: 80%+ of integration scenarios
- **Performance**: All benchmarks within target ranges
- **Reliability**: 95%+ success rate in normal conditions
- **Documentation**: Complete test documentation and examples

---

**The integration test suite ensures the TensorZero Go client is production-ready and fully compatible with the TensorZero platform.** 🚀