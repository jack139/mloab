OUTPUT=mloab

BUILD_TAGS=
LD_FLAGS=
BUILD_FLAGS=-mod=readonly -ldflags "$(LD_FLAGS)"

CGO_ENABLED=1
BUILD_TAGS+=cleveldb

# allow users to pass additional flags via the conventional LDFLAGS variable
LD_FLAGS += $(LDFLAGS)

all: build

build:
	CGO_ENABLED=$(CGO_ENABLED) go build $(BUILD_FLAGS) -tags '$(BUILD_TAGS)' -o $(OUTPUT)

lev:
	CGO_ENABLED=$(CGO_ENABLED) go build $(BUILD_FLAGS) -tags '$(BUILD_TAGS)' misc/levstress.go

test:
	CGO_ENABLED=$(CGO_ENABLED) go build $(BUILD_FLAGS) misc/test.go

clean:
	@rm -f $(OUTPUT) levstress
	@go clean
