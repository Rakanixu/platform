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
    var eventemitter3 = require('eventemitter3'),
        XP            = global.XP || require('expandjs');

    /*********************************************************************/

    /**
     * This class is used to provide event emitting functionalities.
     *
     * @class XPEmitter
     * @description This class is used to provide event emitting functionalities
     */
    module.exports = global.XPEmitter = new XP.Class('XPEmitter', {

        // EXTENDS
        extends: eventemitter3
    });

}(typeof window !== "undefined" ? window : global));