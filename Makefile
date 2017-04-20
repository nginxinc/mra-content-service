tag = ngrefarch/content-service:mesos
ports = -p 8080:8080 # El primer port es l'extern, el segon l'intern

build:
	docker build -t $(tag) .

build-clean:
	docker build --no-cache -t $(tag) .

run:
	docker run -it $(ports) $(tag)

run-v:
	docker run -it ${env} $(ports) $(volumes) $(tag)

shell:
	docker run -it ${env} $(ports) $(volumes) $(tag) bash

push:
	docker push $(tag)

test:
	# Tests not yet implemented
