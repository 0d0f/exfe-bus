GODIR=${GOPATH%%:*}
go get github.com/simonz05/godis
cd $GODIR/src/github.com/simonz05/godis
git checkout stable
cd -
go get github.com/garyburd/go-oauth/oauth
cd $GODIR/src/github.com/garyburd/go-oauth/oauth
git checkout master
cd -
cat build.list | while read l;
do
  GOPATH=$GOPATH:$(pwd) go install $l
done
