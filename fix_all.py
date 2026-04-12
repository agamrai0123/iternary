import os
import re
import glob

go_dir = r"D:\Learn\iternary\itinerary-backend\itinerary"

# Get all non-test Go files
for filepath in glob.glob(os.path.join(go_dir, "*.go")):
    if "_test.go" in filepath:
        continue
    
    print(f"Processing: {os.path.basename(filepath)}")
    
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Fix 1: Replace nil with empty string before closing paren
    content = re.sub(r',\s*nil\)', r', "")', content)
    
    # Fix 2: Replace map[string]string error parameters
    content = re.sub(r'map\[string\]string\{"error":\s*err\.Error\(\)\}', r'err.Error()', content)
    
    # Fix 3: Replace getHTTPStatusCode with GetStatusCode
    content = re.sub(r'getHTTPStatusCode', r'GetStatusCode', content)
    
    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print("  Fixed")

print("\nAll files fixed. Now building...")
os.system('cd D:\\Learn\\iternary\\itinerary-backend && go build -o itinerary-backend.exe . 2>&1')
