import uvicorn

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from starlette.staticfiles import StaticFiles
from utils.logger import setup_logger, LOGGING
from wisp import shutdown_protobuf

from api.conf import settings
from api.middlewares import RequestMiddleware
from api.routers import chat
from jms.poll import setup_poll_jms_event

app = FastAPI()

app.include_router(chat.router, prefix="/kael")

app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.http.cors_allow_origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)
app.mount("/ui", StaticFiles(directory="ui"), name="ui")
app.add_middleware(RequestMiddleware)


def startup_event():
    setup_logger()
    # setup_poll_jms_event()


def shutdown_event():
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
