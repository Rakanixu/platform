!function(){"use strict";var e=require("lodash"),r=require("../assert/assertArgument"),t=require("../tester/isCollection"),i=require("../tester/isFunction");module.exports=function(s,u,n){return r(t(s),1,"Arrayable or Object"),r(i(u),2,"Function"),e.forEachRight(s,u,n)}}();