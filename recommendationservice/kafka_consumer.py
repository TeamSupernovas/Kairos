import os
import json
import psycopg2
from dotenv import load_dotenv
from confluent_kafka import Consumer, KafkaException

# Load .env variables
load_dotenv()

# PostgreSQL connection
def get_db_conn():
    return psycopg2.connect(
        dbname=os.getenv("DB_NAME"),
        user=os.getenv("DB_USER"),
        password=os.getenv("DB_PASSWORD"),
        host=os.getenv("DB_HOST"),
        port=os.getenv("DB_PORT", 5432)
    )

# Insert user interaction (order or rate)
def insert_interaction(user_id, dish_id, action, rating=None):
    conn = get_db_conn()
    cur = conn.cursor()

    if action == "rate" and rating is not None:
        cur.execute(
            """
            INSERT INTO user_dish_interactions (user_id, dish_id, action, rating)
            VALUES (%s, %s, %s, %s)
            """,
            (user_id, dish_id, action, rating)
        )
    elif action == "order":
        cur.execute(
            """
            INSERT INTO user_dish_interactions (user_id, dish_id, action)
            VALUES (%s, %s, %s)
            """,
            (user_id, dish_id, action)
        )
    else:
        print(f"Skipping invalid interaction: {user_id}, {dish_id}, {action}")
        cur.close()
        conn.close()
        return

    conn.commit()
    cur.close()
    conn.close()
    print(f"Inserted {action} for user {user_id}, dish {dish_id}")

# Store dish metadata for content filtering
def store_dish(dish_id, name, features):
    conn = get_db_conn()
    cur = conn.cursor()
    cur.execute("""
        INSERT INTO dishes (dish_id, name, features)
        VALUES (%s, %s, %s)
        ON CONFLICT (dish_id) DO UPDATE
        SET name = EXCLUDED.name,
            features = EXCLUDED.features
    """, (dish_id, name, features))
    conn.commit()
    cur.close()
    conn.close()
    print(f"Stored dish {dish_id}: {name}")

# Main Kafka consumer loop
def start_consumer():
    consumer = Consumer({
        "bootstrap.servers": os.getenv("KAFKA_BOOTSTRAP_SERVERS"),
        "security.protocol": "SASL_SSL",
        "sasl.mechanisms": "PLAIN",
        "sasl.username": os.getenv("KAFKA_USERNAME"),
        "sasl.password": os.getenv("KAFKA_PASSWORD"),
        "group.id": "recommendation-consumer",
        "auto.offset.reset": "earliest"
    })

    topics = [
        "order-service.order.completed",
        "rating-service.rating.created",
        "dish-management-service.dish.created"
    ]

    consumer.subscribe(topics)
    print(f"Subscribed to topics: {topics}")

    try:
        while True:
            msg = consumer.poll(timeout=1.0)
            if msg is None:
                continue
            if msg.error():
                raise KafkaException(msg.error())

            try:
                topic = msg.topic()
                payload = json.loads(msg.value().decode("utf-8"))

                if topic == "order-service.order.completed":
                    insert_interaction(
                        user_id=payload["user_id"],
                        dish_id=payload["dish_id"],
                        action="order"
                    )

                elif topic == "rating-service.rating.created":
                    insert_interaction(
                        user_id=payload["user_id"],
                        dish_id=payload["dish_id"],
                        action="rate",
                        rating=int(payload["rating"])
                    )

                elif topic == "dish-management-service.dish.created":
                    dish_id = payload["dish_id"]
                    name = payload.get("name", "")
                    cuisine = payload.get("cuisine", "")
                    ingredients = payload.get("ingredients", [])
                    tags = payload.get("tags", [])
                    dietary = payload.get("dietary", "")

                    # Build features string
                    features = " ".join([cuisine] + ingredients + tags + [dietary])
                    store_dish(dish_id, name, features)

            except Exception as e:
                print(f"Failed to process message: {e}")
                print(f"Message content: {msg.value().decode('utf-8')}")

    except KeyboardInterrupt:
        print("Kafka consumer interrupted.")
    finally:
        consumer.close()

if __name__ == "__main__":
    start_consumer()
