!function(){"use strict";var r=require("lodash"),e=require("../assert/assertArgument"),t=require("../tester/isCollection"),i=require("../tester/isPredicate"),s=require("../caster/toArray");module.exports=function(u,o,n){return e(t(u=s(u)||u),1,"Arrayable or Object"),e(i(o),2,"Function, Object or string"),r.filter(u,o,n)}}();