!function(){"use strict";var e=require("../assert/assertArgument"),r=require("../tester/isFinite"),t=require("../tester/isNumber");module.exports=function(u,i,n,s){return e(t(u),1,"number"),e(r(i),2,"number"),e(r(n),3,"number"),((s?Math.max(Math.min(u,n),i):u)-i)/(n-i)}}();