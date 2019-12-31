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

function __sysadmin_sk_build() {
  echo "Starting building Sysadmin SK"

  CGO_ENABLED=0 go build -installsuffix cgo -ldflags "$GO_LDFLAGS" -o "bin/sysadmin-sk"
}

# =============================================================================
#  Main Script starts from here
# =============================================================================
__sysadmin_sk_build