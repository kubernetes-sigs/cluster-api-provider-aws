# Hosts running SELinux need :z added to volume mounts
SELINUX_ENABLED := $(shell cat /sys/fs/selinux/enforce 2> /dev/null || echo 0)

ifeq ($(SELINUX_ENABLED),1)
  DOCKER_VOL_OPTS?=:z
endif

.PHONY: docs
docs: diagrams


.PHONY: plantuml-image
plantuml-image:
	docker build . -t "plantuml-builder"

.PHONY: diagrams
diagrams: plantuml-image
	-docker run -u $(UID):$(GID) -v $(abspath .):/docs$(DOCKER_VOL_OPTS) plantuml-builder -tsvg /docs/**/*.md
	-docker run -u $(UID):$(GID) -v $(abspath .):/docs$(DOCKER_VOL_OPTS) plantuml-builder -tsvg /docs/**/*.plantuml


