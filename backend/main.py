import uvicorn
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from utils.logger import setup_logger, get_logger, LOGGING
from api.conf import settings
from api.middlewares import AuthMiddleware
from api.routers import chat

setup_logger()

logger = get_logger(__name__)

app = FastAPI()

app.include_router(chat.router)

# 解决跨站问题
app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.http.cors_allow_origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# TODO 认证
app.add_middleware(AuthMiddleware, my_option='test')


@app.on_event("startup")
async def on_startup():
    print(f"On startup... http://{settings.http.host}/{settings.http.port}")


@app.on_event("shutdown")
async def on_shutdown():
    print("On shutdown...")


if __name__ == "__main__":
    uvicorn.run(
        app,
        host=settings.http.host,
        port=settings.http.port,
        proxy_headers=True,
        forwarded_allow_ips='*',
        log_config=LOGGING,
    )
