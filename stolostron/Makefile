IMG_REG ?= registry.redhat.io
IMG_REPO ?= multicluster
IMG_NAME ?= cluster-api-provider-aws-rhel9
IMG_TAG ?= latest
IMG ?= ${IMG_REG}/${IMG_REPO}/${IMG_NAME}:${IMG_TAG}

docker-build:
	docker build -f Dockerfile.stolostron ../ -t ${IMG}
