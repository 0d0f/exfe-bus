go get $1 github.com/googollee/godis
go get $1 github.com/mrjones/oauth
go get $1 github.com/virushuo/Go-Apns
go get $1 github.com/googollee/go-gcm
go get $1 github.com/googollee/goimap
go get $1 github.com/googollee/go-multiplexer
go get $1 github.com/googollee/go-aws/smtp
go get $1 github.com/googollee/go-aws/s3
go get $1 github.com/googollee/go-aws/smtp
go get $1 github.com/googollee/go-encoding-ex
go get $1 github.com/googollee/go-logger
go get $1 github.com/Go-SQL-Driver/MySQL
go get $1 github.com/gorilla/mux
go get $1 launchpad.net/tomb
cat build.list | while read l;
do
  GOPATH=$GOPATH:$(pwd) go install $1 $l
done
