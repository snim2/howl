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

all:	$(MODEL).png

$(MODEL).png:
	$(DIA) $(DOCS)/$(MODEL).dia $(DIAFLAGS) $(DOCS)/$(MODEL).png

clean:
	-@ $(RM) $(DOCS)/$(MODEL).png