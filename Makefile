VERSION=$(shell git rev-list --count HEAD)-$(shell git describe --always --long)
DEPLOYMENT=bitcoin-exporter

.PHONY: build
build:
	docker-compose -f docker/stack.yml build

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

.PHONY: build-image
build-image:
	docker build -t justinbarrick/bitcoin-exporter:$(VERSION) -f docker/Dockerfile .

.PHONY: push-image
push-image:
	docker push justinbarrick/bitcoin-exporter:$(VERSION)

.PHONY: deploy
deploy:
	helm install --set image.tag=$(VERSION) --name $(DEPLOYMENT) helm

.PHONY: upgrade-deploy
upgrade-deploy:
	helm upgrade --set image.tag=$(VERSION) $(DEPLOYMENT) helm
