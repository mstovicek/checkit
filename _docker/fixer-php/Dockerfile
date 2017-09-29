FROM php:7-alpine

ADD https://github.com/FriendsOfPHP/PHP-CS-Fixer/releases/download/v2.5.0/php-cs-fixer.phar /bin/php-cs-fixer
RUN chmod +x /bin/php-cs-fixer
RUN php-cs-fixer self-update

WORKDIR /app

COPY .php_cs /app/
COPY main /app/

ENTRYPOINT ["./main"]
