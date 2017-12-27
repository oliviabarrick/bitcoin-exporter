VERSION=$(shell git rev-list --count HEAD)-$(shell git describe --always --long)

.PHONY: build
build:
	docker-compose -f docker/stack.yml build

.PHONY: build-image
build-image:
	docker build -t justinbarrick/bitcoin-exporter:$(VERSION) -f docker/Dockerfile .

.PHONY: build-bin
build-bin:
	docker-compose -f docker/docker-compose.yml up --build

.PHONY: up
up:
	docker-compose -f docker/stack.yml up --force-recreate -d
	$(MAKE) status

.PHONY: status
status:
	@echo Nginx URL: https://$(shell docker inspect --format '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $(shell docker ps -f name=nginx -q))/
	@echo Alert Manager URL: http://$(shell docker inspect --format '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $(shell docker ps -f name=alertmanager -q)):9093/
	@echo Prometheus URL: http://$(shell docker inspect --format '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $(shell docker ps -f name=prometheus -q)):9090/
	@echo Grafana URL: http://$(shell docker inspect --format '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $(shell docker ps -f name=grafana -q)):3000/dashboard/db/bitcoins
	@echo Grafana login: admin:admin

.PHONY: logs
logs:
	docker-compose -f docker/stack.yml logs -f bitcoin-exporter

.PHONY: push-image
push-image: build-image
	docker push justinbarrick/bitcoin-exporter:$(VERSION)

.PHONY: initial-deploy
initial-deploy: push-image
	kubectl create configmap bitcoin-exporter-config --from-file=config/bitcoin-exporter.yml
	sed 's/VERSION/$(VERSION)/g' docker/bitcoin-exporter-kubernetes.yml |kubectl apply -f -

.PHONY: update-deploy
update-deploy:
	kubectl set image deployment/bitcoin-exporter-deployment bitcoin-exporter=justinbarrick/bitcoin-exporter:$(VERSION)
