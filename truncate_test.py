file_path = r"D:\Learn\iternary\itinerary-backend\itinerary\group_integration_test.go"

with open(file_path, 'r', encoding='utf-8') as f:
    lines = f.readlines()

# Find where the mockDatabase Close method ends
close_func_index = -1
for i, line in enumerate(lines):
    if "func (m *mockDatabase) Close() error" in line:
        close_func_index = i
        break

# Find the closing brace of the Close function (should be within next few lines)
if close_func_index >= 0:
    for i in range(close_func_index, min(close_func_index + 5, len(lines))):
        if line.strip() == "}":
            close_func_index = i + 1
            break

# Keep only lines up to and including the Close function
if close_func_index > 0:
    with open(file_path, 'w', encoding='utf-8') as f:
        f.writelines(lines[:close_func_index])
    print(f"Truncated file to {close_func_index} lines (removed {len(lines) - close_func_index} corrupted lines)")
else:
    print("Could not find Close function")
