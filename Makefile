
# dependencies installs all required packages
dependencies:
	# Build Dependencies
	go get -u github.com/btcsuite/btcutil
	go get -u github.com/bitspill/bitsig-go
	go get -u github.com/bitspill/json-patch
	go get -u github.com/metacoin/flojson
	# Dev Dependencies
	go get -u github.com/golang/lint/golint
	go get -u github.com/mattn/goveralls

# pkgs changes which packages the makefile calls operate on. run changes which
# tests are run during testing.
run = Test
pkgs = ./messages ./utility

# fmt calls go fmt on all packages.
fmt:
	gofmt -s -l -w $(pkgs)

# vet calls go vet on all packages.
# NOTE: go vet requires packages to be built in order to obtain type info.
vet: release
	go vet $(pkgs)

# will always run on some packages for a while.
lintpkgs = ./messages ./utility
lint:
	@for package in $(lintpkgs); do           \
		golint -min_confidence=1.0 $$package; \
	done

# clean removes all directories that get automatically created during
# development.
clean:
	@rm -rf release cover

test:
	go test -short -timeout=5s $(pkgs) -run=$(run)
test-v:
	go test -race -v -timeout=15s $(pkgs) -run=$(run)
test-long:
	go test -v -race -timeout=500s $(pkgs) -run=$(run)
bench:
	go test -tags='testing' -timeout=500s -run=XXX -bench=. $(pkgs)
cover:
	@set +x
	@mkdir -p cover
	@echo "mode: count" > cover/profile.cov
	@for package in $(pkgs); do                                                               \
		go test -timeout=500s -covermode=count -coverprofile=cover/$$package.out ./$$package  \
		&& go tool cover -html=cover/$$package.out -o=cover/$$package.html;                   \
		if [ -f cover/$$package.out ]; then                                                   \
			cat cover/$$package.out | tail -n +2 >> cover/profile.cov;                        \
		fi;                                                                                   \
	done
	@set -x

coveralls:
	goveralls -coverprofile=cover/profile.cov -service=travis-ci -repotoken=$$COVERALLS_TOKEN

.PHONY: all dependencies fmt clean test test-v test-long release cover coveralls