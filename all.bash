go get $1 github.com/googollee/godis
go get $1 github.com/garyburd/go-oauth/oauth
go get $1 github.com/virushuo/Go-Apns
go get $1 github.com/googollee/go_c2dm
go get $1 github.com/googollee/goimap
go get $1 github.com/googollee/go-encoding-ex
go get $1 github.com/sloonz/go-iconv
cat build.list | while read l;
do
  GOPATH=$GOPATH:$(pwd) go install $1 $l
done
