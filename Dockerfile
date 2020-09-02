FROM golang:1.14

WORKDIR /var/www/app

COPY . .

RUN apt-get update

RUN chmod +x /var/www/app/main.sh

CMD /var/www/app/main.sh