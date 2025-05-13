import pickle
import os
import boto3
from surprise import Dataset, Reader, SVD
from recommender.data_loader import load_interactions
from dotenv import load_dotenv

load_dotenv()

def train_and_save_model(output_path="model.pkl"):
    print("Loading interactions from database...")
    df = load_interactions()

    reader = Reader(rating_scale=(1, 5))
    data = Dataset.load_from_df(df[["user_id", "dish_id", "rating"]], reader)

    print("Training collaborative filtering model (SVD)...")
    trainset = data.build_full_trainset()
    model = SVD()
    model.fit(trainset)

    with open(output_path, "wb") as f:
        pickle.dump(model, f)

    print(f"Model trained and saved to {output_path}")

    return output_path

def upload_model_to_s3(file_path="model.pkl"):
    bucket = os.getenv("S3_BUCKET_NAME")
    key = os.getenv("S3_MODEL_KEY", "models/model.pkl")

    s3 = boto3.client(
        "s3",
        aws_access_key_id=os.getenv("AWS_ACCESS_KEY_ID"),
        aws_secret_access_key=os.getenv("AWS_SECRET_ACCESS_KEY"),
        region_name=os.getenv("AWS_REGION")
    )

    print(f"Uploading {file_path} to s3://{bucket}/{key}...")
    s3.upload_file(file_path, bucket, key)
    print(f"Upload complete: s3://{bucket}/{key}")

if __name__ == "__main__":
    model_path = train_and_save_model()
    upload_model_to_s3(model_path)
