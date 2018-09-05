#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

DIR=$(dirname $(readlink -f $0))

rev=$(git rev-parse HEAD)

tar -C "${DIR}/../.." -zcf ~/rpmbuild/SOURCES/duoldapsync-${rev}.tar.gz --transform s/duoldapsync/duoldapsync-${rev}/ duoldapsync
rpmbuild -ba ${DIR}/duoldapsync.spec
