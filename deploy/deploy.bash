monit -g gobus_workers stop
ls /var/run/gobus/worker_exfe_*.pid | while read l; do sudo kill `cat $l`; done
ls bin/ | while read l; do sudo cp bin/$l /usr/local/bin/exfebus_$l; chmod +x /usr/local/bin/exfebus_$l; done
rm -r /usr/local/etc/exfebus/template /usr/local/etc/exfebus/templates
cp -r template templates /usr/local/etc/exfebus/
mkdir -p /var/run/gobus
monit -g gobus_workers start
