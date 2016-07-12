# TODO
Mobile elements - web-proxy should get device type (we cannot check sreen width/device type relationship anymore) and redirect
Elements for non admin users : kazoup-user-search, this one will be easier to do once we have user api

# Kazoup Web [![CircleCI](https://circleci.com/gh/kazoup/kazoup-web.svg?style=svg&circle-token=1084085b649711ccdac2e6355412dcd9fb259f64)](https://circleci.com/gh/kazoup/kazoup-web)

Packages Kazoup frontend

## DEV environment

Build docker images:

```
make
```

Run docker containers:
```
docker-compose -f docker-compose-dev.yml up
```

SASS files needs to be compiled on the fly. Generated CSS files are mapped into the static-web container.
```
cd frontend
gulp dev
```

## PROD environment


