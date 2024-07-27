.PHONY: test
test:
	iverilog -g2001 cpu_test.v
	./a.out


.PHONY: view
view:
	gtkwave dump.vcd
