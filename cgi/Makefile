OUT = math.cgi
CGIDIR ?= /var/www/cgi-bin
STATICDIR ?= /var/www/math.fent.uk/static
PROBDIR ?= /var/www/math.fent.uk/problems

all:
	go build -buildvcs=false -o $(OUT)

clean:
	rm -f $(OUT)

install:
	install -m 755 $(OUT) $(CGIDIR)/$(OUT)
	mkdir -p $(STATICDIR) $(PROBDIR)
	cp -R ../static/* $(STATICDIR)/
	cp -R ../problems/* $(PROBDIR)/

