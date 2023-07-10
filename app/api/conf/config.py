from typing import Optional, Literal

from pydantic import BaseModel, validator, Field

from utils import singleton_with_lock
from .base import BaseConfig

__all__ = ['settings']

_TYPE_CHECKING = False


class HttpSetting(BaseModel):
    host: str = '0.0.0.0'
    port: int = Field(8083, ge=1, le=65535)
    cors_allow_origins: list[str] = ['*']

    @staticmethod
    @validator("test_host")
    def validate_host(v):
        return v


class ChatGPTSetting(BaseModel):
    openai_base_url: str = 'https://api.openai.com/v1/'
    proxy: Optional[str] = None
    connect_timeout: int = Field(4, ge=1)
    read_timeout: int = Field(8, ge=1)
    api_key: Optional[str] = None


class LogSetting(BaseModel):
    log_level: Literal['INFO', 'DEBUG', 'WARNING'] = 'INFO'


class ConfigModel(BaseModel):
    chat_gpt: ChatGPTSetting = ChatGPTSetting()
    http: HttpSetting = HttpSetting()
    log: LogSetting = LogSetting()

    class Config:
        underscore_attrs_are_private = True


@singleton_with_lock
class Config(BaseConfig[ConfigModel]):
    if _TYPE_CHECKING:
        chat_gpt: ChatGPTSetting = ChatGPTSetting()
        http: HttpSetting = HttpSetting()
        log: LogSetting = LogSetting()

    def __init__(self):
        super().__init__(ConfigModel, "config.yml")


settings = Config()
