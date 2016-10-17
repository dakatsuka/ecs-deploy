GIT := $(shell git rev-parse HEAD)

packages:
	cd cmd/ecs-deploy && gox -os="linux darwin" -arch="amd64" -output "../../pkg/{{.Dir}}-${GIT}-{{.OS}}-{{.Arch}}" -ldflags "-X main.version=${GIT}"

clean:
	rm -f cmd/ecs-deploy/ecs-deploy
	rm -f pkg/*
