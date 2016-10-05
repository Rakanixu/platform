'use strict';

window.Auth = (function() {
  var _profile = {};
  var _customHeaders;
  var _lock = new Auth0Lock('5OCJYuTq5Dog960c3lfVEsBlquDX9Ka2', 'kazoup.eu.auth0.com',{
    auth: {
      redirect : false,
      responseType: 'token',
      redirectUrl: location.href,
      sso:false
    },
    mustAcceptTerms: true,
    languageDictionary: {
      title: "Kazoup",
      signUpTerms: "I agree to the <a href='http://www.kazoup.com/legal/master-agreement/' target='_new'>terms of service</a> and <a href='http://www.kazoup.com/legal/privacy-policy/' target='_new'>privacy policy</a>."
    },
    theme: {
      logo: 'http://www.kazoup.com/img/kazoup-logo-small.png',
      primaryColor: '#2CB4D9',
  	},  
    closable : false
  });

  function _loadWebApp() {
    var e = document.createEvent('Event');
    e.initEvent('web-app-activated', true, true);
    e.data = {}; // Set params if needed
    window.dispatchEvent(e);
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

      // TODO: we need srv to hash the user_id and get it back to send it to Intercom
      // https://segment.com/docs/integrations/intercom/
      console.log(profile)
      window.intercomSettings = {
        name: _profile.name,
        email: _profile.email,
	user_id :_profile.user_id,
	user_hash: _profile.intercom_hash,
        created_at: _profile.created_at
      };
      window.Intercom('update', window.intercomSettings);

      localStorage.setItem('id_token', authResult.idToken);

      _authenticateUser(authResult.idToken).then(function(xhr, response) {
        localStorage.setItem('token', response.id_token);
        _customHeaders = {
          'Authorization': response.id_token
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
    },
    getUserId: function() {
      return _profile.user_id;
    },
    clear: function() {
      _profile = {};
      _customHeaders = {};
    }
  }
}());

