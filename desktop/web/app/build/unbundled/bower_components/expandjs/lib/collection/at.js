!function(){"use strict";var r=require("lodash"),e=require("../assert/assertArgument"),t=require("../tester/isCollection"),a=require("../caster/toArray");module.exports=function(s,u){return e(t(s=a(s)||s),1,"Arrayable or Object"),e(u=a(u),2,"Arrayable"),r.at(s,u)}}();