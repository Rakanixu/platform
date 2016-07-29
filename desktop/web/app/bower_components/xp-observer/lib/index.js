/*jslint browser: true, devel: true, node: true, ass: true, nomen: true, unparam: true, indent: 4 */

/**
 * @license
 * Copyright (c) 2015 The ExpandJS authors. All rights reserved.
 * This code may only be used under the BSD style license found at https://expandjs.github.io/LICENSE.txt
 * The complete set of authors may be found at https://expandjs.github.io/AUTHORS.txt
 * The complete set of contributors may be found at https://expandjs.github.io/CONTRIBUTORS.txt
 */
(function (global) {
    "use strict";

    // Vars
    var XP       = global.XP || require('expandjs'),
        Observer = require('observe-js').ObjectObserver;

    /*********************************************************************/

    /**
     * This class is used to provide object observing functionality.
     *
     * @class XPObserver
     * @description This class is used to provide object observing functionality
     */
    module.exports = global.XPObserver = new XP.Class('XPObserver', {

        /**
         * @constructs
         * @param {Array | Function | Object} value
         * @param {Function} callback
         * @param {boolean} [deep = false]
         */
        initialize: function (value, callback, deep) {

            // Asserting
            XP.assertArgument(XP.isObservable(value), 1, 'Array, Function or Object');
            XP.assertArgument(XP.isFunction(callback), 2, 'Function');

            // Vars
            var self = this;

            // Setting
            self.value      = value;
            self.callback   = callback;
            self.deep       = deep;
            self._observers = [];

            // Observing
            self._addObserver(self.value);

            return self;
        },

        /*********************************************************************/

        /**
         * Disconnects the observer.
         *
         * @method disconnect
         * @returns {Object}
         */
        disconnect: function () {
            var self = this;
            self._removeObserver(self.value);
            return self;
        },

        /*********************************************************************/

        /**
         * Adds the observer for value.
         *
         * @method _addObserver
         * @param {Array | Function | Object} value
         * @param {Array | Object} [wrapper]
         * @returns {Object}
         * @private
         */
        _addObserver: {
            enumerable: false,
            value: function (value, wrapper) {

                // Asserting
                XP.assertArgument(XP.isObservable(value), 1, 'Array, Function or Object');
                XP.assertArgument(XP.isVoid(wrapper) || XP.isCollection(wrapper), 2, 'Array or Object');

                // Vars
                var self     = this,
                    observe  = function (sub) { return XP.isObservable(sub) ? self._addObserver(sub, value) : undefined; },
                    observer = !self._isObserved(value) && (!wrapper || XP.includesDeep(wrapper, value)) && self._connectObserver(new Observer(value));

                // Checking
                if (!observer) { return self; }

                // Adding
                if (value === self.value) { self._observer = observer; } else { XP.push(self._observers, observer); }
                if (self.deep && XP.isCollection(value)) { XP[XP.isArray(value) ? 'forEach' : 'forOwn'](value, observe); }

                return self;
            }
        },

        /**
         * Connects an observer.
         *
         * @method _connectObserver
         * @param {Object} observer
         * @returns {Object}
         * @private
         */
        _connectObserver: {
            enumerable: false,
            value: function (observer) {

                // Asserting
                XP.assertArgument(XP.isObject(observer), 1, 'Object');

                // Vars
                var self     = this,
                    value    = self._getObserved(observer),
                    callback = function (added, removed, changed, getOld) {

                        // Updating
                        XP.forEach(added,   function (sub)      { return XP.isObservable(sub) ? self._addObserver(sub, value) : undefined; });
                        XP.forEach(changed, function (sub, key) { return XP.isObservable(sub) ? self._removeObserver(getOld(key))._addObserver(sub, value) : undefined; });
                        XP.forEach(removed, function (sub, key) { return XP.isObservable(getOld(key)) ? self._removeObserver(getOld(key)) : undefined; });

                        return self.callback(self.value);
                    };

                // Opening
                observer.open(callback);

                return observer;
            }
        },

        /**
         * Returns the value of observer.
         *
         * @method _getObserved
         * @param {Object} observer
         * @returns {Array | Object}
         * @private
         */
        _getObserved: {
            enumerable: false,
            value: function (observer) {
                XP.assertArgument(XP.isObject(observer), 1, 'Object');
                return observer.value_;
            }
        },

        /**
         * Returns the observer of value.
         *
         * @method _getObserver
         * @param {Array | Function | Object} value
         * @returns {Object | undefined}
         * @private
         */
        _getObserver: {
            enumerable: false,
            value: function (value) {
                XP.assertArgument(XP.isObservable(value), 1, 'Array, Function or Object');
                return XP.find(this._observers, function (observer) { return observer.value_ === value; });
            }
        },

        /**
         * Returns true if value is observed.
         *
         * @method _isObserved
         * @param {Array | Function | Object} value
         * @returns {boolean}
         * @private
         */
        _isObserved: {
            enumerable: false,
            value: function (value) {
                XP.assertArgument(XP.isObservable(value), 1, 'Array, Function or Object');
                return value === this.value ? !!this._observer : !!this._getObserver(value);
            }
        },

        /**
         * Removes the observer of value.
         *
         * @method _removeObserver
         * @param {Array | Function | Object} value
         * @returns {Object}
         * @private
         */
        _removeObserver: {
            enumerable: false,
            value: function (value) {

                // Asserting
                XP.assertArgument(XP.isObservable(value), 1, 'Array, Function or Object');

                // Vars
                var self     = this,
                    observe  = function (sub) { return XP.isObservable(sub) ? self._removeObserver(sub) : undefined; },
                    observer = !XP.includesDeep(self.value, value) && self._getObserver(value);

                // Closing
                if (observer) { observer.close(); } else { return self; }

                // Removing
                if (value === self.value) { self._observer = observer; } else { XP.pull(self._observers, observer); }
                if (self.deep && XP.isCollection(value)) { XP[XP.isArray(value) ? 'forEach' : 'forOwn'](value, observe); }

                return self;
            }
        },

        /*********************************************************************/

        /**
         * TODO DOC
         *
         * @property callback
         * @type Function
         */
        callback: {
            set: function (val) { return this.callback || val; },
            validate: function (val) { return !XP.isFunction(val) && 'Function'; }
        },

        /**
         * TODO DOC
         *
         * @property deep
         * @type boolean
         */
        deep: {
            set: function (val) { return !!val; }
        },

        /**
         * TODO DOC
         *
         * @property value
         * @type Array | Function | Object
         */
        value: {
            set: function (val) { return this.value || val; },
            validate: function (val) { return !XP.isObservable(val) && 'Array, Function or Object'; }
        },

        /*********************************************************************/

        /**
         * TODO DOC
         *
         * @property _observer
         * @type Object
         * @private
         */
        _observer: {
            enumerable: false,
            set: function (val) { return this._observer || val; },
            validate: function (val) { return !XP.isObject(val) && 'Object'; }
        },

        /**
         * TODO DOC
         *
         * @property _observers
         * @type Array
         * @private
         */
        _observers: {
            enumerable: false,
            set: function (val) { return this._observers || val; },
            validate: function (val) { return !XP.isArray(val) && 'Array'; }
        }
    });

}(typeof window !== "undefined" ? window : global));
