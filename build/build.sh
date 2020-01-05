#!/usr/bin/env bash
##
# This file is part of the Sysadmin Sidekick Toolkit (Sysadmin-SK) (https://github.com/raffs/sysadmin-sk).
# Copyright (c) 2019 Rafael Oliveira Silva
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, version 3.
#
# This program is distributed in the hope that it will be useful, but
# WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
# General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program. If not, see <http://www.gnu.org/licenses/>.
##

export REPO_PATH="github.com/raffs/sysadmin-sk"
export GO_LDFLAGS="-s" # for building without symbols for debugging.

[[ ! "${BIN_OUTPUT}" ]] && BIN_OUTPUT="${HOME}/.bin/sysadmin-sk"

function __sysadmin_sk_build() {
  echo "Starting building Sysadmin SK"

  if ! CGO_ENABLED=0 go build \
          -gcflags -m \
          -installsuffix cgo \
          -ldflags "$GO_LDFLAGS" \
          -o "${BIN_OUTPUT}"
  then
    echo
    echo "Failed to compile the binary, please see message above"
    exit 1
  else
    echo
    echo "Successfully finish the build"
  fi
}

# =============================================================================
#  Main Script starts from here
# =============================================================================
__sysadmin_sk_build
