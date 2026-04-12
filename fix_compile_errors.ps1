# Script to fix all compilation errors in Go files

$goDir = "D:\Learn\iternary\itinerary-backend\itinerary"

# Get all non-test Go files
Get-ChildItem $goDir -Filter "*.go" -Exclude "*_test.go" | ForEach-Object {
    $file = $_.FullName
    Write-Host "Processing: $($_.Name)"
    
    # Read file content
    $content = Get-Content $file -Raw
    
    # Fix 1: Replace nil with empty string in NewAPIError calls  
    $content = $content -replace 'NewAPIError\(([^,]+),\s*([^,]+),\s*nil\)', 'NewAPIError($1, $2, "")'
    
    # Fix 2: Replace map[string]string error parameters with string
    $content = $content -replace 'map\[string\]string\{\"error\": err\.Error\(\)\}', 'err.Error()'
    
    # Fix 3: Replace getHTTPStatusCode with GetStatusCode (in case it appears)
    $content = $content -replace 'getHTTPStatusCode', 'GetStatusCode'
    
    # Write back to file
    Set-Content $file -Value $content -NoNewline
    Write-Host "  ✓ Fixed"
}

Write-Host ""
Write-Host "All files fixed. Building now..."

# Build the project
cd "D:\Learn\iternary\itinerary-backend"
$buildResult = go build -o itinerary-backend.exe . 2>&1

if ($LASTEXITCODE -eq 0 -or $buildResult -eq "") {
    Write-Host "✓ Build successful!"
    if (Test-Path "D:\Learn\iternary\itinerary-backend\itinerary-backend.exe") {
        $size = (Get-Item "D:\Learn\iternary\itinerary-backend\itinerary-backend.exe").Length / 1MB
        Write-Host "✓ Binary created: itinerary-backend.exe ($([math]::Round($size, 2))MB)"
    }
} else {
    Write-Host "✗ Build failed with errors:"
    $buildResult | Out-Host
}
