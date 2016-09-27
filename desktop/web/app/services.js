const spawn = require("child_process").spawn;

let isDevelopment = process.env.NODE_ENV === "development";

// To manage kazoup platform microservices
module.exports = (function() {
  function archConvert(arch) {
    switch (arch) {
      case "x64":
        arch = "amd64";
        break;
      case "ia32":
        arch = "386";
        break;
    }

    return arch;
  }

  var _startServices = function() {
    running = [];
    resourcesPath = "";

    //Set resources path depending if we running in development - needs to do it better
    if (!isDevelopment) {
      resourcesPath = process.resourcesPath;// Starting all required services FIXME:
      if (process.platform == "win32") {
        elastic = _startService(resourcesPath + "/elasticsearch/bin/elasticsearch.bat", []);
        kazoup = _startService(resourcesPath 		+ "/bin/windows/" + archConvert(process.arch) + "/kazoup.exe", ["desktop"])
      } else {
        elastic = _startService(resourcesPath + "/elasticsearch/bin/elasticsearch", []);
        kazoup = _startService(resourcesPath 		+ "/bin/" + process.platform + "/" + archConvert(process.arch) + "/kazoup", ["desktop"])
      }
      running.push(elastic)
      running.push(kazoup)

    } else {
      //TODO: stupid hack fixme should we point to desktop folder ?
      resourcesPath = __dirname + "/..";
      /*console.log(resourcesPath)

       if (process.platform == "win32") {
       elastic = startService(resourcesPath + "/elasticsearch/bin/elasticsearch.bat", []);
       kazoup = startService(resourcesPath 		+ "/bin/windows/" + archConvert(process.arch) + "/kazoup.exe", ["desktop"])
       } else {
       elastic = startService(resourcesPath + "/elasticsearch/bin/elasticsearch", []);
       kazoup = startService(resourcesPath 		+ "/bin/" + process.platform + "/" + archConvert(process.arch) + "/kazoup", ["desktop"])
       }
       running.push(elastic)
       running.push(kazoup)*/
    }
  };

  var _startService = function(path, args) {
    // Fuck me I hate JS FIXME

    if (path.includes("elastic")) {
      var productionEnv = Object.create(process.env);

      //if (path.includes("win")) {
      if (process.platform.includes("win")) {
        productionEnv.JAVA_HOME = process.resourcesPath + "/java/windows/" + archConvert(process.arch);
      } else {
        productionEnv.JAVA_HOME = process.resourcesPath + "/java/" + process.platform + "/" + archConvert(process.arch);
      }

      console.log("PRODUCTION ENV", productionEnv);

      es = spawn(path, args, {
        wd: path,
        stdio: 'ignore',
        env: productionEnv
      });
    } else {
      es = spawn(path, args, {
        wd: path,
        stdio: 'ignore'
      });
    }

    es.on("close", function(code, signal) {
      console.log("stdout: es process terminated due to receipt of signal " + signal);
    });

    return es;
  };

  return {
    startServices: _startServices,
    startService: _startService
  }
}());
