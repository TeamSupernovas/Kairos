import pickle
from typing import List, Dict
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.metrics.pairwise import cosine_similarity
import pandas as pd
from recommender.data_loader import load_interactions, get_dish_features, load_dishids_for_user
from dotenv import load_dotenv
import os

load_dotenv()

with open("model.pkl", "rb") as f:
    model = pickle.load(f)

def derive_features_from_dish(dish: dict) -> str:
    name = dish.get("DishName", "")
    course = dish.get("MealCourse", "")
    category = dish.get("DietaryCategory", "")
    description = dish.get("Description", "")
    return f"{name} {course} {category} {description}".lower()

def get_user_rated_dishes(user_id: str) -> List[str]:
    df = load_dishids_for_user(user_id)
    return  df["dish_id"].unique().tolist()

def compute_content_scores(user_id: int, nearby_dishes: List[Dict]) -> List[Dict]:
    if not nearby_dishes:
        return []

    nearby_df = pd.DataFrame(nearby_dishes)
    nearby_df["features"] = nearby_df.apply(derive_features_from_dish, axis=1)

    if nearby_df.empty or nearby_df["features"].isnull().all():
        return []

    rated_ids = get_user_rated_dishes(user_id)

    rated_df = get_dish_features(rated_ids)

    if rated_df.empty:
        return [
            {"dish_id": row["DishID"], "name": row["DishName"], "content_score": 0.0}
            for _, row in nearby_df.iterrows()
        ]

    combined = pd.concat([nearby_df[["features"]], rated_df[["features"]]])
    tfidf = TfidfVectorizer()
    tfidf_matrix = tfidf.fit_transform(combined["features"])

    nearby_matrix = tfidf_matrix[:len(nearby_df)]
    rated_matrix = tfidf_matrix[len(nearby_df):]

    cos_sim = cosine_similarity(nearby_matrix, rated_matrix)

    results = []
    for i, row in nearby_df.iterrows():
        sim_scores = cos_sim[i]
        content_score = sum(sim_scores) / len(sim_scores) if len(sim_scores) > 0 else 0.0

        results.append({
            "dish_id": row["DishID"],
            "name": row["DishName"],
            "content_score": round(content_score, 3)
        })

    return results

def hybrid_recommend(user_id: int, nearby_dishes: List[Dict], alpha: float = 0.7) -> List[Dict]:
    if not nearby_dishes:
        return []

    content_scores_list = compute_content_scores(user_id, nearby_dishes)
    content_scores_map = {
        item["dish_id"]: item["content_score"] for item in content_scores_list
    }

    results = []
    for dish in nearby_dishes:
        dish_id = dish["DishID"]
        name = dish["DishName"]
        try:
            collab_score = model.predict(user_id, dish_id).est

        except:
            collab_score = 3.0
        content_score = content_scores_map.get(dish_id, 0.0)

        final_score = alpha * collab_score + (1 - alpha) * content_score

        print(f"Dish: {name}, Collab: {collab_score:.2f}, Content: {content_score:.2f}, Hybrid: {final_score:.2f}")

        results.append({
            "dish_id": dish_id,
            "name": name,
            "score": round(final_score, 3)
        })
    return sorted(results, key=lambda x: x["score"], reverse=True)
