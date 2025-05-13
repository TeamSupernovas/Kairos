import os
import pandas as pd
import psycopg2
from dotenv import load_dotenv

load_dotenv()

def load_interactions():
    conn = psycopg2.connect(
        dbname=os.getenv("DB_NAME"),
        user=os.getenv("DB_USER"),
        password=os.getenv("DB_PASSWORD"),
        host=os.getenv("DB_HOST"),
        port=os.getenv("DB_PORT", 5432)
    )

    query = """
        SELECT user_id, dish_id, action, rating, timestamp
        FROM user_dish_interactions
        ORDER BY user_id, dish_id, 
                 CASE action 
                     WHEN 'rate' THEN 1 
                     WHEN 'order' THEN 2 
                 END
    """

    df = pd.read_sql(query, conn)
    conn.close()

    def resolve_rating(group):
        for _, row in group.iterrows():
            if row["action"] == "rate" and row["rating"]:
                return row["rating"]
            if row["action"] == "order":
                return 4
        return 3

    final_df = df.groupby(["user_id", "dish_id"]).apply(resolve_rating).reset_index()
    final_df.columns = ["user_id", "dish_id", "rating"]

    return final_df
