!function(){"use strict";var e=require("../assert/assertArgument"),r=require("../tester/isElement"),t=require("../tester/isNode"),s=require("../tester/isVoid");module.exports=function(i,u){return e(s(i)||r(i),1,"Element"),e(s(u)||t(u),2,"Node"),i&&u&&i.insertBefore(u,i.firstChild),u}}();