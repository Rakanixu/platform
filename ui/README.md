# Kazoup Desktop

Distributed Kazoup desktop application

# Run 

In development 

Desktop

```

npm run dev

```
Web App access 

```

polymer serve

```

## Build
<!-- 
TODO:Simplify this
-->
To build distribution for all platforms run below
This will build OSX,Linux and Windows packages and save tem in /dist folder
You will need to install dependencies so far it works on Mac see more over [here](https://github.com/electron-userland/electron-builder/wiki/Multi-Platform-Build)
 
```
npm run dist

```

To prepare relase run below.
It will copy neccesery files into /release folder

```

npm run release

```

To deploy packages to github release run below

```

./deploy.sh

```



