import logging
import logging.config
import os

from api.conf import settings
from api import globals as g

LOG_DIR = os.path.join(g.PROJECT_DIR, 'logs')
UNEXPECTED_EXCEPTION_LOG_FILE = os.path.join(LOG_DIR, 'unexpected_exception.log')
LOG_LEVEL = settings.log.log_level

LOGGING = {
    'version': 1,
    'disable_existing_loggers': False,
    'formatters': {
        'verbose': {
            'format': '%(levelname)s %(asctime)s %(module)s %(process)d %(thread)d %(message)s'
        },
        'main': {
            'datefmt': '%Y-%m-%d %H:%M:%S',
            'format': '%(asctime)s [%(module)s %(levelname)s] %(message)s',
        },
        'exception': {
            'datefmt': '%Y-%m-%d %H:%M:%S',
            'format': '\n%(asctime)s [%(levelname)s] %(message)s',
        },
        'simple': {
            'format': '%(levelname)s %(message)s'
        },
        'syslog': {
            'format': 'jumpserver: %(message)s'
        },
        'msg': {
            'format': '%(message)s'
        }
    },
    'handlers': {
        'null': {
            'level': 'DEBUG',
            'class': 'logging.NullHandler',
        },
        'console': {
            'level': 'DEBUG',
            'class': 'logging.StreamHandler',
            'formatter': 'main'
        },
        'unexpected_exception': {
            'encoding': 'utf8',
            'level': 'DEBUG',
            'class': 'logging.handlers.RotatingFileHandler',
            'formatter': 'exception',
            'maxBytes': 1024 * 1024 * 100,
            'backupCount': 7,
            'filename': UNEXPECTED_EXCEPTION_LOG_FILE,
        },
    },
    'loggers': {
        'uvicorn.error:': {
            'handlers': ['console', 'unexpected_exception'],
            'level': LOG_LEVEL,
        },
        'uvicorn.access': {
            'handlers': ['console'],
            'level': LOG_LEVEL,
        }
    }
}


def setup_logger():
    os.makedirs(LOG_DIR, exist_ok=True)
    logging.config.dictConfig(LOGGING)


def get_logger(name):
    return logging.getLogger(f'{name}')
