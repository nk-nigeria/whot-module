PROJECT_NAME=github.com/nakamaFramework/sicbo-module
APP_NAME=sicbo.so
APP_PATH=$(PWD)

update-submodule-dev:
	git submodule update --init
	git submodule update --remote
	cd ./cgp-common && git checkout develop && git pull && cd ..
	go get github.com/nakamaFramework/cgp-common@develop
update-submodule-stg:
	git submodule update --init
	git submodule update --remote
	cd ./cgp-common && git checkout main && cd ..
	go get github.com/nakamaFramework/cgp-common@main
cpdev:
	scp ./bin/${APP_NAME} nakama:/root/cgp-server-dev/dist/data/modules/
build:
	go mod tidy
	go mod vendor
	docker run --rm -w "/app" -v "${APP_PATH}:/app" heroiclabs/nakama-pluginbuilder:3.11.0 build -buildvcs=false --trimpath --buildmode=plugin -o ./bin/${APP_NAME} && cp ./bin/${APP_NAME} ../bin/

sync:
	rsync -aurv --delete ./bin/${APP_NAME} root@cgpdev:/root/cgp-server/dev/data/modules/
	# ssh root@cgpdev 'cd /root/cgp-server && docker restart nakama'

bsync: build sync

dev: update-submodule-dev build cpdev
stg: update-submodule-stg build
proto:
	protoc -I ./ --go_out=$(pwd)/proto  ./proto/sicbo_api.proto
	
