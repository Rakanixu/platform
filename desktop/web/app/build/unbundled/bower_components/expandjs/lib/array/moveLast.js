!function(){"use strict";var e=require("../assert/assertArgument"),r=require("../tester/isIndex"),t=require("../tester/isVoid"),i=require("../caster/toArray");module.exports=function(s,u,a){return e(s=i(s),1,"Arrayable"),e(r(u),2,"a positive number"),e(t(a)||r(a),3,"void or a positive number"),s.push.apply(s,s.splice(u,a)),s}}();