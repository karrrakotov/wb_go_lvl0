build:
	docker build -t wbapp .
build-compose:
	docker compose up --build wildberries
run:
	docker run -d -p 80:1234 --rm --name wbappcont wbapp
stop:
	docker stop wbappcont
rmca:
	docker container prune
rmia:
	docker image prune