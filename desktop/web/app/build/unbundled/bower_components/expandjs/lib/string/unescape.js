!function(){"use strict";var e=require("lodash"),r=require("../assert/assertArgument"),t=require("../tester/isString"),s=require("../tester/isVoid");module.exports=function(i){return r(s(i)||t(i),1,"string"),i?e.unescape(i):""}}();