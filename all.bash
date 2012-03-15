go get github.com/simonz05/godis
go get github.com/googollee/go-gypsy
go get github.com/garyburd/go-oauth
cat build.list | while read l;
do
  GOPATH=$GOPATH:$(pwd) go install $l
done
