#!/usr/bin/env python3
import requests
import json
from datetime import datetime

BASE_URL = "https://itinerary-backend-ikpw.onrender.com"
LOG_FILE = "RENDER_ENDPOINT_TEST_LOG.txt"

def test_endpoints():
    # Initialize log
    with open(LOG_FILE, 'w') as f:
        f.write("="*70 + "\n")
        f.write("RENDER DEPLOYMENT ENDPOINT TEST REPORT\n")
        f.write(f"Date: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n")
        f.write(f"Base URL: {BASE_URL}\n")
        f.write("="*70 + "\n\n")
    
    endpoints = [
        ("/api/health", "Liveness probe - checks if service is running"),
        ("/health", "Alternative health endpoint"),
        ("/api/ready", "Readiness probe - checks if ready for traffic"),
        ("/ready", "Alternative readiness endpoint"),
        ("/api/status", "Detailed status with diagnostics"),
        ("/status", "Alternative status endpoint"),
        ("/api/metrics", "Prometheus-compatible metrics"),
        ("/metrics", "Alternative metrics endpoint"),
    ]
    
    print("\n🚀 Testing Render Deployment Endpoints\n")
    print("="*70)
    
    results = []
    for endpoint, description in endpoints:
        url = BASE_URL + endpoint
        print(f"\n📍 Testing: {endpoint}")
        print(f"📝 {description}")
        print(f"⏱️  Time: {datetime.now().strftime('%H:%M:%S')}")
        print("-" * 70)
        
        try:
            response = requests.get(url, timeout=10)
            
            log_entry = f"\n{'─'*70}\n"
            log_entry += f"📍 ENDPOINT: GET {endpoint}\n"
            log_entry += f"📝 Description: {description}\n"
            log_entry += f"⏱️  Time: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n"
            log_entry += f"{'─'*70}\n"
            log_entry += f"Request URL: {url}\n"
            log_entry += f"HTTP Status Code: {response.status_code}\n"
            log_entry += f"Content-Type: {response.headers.get('Content-Type', 'N/A')}\n"
            log_entry += f"Response Time: {response.elapsed.total_seconds():.3f}s\n"
            log_entry += f"\n📤 Request Headers:\n"
            
            for key, value in response.request.headers.items():
                log_entry += f"  {key}: {value}\n"
            
            log_entry += f"\n📥 Response Headers:\n"
            for key, value in response.headers.items():
                log_entry += f"  {key}: {value}\n"
            
            log_entry += f"\n📦 Response Body:\n"
            
            try:
                json_data = response.json()
                log_entry += json.dumps(json_data, indent=2) + "\n"
                print(f"✅ Status: {response.status_code}")
                print(f"⏱️  Response Time: {response.elapsed.total_seconds():.3f}s")
                print(f"📊 Response: {json.dumps(json_data, indent=2)[:150]}...")
            except:
                log_entry += response.text + "\n"
                print(f"✅ Status: {response.status_code}")
                print(f"⏱️  Response Time: {response.elapsed.total_seconds():.3f}s")
                print(f"📊 Response: {response.text[:100]}...")
            
            log_entry += f"\n✅ Test Status: PASSED\n"
            results.append((endpoint, True, response.status_code))
            
        except requests.exceptions.Timeout:
            log_entry = f"\n{'─'*70}\n"
            log_entry += f"📍 ENDPOINT: GET {endpoint}\n"
            log_entry += f"⏱️  Time: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n"
            log_entry += f"❌ Error: Request Timeout\n"
            print(f"❌ Error: Request Timeout")
            results.append((endpoint, False, "Timeout"))
            
        except Exception as e:
            log_entry = f"\n{'─'*70}\n"
            log_entry += f"📍 ENDPOINT: GET {endpoint}\n"
            log_entry += f"⏱️  Time: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n"
            log_entry += f"❌ Error: {str(e)}\n"
            print(f"❌ Error: {str(e)}")
            results.append((endpoint, False, str(e)))
        
        with open(LOG_FILE, 'a') as f:
            f.write(log_entry + "\n")
    
    print("\n" + "="*70)
    print("📊 TEST SUMMARY")
    print("="*70)
    
    summary = f"\n{'='*70}\n"
    summary += f"TEST SUMMARY - {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n"
    summary += f"{'='*70}\n\n"
    
    passed = sum(1 for _, success, _ in results if success)
    total = len(results)
    
    summary += f"Total Endpoints Tested: {total}\n"
    summary += f"Passed: {passed} ✅\n"
    summary += f"Failed: {total - passed} ❌\n\n"
    
    summary += "DETAILED RESULTS:\n"
    summary += "-" * 70 + "\n"
    for endpoint, success, status in results:
        status_icon = "✅" if success else "❌"
        summary += f"{status_icon} {endpoint:20} - Status: {status}\n"
    
    summary += "\n" + "="*70 + "\n"
    
    print(f"✅ Passed: {passed}/{total}")
    print(f"❌ Failed: {total - passed}/{total}")
    
    with open(LOG_FILE, 'a') as f:
        f.write(summary)
    
    print(f"\n📋 Full logs saved to: {LOG_FILE}")
    return passed == total

if __name__ == "__main__":
    test_endpoints()
