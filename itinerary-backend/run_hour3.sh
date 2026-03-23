#!/bin/bash
# Hour 3: Run Tests with Coverage
# Execute all tests and generate coverage reports

cd /d/Learn/iternary/itinerary-backend

echo "================================"
echo "HOUR 3: TEST EXECUTION"
echo "================================"
echo ""

# Step 1: Run tests with coverage
echo "Running tests with coverage..."
go test ./itinerary -v -cover -coverprofile=coverage.out -timeout 60s

# Step 2: Show test summary
echo ""
echo "================================"
echo "TEST SUMMARY"
echo "================================"
go tool cover -func=coverage.out | tail -10

# Step 3: Generate HTML report
echo ""
echo "Generating coverage HTML report..."
go tool cover -html=coverage.out -o coverage.html

# Step 4: Show coverage percentage
echo ""
coverage_percent=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
echo "✅ Total Coverage: $coverage_percent"

echo ""
echo "================================"
echo "✅ HOUR 3 COMPLETE"
echo "================================"
echo ""
echo "Coverage report: coverage.html"
echo "Test log: test_results.log"
