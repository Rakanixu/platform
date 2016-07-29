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
    var http      = require('http'),
        https     = require('https'),
        adapters  = {'https': https, 'https:': https},
        location  = global.location || {},
        XP        = global.XP || require('expandjs'),
        XPEmitter = global.XPEmitter || require('xp-emitter');

    /*********************************************************************/

    /**
     * A class used to perform XHR requests.
     *
     * @class XPRequest
     * @description A class used to perform XHR requests
     * @extends XPEmitter
     */
    module.exports = global.XPRequest = new XP.Class('XPRequest', {

        // EXTENDS
        extends: XPEmitter,

        /*********************************************************************/

        /**
         * Emitted when a chunk of data is received.
         *
         * @event chunk
         * @param {Buffer | string} chunk
         * @param {Object} emitter
         */

        /**
         * Emitted when data are received.
         *
         * @event data
         * @param {*} data
         * @param {Object} emitter
         */

        /**
         * Emitted when an error is received.
         *
         * @event error
         * @param {string} error
         * @param {Object} emitter
         */

        /**
         * Emitted when a response is received.
         *
         * @event response
         * @param {string} url
         * @param {Object} emitter
         */

        /**
         * Emitted when the request state changes.
         *
         * @event state
         * @param {string} state
         * @param {Object} emitter
         */

        /**
         * Emitted when the request is submitted.
         *
         * @event submit
         * @param {Buffer | string} data
         * @param {Object} emitter
         */

        /*********************************************************************/

        /**
         * @constructs
         * @param {Object | string} options The request url or options.
         *   @param {string} [options.contentType] A shortcut for the "Content-Type" header.
         *   @param {string} [options.encoding] The response encoding.
         *   @param {Object} [options.headers] An object containing request headers.
         *   @param {string} [options.hostname] The request hostname, usable in alternative to url.
         *   @param {number} [options.keepAlive = 0] How often to submit TCP KeepAlive packets over sockets being kept alive.
         *   @param {string} [options.method = "GET"] A string specifying the HTTP request method.
         *   @param {string} [options.path] The request path, usable in alternative to url.
         *   @param {number} [options.port] The request port, usable in alternative to url.
         *   @param {number} [options.protocol = "http:"] The request protocol, usable in alternative to url.
         *   @param {string} [options.responseType] The type of data expected back from the server.
         *   @param {string} [options.url] The request url.
         */
        initialize: {
            promise: true,
            value: function (options, resolver) {

                // Vars
                var self    = this,
                    adaptee = null;

                // Super
                XPEmitter.call(self);

                // Overriding
                if (!XP.isObject(options)) { options = {url: options}; }
                if (!XP.isFalsy(options.url)) { XP.assign(options, XP.pick(XP.parseURL(options.url), ['hostname', 'path', 'port', 'protocol'])); }

                // Setting
                self.chunks       = [];
                self.state        = 'idle';
                self.options      = options;
                self.contentType  = self.options.contentType || null;
                self.encoding     = self.options.encoding || null;
                self.headers      = self.options.headers || {};
                self.hostname     = self.options.hostname || location.hostname || null;
                self.keepAlive    = self.options.keepAlive || 0;
                self.method       = self.options.method || 'GET';
                self.path         = self.options.path || null;
                self.port         = self.options.port || (!self.options.hostname && XP.toNumber(location.port)) || null;
                self.protocol     = self.options.protocol || (!self.options.hostname && location.protocol) || 'http:';
                self.responseType = self.options.responseType || null;
                self.url          = XP.toURL({protocol: self.protocol, hostname: self.hostname, port: self.port, pathname: self.pathname, search: self.search});

                // Adapting
                adaptee = (adapters[self.protocol] || http).request({
                    headers: self.headers,
                    hostname: self.hostname,
                    keepAlive: self.keepAlive > 0,
                    keepAliveMsecs: self.keepAlive,
                    method: self.method,
                    path: self.path,
                    port: self.port,
                    protocol: self.protocol,
                    withCredentials: false
                });

                // Listening
                adaptee.on('error', self._handleError.bind(self, resolver));
                adaptee.on('response', self._handleResponse.bind(self, resolver));

                // Overriding
                self.abort  = self.abort.bind(self, adaptee);
                self.header = self.header.bind(self, adaptee);
                self.submit = self.submit.bind(self, adaptee);
            }
        },

        /*********************************************************************/

        /**
         * Aborts the request.
         *
         * @method abort
         * @returns {Object}
         */
        abort: function (adaptee) {

            // Vars
            var self = this;

            // Checking
            if (self.tsAbort) { return self; }

            // Aborting
            adaptee.abort();

            // Setting
            self.state   = 'aborted';
            self.tsAbort = Date.now();

            return self;
        },

        /**
         * Get or set a header.
         *
         * @method header
         * @param {string} name
         * @param {number | string} [value]
         * @returns {number | string}
         */
        header: function (adaptee, name, value) {

            // Asserting
            XP.assertArgument(XP.isString(name, true), 1, 'string');
            XP.assertArgument(XP.isVoid(value) || XP.isFalse(value) || XP.isInput(value, true), 2, 'string');

            // Vars
            var self = this;

            // Getting
            if (!XP.isDefined(value) || self.state !== 'idle') { return self.headers[name]; }

            // Setting
            if (value) { adaptee.setHeader(name, self.headers[name] = value); return value; }

            // Removing
            adaptee.removeHeader(name);

            // Deleting
            delete self.headers[name];
        },

        /**
         * Submits the request, using data for the request body.
         *
         * @method submit
         * @param {*} [data]
         * @param {Function} [resolver]
         * @returns {Promise}
         */
        submit: {
            promise: true,
            value: function (adaptee, data, resolver) {

                // Asserting
                XP.assertArgument(XP.isVoid(resolver) || XP.isFunction(resolver), 2, 'Function');

                // Vars
                var self = this;

                // Checking
                if (self.tsSubmit) { return self; }

                // Serializing
                if (self.method === 'GET') { data = undefined; }
                if (self.method !== 'GET') { data = (XP.isInput(data, true) || XP.isBuffer(data) ? data : (XP.isCollection(data) ? XP.toJSON(data) : undefined)); }

                // Catching
                self.catch(function (err) { resolver(err, null); });
                self.then(function (data) { resolver(null, data); });

                // Ending
                adaptee.end(data);

                // Setting
                self.state    = 'pending';
                self.tsSubmit = Date.now();

                // Emitting
                self.emit('submit', data, self);
            }
        },

        /*********************************************************************/

        /**
         * TODO DOC
         *
         * @property chunks
         * @type Array
         */
        chunks: {
            set: function (val) { return this.chunks || val; },
            validate: function (val) { return !XP.isArray(val) && 'Array'; }
        },

        /**
         * TODO DOC
         *
         * @property contentType
         * @type string
         */
        contentType: {
            set: function (val) { return XP.isDefined(this.contentType) ? this.contentType : val; },
            validate: function (val) { return !XP.isVoid(val) && !XP.isString(val, true) && 'string'; }
        },

        /**
         * TODO DOC
         *
         * @property data
         * @type *
         * @readonly
         */
        data: {
            set: function (val) { return XP.isDefined(this.data) ? this.data : val; }
        },

        /**
         * TODO DOC
         *
         * @property encoding
         * @type string
         */
        encoding: {
            set: function (val) { return XP.isDefined(this.encoding) ? this.encoding : val; },
            validate: function (val) { return !XP.isVoid(val) && !XP.isString(val, true) && 'string'; }
        },

        /**
         * TODO DOC
         *
         * @property error
         * @type string
         * @readonly
         */
        error: {
            set: function (val) { return XP.isDefined(this.error) ? this.error : val; },
            validate: function (val) { return !XP.isVoid(val) && !XP.isString(val) && 'string'; }
        },

        /**
         * TODO DOC
         *
         * @property headers
         * @type Object
         */
        headers: {
            set: function (val) { return this.headers || (XP.isObject(val) && XP.cloneDeep(val)); },
            then: function () { if (this.contentType) { this.headers['content-type'] = this.contentType; } },
            validate: function (val) { return !XP.isObject(val) && 'Object'; }
        },

        /**
         * TODO DOC
         *
         * @property hostname
         * @type string
         */
        hostname: {
            set: function (val) { return this.hostname || val; },
            validate: function (val) { return !XP.isString(val, true) && 'string'; }
        },

        /**
         * TODO DOC
         *
         * @property keepAlive
         * @type number
         * @default 0
         */
        keepAlive: {
            set: function (val) { return XP.isDefined(this.keepAlive) ? this.keepAlive : val; },
            validate: function (val) { return !XP.isInt(val, true) && 'number'; }
        },

        /**
         * TODO DOC
         *
         * @property method
         * @type string
         * @default "GET"
         */
        method: {
            set: function (val) { return this.method || XP.upperCase(val); },
            validate: function (val) { return !XP.isString(val, true) && 'string'; }
        },

        /**
         * TODO DOC
         *
         * @property path
         * @type string
         */
        path: {
            set: function (val) { return XP.isDefined(this.path) ? this.path : val; },
            then: function (post) { var parts = XP.split(post, '?', true); this.pathname = parts[0] || null; this.search = parts[1] ? '?' + parts[1] : null; },
            validate: function (val) { return !XP.isVoid(val) && !XP.isString(val, true) && 'string'; }
        },

        /**
         * TODO DOC
         *
         * @property pathname
         * @type string
         */
        pathname: {
            set: function (val) { return XP.isDefined(this.path) ? this.path : val; },
            validate: function (val) { return !XP.isVoid(val) && !XP.isString(val, true) && 'string'; }
        },

        /**
         * TODO DOC
         *
         * @property port
         * @type number
         */
        port: {
            set: function (val) { return XP.isDefined(this.port) ? this.port : val; },
            validate: function (val) { return !XP.isVoid(val) && !XP.isInt(val, true) && 'number'; }
        },

        /**
         * TODO DOC
         *
         * @property protocol
         * @type string
         */
        protocol: {
            set: function (val) { return this.protocol || val; },
            validate: function (val) { return !XP.isString(val, true) && 'string'; }
        },

        /**
         * TODO DOC
         *
         * @property responseType
         * @type string
         */
        responseType: {
            set: function (val) { return XP.isDefined(this.responseType) ? this.responseType : val; },
            validate: function (val) { return !XP.isVoid(val) && this.responseTypes.indexOf(val) < 0 && 'string'; }
        },

        /**
         * TODO DOC
         *
         * @property responseTypes
         * @type Array
         * @default ["json"]
         * @readonly
         */
        responseTypes: {
            frozen: true,
            writable: false,
            value: ['json']
        },

        /**
         * TODO DOC
         *
         * @property pathname
         * @type string
         */
        search: {
            set: function (val) { return XP.isDefined(this.search) ? this.search : val; },
            validate: function (val) { return !XP.isVoid(val) && !XP.isString(val, true) && 'string'; }
        },

        /**
         * TODO DOC
         *
         * @property state
         * @type string
         * @readonly
         */
        state: {
            set: function (val) { return val; },
            then: function (post) { this.emit('state', post, this); },
            validate: function (val) { return this.states.indexOf(val) < 0 && 'string'; }
        },

        /**
         * TODO DOC
         *
         * @property states
         * @type Array
         * @default ["aborted", "failed", "idle", "pending", "received", "receiving"]
         * @readonly
         */
        states: {
            frozen: true,
            writable: false,
            value: ['aborted', 'failed', 'idle', 'pending', 'received', 'receiving']
        },

        /**
         * TODO DOC
         *
         * @property statusCode
         * @type number
         * @readonly
         */
        statusCode: {
            set: function (val) { return this.statusCode || val; },
            validate: function (val) { return !XP.isInt(val, true) && 'number'; }
        },

        /**
         * TODO DOC
         *
         * @property statusMessage
         * @type string
         * @readonly
         */
        statusMessage: {
            set: function (val) { return this.statusMessage || val; },
            validate: function (val) { return !XP.isString(val) && 'string'; }
        },

        /**
         * TODO DOC
         *
         * @property tsAbort
         * @type number
         * @readonly
         */
        tsAbort: {
            set: function (val) { return this.tsAbort || val; },
            validate: function (val) { return !XP.isInt(val, true) && 'number'; }
        },

        /**
         * TODO DOC
         *
         * @property tsData
         * @type number
         * @readonly
         */
        tsData: {
            set: function (val) { return this.tsData || val; },
            validate: function (val) { return !XP.isInt(val, true) && 'number'; }
        },

        /**
         * TODO DOC
         *
         * @property tsResponse
         * @type number
         * @readonly
         */
        tsResponse: {
            set: function (val) { return this.tsResponse || val; },
            validate: function (val) { return !XP.isInt(val, true) && 'number'; }
        },

        /**
         * TODO DOC
         *
         * @property tsSubmit
         * @type number
         * @readonly
         */
        tsSubmit: {
            set: function (val) { return this.tsSubmit || val; },
            validate: function (val) { return !XP.isInt(val, true) && 'number'; }
        },

        /**
         * TODO DOC
         *
         * @property url
         * @type string
         */
        url: {
            set: function (val) { return this.url || val; },
            validate: function (val) { return !XP.isString(val, true) && 'string'; }
        },

        /*********************************************************************/

        // HANDLER
        _handleData: function (chunk) {

            // Vars
            var self = this;

            // Setting
            self.chunks.push(chunk);

            // Emitting
            self.emit('chunk', chunk, self);
        },

        // HANDLER
        _handleEnd: function (resolver) {

            // Vars
            var self   = this,
                failed = self.statusCode >= 400,
                data   = failed ? null : XP.join(self.chunks),
                error  = failed ? XP.join(self.chunks).toString() || self.statusMessage : null,
                state  = failed ? 'failed' : 'received',
                type   = failed ? 'error' : self.responseType;

            // Parsing
            if (type === 'json') { data = XP.parseJSON(data.toString()); }

            // Setting
            self.data   = data;
            self.error  = error;
            self.state  = state;
            self.tsData = Date.now();

            // Resolving
            resolver(self.error, self.data);

            // Emitting
            self.emit(failed ? 'error' : 'data', failed ? self.error : self.data, self);
        },

        // HANDLER
        _handleError: function (resolver, error) {

            // Vars
            var self = this;

            // Setting
            self.error = error.message || 'Unknown';
            self.state = 'failed';

            // Resolving
            resolver(self.error, null);

            // Emitting
            self.emit('error', self.error, self);
        },

        // HANDLER
        _handleResponse: function (resolver, response) {

            // Vars
            var self = this;

            // Encoding
            if (self.encoding) { response.setEncoding(self.encoding); }

            // Setting
            self.state         = 'receiving';
            self.statusCode    = response.statusCode;
            self.statusMessage = response.statusMessage || http.STATUS_CODES[self.statusCode] || 'Unknown';
            self.tsResponse    = Date.now();

            // Listening
            response.on('data', self._handleData.bind(self));
            response.on('end', self._handleEnd.bind(self, resolver));

            // Emitting
            self.emit('response', self.url, self);
        }
    });

}(typeof window !== "undefined" ? window : global));
