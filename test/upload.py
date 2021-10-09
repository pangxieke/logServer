#!/usr/bin/python3
# -*- coding: utf-8 -*-

import json
import requests

def run(host, port):
    params = {
        "mac_addr": "00:00:00:11",
        "event_time": 1633682385000
    }
    params = json.dumps(params)

    url = "%s:%d" % (host, port)
    print(url)

    response = requests.post(url, data=params)

    """响应"""

    print(response.status_code)
    print(response.content)


if __name__ == "__main__":
    run("http://localhost", 8080)
