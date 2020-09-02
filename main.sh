go run /var/www/app/migrations/*.go migrate

if [ "$API_MODE" = 'prod' ]; then
  go build /var/www/app/*.go 
  ./main & ./main -server=queue -workers=30 & ./main -server=stream -workers=15 & ./main -server=runner -restart=true
else
  go get github.com/codegangsta/gin
  /go/bin/gin run /var/www/app/*.go
fi
