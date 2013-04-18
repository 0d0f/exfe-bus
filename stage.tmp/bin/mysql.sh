#!/bin/sh

DIR=`dirname $0`
cd ${DIR}

DB=`./config.py '["db"]["db_name"]'`
DB_NAME=`./config.py '["db"]["username"]'`
DB_PASS=`./config.py '["db"]["password"]'`

PASS_ARG="-p${DB_PASS}"
if [ "${DB_PASS}" = "" ]
then
	PASS_ARG="-u${DB_NAME}"
fi

cd mysql
LAST=`cat last`

ls *.sql | sort | awk -F'_' -v last="$LAST" '{if ($1 > last) print $0}' | while read l; do mysql -u"${DB_NAME}" "${PASS_ARG}" ${DB} < "$l"; done 

ls *.sql | sort | tail -1 | sed 's/^\([0-9]*\)_.*$/\1/' > last
