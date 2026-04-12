#!/usr/bin/env python3
"""
Test script to verify the itinerary API endpoints are working correctly.
"""

import json
import urllib.request
import urllib.error
import sys

BASE_URL = "http://localhost:8080"

def test_endpoint(method, endpoint, data=None, expected_status=200):
    """Test a single API endpoint"""
    url = f"{BASE_URL}{endpoint}"
    
    try:
        if data:
            data_bytes = json.dumps(data).encode('utf-8')
            req = urllib.request.Request(url, data=data_bytes, method=method)
            req.add_header('Content-Type', 'application/json')
        else:
            req = urllib.request.Request(url, method=method)
        
        response = urllib.request.urlopen(req)
        status = response.status
        content = response.read().decode('utf-8')
        
        if status == expected_status:
            print(f"✓ {method} {endpoint} - Status {status}")
            try:
                resp_json = json.loads(content)
                print(f"  Response: {json.dumps(resp_json, indent=2)[:200]}...")
            except:
                print(f"  Response: {content[:100]}...")
            return True
        else:
            print(f"✗ {method} {endpoint} - Expected {expected_status}, got {status}")
            return False
            
    except urllib.error.HTTPError as e:
        print(f"✗ {method} {endpoint} - HTTP Error {e.code}")
        try:
            error_content = e.read().decode('utf-8')
            print(f"  Error: {error_content[:200]}")
        except:
            pass
        return False
    except Exception as e:
        print(f"✗ {method} {endpoint} - Error: {str(e)}")
        return False

def main():
    print("=" * 60)
    print("Testing Itinerary API Endpoints")
    print("=" * 60)
    
    tests = [
        ("GET", "/api/destinations", None),
        ("GET", "/api/destinations?page=1&pageSize=10", None),
        ("GET", "/api/destinations/dest-001", None),
        ("GET", "/api/itineraries/destination/dest-001", None),
        ("GET", "/api/itineraries/itin-001", None),
        ("GET", "/api/itineraries/itin-001/items", None),
        ("GET", "/api/itineraries/itin-001/comments", None),
        ("POST", "/api/itineraries/itin-001/like", {}, 200),
        ("POST", "/api/itineraries/itin-001/comments", 
         {"content": "Test comment", "rating": 4.5}, 201),
    ]
    
    passed = 0
    failed = 0
    
    for method, endpoint, data in tests:
        if test_endpoint(method, endpoint, data):
            passed += 1
        else:
            failed += 1
        print()
    
    print("=" * 60)
    print(f"Results: {passed} passed, {failed} failed")
    print("=" * 60)
    
    return 0 if failed == 0 else 1

if __name__ == "__main__":
    sys.exit(main())
