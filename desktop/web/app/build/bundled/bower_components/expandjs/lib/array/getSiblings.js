!function(){"use strict";var r=require("../assert/assertArgument"),e=require("../array/concat"),a=require("../array/getNext"),t=require("../array/getPrevious"),u=require("../caster/toArray");module.exports=function(i,s){r(i=u(i),1,"Arrayable");var n=t(i,s),o=a(i,s);return e(n?[n]:[],o?[o]:[])}}();