monit -g gobus stop
ls /var/run/gobus/exfe_*.pid | while read l; do sudo kill `cat $l`; done
ls bin/ | while read l; do sudo cp bin/$l /usr/local/bin/exfe_$l; chmod +x /usr/local/bin/exfe_$l; done
rm -r /usr/local/etc/gobus/templates
cp -r templates /usr/local/etc/gobus/
mkdir -p /var/run/gobus
monit -g gobus start
