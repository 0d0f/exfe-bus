#!/usr/bin/env python

import json
import sys

with open('/usr/local/etc/gobus/exfe.json') as fp:
    config = json.load(fp)

if len(sys.argv) > 1:
    cmd = "config%s" % (sys.argv[1])
    print eval(cmd)