.PHONY: docs

docs:
	go list ./... | xargs -n1 go doc -all
	