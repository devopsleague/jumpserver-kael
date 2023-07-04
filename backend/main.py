import uvicorn
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from utils.logger import setup_logger, LOGGING
from wisp import setup_protobuf, shutdown_protobuf
from api.conf import settings
from api.middlewares import AuthMiddleware
from api.routers import chat

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


def startup_event():
    setup_logger()
    setup_protobuf()
    print(f"On startup... http://{settings.http.host}:{settings.http.port}")


def shutdown_event():
    # 在程序关闭时执行的操作，例如释放资源、关闭连接等
    print("应用程序关闭，执行清理操作")
    shutdown_protobuf()


@app.on_event("startup")
async def on_startup():
    startup_event()


@app.on_event("shutdown")
async def on_shutdown():
    shutdown_event()


if __name__ == "__main__":
    uvicorn.run(
        app,
        host=settings.http.host,
        port=settings.http.port,
        proxy_headers=True,
        forwarded_allow_ips='*',
        log_config=LOGGING,
    )
