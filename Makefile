#
# Simple Makefile to generate some docs.
# 
# TODO: Add support for unit testing, etc.
#

SHELL=/bin/sh
RM=/bin/rm
DIA=/usr/bin/dia
DIAFLAGS=-e

DOCS=docs
MODEL=data_model

GODOC=/usr/bin/godoc
GODOC_FLAGS=-timestamps=true -index -html

all:	docs $(MODEL).png 

docs:	$(DOCS)/howl.html $(DOCS)/model.html $(DOCS)/view.html $(DOCS)/controller.html

$(MODEL).png:
	$(DIA) $(DOCS)/$(MODEL).dia $(DIAFLAGS) $(DOCS)/$(MODEL).png

$(DOCS)/%.html: % 
	$(GODOC) $(GODOC_FLAGS) ./$</ > $@

clean:
	-@ $(RM) $(DOCS)/$(MODEL).png
	-@ $(RM) $(DOCS)/*.html
	-@ $(RM) $(DOCS)/*~