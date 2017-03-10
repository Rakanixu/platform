'use strict';

window.Auth = (function() {
  var _profile = {};
  var _customHeaders;
  var _lock = new Auth0Lock('5OCJYuTq5Dog960c3lfVEsBlquDX9Ka2', 'kazoup.eu.auth0.com',{
    auth: {
      redirect : false,
      responseType: 'token',
      redirectUrl: location.href,
      sso:true
    },
    mustAcceptTerms: true,
    languageDictionary: {
      title: "Kazoup",
      signUpTerms: "I agree to the <a href='http://www.kazoup.com/legal/kazoup-term-of-service/' target='_new'>terms of service</a> and <a href='http://www.kazoup.com/legal/privacy-policy/' target='_new'>privacy policy</a>."
    },
    theme: {
      logo: 'https://app.kazoup.io/src/static/kazoup-logo-small.png',
      primaryColor: '#2CB4D9',
  	},  
    closable : false
  });

  function _loadSearchPage() {
    document.querySelector('kazoup-app').routeTo(
      MapBehaviorImp.properties.pagesMap.value().search.path,
      {
        index: Auth.getMD5UserId()
      }
    );
  }

  function _authenticateUser(code) {
    return qwest.post('https://kazoup.eu.auth0.com/delegation', {
      grant_type: 'urn:ietf:params:oauth:grant-type:jwt-bearer',
      target: '6zIDm8InhbTRp1bL2C4m1TK4Llr4arTy',
      client_id: '5OCJYuTq5Dog960c3lfVEsBlquDX9Ka2',
      scope: 'openid roles',
      api_type: 'app',
      id_token: code
    }, {
      headers: {
        'Cache-Control': ''
      },
      cache:true
    });
  }

  _lock.on("authenticated", function(authResult) {
    _lock.getProfile(authResult.idToken, function(error, profile) {
      if (error) {
        return;
      }
      _profile = profile;
      localStorage.setItem('profile', JSON.stringify(profile));

      window.intercomSettings = {
        name: _profile.name,
        email: _profile.email,
        user_id :_profile.user_id,
        user_hash: _profile.intercom_hash,
        created_at: _profile.created_at
      };
      window.Intercom('update', window.intercomSettings);

      window.dispatchEvent(new CustomEvent('init-notify-messages', {
        detail: _profile
      }));

      localStorage.setItem('id_token', authResult.idToken);

      _authenticateUser(authResult.idToken).then(function(xhr, response) {
        localStorage.setItem('token', response.id_token);
        _customHeaders = {
          'Authorization': response.id_token
        };

        _lock.hide();
        _loadSearchPage();
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
      if (_.isEmpty(_customHeaders)) {
        // https://github.com/Polymer/polymer/issues/4012
        // Safe way to redirect to login if unauthorize
        return {
          Authorization: localStorage.getItem('token') ?
            localStorage.getItem('token') :
            'not_valid_bitch'
        }
      }

      return _customHeaders;
    },
    getJWT: function() {
      return localStorage.getItem('token');
    },
    getProfile: function() {
      if (_.isEmpty(_profile)) {
        return JSON.parse(localStorage.getItem('profile'));
      }

      return _profile;
    },
    getUserId: function() {
      try {
        return this.getProfile().user_id;
      } catch(e) {
        return '';
      }
    },
    getMD5UserId: function() {
      var profile = this.getProfile();

      if (!_.isEmpty(profile) && profile.user_id) {
        return md5(profile.user_id);
      }

      return '';
    },
    clear: function() {
      localStorage.removeItem('profile');
      localStorage.removeItem('token');
      localStorage.removeItem('id_token');

      _profile = {};
      _customHeaders = {};
    }
  }
}());

