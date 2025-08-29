REPO=blacktop
NAME=go-macho
VERSION=$(shell svu current)
NEXT_VERSION:=$(shell svu patch)

GIT_COMMIT=$(git rev-parse HEAD)
GIT_DIRTY=$(test -n "`git status --porcelain`" && echo "+CHANGES" || true)
GIT_DESCRIBE=$(git describe --tags)

.PHONY: dev-deps
dev-deps: ## Install the dev dependencies
	@brew install gh
	@go install github.com/goreleaser/chglog/cmd/chglog@latest
	@go install github.com/caarlos0/svu@v1.4.1

.PHONY: bump
bump: ## Incriment version patch number
	@echo " > Bumping VERSION"
	@chglog add --version ${NEXT_VERSION}

.PHONY: changelog
changelog: bump ## Create a new CHANGELOG.md
	@echo " > Creating CHANGELOG.md"
	@chglog format --template release > CHANGELOG.md

.PHONY: release
release: changelog ## Create a new release from the VERSION
	@echo " > Creating Release"
	@gh release create ${NEXT_VERSION} -F CHANGELOG.md

.PHONY: destroy
destroy: ## Remove release from the VERSION
	@echo " > Deleting Release"
	git tag -d ${VERSION}
	git push origin :refs/tags/${VERSION}

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help