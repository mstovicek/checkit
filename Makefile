all: ;

start-fixer-php:
	docker build --rm -t checkit-php ./fixer/php
	docker run --rm -d -p 8080:80 -v `pwd`/files:/files/ --name checkit-php checkit-php

stop-fixer-php:
	docker kill checkit-php
