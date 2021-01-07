.PHONY: build push
build:
	docker build -t bjjb/urleen .
push:
	docker push bjjb/urleen
