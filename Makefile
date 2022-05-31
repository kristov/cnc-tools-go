TOOLS := cnc cnc-stl-view cnc-view2d svg2wkt

all: $(TOOLS)

%: %.go cnclib/cnclib.go
	go build -o $@ $<

clean:
	rm -f $(TOOLS)
