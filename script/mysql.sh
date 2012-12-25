#!/bin/sh

DIR=`dirname $0`
cd ${DIR}

DB=`./config.py '["db"]["db_name"]'`
DB_NAME=`./config.py '["db"]["username"]'`
DB_PASS=`./config.py '["db"]["password"]'`

PASS_ARG="-p${DB_PASS}"
if [ "${DB_PASS}" = "" ]
then
	PASS_ARG=""
fi

ls ${DIR}/mysql/*.sql | sort | while read l; do mysql -u"{DB_NAME}" "${PASS_ARG}" ${DB} < "$l"; done 