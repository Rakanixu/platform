!function(){"use strict";var e=require("../assert/assertArgument"),t=require("../tester/isElement"),r=require("../tester/isString");module.exports=function(s,i){return e(t(s),1,"Element"),e(r(i,!0),2,"string"),s.hasAttribute(i)}}();