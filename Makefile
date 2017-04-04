.PHONY: docker-image docker-push
docker-image:
	docker build -t bjjb/urleen .
docker-push:
	docker push bjjb/urleen
