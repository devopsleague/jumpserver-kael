import os

import yaml
from typing import TypeVar, Generic, Type
from pydantic import BaseModel

from api import globals as g

__all__ = ['BaseConfig']

T = TypeVar("T", bound=BaseModel)


class BaseConfig(Generic[T]):
    _config: T = None
    _config_path = None
    _config_type = None

    def __init__(self, _config_type: Type, config_filename: str):
        self._config_type = _config_type
        config_dir = os.environ.get('KAEL_CONFIG_DIR', g.PROJECT_DIR)
        self._config_path = os.path.join(config_dir, config_filename)
        self.load()

    def __getattr__(self, key):
        return getattr(self._config, key)

    def __setattr__(self, key, value):
        if key in ('_config', '_config_type', '_config_path'):
            super().__setattr__(key, value)
        else:
            setattr(self._config, key, value)

    def load(self):
        if not os.path.exists(self._config_path):
            config_dict = self._config_type().dict()
            self._config = self._config_type(**config_dict)
            return
        try:
            with open(self._config_path, encoding='utf8') as f:
                config_dict = yaml.safe_load(f) or {}
                self._config = self._config_type(**config_dict)
        except Exception as e:
            raise Exception(f'cannot read config ({self._config_path}), error: {str(e)}')
