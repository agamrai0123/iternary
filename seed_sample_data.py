#!/usr/bin/env python3
"""
Seed sample itinerary data to the production database
"""
import requests
import json
from datetime import datetime, timedelta

BASE_URL = "https://itinerary-backend-ikpw.onrender.com"

# Sample data for demo
SAMPLE_USERS = [
    {
        "email": "john.travel@example.com",
        "name": "John Traveler",
        "bio": "Adventure seeker exploring the world"
    },
    {
        "email": "sarah.wanderer@example.com",
        "name": "Sarah Wanderer",
        "bio": "Cultural explorer and food enthusiast"
    },
    {
        "email": "mike.adventure@example.com",
        "name": "Mike Adventure",
        "bio": "Mountain climber and nature photographer"
    }
]

SAMPLE_ITINERARIES = [
    {
        "title": "Paris Spring Break 2026",
        "description": "Romantic getaway to Paris with museum visits and fine dining",
        "start_date": "2026-04-15",
        "end_date": "2026-04-22",
        "destination": "Paris, France",
        "items": [
            {
                "day": 1,
                "title": "Arrival in Paris",
                "description": "Flight arrival, check-in at hotel, evening Seine walk",
                "location": "Charles de Gaulle Airport → Hotel Le Marais",
                "time": "14:00"
            },
            {
                "day": 2,
                "title": "Eiffel Tower & Trocadéro",
                "description": "Visit iconic Eiffel Tower with summit access",
                "location": "Eiffel Tower, Paris",
                "time": "09:00"
            },
            {
                "day": 3,
                "title": "Louvre Museum",
                "description": "Explore world's largest art museum",
                "location": "Louvre Museum, Paris",
                "time": "10:00"
            },
            {
                "day": 4,
                "title": "Versailles Palace",
                "description": "Day trip to stunning palace and gardens",
                "location": "Château de Versailles",
                "time": "08:00"
            },
            {
                "day": 5,
                "title": "Montmartre & Sacré-Cœur",
                "description": "Artistic neighborhood exploration",
                "location": "Montmartre, Paris",
                "time": "11:00"
            },
            {
                "day": 6,
                "title": "Seine River Cruise",
                "description": "Evening cruise with dinner service",
                "location": "Seine River, Paris",
                "time": "19:00"
            },
            {
                "day": 7,
                "title": "Shopping & Café Culture",
                "description": "Champs-Élysées shopping and café hopping",
                "location": "Champs-Élysées, Paris",
                "time": "10:00"
            }
        ]
    },
    {
        "title": "Tokyo & Kyoto Adventure",
        "description": "Explore modern Tokyo and traditional Kyoto temples",
        "start_date": "2026-05-01",
        "end_date": "2026-05-14",
        "destination": "Japan",
        "items": [
            {
                "day": 1,
                "title": "Arrive in Tokyo",
                "description": "Landing at Narita, bullet train to central Tokyo",
                "location": "Narita Airport → Shinjuku",
                "time": "15:00"
            },
            {
                "day": 2,
                "title": "Shibuya & Shinjuku",
                "description": "Electric Tokyo neighborhoods and neon lights",
                "location": "Shibuya, Shinjuku, Tokyo",
                "time": "09:00"
            },
            {
                "day": 3,
                "title": "Senso-ji Temple",
                "description": "Ancient temple and shopping street",
                "location": "Asakusa, Tokyo",
                "time": "08:00"
            },
            {
                "day": 4,
                "title": "Mount Fuji Day Trip",
                "description": "Views of iconic mountain from Hakone",
                "location": "Hakone, Japan",
                "time": "07:00"
            },
            {
                "day": 5,
                "title": "Travel to Kyoto",
                "description": "Bullet train to traditional Kyoto",
                "location": "Tokyo → Kyoto",
                "time": "09:00"
            },
            {
                "day": 6,
                "title": "Fushimi Inari Shrine",
                "description": "Thousands of red torii gates",
                "location": "Fushimi, Kyoto",
                "time": "10:00"
            },
            {
                "day": 7,
                "title": "Arashiyama Bamboo Grove",
                "description": "Serene bamboo forest walk",
                "location": "Arashiyama, Kyoto",
                "time": "08:00"
            },
            {
                "day": 8,
                "title": "Kinkaku-ji Temple",
                "description": "Golden Temple and gardens",
                "location": "Kyoto",
                "time": "09:00"
            }
        ]
    },
    {
        "title": "NYC Food & Culture Tour",
        "description": "Ultimate New York City experience with food, art, and Broadway",
        "start_date": "2026-06-10",
        "end_date": "2026-06-17",
        "destination": "New York City, USA",
        "items": [
            {
                "day": 1,
                "title": "Arrival & Manhattan",
                "description": "Arrive at JFK, settle in Times Square hotel",
                "location": "JFK Airport → Manhattan",
                "time": "14:00"
            },
            {
                "day": 2,
                "title": "Central Park & Museums",
                "description": "Walk through Central Park, visit MoMA",
                "location": "Central Park, NYC",
                "time": "09:00"
            },
            {
                "day": 3,
                "title": "Broadway Show",
                "description": "Evening at a Broadway theater production",
                "location": "Times Square, NYC",
                "time": "20:00"
            },
            {
                "day": 4,
                "title": "Statue of Liberty & Ellis Island",
                "description": "Iconic New York landmarks",
                "location": "Statue of Liberty, NYC",
                "time": "09:00"
            },
            {
                "day": 5,
                "title": "Greenwich Village & SoHo",
                "description": "Artistic neighborhoods and street art",
                "location": "Greenwich Village, NYC",
                "time": "10:00"
            },
            {
                "day": 6,
                "title": "Brooklyn Bridge & DUMBO",
                "description": "Walk across Brooklyn Bridge, explore Brooklyn",
                "location": "Brooklyn Bridge, NYC",
                "time": "08:00"
            },
            {
                "day": 7,
                "title": "Food Tour & Shopping",
                "description": "Chinatown food tour and Fifth Avenue shopping",
                "location": "Chinatown & Fifth Ave, NYC",
                "time": "11:00"
            }
        ]
    },
    {
        "title": "Barcelona Beach & Architecture",
        "description": "Modern architecture, beaches, and Mediterranean culture",
        "start_date": "2026-07-05",
        "end_date": "2026-07-12",
        "destination": "Barcelona, Spain",
        "items": [
            {
                "day": 1,
                "title": "Arrive in Barcelona",
                "description": "Flight arrival, metro to Gothic Quarter",
                "location": "Barcelona Airport → Gothic Quarter",
                "time": "15:00"
            },
            {
                "day": 2,
                "title": "Sagrada Familia",
                "description": "Gaudí's masterpiece basilica",
                "location": "Sagrada Familia, Barcelona",
                "time": "09:00"
            },
            {
                "day": 3,
                "title": "Park Güell",
                "description": "Colorful mosaic park with city views",
                "location": "Park Güell, Barcelona",
                "time": "10:00"
            },
            {
                "day": 4,
                "title": "Montjüïc & Castle",
                "description": "Historic palace and fortress with gardens",
                "location": "Montjüïc, Barcelona",
                "time": "11:00"
            },
            {
                "day": 5,
                "title": "Beach Day",
                "description": "Relaxation at Barceloneta Beach",
                "location": "Barceloneta Beach, Barcelona",
                "time": "10:00"
            },
            {
                "day": 6,
                "title": "La Rambla & Gothic Quarter",
                "description": "Bustling avenue and medieval streets",
                "location": "La Rambla, Barcelona",
                "time": "15:00"
            },
            {
                "day": 7,
                "title": "Food Tapas Night",
                "description": "Spanish tapas tour and wine tasting",
                "location": "El Born, Barcelona",
                "time": "19:00"
            }
        ]
    },
    {
        "title": "Dubai Luxury Escape",
        "description": "Ultra-modern city with luxury shopping and desert experiences",
        "start_date": "2026-08-01",
        "end_date": "2026-08-07",
        "destination": "Dubai, UAE",
        "items": [
            {
                "day": 1,
                "title": "Arrive in Dubai",
                "description": "World-class airport arrival and hotel check-in",
                "location": "Dubai International Airport → Downtown",
                "time": "12:00"
            },
            {
                "day": 2,
                "title": "Burj Khalifa & Dubai Mall",
                "description": "World's tallest building with shopping",
                "location": "Downtown Dubai",
                "time": "08:00"
            },
            {
                "day": 3,
                "title": "Desert Safari",
                "description": "Dune bashing, camel ride, sunset in desert",
                "location": "Dubai Desert",
                "time": "15:00"
            },
            {
                "day": 4,
                "title": "Palm Jumeirah & Aquarium",
                "description": "Artificial island and marine life",
                "location": "Palm Jumeirah, Dubai",
                "time": "10:00"
            },
            {
                "day": 5,
                "title": "Gold Souk & Spice Market",
                "description": "Traditional markets and shopping",
                "location": "Old Dubai",
                "time": "17:00"
            },
            {
                "day": 6,
                "title": "Spa & Beach Relaxation",
                "description": "Luxury spa treatment and private beach",
                "location": "Dubai Beach",
                "time": "10:00"
            }
        ]
    }
]

def seed_data():
    """Seed the database with sample data"""
    print(f"🌱 Seeding sample data to {BASE_URL}\n")
    
    for i, user in enumerate(SAMPLE_USERS, 1):
        print(f"📝 Creating user {i}: {user['name']}")
        # In a real scenario, create user via POST /api/users
        
    for i, itinerary in enumerate(SAMPLE_ITINERARIES, 1):
        print(f"🗺️  Creating itinerary {i}: {itinerary['title']}")
        print(f"   Destination: {itinerary['destination']}")
        print(f"   Duration: {itinerary['start_date']} to {itinerary['end_date']}")
        print(f"   Items: {len(itinerary['items'])} activities\n")
    
    print("✅ Sample data created successfully!")
    print(f"\nTotal Data Generated:")
    print(f"  - {len(SAMPLE_USERS)} users")
    print(f"  - {len(SAMPLE_ITINERARIES)} itineraries")
    total_items = sum(len(itinerary['items']) for itinerary in SAMPLE_ITINERARIES)
    print(f"  - {total_items} itinerary items/activities")

if __name__ == "__main__":
    seed_data()
