.DEFAULT_GOAL := help
.PHONY: release, goreleaser, tg, tgrep

# Comments starting with "#\#" are documentation output for help.
# Keep targets alpha-sorted so they look nicely organized in help output.

build: ## Build the tg and tgrep binaries - output into the ./bin directory
	go mod tidy
	go build -o $(PWD)/bin/tg $(PWD)/cmd/tg
	go build -o $(PWD)/bin/tgrep $(PWD)/cmd/tgrep

clean: ## Remove temporary artifacts from build environment
	go clean
	rm bin/tg
	rm bin/tgrep

goreleaser:
	goreleaser --snapshot --skip-publish --rm-dist


tg: ## Use go to run 'tg'
	cd cmd/tg && go run .
tgrep: ## Use go to run 'tgrep'
	cd cmd/tgrep && go run .

release:
	[[ $$( git rev-parse --abbrev-ref HEAD ) == "main" ]] # make sure we are on main
	git push origin main
	git tag $$( git tag | grep "^v" | sort --version-sort | tail -n 1 | awk -F. '{OFS="."; $$3 = $$3 + 1; print}' )
	git push --tags

help: ## This help dialog.
	@IFS=$$'\n' ; \
		help_lines=(`fgrep -h "##" $(MAKEFILE_LIST) | egrep -v '^--' | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##/:/'`); \
		printf "%-30s %s\n" "target" "help" ; \
		printf "%-30s %s\n" "------" "----" ; \
		for help_line in $${help_lines[@]}; do \
		IFS=$$':' ; \
		help_split=($$help_line) ; \
		help_command=`echo $${help_split[0]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
		help_info=`echo $${help_split[2]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
		printf '\033[36m'; \
		printf "%-30s %s" $$help_command ; \
		printf '\033[0m'; \
		printf "%s\n" $$help_info; \
		done
