import logging
import logging.config
import os

from api.conf import settings
from api import globals as g

LOG_DIR = os.path.join(g.PROJECT_DIR, 'logs')
UNEXPECTED_EXCEPTION_LOG_FILE = os.path.join(LOG_DIR, 'kael.log')
LOG_LEVEL = settings.log.log_level

LOGGING = {
    'version': 1,
    'formatters': {
        'simple': {
            'format': '%(asctime)s.%(msecs)03d %(levelname)8s: [%(name)s]\t%(message)s',
            'datefmt': '%Y/%m/%d %H:%M:%S'
        },
        'proxy-output': {
            'format': '%(message)s',
            'datefmt': '%Y/%m/%d %H:%M:%S'
        },
        'colored': {
            '()': 'colorlog.ColoredFormatter',
            'datefmt': '%Y/%m/%d %H:%M:%S',
            'format': '%(asctime)s.%(msecs)03d %(log_color)s%(levelname)8s%(reset)s: %(cyan)s[%(name)s]%(reset)s %(message)s',
        }
    },
    'handlers': {
        'file_handler': {
            'class': 'logging.handlers.RotatingFileHandler',
            'formatter': 'simple',
            'encoding': 'utf8',
            'level': 'DEBUG',
            'filename': UNEXPECTED_EXCEPTION_LOG_FILE,
            'maxBytes': 1024 * 1024 * 100
        },
        'console_handler': {
            'class': 'logging.StreamHandler',
            'formatter': 'colored',
            'level': 'DEBUG'
        },
    },
    'root': {
        'level': 'DEBUG',
        'handlers': []
    },
    'loggers': {
        'uvicorn.error': {
            'level': 'INFO',
            'handlers': ['file_handler']
        },
        'uvicorn.access': {
            'level': 'INFO',
            'handlers': ['console_handler']
        },
        'kael': {
            'level': 'INFO',
            'handlers': ['console_handler', 'file_handler']
        },
    }
}


def setup_logger():
    os.makedirs(LOG_DIR, exist_ok=True)
    logging.config.dictConfig(LOGGING)


def get_logger(name):
    return logging.getLogger(f'kael.{name}')
