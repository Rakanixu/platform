!function(){"use strict";var r=require("../error/ArgumentError"),e=require("../tester/isElement"),t=require("../tester/isFunction"),i=require("../tester/isString"),n=require("../tester/isVoid"),u=require("../function/mock"),o=require("../caster/toDOMPredicate");module.exports=function(s){if(e(s))return function(r){return r===s};if(t(s)||i(s,!0))return o(s);if(n(s)||i(s,!1))return u();throw new r(1,"Element, Function or string")}}();