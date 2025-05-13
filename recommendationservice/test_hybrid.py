from recommender.hybrid_recommender import hybrid_recommend

nearby_dishes = [
    {"dish_id": 1, "name": "Chicken Biryani", "features": "indian spicy rice chicken"},
    {"dish_id": 2, "name": "Paneer Tikka", "features": "indian grilled paneer spicy"},
    {"dish_id": 3, "name": "Sushi", "features": "japanese rice fish seaweed"},
]

results = hybrid_recommend(user_id=1, nearby_dishes=nearby_dishes)
print(results)

