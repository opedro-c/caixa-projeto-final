FROM python:3.14.3-slim

WORKDIR /app

RUN apt-get update && apt-get install -y \
    libpq-dev \
    gcc \
    && rm -rf /var/lib/apt/lists/*

COPY requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt

COPY to_postgres.py .
COPY movies.csv .
COPY users.csv .
COPY ratings.csv .

CMD ["python", "to_postgres.py"]