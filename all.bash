go get github.com/simonz05/godis
go get github.com/garyburd/go-oauth/oauth
go get github.com/virushuo/Go-Apns
go get github.com/googollee/go_c2dm
cat build.list | while read l;
do
  GOPATH=$GOPATH:$(pwd) go install $l
done
