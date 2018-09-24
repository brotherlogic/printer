#!/usr/bin/python

from Adafruit_Thermal import *
import sys

printer = Adafruit_Thermal("/dev/serial0", 19200, timeout=5)

printer.println(sys.argv[1])
printer.feed(4)
