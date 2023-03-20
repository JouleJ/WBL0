FROM python:3.11

COPY . /app
RUN ["pip3", "install", "nats-py"]

ENTRYPOINT ["python3", "/app/main.py"]
