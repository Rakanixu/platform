!function(){"use strict";var r=require("../assert/assertArgument"),e=require("../tester/isString");module.exports=function(t,s){return r(e(t,!0),1,"string"),r(e(s,!0),2,"string"),"Basic "+new Buffer(t+":"+s).toString("base64")}}();