TOOLS := cnc cnc-view2d svg2wkt

all: $(TOOLS)

install: all
	cp $(TOOLS) ~/bin/

%: %.go cnclib/cnclib.go
	go build -o $@ $<

clean:
	rm -f $(TOOLS)
