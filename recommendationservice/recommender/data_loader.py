import os
import pandas as pd
import psycopg2
from dotenv import load_dotenv
from typing import List, Dict

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


def load_dishids_for_user(user_id):
    conn = psycopg2.connect(
        dbname=os.getenv("DB_NAME"),
        user=os.getenv("DB_USER"),
        password=os.getenv("DB_PASSWORD"),
        host=os.getenv("DB_HOST"),
        port=os.getenv("DB_PORT", 5432)
    )

    query = """
        SELECT distinct dish_id
        FROM user_dish_interactions
        WHERE user_id = %s
    """

    df = pd.read_sql(query, conn, params=(str(user_id),))
    conn.close()

    return df


def get_dish_features(dish_ids: List[str]) -> pd.DataFrame:
    if not dish_ids:
        return pd.DataFrame(columns=["dish_id", "features"])

    conn = psycopg2.connect(
        dbname=os.getenv("DB_NAME"),
        user=os.getenv("DB_USER"),
        password=os.getenv("DB_PASSWORD"),
        host=os.getenv("DB_HOST"),
        port=os.getenv("DB_PORT", 5432)
    )

    query = """
        SELECT
            dish_id,
            features
        FROM dishes
        WHERE dish_id = ANY(%s)
    """

    df = pd.read_sql(query, conn, params=(dish_ids,))
    conn.close()

    return df[["dish_id", "features"]]