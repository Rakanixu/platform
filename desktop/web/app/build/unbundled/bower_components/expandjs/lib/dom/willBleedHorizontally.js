!function(){"use strict";var e=require("../assert/assertArgument"),t=require("../dom/getWidth"),r=require("../tester/isObject");module.exports=function(i,s){return e(r(i),1,"Object"),e(r(s),2,"Object"),s.left+i.width+s.right>t()}}();