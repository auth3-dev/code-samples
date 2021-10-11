// configure
const clientID = 'tutorialgolangserver%405ddafaf2-044e-46a4-9b79-e6d002687170.developer.apps.auth3.dev';
const auth3ProjectId = '5ddafaf2-044e-46a4-9b79-e6d002687170';

// do not edit
const endpoint = 'https://' + auth3ProjectId + '.as.auth3.dev/';
const authEndpoint = endpoint + 'oauth2/auth';

export const auth = () => {
  var state = 'ntwfZrndRp'; // TODO generate and store a valid state and verify it later 
  var nonce = 'eobpHczU9G'; // TODO generate and store a valid nonce and verify it later

  window.location.href = authEndpoint + '?' +
    'response_type=token&' + // set: id_token+token if you also need an ID Token
    'client_id=' + clientID + '&' + 
    'redirect_uri=' + window.location.href + '&' +
    'state=' + state + '&' +
    'nonce=' + nonce;

  // you might need to change this if you don't use React.
  return <>Logging in...</>;
}

export const getParams = () => {
    var params = {};

    if (window.location.search.length) {
        var query = window.location.search.substring(1);
        params = extractParams({query, params});
    }

    if (window.location.hash.length) {
        var query = window.location.hash.substring(1);
        params = extractParams({query, params});
    }

    return params;
}

const extractParams = ({query, params = {}}) => {
    var pairs = query.split("&");
    console.log(pairs);

    for (var i = 0; i < pairs.length; i++) {
        var pair = pairs[i].split("=");
        console.log(pairs);
        params[pair[0]] = decodeURIComponent(pair[1].replace('+', ' '));
    }

    return params;
}