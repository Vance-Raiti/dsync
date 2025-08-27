#!/bin/bash

set -euxo pipefail

SPEC="dsync.spec"
NAME=$(rpmspec -q --qf "%{name}" $SPEC)
VERSION=$(rpmspec -q --qf "%{version}" $SPEC)
REV=$(git rev-parse --short HEAD)

mkdir -p ~/rpmbuild/{BUILD,RPMS,SOURCES,SPECS,SRPMS}
cp src/* ~/rpmbuild/SOURCES

git archive --format=tar.gz --prefix="${NAME}-${VERSION}/" -o ~/rpmbuild/SOURCES/${NAME}-${VERSION}.tar.gz HEAD
cp dsync.spec ~/rpmbuild/SPECS/

rpmbuild -ba --define "REV $REV" $SPEC
rm -r rpms
mv ~/rpmbuild/RPMS rpms
rm -r ~/rpmbuild
