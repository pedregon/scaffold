set dotenv-load

lint:
	golangci-lint run

format:
	go fmt

markdown:
	npm run markdownlint -- *.md

test:
	go test

pre-commit:
	pre-commit run --all-files

bump:
	git cliff --bump --unreleased