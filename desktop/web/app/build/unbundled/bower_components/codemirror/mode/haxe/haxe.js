!function(e){"object"==typeof exports&&"object"==typeof module?e(require("../../lib/codemirror")):"function"==typeof define&&define.amd?define(["../../lib/codemirror"],e):e(CodeMirror)}(function(e){"use strict";e.defineMode("haxe",function(e,t){function n(e){return{type:e,style:"keyword"}}function r(e,t,n){return t.tokenize=n,n(e,t)}function a(e,t){for(var n,r=!1;null!=(n=e.next());){if(n==t&&!r)return!0;r=!r&&"\\"==n}}function i(e,t,n){return J=e,K=n,t}function o(e,t){var n=e.next();if('"'==n||"'"==n)return r(e,t,l(n));if(/[\[\]{}\(\),;\:\.]/.test(n))return i(n);if("0"==n&&e.eat(/x/i))return e.eatWhile(/[\da-f]/i),i("number","number");if(/\d/.test(n)||"-"==n&&e.eat(/\d/))return e.match(/^\d*(?:\.\d*(?!\.))?(?:[eE][+\-]?\d+)?/),i("number","number");if(t.reAllowed&&"~"==n&&e.eat(/\//))return a(e,"/"),e.eatWhile(/[gimsu]/),i("regexp","string-2");if("/"==n)return e.eat("*")?r(e,t,u):e.eat("/")?(e.skipToEnd(),i("comment","comment")):(e.eatWhile(ne),i("operator",null,e.current()));if("#"==n)return e.skipToEnd(),i("conditional","meta");if("@"==n)return e.eat(/:/),e.eatWhile(/[\w_]/),i("metadata","meta");if(ne.test(n))return e.eatWhile(ne),i("operator",null,e.current());var o;if(/[A-Z]/.test(n))return e.eatWhile(/[\w_<>]/),o=e.current(),i("type","variable-3",o);e.eatWhile(/[\w_]/);var o=e.current(),c=te.propertyIsEnumerable(o)&&te[o];return c&&t.kwAllowed?i(c.type,c.style,o):i("variable","variable",o)}function l(e){return function(t,n){return a(t,e)&&(n.tokenize=o),i("string","string")}}function u(e,t){for(var n,r=!1;n=e.next();){if("/"==n&&r){t.tokenize=o;break}r="*"==n}return i("comment","comment")}function c(e,t,n,r,a,i){this.indented=e,this.column=t,this.type=n,this.prev=a,this.info=i,null!=r&&(this.align=r)}function f(e,t){for(var n=e.localVars;n;n=n.next)if(n.name==t)return!0}function s(e,t,n,r,a){var i=e.cc;for(ae.state=e,ae.stream=a,ae.marked=null,ae.cc=i,e.lexical.hasOwnProperty("align")||(e.lexical.align=!0);;){var o=i.length?i.pop():V;if(o(n,r)){for(;i.length&&i[i.length-1].lex;)i.pop()();return ae.marked?ae.marked:"variable"==n&&f(e,r)?"variable-2":"variable"==n&&d(e,r)?"variable-3":t}}}function d(e,t){if(/[a-z]/.test(t.charAt(0)))return!1;for(var n=e.importedtypes.length,r=0;r<n;r++)if(e.importedtypes[r]==t)return!0}function p(e){for(var t=ae.state,n=t.importedtypes;n;n=n.next)if(n.name==e)return;t.importedtypes={name:e,next:t.importedtypes}}function m(){for(var e=arguments.length-1;e>=0;e--)ae.cc.push(arguments[e])}function v(){return m.apply(null,arguments),!0}function b(e,t){for(var n=t;n;n=n.next)if(n.name==e)return!0;return!1}function y(e){var t=ae.state;if(t.context){if(ae.marked="def",b(e,t.localVars))return;t.localVars={name:e,next:t.localVars}}else if(t.globalVars){if(b(e,t.globalVars))return;t.globalVars={name:e,next:t.globalVars}}}function x(){ae.state.context||(ae.state.localVars=ie),ae.state.context={prev:ae.state.context,vars:ae.state.localVars}}function h(){ae.state.localVars=ae.state.context.vars,ae.state.context=ae.state.context.prev}function k(e,t){var n=function(){var n=ae.state;n.lexical=new c(n.indented,ae.stream.column(),e,null,n.lexical,t)};return n.lex=!0,n}function w(){var e=ae.state;e.lexical.prev&&(")"==e.lexical.type&&(e.indented=e.lexical.indented),e.lexical=e.lexical.prev)}function g(e){function t(n){return n==e?v():";"==e?m():v(t)}return t}function V(e){return"@"==e?v(z):"var"==e?v(k("vardef"),j,g(";"),w):"keyword a"==e?v(k("form"),A,V,w):"keyword b"==e?v(k("form"),V,w):"{"==e?v(k("}"),x,_,w,h):";"==e?v():"attribute"==e?v(W):"function"==e?v(F):"for"==e?v(k("form"),g("("),k(")"),q,g(")"),w,V,w):"variable"==e?v(k("stat"),Z):"switch"==e?v(k("form"),A,k("}","switch"),g("{"),_,w,w):"case"==e?v(A,g(":")):"default"==e?v(g(":")):"catch"==e?v(k("form"),x,g("("),H,g(")"),V,w,h):"import"==e?v(C,g(";")):"typedef"==e?v(T):m(k("stat"),A,g(";"),w)}function A(e){return re.hasOwnProperty(e)?v(E):"type"==e?v(E):"function"==e?v(F):"keyword c"==e?v(S):"("==e?v(k(")"),S,g(")"),w,E):"operator"==e?v(A):"["==e?v(k("]"),P(S,"]"),w,E):"{"==e?v(k("}"),P(O,"}"),w,E):v()}function S(e){return e.match(/[;\}\)\],]/)?m():m(A)}function E(e,t){if("operator"==e&&/\+\+|--/.test(t))return v(E);if("operator"==e||":"==e)return v(A);if(";"!=e)return"("==e?v(k(")"),P(A,")"),w,E):"."==e?v(I,E):"["==e?v(k("]"),A,g("]"),w,E):void 0}function W(e){return"attribute"==e?v(W):"function"==e?v(F):"var"==e?v(j):void 0}function z(e){return":"==e?v(z):"variable"==e?v(z):"("==e?v(k(")"),P(M,")"),w,V):void 0}function M(e){if("variable"==e)return v()}function C(e,t){return"variable"==e&&/[A-Z]/.test(t.charAt(0))?(p(t),v()):"variable"==e||"property"==e||"."==e||"*"==t?v(C):void 0}function T(e,t){return"variable"==e&&/[A-Z]/.test(t.charAt(0))?(p(t),v()):"type"==e&&/[A-Z]/.test(t.charAt(0))?v():void 0}function Z(e){return":"==e?v(w,V):m(E,g(";"),w)}function I(e){if("variable"==e)return ae.marked="property",v()}function O(e){if("variable"==e&&(ae.marked="property"),re.hasOwnProperty(e))return v(g(":"),A)}function P(e,t){function n(r){return","==r?v(e,n):r==t?v():v(g(t))}return function(r){return r==t?v():m(e,n)}}function _(e){return"}"==e?v():m(V,_)}function j(e,t){return"variable"==e?(y(t),v(U,D)):v()}function D(e,t){return"="==t?v(A,D):","==e?v(j):void 0}function q(e,t){return"variable"==e?(y(t),v(B,A)):m()}function B(e,t){if("in"==t)return v()}function F(e,t){return"variable"==e||"type"==e?(y(t),v(F)):"new"==t?v(F):"("==e?v(k(")"),x,P(H,")"),w,U,V,h):void 0}function U(e){if(":"==e)return v($)}function $(e){return"type"==e?v():"variable"==e?v():"{"==e?v(k("}"),P(G,"}"),w):void 0}function G(e){if("variable"==e)return v(U)}function H(e,t){if("variable"==e)return y(t),v(U)}var J,K,L=e.indentUnit,N=n("keyword a"),Q=n("keyword b"),R=n("keyword c"),X=n("operator"),Y={type:"atom",style:"atom"},ee={type:"attribute",style:"attribute"},J=n("typedef"),te={if:N,while:N,else:Q,do:Q,try:Q,return:R,break:R,continue:R,new:R,throw:R,var:n("var"),inline:ee,static:ee,using:n("import"),public:ee,private:ee,cast:n("cast"),import:n("import"),macro:n("macro"),function:n("function"),catch:n("catch"),untyped:n("untyped"),callback:n("cb"),for:n("for"),switch:n("switch"),case:n("case"),default:n("default"),in:X,never:n("property_access"),trace:n("trace"),class:J,abstract:J,enum:J,interface:J,typedef:J,extends:J,implements:J,dynamic:J,true:Y,false:Y,null:Y},ne=/[+\-*&%=<>!?|]/,re={atom:!0,number:!0,variable:!0,string:!0,regexp:!0},ae={state:null,column:null,marked:null,cc:null},ie={name:"this",next:null};return h.lex=!0,w.lex=!0,{startState:function(e){var n=["Int","Float","String","Void","Std","Bool","Dynamic","Array"],r={tokenize:o,reAllowed:!0,kwAllowed:!0,cc:[],lexical:new c((e||0)-L,0,"block",(!1)),localVars:t.localVars,importedtypes:n,context:t.localVars&&{vars:t.localVars},indented:0};return t.globalVars&&"object"==typeof t.globalVars&&(r.globalVars=t.globalVars),r},token:function(e,t){if(e.sol()&&(t.lexical.hasOwnProperty("align")||(t.lexical.align=!1),t.indented=e.indentation()),e.eatSpace())return null;var n=t.tokenize(e,t);return"comment"==J?n:(t.reAllowed=!("operator"!=J&&"keyword c"!=J&&!J.match(/^[\[{}\(,;:]$/)),t.kwAllowed="."!=J,s(t,n,J,K,e))},indent:function(e,t){if(e.tokenize!=o)return 0;var n=t&&t.charAt(0),r=e.lexical;"stat"==r.type&&"}"==n&&(r=r.prev);var a=r.type,i=n==a;return"vardef"==a?r.indented+4:"form"==a&&"{"==n?r.indented:"stat"==a||"form"==a?r.indented+L:"switch"!=r.info||i?r.align?r.column+(i?0:1):r.indented+(i?0:L):r.indented+(/^(?:case|default)\b/.test(t)?L:2*L)},electricChars:"{}",blockCommentStart:"/*",blockCommentEnd:"*/",lineComment:"//"}}),e.defineMIME("text/x-haxe","haxe"),e.defineMode("hxml",function(){return{startState:function(){return{define:!1,inString:!1}},token:function(e,t){var n=e.peek(),r=e.sol();if("#"==n)return e.skipToEnd(),"comment";if(r&&"-"==n){var a="variable-2";return e.eat(/-/),"-"==e.peek()&&(e.eat(/-/),a="keyword a"),"D"==e.peek()&&(e.eat(/[D]/),a="keyword c",t.define=!0),e.eatWhile(/[A-Z]/i),a}var n=e.peek();return 0==t.inString&&"'"==n&&(t.inString=!0,n=e.next()),1==t.inString?(e.skipTo("'")||e.skipToEnd(),"'"==e.peek()&&(e.next(),t.inString=!1),"string"):(e.next(),null)},lineComment:"#"}}),e.defineMIME("text/x-hxml","hxml")});