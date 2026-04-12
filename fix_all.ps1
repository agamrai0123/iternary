$goDir = "D:\Learn\iternary\itinerary-backend\itinerary"
Get-ChildItem $goDir -Filter "*.go" -Exclude "*_test.go" | ForEach-Object {
    $file = $_.FullName
    Write-Host "Processing: $($_.Name)"
    $content = Get-Content $file -Raw
    $content = $content -replace 'NewAPIError\(([^,]+),\s*([^,]+),\s*nil\)', 'NewAPIError($1, $2, "")'
    $content = $content -replace 'map\[string\]string\{\"error\": err\.Error\(\)\}', 'err.Error()'
    $content = $content -replace 'getHTTPStatusCode', 'GetStatusCode'
    Set-Content $file -Value $content -NoNewline
}
Write-Host "Building..."
cd "D:\Learn\iternary\itinerary-backend"
go build -o itinerary-backend.exe . 2>&1
