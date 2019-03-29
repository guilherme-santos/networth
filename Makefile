mkfile_path := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

run:
	@export `cat ${mkfile_path}.env | xargs`; go run cmd/networth/*.go
