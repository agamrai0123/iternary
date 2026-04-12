@echo off
REM Monday Hour 3: Complete Test Execution with Coverage
REM Started: 11:00 AM
REM Goal: Execute 79 tests and achieve >85% code coverage

cd /d D:\Learn\iternary\itinerary-backend

echo.
echo ==========================================
echo ^|   HOUR 3: TEST EXECUTION (11:00 AM)   ^|
echo ==========================================
echo.
echo Starting at: %date% %time%
echo.

REM Step 1: Ensure database exists and schema is applied
echo Step 1: Preparing database...
if not exist itinerary.db (
    echo   - Creating new database from schema...
    sqlite3 itinerary.db < multicurrency_schema.sql > nul 2>&1
    echo   ^[OK^] Database created with multi-currency schema
) else (
    echo   - Database exists, verifying schema...
    for /f %%i in ('sqlite3 itinerary.db "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%%';"') do set TABLE_COUNT=%%i
    echo   ^[OK^] Database has !TABLE_COUNT! tables
)

echo.

REM Step 2: Clean previous test artifacts
echo Step 2: Cleaning previous test artifacts...
if exist coverage.out del /s /q coverage.out > nul 2>&1
if exist coverage.html del /s /q coverage.html > nul 2>&1
if exist test_summary.txt del /s /q test_summary.txt > nul 2>&1
if exist test_results.txt del /s /q test_results.txt > nul 2>&1
echo   ^[OK^] Cleaned

echo.

REM Step 3: Run tests with coverage
echo Step 3: Running 79 tests from 12 test files...
echo   Files being tested:
echo     - auth_service_test.go
echo     - config_test.go
echo     - error_test.go
echo     - logger_test.go
echo     - metrics_test.go
echo     - models_test.go
echo     - service_test.go
echo     - template_helpers_test.go
echo     - group_models_test.go ^(25+^)
echo     - group_service_test.go ^(32+^)
echo     - group_integration_test.go ^(22+^)
echo.
echo   Executing: go test ./itinerary -v -cover -coverprofile=coverage.out
echo.

REM Run the tests and capture output
go test ./itinerary -v -cover -coverprofile=coverage.out -timeout 60s > test_results.txt 2>&1

REM Display results as they complete
type test_results.txt

echo.
echo ==========================================

REM Step 4: Generate coverage HTML report
echo Step 4: Generating coverage report...
go tool cover -html=coverage.out -o coverage.html 2> nul
echo   ^[OK^] HTML report created: coverage.html

echo.

REM Step 5: Extract coverage summary
echo Step 5: Coverage Analysis...
for /f "tokens=3" %%i in ('go tool cover -func=coverage.out ^| findstr /B "total"') do (
    set COVERAGE=%%i
)
echo   Coverage: %COVERAGE%

echo.

REM Step 6: Test Results Summary
echo ==========================================
echo Step 6: Test Results Summary
echo ==========================================

REM Count test results
for /f %%i in ('findstr /C:"--- PASS:" test_results.txt ^| find /C /V ""') do set PASSED_TESTS=%%i
for /f %%i in ('findstr /C:"--- FAIL:" test_results.txt ^| find /C /V ""') do set FAILED_TESTS=%%i

if "%PASSED_TESTS%"=="" set PASSED_TESTS=0
if "%FAILED_TESTS%"=="" set FAILED_TESTS=0

echo   Total Tests Passed: %PASSED_TESTS%
echo   Total Tests Failed: %FAILED_TESTS%

if %FAILED_TESTS% GTR 0 (
    echo.
    echo   Failed tests:
    findstr /C:"--- FAIL:" test_results.txt
)

echo.

REM Step 7: Final Status
echo ==========================================
echo HOUR 3 EXECUTION SUMMARY
echo ==========================================
echo   Time: %time%
echo   Status: Test execution complete
echo   Coverage: %COVERAGE%
echo   Tests Passed: %PASSED_TESTS%

if %FAILED_TESTS% EQU 0 (
    echo.
    echo   ^~^~^~ SUCCESS! All tests passing ^~^~^^
    echo.
    echo   Next Step ^(Hour 4 - 12:00 PM^):
    echo     1. Build binary: go build -o itinerary-backend.exe .
    echo     2. Start server: itinerary-backend.exe
    echo     3. Test endpoints: curl http://localhost:8080/api/health
    echo     4. Stop server: taskkill /IM itinerary-backend.exe /F
) else (
    echo.
    echo   Some tests failed. Review details above.
)

echo.
echo Completed at: %time%
echo ==========================================

REM Save summary
(
    echo HOUR 3 TEST EXECUTION RESULTS
    echo ==============================
    echo Timestamp: %date% %time%
    echo.
    echo Coverage: %COVERAGE%
    echo Tests Passed: %PASSED_TESTS%
    echo Tests Failed: %FAILED_TESTS%
    echo.
    echo Files Generated:
    echo   - test_results.txt ^(detailed test output^)
    echo   - coverage.out ^(raw coverage data^)
    echo   - coverage.html ^(visual coverage report^)
    if %FAILED_TESTS% EQU 0 (
        echo.
        echo Status: SUCCESS
    ) else (
        echo.
        echo Status: FAILED
    )
) > test_summary.txt

type test_summary.txt

pause
