go run /var/www/app/migrations/*.go migrate

export GO111MODULE=on
export PATH="$PATH:$(go env GOPATH)/bin"
apt install -y protobuf-compiler
go get google.golang.org/protobuf/cmd/protoc-gen-go google.golang.org/grpc/cmd/protoc-gen-go-grpc

if [ "$API_MODE" = 'prod' ]; then
  go build /var/www/app/*.go 
  ./main & ./main -server=queue -workers=30 & ./main -server=stream -workers=15 & ./main -server=runner -workers=250 -restart=true
else
  go get github.com/codegangsta/gin
  /go/bin/gin run /var/www/app/*.go
fi
