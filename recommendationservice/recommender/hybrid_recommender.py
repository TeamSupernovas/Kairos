import pickle
from typing import List, Dict
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.metrics.pairwise import cosine_similarity
import pandas as pd
from recommender.data_loader import load_interactions
from dotenv import load_dotenv
import os

load_dotenv()

# Load collaborative model
with open("model.pkl", "rb") as f:
    model = pickle.load(f)

# Derive feature string from dish object
def derive_features_from_dish(dish: dict) -> str:
    name = dish.get("DishName", "")
    course = dish.get("MealCourse", "")
    category = dish.get("DietaryCategory", "")
    description = dish.get("Description", "")
    return f"{name} {course} {category} {description}".lower()

def get_user_rated_dishes(user_id: int) -> List[str]:
    df = load_interactions()
    return df[df["user_id"] == user_id]["dish_id"].unique().tolist()

def hybrid_recommend(user_id: int, nearby_dishes: List[Dict], alpha: float = 0.7) -> List[Dict]:
    if not nearby_dishes:
        return []

    # Derive features dynamically from dish metadata
    for dish in nearby_dishes:
        dish["features"] = derive_features_from_dish(dish)

    dish_df = pd.DataFrame(nearby_dishes)
    dish_df = dish_df.dropna(subset=["features"])

    if dish_df.empty:
        return []

    user_rated = get_user_rated_dishes(user_id)
    print(f"User {user_id} rated dishes: {user_rated}")

    # TF-IDF similarity
    tfidf = TfidfVectorizer()
    tfidf_matrix = tfidf.fit_transform(dish_df["features"])
    cos_sim = cosine_similarity(tfidf_matrix)

    results = []
    for idx, row in dish_df.iterrows():
        dish_id = row["DishID"]
        name = row["DishName"]

        # 1. Collaborative Score
        try:
            collab_score = model.predict(user_id, dish_id).est
        except:
            collab_score = 3.0

        # 2. Content Score
        sim_scores = []
        for rated_id in user_rated:
            if rated_id in dish_df["DishID"].values:
                i = dish_df[dish_df["DishID"] == dish_id].index[0]
                j = dish_df[dish_df["DishID"] == rated_id].index[0]
                sim_scores.append(cos_sim[i][j])
        content_score = sum(sim_scores) / len(sim_scores) if sim_scores else 0.0

        # 3. Hybrid Score
        final_score = alpha * collab_score + (1 - alpha) * content_score

        print(f"Dish: {name}, Collab: {collab_score:.2f}, Content: {content_score:.2f}, Hybrid: {final_score:.2f}")

        results.append({
            "dish_id": dish_id,
            "name": name,
            "score": round(final_score, 3)
        })

    return sorted(results, key=lambda x: x["score"], reverse=True)
