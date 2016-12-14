const spawn = require("child_process").spawn;

// To manage kazoup platform microservices
module.exports = (function() {
  let isDevelopment = process.env.NODE_ENV === "development";
  var running = [];



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

  var _getServices = function() {
    return running;
  };

  var _startServices = function() {
    let resourcesPath = process.env.PWD;

    // Set resources path depending if we running in development
    if (!isDevelopment) {
      //resourcesPath = process.resourcesPath;
      if (process.platform == "win32") {
        elastic = _startService(resourcesPath + "/elasticsearch/bin/elasticsearch.bat", []);
        kazoup = _startService(resourcesPath 		+ "/bin/windows/" + archConvert(process.arch) + "/kazoup.exe", ["desktop"])
      } else {
        elastic = _startService(resourcesPath + "/elasticsearch/bin/elasticsearch", []);
        kazoup = _startService(resourcesPath 		+ "/bin/" + process.platform + "/" + archConvert(process.arch) + "/kazoup", ["desktop"]);
      }

      console.log(resourcesPath 		+ "/bin/" + process.platform + "/" + archConvert(process.arch) + "/kazoup")
      running.push(elastic)
      running.push(kazoup)

    } else {
      //TODO: stupid hack fixme should we point to desktop folder ?
      resourcesPath = __dirname + "/..";
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
    getServices: _getServices,
    startServices: _startServices,
    startService: _startService
  }
}());
