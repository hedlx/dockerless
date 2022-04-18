.PHONY: api

api:
	rm -rf client client-tmp
	openapi-generator generate \
		-i openapi.yaml \
		-g go \
		-o client-tmp \
		--additional-properties=packageName="api" \
		--git-user-id "hedlx" \
		--git-repo-id "doless/client"
	mkdir -p client
	cp -r client-tmp/{*.go,go.mod,go.sum,docs,README.md} client/
	rm -rf client-tmp
	cd client && go mod tidy
