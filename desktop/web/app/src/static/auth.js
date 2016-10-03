'use strict';

window.Auth = (function() {
  var _profile = {};
  var _customHeaders;
  var _lock = new Auth0Lock('5OCJYuTq5Dog960c3lfVEsBlquDX9Ka2', 'kazoup.eu.auth0.com',{
    auth: {
      redirect :false
    }
  });

  function _loadWebApp() {
    var e = document.createEvent('Event');
    e.initEvent('web-app-activated', true, true);
    e.data = {}; // Set params if needed
    window.dispatchEvent(e);
    document.querySelector('#ironPages').selected = 1;
  }

  function _authenticateUser(code) {
    return qwest.post('https://kazoup.eu.auth0.com/delegation', {
      grant_type: 'urn:ietf:params:oauth:grant-type:jwt-bearer',
      target: '6zIDm8InhbTRp1bL2C4m1TK4Llr4arTy',
      client_id: '5OCJYuTq5Dog960c3lfVEsBlquDX9Ka2',
      scope: 'openid',
      api_type: 'app',
      id_token: code
    });
  }

  _lock.on("authenticated", function(authResult) {
    _lock.getProfile(authResult.idToken, function(error, profile) {
      if (error) {
        return;
      }
      _profile = profile;

      localStorage.setItem('id_token', authResult.idToken);

      _authenticateUser(authResult.idToken).then(function(xhr, response) {
        localStorage.setItem('token', response.id_token);
        _customHeaders = {
          Token: response.id_token
        };

        _loadWebApp();
        _lock.hide();
      }).catch(function(e, xhr, response) {

      });
    });
  });

  return {
    showLogin: function() {
      _lock.show();
    },
    hideLogin: function() {
      _lock.hide();
    },
    setHeaders: function(headers) {
      _customHeaders = headers;
    },
    getHeaders: function() {
      return _customHeaders;
    },
    getProfile: function() {
      return _profile;
    }
  }
}());

