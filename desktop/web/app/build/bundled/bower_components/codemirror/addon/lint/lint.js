!function(t){"object"==typeof exports&&"object"==typeof module?t(require("../../lib/codemirror")):"function"==typeof define&&define.amd?define(["../../lib/codemirror"],t):t(CodeMirror)}(function(t){"use strict";function e(e,n){function o(e){return r.parentNode?(r.style.top=Math.max(0,e.clientY-r.offsetHeight-5)+"px",void(r.style.left=e.clientX+5+"px")):t.off(document,"mousemove",o)}var r=document.createElement("div");return r.className="CodeMirror-lint-tooltip",r.appendChild(n.cloneNode(!0)),document.body.appendChild(r),t.on(document,"mousemove",o),o(e),null!=r.style.opacity&&(r.style.opacity=1),r}function n(t){t.parentNode&&t.parentNode.removeChild(t)}function o(t){t.parentNode&&(null==t.style.opacity&&n(t),t.style.opacity=0,setTimeout(function(){n(t)},600))}function r(n,r,i){function a(){t.off(i,"mouseout",a),l&&(o(l),l=null)}var l=e(n,r),u=setInterval(function(){if(l)for(var t=i;;t=t.parentNode){if(t&&11==t.nodeType&&(t=t.host),t==document.body)return;if(!t){a();break}}if(!l)return clearInterval(u)},400);t.on(i,"mouseout",a)}function i(t,e,n){this.marked=[],this.options=e,this.timeout=null,this.hasGutter=n,this.onMouseOver=function(e){g(t,e)},this.waitingFor=0}function a(t,e){return e instanceof Function?{getAnnotations:e}:(e&&e!==!0||(e={}),e)}function l(t){var e=t.state.lint;e.hasGutter&&t.clearGutter(y);for(var n=0;n<e.marked.length;++n)e.marked[n].clear();e.marked.length=0}function u(e,n,o,i){var a=document.createElement("div"),l=a;return a.className="CodeMirror-lint-marker-"+n,o&&(l=a.appendChild(document.createElement("div")),l.className="CodeMirror-lint-marker-multiple"),0!=i&&t.on(l,"mouseover",function(t){r(t,e,l)}),a}function s(t,e){return"error"==t?t:e}function c(t){for(var e=[],n=0;n<t.length;++n){var o=t[n],r=o.from.line;(e[r]||(e[r]=[])).push(o)}return e}function f(t){var e=t.severity;e||(e="error");var n=document.createElement("div");return n.className="CodeMirror-lint-message-"+e,n.appendChild(document.createTextNode(t.message)),n}function m(e,n,o){function r(){a=-1,e.off("change",r)}var i=e.state.lint,a=++i.waitingFor;e.on("change",r),n(e.getValue(),function(n,o){e.off("change",r),i.waitingFor==a&&(o&&n instanceof t&&(n=o),p(e,n))},o,e)}function d(e){var n=e.state.lint,o=n.options,r=o.options||o,i=o.getAnnotations||e.getHelper(t.Pos(0,0),"lint");i&&(o.async||i.async?m(e,i,r):p(e,i(e.getValue(),r,e)))}function p(t,e){l(t);for(var n=t.state.lint,o=n.options,r=c(e),i=0;i<r.length;++i){var a=r[i];if(a){for(var m=null,d=n.hasGutter&&document.createDocumentFragment(),p=0;p<a.length;++p){var h=a[p],v=h.severity;v||(v="error"),m=s(m,v),o.formatAnnotation&&(h=o.formatAnnotation(h)),n.hasGutter&&d.appendChild(f(h)),h.to&&n.marked.push(t.markText(h.from,h.to,{className:"CodeMirror-lint-mark-"+v,__annotation:h}))}n.hasGutter&&t.setGutterMarker(i,y,u(d,m,a.length>1,n.options.tooltips))}}o.onUpdateLinting&&o.onUpdateLinting(e,r,t)}function h(t){var e=t.state.lint;e&&(clearTimeout(e.timeout),e.timeout=setTimeout(function(){d(t)},e.options.delay||500))}function v(t,e){for(var n=e.target||e.srcElement,o=document.createDocumentFragment(),i=0;i<t.length;i++){var a=t[i];o.appendChild(f(a))}r(e,o,n)}function g(t,e){var n=e.target||e.srcElement;if(/\bCodeMirror-lint-mark-/.test(n.className)){for(var o=n.getBoundingClientRect(),r=(o.left+o.right)/2,i=(o.top+o.bottom)/2,a=t.findMarksAt(t.coordsChar({left:r,top:i},"client")),l=[],u=0;u<a.length;++u){var s=a[u].__annotation;s&&l.push(s)}l.length&&v(l,e)}}var y="CodeMirror-lint-markers";t.defineOption("lint",!1,function(e,n,o){if(o&&o!=t.Init&&(l(e),e.state.lint.options.lintOnChange!==!1&&e.off("change",h),t.off(e.getWrapperElement(),"mouseover",e.state.lint.onMouseOver),clearTimeout(e.state.lint.timeout),delete e.state.lint),n){for(var r=e.getOption("gutters"),u=!1,s=0;s<r.length;++s)r[s]==y&&(u=!0);var c=e.state.lint=new i(e,a(e,n),u);c.options.lintOnChange!==!1&&e.on("change",h),0!=c.options.tooltips&&t.on(e.getWrapperElement(),"mouseover",c.onMouseOver),d(e)}}),t.defineExtension("performLint",function(){this.state.lint&&d(this)})});