!function(){"use strict";var r=require("../assert/assertArgument"),e=require("../tester/isString"),t=require("../tester/isVoid"),i=require("../string/trim");module.exports=function(s){return r(t(s)||e(s),1,"string"),s?i(s.replace(/[ ]+/g," ")):""}}();