include $(GOROOT)/src/Make.inc

TARG=github.com/mkrautz/goar

GOFILES=\
	common.go\
	writer.go\
	reader.go

include $(GOROOT)/src/Make.pkg
