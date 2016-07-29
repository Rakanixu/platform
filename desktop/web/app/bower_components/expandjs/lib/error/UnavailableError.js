/*jslint browser: true, devel: true, node: true, ass: true, nomen: true, unparam: true, indent: 4 */

/**
 * @license
 * Copyright (c) 2015 The ExpandJS authors. All rights reserved.
 * This code may only be used under the BSD style license found at https://expandjs.github.io/LICENSE.txt
 * The complete set of authors may be found at https://expandjs.github.io/AUTHORS.txt
 * The complete set of contributors may be found at https://expandjs.github.io/CONTRIBUTORS.txt
 */
(function () {
    "use strict";

    var CustomError = require('../error/CustomError');

    /**
     * Creates a custom error with the `UnavailableError` type and a predefined message.
     *
     * ```js
     *  console.log(new UnavailableError('myVar'));
     *  // => UnavailableError{message: 'myVar is not available', stack: '...'}
     * ```
     *
     * @function UnavailableError
     * @param {string} key The key to be used in the error message
     * @param {number} [code] Additional status code.
     * @constructor
     */
    module.exports = function UnavailableError(key, code) {
        CustomError.call(this, 'UnavailableError', key + ' is not available', code);
    };

}());