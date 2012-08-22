U_NAME="$1"
G_NAME="$2"

mkdir -p /var/log/gobus
mkdir -p /var/run/gobus
mkdir -p /usr/local/etc/exfebus

chown %{U_NAME}:%{G_NAME} /var/log/gobus /var/run/gobus