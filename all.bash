go get github.com/simonz05/godis
go get github.com/googollee/go-gypsy
go get github.com/garyburd/go-oauth
GOPATH=$GOPATH:$(pwd) go install config
GOPATH=$GOPATH:$(pwd) go install gobus
GOPATH=$GOPATH:$(pwd) go install gosque
GOPATH=$GOPATH:$(pwd) go install oauth
GOPATH=$GOPATH:$(pwd) go install twitter/job
GOPATH=$GOPATH:$(pwd) go install twitter/service
GOPATH=$GOPATH:$(pwd) go install twitter
GOPATH=$GOPATH:$(pwd) go install mail
