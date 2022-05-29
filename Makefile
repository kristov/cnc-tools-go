TOOLS := cnc cnc-stl-view cnc-path-view

all: $(TOOLS)

%: %.go
	go build -o $@ $<

clean:
	rm -f $(TOOLS)
