!function(){"use strict";var t=require("../tester/isFunction"),r=require("../tester/isVoid");module.exports=function(i,n,e){return r(i)?"null":JSON.stringify(i,function(i,e){var o=e&&e.toJSON?e.toJSON():e;return t(o)?o.toString():r(o)?n?void 0:null:e},e?"  ":void 0)}}();