#!/bin/bash

BASEDIR=$(dirname $(pwd))
PACKAGENAME=$1
echo Building RPMs...
VERSION?=v0.0.2
RELEASE?=rel01
BIN_NAME=${PACKAGENAME}_${VERSION}_${RELEASE}
echo "BIN_NAME: $BIN_NAME"


mkdir -p ${BASEDIR}/rpmbuild/SOURCES/${BIN_NAME}
cp ../$(PACKAGENAME)*  $(BASEDIR)/rpmbuild/SOURCES/$(BIN_NAME)
cp ../package/*  $(BASEDIR)/rpmbuild/SOURCES/$(BIN_NAME)
cd $(BASEDIR)/rpmbuild/SOURCES && tar cvfz $(BIN_NAME).tar.gz $(BIN_NAME)
rpmbuild --define '_rpmfilename $(BIN_NAME).rpm' \
         --define '_topdir $(BASEDIR)/rpmbuild' \
         --define 'version $(VERSION)' \
         --define 'release $(RELEASE)' \
         -ba --clean  $(BASEDIR)/package/$(REPO).spec
