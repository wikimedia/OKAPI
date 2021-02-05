FROM golang:1.15

WORKDIR /var/www/app

COPY . .

RUN ["chmod", "+x", "/var/www/app/server/main.sh"]

CMD ["/var/www/app/server/main.sh"]
