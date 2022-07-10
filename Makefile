.PHONY: api redis-cli

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
	cd client-tmp && cp -r *.go go.mod go.sum docs README.md ../client/
	rm -rf client-tmp
	cd client && go mod tidy

redis-cli:
	docker run -it --network doless_default --rm redis:6-alpine redis-cli -h redis
