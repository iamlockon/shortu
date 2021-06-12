#!bin/sh
set -ex

mockgen -source internal/cache/model.go -package mock -destination mock/cache_mock.go
mockgen -source internal/config/model.go -package mock -destination mock/config_mock.go
mockgen -source internal/db/model.go -package mock -destination mock/db_mock.go