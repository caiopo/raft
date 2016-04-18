#! /usr/bin/env python3

import sys
import json

pipe = sys.stdin.readlines()

pipe = json.loads(pipe)

print(pipe['spec']['ports']['nodePort'])