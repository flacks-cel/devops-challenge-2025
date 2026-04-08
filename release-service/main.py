from fastapi import FastAPI
from datetime import datetime, timezone

app = FastAPI()

@app.get("/")
def release_info():
    return {
        "service": "release-service",
        "version": "1.4.2",
        "environment": "production",
        "status": "healthy"
    }

@app.get("/time")
def last_deployed():
    return {
        "last_deployed_at": datetime.now(timezone.utc).isoformat()
    }