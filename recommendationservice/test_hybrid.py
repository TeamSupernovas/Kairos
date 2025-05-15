import unittest
from recommender.hybrid_recommender import hybrid_recommend

class TestHybridRecommender(unittest.TestCase):

    def setUp(self):
        self.user_id = "117899928711855949682"
        self.nearby_dishes = [
            {
                "DishID": "0f04a781-37cc-44ac-916e-e462557aef84",
                "DishName": "Sushi",
                "MealCourse": "Main Course",
                "DietaryCategory": "Non-Vegetarian",
                "Description": "A delicate Japanese dish featuring seasoned rice paired with fresh seafood, vegetables, or egg, artfully rolled or layered and served with soy sauce, wasabi, and pickled ginger."
            },
            {
                "DishID": "e81495ef-bced-4d1e-828a-19a388b7ee45",
                "DishName": "Chicken Biryani",
                "MealCourse": "Main Course",
                "DietaryCategory": "Non-Vegetarian",
                "Description": "A fragrant and flavorful rice dish layered with tender marinated chicken, aromatic basmati rice, and bold spices, slow-cooked to perfection and served with raita and salan."
            },
            {
                "DishID": "5391f6da-e830-4536-955a-3f277ea2eaab",
                "DishName": "Dosa",
                "MealCourse": "Starter",
                "DietaryCategory": "Vegetarian",
                "Description": "A dosa batter that turns crisp, chewy, and light on the griddle, with a sourdough-like tang thanks to a double fermentation."
            }
        ]

    def test_hybrid_recommend_returns_scores(self):
        results = hybrid_recommend(user_id=self.user_id, nearby_dishes=self.nearby_dishes)
        
        self.assertIsInstance(results, list)
        self.assertGreater(len(results), 0)

        for item in results:
            self.assertIn("dish_id", item)
            self.assertIn("name", item)
            self.assertIn("score", item)
            print(item)  # Optional: to see output during test run

if __name__ == "__main__":
    unittest.main()
