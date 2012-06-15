go get $1 github.com/googollee/godis
go get $1 github.com/garyburd/go-oauth/oauth
go get $1 github.com/virushuo/Go-Apns
go get $1 github.com/googollee/go_c2dm
cat build.list | while read l;
do
  GOPATH=$GOPATH:$(pwd) go install $1 $l
done
