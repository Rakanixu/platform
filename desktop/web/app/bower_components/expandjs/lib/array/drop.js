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

    var _              = require('lodash'),
        assertArgument = require('../assert/assertArgument'),
        isVoid         = require('../tester/isVoid'),
        isIndex        = require('../tester/isIndex'),
        toArray        = require('../caster/toArray');

    /**
     * Creates a slice of `array` with `n` elements dropped from the beginning.
     *
     * ```js
     * XP.drop([1, 2, 3]);
     * // => [2, 3]
     *
     * XP.drop([1, 2, 3], 2);
     * // => [3]
     *
     * XP.drop([1, 2, 3], 5);
     * // => []
     *
     * XP.drop([1, 2, 3], 0);
     * // => [1, 2, 3]
     * ```
     *
     * @function drop
     * @param {Array} array The array to query.
     * @param {number} [n = 1] The number of elements to drop.
     * @returns {Array} Returns the slice of `array`.
     */
    module.exports = function drop(array, n) {
        assertArgument(array = toArray(array), 1, 'Arrayable');
        assertArgument(isVoid(n) || isIndex(n), 2, 'number');
        return _.drop(array, n);
    };

}());