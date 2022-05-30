TOOLS := cnc cnc-stl-view cnc-view2d

all: $(TOOLS)

%: %.go
	go build -o $@ $<

clean:
	rm -f $(TOOLS)
