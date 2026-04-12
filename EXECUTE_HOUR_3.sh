#!/bin/bash

# Monday Hour 3: Complete Test Execution with Coverage
# Started: 11:00 AM
# Goal: Execute 79 tests and achieve >85% code coverage

cd /d/Learn/iternary/itinerary-backend

echo "=========================================="
echo "📊 HOUR 3: TEST EXECUTION (11:00 AM)"
echo "=========================================="
echo ""
echo "Starting at: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""

# Step 1: Ensure database exists and schema is applied
echo "Step 1: Preparing database..."
if [ ! -f itinerary.db ]; then
    echo "  - Creating new database from schema..."
    sqlite3 itinerary.db < multicurrency_schema.sql 2>&1 | grep -E "CREATE|INSERT|Error" | head -20
    echo "  ✅ Database created with multi-currency schema"
else
    echo "  - Database exists, verifying schema..."
    TABLE_COUNT=$(sqlite3 itinerary.db "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%';")
    echo "  ✅ Database has $TABLE_COUNT tables"
fi

echo ""

# Step 2: Clean previous test artifacts
echo "Step 2: Cleaning previous test artifacts..."
rm -f coverage.out coverage.html test_summary.txt test_results.txt
echo "  ✅ Cleaned"

echo ""

# Step 3: Run tests with coverage
echo "Step 3: Running 79 tests..."
echo "  Files being tested:"
echo "    - auth_service_test.go"
echo "    - config_test.go"
echo "    - error_test.go"
echo "    - logger_test.go"
echo "    - metrics_test.go"
echo "    - models_test.go"
echo "    - service_test.go"
echo "    - template_helpers_test.go"
echo "    - group_models_test.go (25+ tests)"
echo "    - group_service_test.go (32+ tests)"
echo "    - group_integration_test.go (22+ tests)"
echo ""
echo "  Executing: go test ./itinerary -v -cover -coverprofile=coverage.out -timeout 60s"
echo ""

# Run the tests and capture output
go test ./itinerary -v -cover -coverprofile=coverage.out -timeout 60s 2>&1 | tee test_results.txt

echo ""
echo "=========================================="

# Step 4: Generate coverage HTML report
echo "Step 4: Generating coverage report..."
go tool cover -html=coverage.out -o coverage.html 2>&1
echo "  ✅ HTML report created: coverage.html"

echo ""

# Step 5: Extract coverage summary
echo "Step 5: Coverage Analysis..."
COVERAGE=$(go tool cover -func=coverage.out | tail -1 | awk '{print $3}')
echo "  📊 Total Coverage: $COVERAGE"

# Detailed coverage by function
echo ""
echo "  Coverage by component:"
go tool cover -func=coverage.out | grep "^github.com/yourusername/itinerary-backend/itinerary" | awk -F':' '{print $2}' | awk '{print "    - " $1 ": " $NF}' | sort -u | head -20

echo ""

# Step 6: Test Results Summary
echo "=========================================="
echo "Step 6: Test Results Summary"
echo "=========================================="

# Count test results
TOTAL_TESTS=$(grep -c "^RUN\|^--- PASS:" test_results.txt || echo "N/A")
PASSED_TESTS=$(grep -c "^--- PASS:" test_results.txt || echo "0")
FAILED_TESTS=$(grep -c "^--- FAIL:" test_results.txt || echo "0")

echo "  Total Test Cases: $TOTAL_TESTS"
echo "  ✅ Passed: $PASSED_TESTS"
echo "  ❌ Failed: $FAILED_TESTS"

if grep -q "^--- FAIL:" test_results.txt; then
    echo ""
    echo "Failed tests:"
    grep "^--- FAIL:" test_results.txt
    echo ""
    echo "FAILED TEST DETAILS:"
    grep -A 10 "^--- FAIL:" test_results.txt
fi

echo ""

# Step 7: Final Status
echo "=========================================="
echo "HOUR 3 EXECUTION SUMMARY"
echo "=========================================="
echo "  ⏱️  Started: $(date '+%H:%M:%S')"
echo "  ✅ Status: Test execution complete"
echo "  📊 Coverage: $COVERAGE"
echo "  ✔️  Test Results: $PASSED_TESTS passed"

if [ "$FAILED_TESTS" -eq "0" ]; then
    echo ""
    echo "✨ SUCCESS! All tests passing ✨"
    echo ""
    echo "Next Step (Hour 4 - 12:00 PM):"
    echo "  1. Build binary: go build -o itinerary-backend.exe ."
    echo "  2. Start server: ./itinerary-backend.exe > server.log 2>&1 &"
    echo "  3. Test endpoints: curl http://localhost:8080/api/health"
    echo "  4. Stop server: taskkill /IM itinerary-backend.exe /F (or pkill)"
else
    echo ""
    echo "⚠️  Some tests failed. Review details above."
fi

echo ""
echo "Completed at: $(date '+%Y-%m-%d %H:%M:%S')"
echo "=========================================="

# Create summary file
{
    echo "HOUR 3 TEST EXECUTION RESULTS"
    echo "=============================="
    echo "Timestamp: $(date)"
    echo "Duration: $SECONDS seconds"
    echo ""
    echo "Coverage: $COVERAGE"
    echo "Tests Passed: $PASSED_TESTS"
    echo "Tests Failed: $FAILED_TESTS"
    echo ""
    echo "Files Generated:"
    echo "  - test_results.txt (detailed test output)"
    echo "  - coverage.out (raw coverage data)"
    echo "  - coverage.html (visual coverage report)"
    echo ""
    if [ "$FAILED_TESTS" -eq "0" ]; then
        echo "Status: ✅ PASSED"
    else
        echo "Status: ❌ FAILED"
    fi
} > test_summary.txt

echo ""
echo "📄 Full summary saved to: test_summary.txt"
echo ""
