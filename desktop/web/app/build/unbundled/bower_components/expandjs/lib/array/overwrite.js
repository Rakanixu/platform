!function(){"use strict";var r=require("../assert/assertArgument"),e=require("../array/concat"),t=require("../collection/reduce"),a=require("../tester/isArray"),n=require("../tester/isArrayable");module.exports=function(u,i){r(a(u),1,"Array"),r(n(i),2,"Arrayable");var s=u.length!==i.length||t(u,function(r,e,t){return r||e!==i[t]});return s&&Array.prototype.splice.apply(u,e([0,u.length],i)),u}}();