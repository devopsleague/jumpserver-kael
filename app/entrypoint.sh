#!/bin/bash
while [ "$(curl -I -m 10 -o /dev/null -s -w %{http_code} ${CORE_HOST}/api/health/)" != "200" ]
do
    echo "wait for jms_core $CORE_HOST ready"
    sleep 2
done

echo
date
echo "KAEL Version $VERSION, more see https://www.jumpserver.org"
echo "Quit the server with CONTROL-C."
echo

wisp

cd /opt/kael
uvicorn main:app --host 0.0.0.0 --port 8083