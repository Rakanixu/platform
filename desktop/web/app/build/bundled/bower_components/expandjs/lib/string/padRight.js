!function(){"use strict";var e=require("lodash"),r=require("../assert/assertArgument"),t=require("../tester/isVoid"),i=require("../tester/isIndex"),s=require("../tester/isString");module.exports=function(u,n,o){return r(t(u)||s(u),1,"string"),r(t(n)||i(n),2,"number"),r(t(o)||s(o),3,"string"),e.padRight(u,n,o)}}();