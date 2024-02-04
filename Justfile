set dotenv-load

current := '1.0.0'

tidy:
	go mod tidy

lint:
	golangci-lint run

format:
	go fmt

markdown:
	npm run markdownlint -- *.md

test:
	go test -v -covermode=count -coverprofile=coverage.out

report:
	go tool cover -html=coverage.out

pre-commit:
	pre-commit run --all-files

bump version=current:
	git cliff --unreleased --tag {{version}} -o CHANGELOG.md

release:
	@echo 'Release!'
	# git tag v{{current}} && git push origin v{{current}}

publish:
	@echo 'Publish!'
	# GOPROXY=proxy.golang.org go list -m github.com/pedregon/scaffold@v{{current}}

deploy: tidy lint test bump markdown release publish
	@echo v{{current}}