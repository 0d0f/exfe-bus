rm -Rf ./deploy

mkdir -p ./deploy/usr/local/bin
mkdir -p ./deploy/usr/local/etc

ls ./bin | while read l; do cp -Rf ./bin/$l ./deploy/usr/local/bin/exfebus_$l; done
cp -Rf ./template ./deploy/usr/local/etc/

cd ./deploy
rsync --progress -av * root@exfe.com:/
ssh root@exfe.com 'ls /var/run/resque/worker_exfe_*.pid | while read l; do kill `cat $l`; done'
