go get $1 github.com/mrjones/oauth
go get $1 code.google.com/p/go-imap/go1/imap
go get $1 github.com/virushuo/Go-Apns
go get $1 github.com/googollee/go-socket.io
go get $1 github.com/googollee/go-rest
go get $1 github.com/googollee/godis
go get $1 github.com/googollee/go-gcm
go get $1 github.com/googollee/go-multiplexer
go get $1 github.com/googollee/go-aws/smtp
go get $1 github.com/googollee/go-aws/s3
go get $1 github.com/googollee/go-aws/smtp
go get $1 github.com/googollee/go-encoding
go get $1 github.com/googollee/go-logger
go get $1 github.com/go-sql-driver/mysql
go get $1 github.com/gorilla/mux
go get $1 launchpad.net/tomb
cat build.list | while read l;
do
  GOPATH=$GOPATH:$(pwd) go install $l
done
