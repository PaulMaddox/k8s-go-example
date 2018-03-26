all: deploy

name := nlb-test
registry := 860498507463.dkr.ecr.us-west-2.amazonaws.com
tag := $(shell git describe --tags --always --dirty)
podfile := kubernetes.yaml

build: 
	$(call blue, "Building container...")
	docker build -t ${name} .

tag: build
	$(call blue, "Tagging image (tag: ${tag})...")
	docker tag ${name} ${name}:${tag}
	docker tag ${name} ${registry}/${name}:latest
	docker tag ${name} ${registry}/${name}:${tag}

publish: tag  
	$(call blue, "Publishing Docker image to registry...")
	docker push ${registry}/${name}:${tag} 
	docker push ${registry}/${name}:latest 

deploy: publish
	$(call blue, "Deploying to Kubernetes...")
	sed -i 's/image: .*/image: ${registry}\/${name}:${tag}/g' ${podfile}
	git add ${podfile}
	git commit -m 'Updated Kubernetes podfile'
	kubectl apply -f ${podfile}

run: build
	$(call blue, "Running Docker image locally...")
	docker run -i -t --rm -P ${name}:latest

define blue
	@tput setaf 6
	@echo $1
	@tput sgr0
endef
