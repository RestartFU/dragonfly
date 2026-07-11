.PHONY: lint mocks

MOCKGEN_VERSION ?= v0.6.0
MOCKGEN := go run go.uber.org/mock/mockgen@$(MOCKGEN_VERSION)

lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.4 run ./...

mocks:
	mkdir -p mocks
	$(MOCKGEN) -source=server/session/session.go -destination=mocks/session_conn.go -package=mocks -mock_names=Conn=MockSessionConn
	$(MOCKGEN) -source=server/player/provider.go -destination=mocks/player_provider.go -package=mocks -mock_names=Provider=MockPlayerProvider
	$(MOCKGEN) -source=server/world/provider.go -destination=mocks/world_provider.go -package=mocks -mock_names=Provider=MockWorldProvider
	$(MOCKGEN) -source=server/world/viewer.go -destination=mocks/world_viewer.go -package=mocks -mock_names=Viewer=MockWorldViewer
