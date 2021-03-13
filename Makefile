VERSION=$(shell git describe --no-match --dirty --always --abbrev=8)
REPOSITORY=quay.io/conprof/perfessor
docker: bin
	docker build . -t perfessor
	docker tag perfessor:latest $(REPOSITORY):$(VERSION)-linux-amd64

push: docker
	docker push $(REPOSITORY):$(VERSION)-linux-amd64

bin:
	go build -o perfessor ./cmd/perfessor/main.go

run: docker
	docker run --rm -it --cap-add SYS_ADMIN perfessor:latest

clean: 
	rm -f perfessor
