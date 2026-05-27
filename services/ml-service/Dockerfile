FROM python:3.12-slim

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY app ./app
COPY data ./data

RUN mkdir -p /app/models

ENV MODEL_PATH=/app/models/tag_model.joblib
ENV TRAINING_DATA_PATH=/app/data/tag_training.csv

EXPOSE 8091

CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8091"]
