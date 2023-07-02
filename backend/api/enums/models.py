from enum import auto

from strenum import StrEnum


class ChatSourceTypes(StrEnum):
    chat_gpt = auto()


chat_model_definitions = {
    "chat_gpt": {
        "gpt_3_5": "gpt-3.5-turbo",
        "gpt_4": "gpt-4",
    }
}

cls_to_source = {
    "ChatGPTModels": ChatSourceTypes.chat_gpt,
}


class BaseChatModelEnum(StrEnum):
    def code(self):
        source = cls_to_source.get(self.__class__.__name__, None)
        result = chat_model_definitions[source].get(self.name, None)
        assert result, f"model name not found: {self.name}"
        return result

    @classmethod
    def from_code(cls, code: str):
        source = cls_to_source.get(cls.__name__, None)
        for name, value in chat_model_definitions[source].items():
            if value == code:
                return cls[name]
        return None


class ChatGPTModels(BaseChatModelEnum):
    gpt_3_5 = auto()
    gpt_4 = auto()
