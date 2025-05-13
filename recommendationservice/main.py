from fastapi import FastAPI
from pydantic import BaseModel
from typing import List, Optional
from recommender.hybrid_recommender import hybrid_recommend

app = FastAPI()

class DishInput(BaseModel):
    DishID: str
    DishName: str
    Description: Optional[str] = ""
    MealCourse: Optional[str] = ""
    DietaryCategory: Optional[str] = ""

class RecommendationRequest(BaseModel):
    user_id: int
    nearby_dishes: List[DishInput]

@app.post("/recommendations")
def get_recommendations(req: RecommendationRequest):
    dishes = [dish.model_dump() for dish in req.nearby_dishes]

    results = hybrid_recommend(
        user_id=req.user_id,
        nearby_dishes=dishes
    )

    return {
        "user_id": req.user_id,
        "recommendations": results
    }
