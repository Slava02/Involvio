package Involvio

//go:generate go run github.com/go-swagger/go-swagger/cmd/swagger generate server --name=Involvio  --spec=docs/swagger3.yaml --api-package=api --model-package=internal/entity --default-scheme=http --main-package=cmd/Involvio --server-package=internal/route --implementation-package=github.com/Slava02/Involvio/internal/app --regenerate-configureapi
