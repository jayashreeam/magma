/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @format
 */
'use strict';

import 'babel-polyfill';
import {BrowserRouter} from 'react-router-dom';
import ReactDOM from 'react-dom';
import React from 'react';
import Main from './components/Main';

import {} from './common/axiosConfig';

ReactDOM.render(
  <BrowserRouter>
    <Main />
  </BrowserRouter>,
  document.getElementById('root'),
);
