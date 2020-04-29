bans-sidecar: #
	docker build -t redditopenttd/k8s-helpers:latest -f docker/bans-sidecar/Dockerfile .

k8s-preinit:
	docker build -t redditopenttd/k8s-helpers:latest -f docker/k8s-preinit/Dockerfile .