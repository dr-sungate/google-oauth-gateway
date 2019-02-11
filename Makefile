.PHONY: install,run

install:
	cd /go/srcgithub.com/dr-sungate/google-oauth-gateway/gateway &&  export GO111MODULE=on go install

run:
	export GO111MODULE=on && go run ./gateway/gateway.go
