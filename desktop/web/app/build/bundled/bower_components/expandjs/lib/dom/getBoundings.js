!function(){"use strict";var t=require("../assert/assertArgument"),e=require("../tester/isElement");module.exports=function(r){t(e(r),1,"Element");var i=r.getBoundingClientRect();return{bottom:i.bottom,height:i.height,left:i.left,right:i.right,top:i.top,width:i.width}}}();