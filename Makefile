.PHONY: build
build:
	go install github.com/goreleaser/goreleaser@latest
	goreleaser --snapshot --skip-publish --rm-dist

.PHONY: release
release: ## example: make release V=0.0.0
	@read -p "Press enter to confirm and push to origin ..."
	git tag v$(V)
	git push origin v$(V)
