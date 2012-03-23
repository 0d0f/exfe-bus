rm -Rf ./deploy

mkdir -p ./deploy/usr/local/bin
mkdir -p ./deploy/usr/local/etc

cp -Rf ./bin/twitter ./deploy/usr/local/bin/exfebus_twitter
cp -Rf ./template ./deploy/usr/local/etc/

cd ./deploy
rsync --progress -av * root@exfe.com:/
ssh root@exfe.com 'kill `cat /var/run/resque/worker_exfe_twitter.pid`'
