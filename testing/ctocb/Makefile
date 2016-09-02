.PHONY: go16 go17 push

push: go16 go17
	gcloud docker push gcr.io/golang-samples-tests/go16
	gcloud docker push gcr.io/golang-samples-tests/go17

go16:
	docker build -t gcr.io/golang-samples-tests/go16 -f Dockerfile-1.6 .

go17:
	docker build -t gcr.io/golang-samples-tests/go17 -f Dockerfile-1.7 .
