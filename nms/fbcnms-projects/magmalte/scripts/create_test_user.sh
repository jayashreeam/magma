#!/usr/bin/env bash
#
# Copyright 2004-present Facebook. All Rights Reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

cd "$(dirname "$0")/.." && docker-compose run magmalte yarn run setAdminPassword admin@magma.test password1234
