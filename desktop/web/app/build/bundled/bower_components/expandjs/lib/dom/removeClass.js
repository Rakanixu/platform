!function(){"use strict";var e=require("../assert/assertArgument"),r=require("../tester/isElement"),t=require("../tester/isString"),s=require("../tester/isVoid");module.exports=function(i,n){return e(s(i)||r(i),1,"Element"),e(s(n)||t(n),2,"string"),i&&n&&i.classList.remove(n),i}}();