!function(){"use strict";var r=require("lodash"),e=require("../assert/assertArgument"),t=require("../caster/toArray");module.exports=function(a){return e(a=t(a),1,"Arrayable"),r.compact(a)}}();