!function(e){"object"==typeof exports&&"object"==typeof module?e(require("../../lib/codemirror")):"function"==typeof define&&define.amd?define(["../../lib/codemirror"],e):e(CodeMirror)}(function(e){"use strict";e.defineMode("asterisk",function(){function e(e,n){var a="",r=e.next();if(";"==r)return e.skipToEnd(),"comment";if("["==r)return e.skipTo("]"),e.eat("]"),"header";if('"'==r)return e.skipTo('"'),"string";if("'"==r)return e.skipTo("'"),"string-2";if("#"==r&&(e.eatWhile(/\w/),a=e.current(),i.indexOf(a)!==-1))return e.skipToEnd(),"strong";if("$"==r){var o=e.peek();if("{"==o)return e.skipTo("}"),e.eat("}"),"variable-3"}if(e.eatWhile(/\w/),a=e.current(),t.indexOf(a)!==-1){switch(n.extenStart=!0,a){case"same":n.extenSame=!0;break;case"include":case"switch":case"ignorepat":n.extenInclude=!0}return"atom"}}var t=["exten","same","include","ignorepat","switch"],i=["#include","#exec"],n=["addqueuemember","adsiprog","aelsub","agentlogin","agentmonitoroutgoing","agi","alarmreceiver","amd","answer","authenticate","background","backgrounddetect","bridge","busy","callcompletioncancel","callcompletionrequest","celgenuserevent","changemonitor","chanisavail","channelredirect","chanspy","clearhash","confbridge","congestion","continuewhile","controlplayback","dahdiacceptr2call","dahdibarge","dahdiras","dahdiscan","dahdisendcallreroutingfacility","dahdisendkeypadfacility","datetime","dbdel","dbdeltree","deadagi","dial","dictate","directory","disa","dumpchan","eagi","echo","endwhile","exec","execif","execiftime","exitwhile","extenspy","externalivr","festival","flash","followme","forkcdr","getcpeid","gosub","gosubif","goto","gotoif","gotoiftime","hangup","iax2provision","ices","importvar","incomplete","ivrdemo","jabberjoin","jabberleave","jabbersend","jabbersendgroup","jabberstatus","jack","log","macro","macroexclusive","macroexit","macroif","mailboxexists","meetme","meetmeadmin","meetmechanneladmin","meetmecount","milliwatt","minivmaccmess","minivmdelete","minivmgreet","minivmmwi","minivmnotify","minivmrecord","mixmonitor","monitor","morsecode","mp3player","mset","musiconhold","nbscat","nocdr","noop","odbc","odbc","odbcfinish","originate","ospauth","ospfinish","osplookup","ospnext","page","park","parkandannounce","parkedcall","pausemonitor","pausequeuemember","pickup","pickupchan","playback","playtones","privacymanager","proceeding","progress","queue","queuelog","raiseexception","read","readexten","readfile","receivefax","receivefax","receivefax","record","removequeuemember","resetcdr","retrydial","return","ringing","sayalpha","saycountedadj","saycountednoun","saycountpl","saydigits","saynumber","sayphonetic","sayunixtime","senddtmf","sendfax","sendfax","sendfax","sendimage","sendtext","sendurl","set","setamaflags","setcallerpres","setmusiconhold","sipaddheader","sipdtmfmode","sipremoveheader","skel","slastation","slatrunk","sms","softhangup","speechactivategrammar","speechbackground","speechcreate","speechdeactivategrammar","speechdestroy","speechloadgrammar","speechprocessingsound","speechstart","speechunloadgrammar","stackpop","startmusiconhold","stopmixmonitor","stopmonitor","stopmusiconhold","stopplaytones","system","testclient","testserver","transfer","tryexec","trysystem","unpausemonitor","unpausequeuemember","userevent","verbose","vmauthenticate","vmsayname","voicemail","voicemailmain","wait","waitexten","waitfornoise","waitforring","waitforsilence","waitmusiconhold","waituntil","while","zapateller"];return{startState:function(){return{extenStart:!1,extenSame:!1,extenInclude:!1,extenExten:!1,extenPriority:!1,extenApplication:!1}},token:function(t,i){var a="";return t.eatSpace()?null:i.extenStart?(t.eatWhile(/[^\s]/),a=t.current(),/^=>?$/.test(a)?(i.extenExten=!0,i.extenStart=!1,"strong"):(i.extenStart=!1,t.skipToEnd(),"error")):i.extenExten?(i.extenExten=!1,i.extenPriority=!0,t.eatWhile(/[^,]/),i.extenInclude&&(t.skipToEnd(),i.extenPriority=!1,i.extenInclude=!1),i.extenSame&&(i.extenPriority=!1,i.extenSame=!1,i.extenApplication=!0),"tag"):i.extenPriority?(i.extenPriority=!1,i.extenApplication=!0,t.next(),i.extenSame?null:(t.eatWhile(/[^,]/),"number")):i.extenApplication?(t.eatWhile(/,/),a=t.current(),","===a?null:(t.eatWhile(/\w/),a=t.current().toLowerCase(),i.extenApplication=!1,n.indexOf(a)!==-1?"def strong":null)):e(t,i)}}}),e.defineMIME("text/x-asterisk","asterisk")});