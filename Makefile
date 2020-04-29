bans-sidecar: #
	docker build -t redditopenttd/k8s-helper:latest -f docker/bans-sidecar/Dockerfile .

k8s-preinit:
	docker build -t redditopenttd/k8s-helper:latest -f docker/k8s-preinit/Dockerfile .