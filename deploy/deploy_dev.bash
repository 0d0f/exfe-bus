sudo /usr/local/etc/rc.d/monit stop
ls /var/run/gobus/worker_exfe_*.pid | while read l; do sudo kill `cat $l`; done
ls bin/ | while read l; do sudo cp bin/$l /usr/local/bin/exfebus_$l; chmod +x /usr/local/bin/exfebus_$l; done
sudo rm -r /usr/local/etc/exfebus/template
sudo cp -r template /usr/local/etc/exfebus/
sudo /usr/local/etc/rc.d/monit start
