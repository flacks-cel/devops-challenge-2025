from fastapi import FastAPI
from datetime import datetime, timezone
from prometheus_client import Counter, Histogram, make_asgi_app
import time

app = FastAPI()

REQUEST_COUNT = Counter(
    "release_requests_total",
    "Total de requests no release-service",
    ["method", "endpoint", "status"]
)

REQUEST_LATENCY = Histogram(
    "release_request_duration_seconds",
    "Latência das requisições no release-service",
    ["endpoint"]
)

metrics_app = make_asgi_app()
app.mount("/metrics", metrics_app)

@app.get("/")
def release_info():
    start = time.time()
    REQUEST_COUNT.labels(method="GET", endpoint="/", status="200").inc()
    REQUEST_LATENCY.labels(endpoint="/").observe(time.time() - start)
    return {
        "service": "release-service",
        "version": "1.4.2",
        "environment": "production",
        "status": "healthy"
    }

@app.get("/time")
def last_deployed():
    start = time.time()
    REQUEST_COUNT.labels(method="GET", endpoint="/time", status="200").inc()
    REQUEST_LATENCY.labels(endpoint="/time").observe(time.time() - start)
    return {
        "last_deployed_at": datetime.now(timezone.utc).isoformat()
    }