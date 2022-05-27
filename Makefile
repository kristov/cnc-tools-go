TOOLS := cnc-cutting-path cnc-stl-view cnc-view-paths

all: $(TOOLS)

%: %.go
	go build -o $@ $<

clean:
	rm -f $(TOOLS)
