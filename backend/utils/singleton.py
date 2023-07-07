from threading import Lock
from typing import Type, TypeVar

T = TypeVar("T")

__all__ = ['singleton_with_lock']


def singleton_with_lock(cls: Type[T]):
    instances = {}
    lock = Lock()

    def get_instance(*args, **kwargs) -> T:
        with lock:
            if cls not in instances:
                instances[cls] = cls(*args, **kwargs)
            return instances[cls]

    return get_instance
