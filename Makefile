#
# Simple Makefile to generate some docs.
# 
# TODO: Add support for unit testing, etc.
#

APP=fieldsensing

SHELL=/bin/sh
RM=/bin/rm
DIA=/usr/bin/dia
DIAFLAGS=-e

DOCS=docs
MODEL=data_model

GODOC=/usr/bin/godoc
GODOC_FLAGS=-timestamps=true -index -html

DEVSERVER=dev_appserver.py 
APPCFG=appcfg.py
APPENGINE_PATH=~/inst/go_appengine
DEVSERVER_FLAGS=

all:	docs $(MODEL).png dev

dev:
	$(APPENGINE_PATH)/$(DEVSERVER) $(DEVSERVER_FLAGS) .

docs:	$(DOCS)/howl.html $(DOCS)/model.html $(DOCS)/view.html $(DOCS)/controller.html

$(MODEL).png:
	$(DIA) $(DOCS)/$(MODEL).dia $(DIAFLAGS) $(DOCS)/$(MODEL).png

$(DOCS)/%.html: % 
	$(GODOC) $(GODOC_FLAGS) ./$</ > $@

upload:
	$(APPENGINE_PATH)/$(APPCFG) update .

clean:
	-@ $(RM) $(DOCS)/$(MODEL).png
	-@ $(RM) $(DOCS)/*.html
	-@ $(RM) $(DOCS)/*~