'use strict';

window.Kazoup = window.Kazoup || {};
Kazoup.SharedData = Kazoup.SharedData || {};

Kazoup.SharedData = {
  config: {},
  currentUser: {},
  getConfig: function() {
    var promise = pinkySwear();
    var _self = this;

    // Check if config already available to avoid request AND
    // APPLIANCE_IS_REGISTER is false, because state could change in meantime
    if (_.isEmpty(this.config) || !this.config.APPLIANCE_IS_REGISTERED) {
      qwest.get(
        'http://localhost:8080/v2/settings/public/read',
        {},
        {
          dataType: 'json',
          cache: true,
          headers: {
            'Accept': 'application/json',
            'Content-type': 'application/json'
          }
        }
      ).then(function(xhr, response) {
        _self.config = response;

        promise(true, [_self.config]);
        _self._fireEvent('config-ready');
      }).catch(function(xhr, response, e) {
        promise(false, [e]);
      });
    } else {
      promise(true, [this.config]);
      _self._fireEvent('config-ready');
    }

    return promise;
  },
  getCurrentUser: function() {
    var promise = pinkySwear();
    var _self = this;

    // Check if currentUser is already available to avoid request
    // In any case promise should be resolved or rejected
    if (_.isEmpty(this.currentUser)) {
      qwest.get(
        'http://localhost:8080/v2/settings/user/read',
        {},
        {
          dataType: 'json',
          cache: true,
          headers: {
            Accept: 'application/json',
            Authorization: 'JWT ' + localStorage.getItem('authToken')
          }
        }
      ).then(function(xhr, response) {
        _self.currentUser = response;
        _self._fireEvent('user-ready');
        promise(true, [_self.currentUser]);
      }).catch(function(xhr, response) {
        // unauthorize
        if (xhr.status === 401) {
          localStorage.removeItem('authToken');
          _self._fireEvent('user-ready');
        }
        promise(false, [response]);
      });
    } else {
      _self._fireEvent('user-ready');
      promise(true, [this.currentUser]);
    }

    return promise;
  },
  _fireEvent: function(eventName, obj) {
    // Support for IE
    // https://developer.mozilla.org/en-US/docs/Web/Guide/Events/Creating_and_triggering_events
    this._customEvent = document.createEvent('Event');
    this._customEvent.initEvent(eventName, true, true);
    this._customEvent.data = obj;

    window.dispatchEvent(this._customEvent);
  },
  fireEvent: function(eventName, obj) {
    this._fireEvent(eventName, obj);
  }
};
