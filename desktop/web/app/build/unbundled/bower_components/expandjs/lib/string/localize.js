!function(){"use strict";var e=require("../assert/assertArgument"),r=require("../tester/isArrayable"),t=require("../tester/isCollection"),i=require("../tester/isObject"),u=require("../tester/isString"),s=require("../tester/isVoid"),n=require("../collection/map"),o=require("../object/mapValues"),c=require("../object/value");module.exports=function a(q,l){return e(s(q)||u(q)||t(q),1,"Array, Object or string"),e(s(l)||i(l),2,"Object"),q&&l?u(q)?c(l,q,q):r(q)?n(q,function(e){return a(l,e)}):i(q)?o(q,function(e,r){return a(l,r)}):void 0:q||""}}();