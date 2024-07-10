.PHONY: all build-server build-client run-server run-client

all: clean build-server build-client run-server run-client

build-server:
	@echo "================ BUILD TCP SERVER ===================="
	docker build -f server.Dockerfile -t my-tcp-server .
	@echo "============== COMPLETED TCP SERVER =================="

build-client:
	@echo "================ BUILD TCP CLIENT ===================="
	docker build -f client.Dockerfile -t my-tcp-client .
	@echo "============== COMPLETED TCP CLIENT =================="

run-server:
	@echo "================= RUN TCP SERVER ====================="
	#docker stop tcp-server
	#docker rm tcp-server
	docker run -d --name tcp-server -p12345:12345 my-tcp-server

run-client:
	@echo "================= RUN TCP CLIENT ====================="
	#docker stop tcp-client
	#docker rm tcp-client
	docker run --network host my-tcp-client

clean:
	@echo "================= CLEAN UP ==========================="
	docker rm -f tcp-server || true
	docker rmi my-tcp-server my-tcp-client || true
	@echo "============== COMPLETED CLEAN UP ===================="

