deps:
	go get github.com/julienschmidt/httprouter
	go get github.com/kylelemons/go-gypsy/yaml
test: deps
	go test
build: test
	go build -i -o bin/queuekeeper
docker: build
	sudo docker build -t queuekeeper .
run:
	sudo docker run -d --rm --name queuekeeper_instance -p 8086:8086 queuekeeper
	sudo docker logs -f queuekeeper_instance
kill:
	sudo docker kill --signal 2 queuekeeper_instance

