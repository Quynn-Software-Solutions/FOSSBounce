.PHONY: run build_docker run_docker


run:
	go build . && ./fossbounce

build_docker:
	docker build -t fossbounce .

run_docker:
	docker run -p 8080:8080 -t fossbounce