OUT = math.cgi
CGIDIR ?= /var/www/cgi-bin
STATICDIR ?= /var/www/math.fent.uk

all:
	go build -buildvcs=false -o $(OUT)

clean:
	rm -f $(OUT)

install:
	install -m 755 $(OUT) $(CGIDIR)/$(OUT)
	mkdir -p $(STATICDIR)
	cp -R ../static/* $(STATICDIR)/

