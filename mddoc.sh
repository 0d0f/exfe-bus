#!/bin/sh

cat - | grep -A1 . | grep -v "^--$" | sed 's/^\([ftpc })]\)/    \1/g' | sed -E 's/^([a-zA-Z]+)$/\1\
-------/g' | sed -E 's/^    type ([A-Z][a-zA-Z]*)/type <span id="\1">\1<\/span>\
\
    type \1/g' | sed -E 's/    func ([A-Z][a-zA-Z]*)/func <span id="\1">\1<\/span>\
\
    func \1/g'
