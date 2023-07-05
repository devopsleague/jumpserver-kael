from jms.base import BaseHandler


class CommandHandler(BaseHandler):

    async def test(self):
        from time import sleep

        sleep(3)
        print('-----------------')
        pass
