#!/bin/bash

# imagec wrapper file - furnish the pretence of rpctool based default args

component=/sbin/imagec
args="$(rpctool -get vch$component) $@" 
${component}.bin $args
