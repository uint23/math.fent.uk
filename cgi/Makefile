OUT = math.cgi
CGIDIR ?= /var/www/cgi-bin
GOFLAGS = -buildvcs=false

all:
	go build $(GOFLAGS) -o $(OUT)

clean:
	rm -f $(OUT)

install:
	install -m 755 $(OUT) $(CGIDIR)/$(OUT)

