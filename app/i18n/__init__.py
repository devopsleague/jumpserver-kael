from babel.support import Translations
from api import globals

TRANSLATIONS = {
    "zh-hans": Translations.load(globals.I18N_DIR, locales=["zh_CN"]),
}


def gettext(msg: str):
    language = globals.language
    t = TRANSLATIONS.get(language, 'zh-hans')
    return t.gettext(msg)
