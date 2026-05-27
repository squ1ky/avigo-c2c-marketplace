from contextlib import asynccontextmanager
from typing import List, Optional

from fastapi import FastAPI, HTTPException
from pydantic import BaseModel, Field

from .model import TagModel

tag_model = TagModel()


class PredictTagsRequest(BaseModel):
    title: str = Field(..., min_length=3, max_length=200)
    category: Optional[str] = Field(default="", max_length=80)


class PredictTagsResponse(BaseModel):
    tags: List[str]
    model: str
    labels_count: int


@asynccontextmanager
async def lifespan(app: FastAPI):
    tag_model.load_or_train()
    yield


app = FastAPI(title="AviGo ML Service", version="0.1.0", lifespan=lifespan)


@app.get("/health")
def health():
    return {
        "status": "ok",
        "trained": tag_model.is_ready,
        "labels_count": tag_model.labels_count,
    }


@app.post("/train")
def train():
    tag_model.train(force=True)
    return {
        "status": "trained",
        "labels_count": tag_model.labels_count,
        "samples_count": tag_model.samples_count,
    }


@app.post("/predict-tags", response_model=PredictTagsResponse)
def predict_tags(payload: PredictTagsRequest):
    if not tag_model.is_ready:
        raise HTTPException(status_code=503, detail="model is not trained")

    tags = tag_model.predict(payload.title, payload.category or "")
    return PredictTagsResponse(
        tags=tags,
        model="tfidf-one-vs-rest-logreg",
        labels_count=tag_model.labels_count,
    )
