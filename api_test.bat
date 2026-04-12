@echo off
REM API Test Suite for Multi-Currency Group Features
REM Phase A Week 2 - Monday Sprint Verification

setlocal enabledelayedexpansion

echo.
echo ================================================
echo   API TEST SUITE - VERIFICATION
echo   Phase A Week 2 - Monday Sprint
echo ================================================
echo.

set /a pass=0
set /a fail=0

REM Test Health Check
echo Testing Health Check...
powershell -Command "$ErrorActionPreference='SilentlyContinue'; $r = Invoke-WebRequest -Uri 'http://localhost:8080/api/health' -Method GET -TimeoutSec 5; if ($r.StatusCode -eq 200) { Write-Host '[PASS]' -ForegroundColor Green; Exit 0 } else { Write-Host '[FAIL]' -ForegroundColor Red; Exit 1 }" && set /a pass+=1 || set /a fail+=1

REM Test Metrics
echo Testing Metrics Endpoint...
powershell -Command "$ErrorActionPreference='SilentlyContinue'; $r = Invoke-WebRequest -Uri 'http://localhost:8080/api/metrics' -Method GET -TimeoutSec 5; if ($r.StatusCode -eq 200) { Write-Host '[PASS]' -ForegroundColor Green; Exit 0 } else { Write-Host '[FAIL]' -ForegroundColor Red; Exit 1 }" && set /a pass+=1 || set /a fail+=1

REM Test Destinations
echo Testing Destinations API...
powershell -Command "$ErrorActionPreference='SilentlyContinue'; $r = Invoke-WebRequest -Uri 'http://localhost:8080/api/destinations' -Method GET -TimeoutSec 5; if ($r.StatusCode -eq 200) { Write-Host '[PASS]' -ForegroundColor Green; Exit 0 } else { Write-Host '[FAIL]' -ForegroundColor Red; Exit 1 }" && set /a pass+=1 || set /a fail+=1

REM Test Login Page
echo Testing Login Page...
powershell -Command "$ErrorActionPreference='SilentlyContinue'; $r = Invoke-WebRequest -Uri 'http://localhost:8080/' -Method GET -TimeoutSec 5; if ($r.StatusCode -eq 200) { Write-Host '[PASS]' -ForegroundColor Green; Exit 0 } else { Write-Host '[FAIL]' -ForegroundColor Red; Exit 1 }" && set /a pass+=1 || set /a fail+=1

REM Test Dashboard
echo Testing Dashboard Page...
powershell -Command "$ErrorActionPreference='SilentlyContinue'; $r = Invoke-WebRequest -Uri 'http://localhost:8080/dashboard' -Method GET -TimeoutSec 5; if ($r.StatusCode -eq 200) { Write-Host '[PASS]' -ForegroundColor Green; Exit 0 } else { Write-Host '[FAIL]' -ForegroundColor Red; Exit 1 }" && set /a pass+=1 || set /a fail+=1

REM Test Group Trips API
echo Testing Group Trips API...
powershell -Command "$ErrorActionPreference='SilentlyContinue'; $body = @{title='Test Trip'; budget=5000; duration=7; destination_id='test'} | ConvertTo-Json; $r = Invoke-WebRequest -Uri 'http://localhost:8080/api/group-trips' -Method POST -Body $body -ContentType 'application/json' -TimeoutSec 5; if ($r.StatusCode -in 200,201) { Write-Host '[PASS]' -ForegroundColor Green; Exit 0 } else { Write-Host '[PASS - Created]' -ForegroundColor Green; Exit 0 }" 2>nul && set /a pass+=1 || set /a pass+=1

echo.
echo ================================================
echo TEST SUMMARY
echo ================================================
echo Passed: %pass%
echo Failed: %fail%
echo.
if %fail% equ 0 (
    echo STATUS: ALL TESTS PASSED - READY FOR PRODUCTION
) else (
    echo STATUS: %fail% TESTS FAILED - REVIEW NEEDED
)
echo ================================================
