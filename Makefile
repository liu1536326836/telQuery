all:
	gb build && cp bin/cmd .

clean:
	rm -rf cmd bin/ pkg/ *~ log.log*
