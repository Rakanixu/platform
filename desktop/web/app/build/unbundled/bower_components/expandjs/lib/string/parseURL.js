!function(){"use strict";var r=require("../assert/assertArgument"),e=require("../object/assign"),t=require("../tester/isString"),i=require("../tester/isVoid"),s=require("../caster/toNumber"),u=require("url");module.exports=function(n,o,a){r(i(n)||t(n),1,"string");var l=n?u.parse(n,!!o,!!a):null;if(l)return e(l,{port:s(l.port)||null})}}();