$file = "D:\Learn\iternary\itinerary-backend\itinerary\group_database.go"
$content = Get-Content $file -Raw

# Fix method calls
$content = $content -replace 'db\.exec\(', 'db.conn.Exec('
$content = $content -replace 'db\.query\(', 'db.conn.QueryRow('
$content = $content -replace 'db\.queryRows\(', 'db.conn.Query('

# Fix NewAPIError calls with map[string]string
$content = $content -replace ', map\[string\]string\{"error": err\.Error\(\)\}\)', ', err.Error())'

$content | Set-Content $file
Write-Host "Fixed all issues in group_database.go"
