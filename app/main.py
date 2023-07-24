import uvicorn
from fastapi import FastAPI
from fastapi.responses import RedirectResponse
from fastapi.middleware.cors import CORSMiddleware
from starlette.staticfiles import StaticFiles
from utils.logger import setup_logger, LOGGING, get_logger
from wisp import shutdown_protobuf
from api.conf import settings
from api.middlewares import I18nMiddleware
from api.routers import chat, health, handlers
from api.jms.poll import setup_poll_jms_event

logger = get_logger(__name__)

app = FastAPI()

app.include_router(chat.router, prefix="/kael")
app.include_router(health.router, prefix="/kael")
app.include_router(handlers.router, prefix="/kael")

app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.http.cors_allow_origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)
app.add_middleware(I18nMiddleware)


@app.get("/kael/connect")
async def connect(token: str):
    redirect_url = f"/kael/?token={token}"
    return RedirectResponse(redirect_url)


app.mount("/kael", StaticFiles(directory="ui", html=True), name="ui")


@app.on_event("startup")
async def on_startup():
    setup_logger()
    logger.info("应用程序启动，执行初始化操作")
    setup_poll_jms_event()


@app.on_event("shutdown")
async def on_shutdown():
    logger.info("应用程序关闭，执行清理操作")
    shutdown_protobuf()


if __name__ == "__main__":
    uvicorn.run(
        app,
        host=settings.http.host,
        port=settings.http.port,
        proxy_headers=True,
        forwarded_allow_ips='*',
        log_config=LOGGING,
    )
