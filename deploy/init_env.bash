U_NAME="$1"
G_NAME="$2"

mkdir -p /var/log/gobus
mkdir -p /var/run/gobus
mkdir -p /usr/local/etc/gobus

chown ${U_NAME}:${G_NAME} /var/log/gobus /var/run/gobus

echo RUN="/var/run/gobus" > /etc/init.d/exfe
echo mkdir -p '$RUN' >> /etc/init.d/exfe
echo chown ${U_NAME}:${G_NAME} '$RUN' >> /etc/init.d/exfe

chmod +x /etc/init.d/exfe

update-rc.d exfe defaults
