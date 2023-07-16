import argparse
import subprocess


class FastapiBabel:
    """
        pybabel extract -F i18n/babel.cfg -o i18n/messages.pot api
        pybabel init -i i18n/messages.pot -d i18n -l zh_CN
        pybabel init -i i18n/messages.pot -d i18n -l ja_JP
        pybabel update -i i18n/messages.pot -d i18n
        pybabel compile -d i18n
    """

    def __init__(self):
        self.extract_command = [
            'pybabel', 'extract',
            '-F', 'i18n/babel.cfg',
            '-o', 'i18n/messages.pot',
            'api'
        ]

        self.update_command = [
            'pybabel', 'update',
            '-i', 'i18n/messages.pot',
            '-d', 'i18n',
        ]

        self.compile_command = ['pybabel', 'compile', '-d', 'i18n']

    def init_language(self, language_name: str):
        self.run_command(self.extract_command)
        command = [
            'pybabel', 'init', '-i', 'i18n/messages.pot',
            '-d', 'i18n', '-l', language_name
        ]
        self.run_command(command)

    def makemessages(self):
        self.run_command(self.extract_command)
        self.run_command(self.update_command)

    def compilemessages(self):
        self.run_command(self.compile_command)

    @staticmethod
    def run_command(command):
        result = subprocess.run(
            command, capture_output=True, text=True
        )
        output = result.stdout + result.stderr
        print(output)


if __name__ == '__main__':
    parser = argparse.ArgumentParser(
        description="""
        Kael console is currently engaged in internationalization related operations;

        Example: \r\n

        %(prog)s makemessages;
        """
    )
    parser.add_argument(
        'action', type=str,
        choices=("makemessages", "compilemessages"),
        help="Action to run"
    )

    args = parser.parse_args()

    action = args.action
    if action == "makemessages":
        FastapiBabel().makemessages()
    elif action == "compilemessages":
        FastapiBabel().compilemessages()
    else:
        print("Action not found: %s" % action)
