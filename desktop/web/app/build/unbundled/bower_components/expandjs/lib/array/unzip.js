!function(){"use strict";var r=require("lodash"),e=require("../assert/assertArgument"),t=require("../caster/toArray");module.exports=function(u){return e(u=t(u),1,"Arrayable"),r.unzip(u)}}();