TOOLS := cnc-path-cutting cnc-stl-view cnc-path-view cnc-path-translate

all: $(TOOLS)

%: %.go
	go build -o $@ $<

clean:
	rm -f $(TOOLS)
