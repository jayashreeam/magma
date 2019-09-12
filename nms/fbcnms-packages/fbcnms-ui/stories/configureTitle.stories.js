/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import ConfigureTitle from '../components/ConfigureTitle.react';
import React from 'react';
import {storiesOf} from '@storybook/react';

storiesOf('ConfigureTitle', module).add('default', () => (
  <ConfigureTitle
    title={'This is a title'}
    subtitle={'This is a descriptive subtitle'}
  />
));
