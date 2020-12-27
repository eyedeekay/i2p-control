
USER_GH=eyedeekay
VERSION=0.1.04
packagename=i2p-control

GO_COMPILER_OPTS = -a -tags netgo -ldflags '-w -extldflags "-static"'

echo:
	@echo "type make version to do release $(VERSION)"

version:
	gothub release -s $(GITHUB_TOKEN) -u $(USER_GH) -r $(packagename) -t v$(VERSION) -d "version $(VERSION)"

del:
	gothub delete -s $(GITHUB_TOKEN) -u $(USER_GH) -r $(packagename) -t v$(VERSION)

tar:
	tar --exclude .git \
		--exclude .go \
		--exclude bin \
		--exclude examples \
		-cJvf ../$(packagename)_$(VERSION).orig.tar.xz .

all: windows osx linux

windows: fmt
	GOOS=windows go build $(GO_COMPILER_OPTS) -o $(packagename).exe

osx: fmt
	GOOS=darwin go build $(GO_COMPILER_OPTS) -o $(packagename)-darwin

linux: fmt
	GOOS=linux go build $(GO_COMPILER_OPTS) -o $(packagename)

sumwindows=`sha256sum $(packagename).exe`
sumlinux=`sha256sum $(packagename)`
sumdarwin=`sha256sum $(packagename)-darwin`

upload-windows:
	gothub upload -R -u eyedeekay -r "$(packagename)" -t v$(VERSION) -l "$(sumwindows)" -n "$(packagename).exe" -f "$(packagename).exe"

upload-darwin:
	gothub upload -R -u eyedeekay -r "$(packagename)" -t v$(VERSION) -l "$(sumdarwin)" -n "$(packagename)-darwin" -f "$(packagename)-darwin"

upload-linux:
	gothub upload -R -u eyedeekay -r "$(packagename)" -t v$(VERSION) -l "$(sumlinux)" -n "$(packagename)" -f "$(packagename)"

upload: upload-windows upload-darwin upload-linux

release: version upload

fmt:
	gofmt -w -s main.go

