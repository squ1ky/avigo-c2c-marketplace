import csv
import os
import re
from dataclasses import dataclass
from pathlib import Path
from typing import Iterable, List, Sequence

import joblib
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.linear_model import LogisticRegression
from sklearn.multiclass import OneVsRestClassifier
from sklearn.preprocessing import MultiLabelBinarizer
from sklearn.pipeline import Pipeline


DEFAULT_DATA_PATH = "/app/data/tag_training.csv"
DEFAULT_MODEL_PATH = "/app/models/tag_model.joblib"


@dataclass
class TrainingSet:
    texts: List[str]
    labels: List[List[str]]


class TagModel:
    def __init__(self) -> None:
        self.data_path = Path(os.getenv("TRAINING_DATA_PATH", DEFAULT_DATA_PATH))
        self.model_path = Path(os.getenv("MODEL_PATH", DEFAULT_MODEL_PATH))
        self.pipeline: Pipeline | None = None
        self.binarizer: MultiLabelBinarizer | None = None
        self.samples_count = 0

    @property
    def is_ready(self) -> bool:
        return self.pipeline is not None and self.binarizer is not None

    @property
    def labels_count(self) -> int:
        if self.binarizer is None:
            return 0
        return len(self.binarizer.classes_)

    def load_or_train(self) -> None:
        if self.model_path.exists():
            payload = joblib.load(self.model_path)
            self.pipeline = payload["pipeline"]
            self.binarizer = payload["binarizer"]
            self.samples_count = payload.get("samples_count", 0)
            return

        self.train(force=True)

    def train(self, force: bool = False) -> None:
        if not force and self.model_path.exists():
            self.load_or_train()
            return

        training_set = self._load_training_set()
        if not training_set.texts:
            raise ValueError("training dataset is empty")

        binarizer = MultiLabelBinarizer()
        y = binarizer.fit_transform(training_set.labels)

        pipeline = Pipeline(
            steps=[
                (
                    "tfidf",
                    TfidfVectorizer(
                        analyzer="char_wb",
                        ngram_range=(2, 5),
                        lowercase=True,
                        min_df=1,
                    ),
                ),
                (
                    "classifier",
                    OneVsRestClassifier(
                        LogisticRegression(
                            max_iter=1000,
                            class_weight="balanced",
                        )
                    ),
                ),
            ]
        )

        pipeline.fit(training_set.texts, y)

        self.pipeline = pipeline
        self.binarizer = binarizer
        self.samples_count = len(training_set.texts)

        self.model_path.parent.mkdir(parents=True, exist_ok=True)
        joblib.dump(
            {
                "pipeline": self.pipeline,
                "binarizer": self.binarizer,
                "samples_count": self.samples_count,
            },
            self.model_path,
        )

    def predict(self, title: str, category: str = "", limit: int = 8) -> List[str]:
        if not self.is_ready:
            raise RuntimeError("model is not trained")

        text = self._feature_text(title, category)
        probabilities = self.pipeline.predict_proba([text])[0]
        labels = list(self.binarizer.classes_)
        ranked = sorted(zip(labels, probabilities), key=lambda item: item[1], reverse=True)

        tags = [label for label, score in ranked if score >= 0.55]
        if len(tags) < 4:
            tags = [label for label, _ in ranked[:6]]

        tags.extend(extract_title_tags(title))
        return normalize_tags(tags, limit=limit)

    def _load_training_set(self) -> TrainingSet:
        texts: List[str] = []
        labels: List[List[str]] = []

        with self.data_path.open("r", encoding="utf-8", newline="") as file:
            reader = csv.DictReader(file)
            for row in reader:
                title = row["title"].strip()
                category = row["category"].strip()
                row_tags = split_tags(row["tags"])
                if title and row_tags:
                    texts.append(self._feature_text(title, category))
                    labels.append(row_tags)

        return TrainingSet(texts=texts, labels=labels)

    @staticmethod
    def _feature_text(title: str, category: str) -> str:
        return f"{category} {title}".strip().lower()


def split_tags(value: str) -> List[str]:
    return normalize_tags(value.split(";"), limit=40)


def normalize_tags(tags: Iterable[str], limit: int) -> List[str]:
    blocked = {
        "новый",
        "новое",
        "новая",
        "бу",
        "б/у",
        "отличное состояние",
    }
    result: List[str] = []
    seen = set()

    for tag in tags:
        normalized = " ".join(str(tag).strip().lower().split())
        if not normalized or len(normalized) > 32 or normalized in blocked:
            continue
        if normalized in seen:
            continue

        seen.add(normalized)
        result.append(normalized)
        if len(result) >= limit:
            break

    return result


def extract_title_tags(title: str) -> Sequence[str]:
    lower = title.lower()
    tags: List[str] = []

    for amount, unit in re.findall(r"(\d+)\s*(gb|гб|tb|тб)", lower):
        unit = "gb" if unit in {"gb", "гб"} else "tb"
        tags.append(f"{amount}{unit}")
        if unit == "gb":
            tags.append(f"{amount}гб")

    iphone_match = re.search(r"(?:iphone|айфон)\s*(\d{1,2})(?:\s*(pro max|pro|plus|max))?", lower)
    if iphone_match:
        model = iphone_match.group(1)
        suffix = iphone_match.group(2) or ""
        model_tag = f"iphone {model} {suffix}".strip()
        tags.extend(["iphone", "айфон", model_tag])

    if "macbook" in lower or "макбук" in lower:
        tags.extend(["macbook", "макбук", "ноутбук"])
        if "air" in lower:
            tags.append("macbook air")
        if "pro" in lower:
            tags.append("macbook pro")

    return tags
