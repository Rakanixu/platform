!function(e){"use strict";var t=require("../assert/assertArgument"),r=require("../tester/isElement"),n=require("../tester/isString");module.exports=function(i,s){return t(r(i),1,"Element"),t(n(s,!0),2,"string"),e.getComputedStyle(i)[s]}}("undefined"!=typeof window?window:global);