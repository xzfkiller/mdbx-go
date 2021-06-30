
.PHONY: all test

BRANCH=`git rev-parse --abbrev-ref HEAD`
COMMIT=`git rev-parse --short HEAD`
MASTER_COMMIT=`git rev-parse --short origin/master`
GOLDFLAGS="-X main.branch $(BRANCH) -X main.commit $(COMMIT)"

all:
test:
	go test ./mdbx

clean:
	@rm -rf *~
