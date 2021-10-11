import React from 'react';
import {auth, getParams} from './auth';

export default () => {
  var token;
  var params = getParams({location: window.location})

  // error
  if (params['error_description']) {
      return <>
          <p style={{'max-width': '100%'}}>{params['error_description']}</p>
      </>
  }

  if (params['access_token']) {
      // TODO: store in state or in a persistence storage
      // TODO: you should validate state and nonce here
      token = params['access_token'];
  }

  if(!token) {
      return auth();
  }

  return <>
    <h3>Protected</h3>
    <p>Here, we would load some data via a non-public api, or whatever.</p>
    <p>Token: {token}</p>
  </>
}