go get github.com/simonz05/godis
go get github.com/garyburd/go-oauth/oauth
cat build.list | while read l;
do
  GOPATH=$GOPATH:$(pwd) go install $l
done
