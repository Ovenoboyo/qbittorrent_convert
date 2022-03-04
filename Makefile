# PREFIX is environment variable, but if it is not set, then set default value
ifeq ($(PREFIX),)
    PREFIX := /usr/local
endif

all: install clean

install:
	go build -o qbittorrent_convert main.go
	install -m 0755 qbittorrent_convert $(PREFIX)/bin

clean:
	rm -rf qbittorrent_convert