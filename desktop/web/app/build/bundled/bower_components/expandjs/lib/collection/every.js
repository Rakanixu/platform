!function(){"use strict";var e=require("lodash"),r=require("../assert/assertArgument"),t=require("../tester/isCollection"),i=require("../tester/isPredicate"),s=require("../caster/toArray");module.exports=function(u,o,n){return r(t(u=s(u)||u),1,"Arrayable or Object"),r(i(o),2,"Function, Object or string"),e.every(u,o,n)}}();