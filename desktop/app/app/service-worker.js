/**
 * Copyright 2016 Google Inc. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// This generated service worker JavaScript will precache your site's resources.
// The code needs to be saved in a .js file at the top-level of your site, and registered
// from your pages in order to be used. See
// https://github.com/googlechrome/sw-precache/blob/master/demo/app/js/service-worker-registration.js
// for an example of how you can register this script and handle various service worker events.

/* eslint-env worker, serviceworker */
/* eslint-disable indent, no-unused-vars, no-multiple-empty-lines, max-nested-callbacks, space-before-function-paren */
'use strict';





/* eslint-disable quotes, comma-spacing */
var PrecacheConfig = [["/bower_components/expandjs/dist/expandjs.min.js","5ae1d7e9a0967aed8480cb9f80258f74"],["/bower_components/expandjs/expandjs.html","d6d82e14c9996cb71df9d61107a9ddfd"],["/bower_components/font-roboto/roboto.html","09500fd5adfad056ff5aa05e2aae0ec5"],["/bower_components/iron-a11y-keys-behavior/iron-a11y-keys-behavior.html","b9a8e766d0ab03a5d13e275754ec3d54"],["/bower_components/iron-behaviors/iron-button-state.html","6565a80d1af09299c1201f8286849c3b"],["/bower_components/iron-behaviors/iron-control-state.html","1c12ee539b1dbbd0957ae26b3549cc13"],["/bower_components/iron-flex-layout/iron-flex-layout.html","f6d0b5075d5f70cac9b4bd66bd88c459"],["/bower_components/iron-icon/iron-icon.html","923e727686cfe7c6e58cc32c2a49717d"],["/bower_components/iron-icons/iron-icons.html","c8f9154ae89b94e658e4a52eee690a16"],["/bower_components/iron-iconset-svg/iron-iconset-svg.html","8fb45b1b4668dae069f5efb5004c2af4"],["/bower_components/iron-media-query/iron-media-query.html","7436f9608ebd2d31e4b346921651f84b"],["/bower_components/iron-menu-behavior/iron-menu-behavior.html","8c4fc9ccbb28f3bf68c621ebc3859fb7"],["/bower_components/iron-meta/iron-meta.html","dd4ef14e09c5771e70292d70963f6718"],["/bower_components/iron-resizable-behavior/iron-resizable-behavior.html","e93449ccd4312e4e30060c87bd52ed25"],["/bower_components/iron-selector/iron-multi-selectable.html","46d6620acd7bad986d81097d9ca91692"],["/bower_components/iron-selector/iron-selectable.html","65b04f3f5f1b551d91a82b36f916f6b6"],["/bower_components/iron-selector/iron-selection.html","83545b7d1eae4020594969e6b9790b65"],["/bower_components/iron-selector/iron-selector.html","4d2657550768bec0788eba5190cddc66"],["/bower_components/mat-avatar/mat-avatar.html","89b0a388158390e302cf980d24ccf740"],["/bower_components/mat-breadcrumb/mat-breadcrumb-step.html","5c3957d6f5884b30cdb0ae0af0e5d0ef"],["/bower_components/mat-breadcrumb/mat-breadcrumb.html","b2059c2464c6c3b72dfa000753d1f856"],["/bower_components/mat-divider/mat-divider.html","e2847468a88ef43ed7374ff0b3da7b33"],["/bower_components/mat-icon-button/mat-icon-button.html","5735fd85dab3f3a5f0ef5809af978e0d"],["/bower_components/mat-icon/mat-icon.html","717bb03f6435b9d6c106736605bd4354"],["/bower_components/mat-icons/mat-icons.html","3aba3e418f9b4fb586c77c68c3552727"],["/bower_components/mat-ink/mat-ink-behavior.html","1984ec8e24c17c86d3430f6f32c8fe52"],["/bower_components/mat-ink/mat-ink-styles.html","bd5e50323c087832411c07b7bf974881"],["/bower_components/mat-ink/mat-ink.html","8d1639516cd8e79dc92764efc9da3287"],["/bower_components/mat-item/mat-item.html","83ded5653109cf9ac5d0860959b6e99f"],["/bower_components/mat-list/mat-list.html","da395dbea2a87bf0436326bd93f80cb8"],["/bower_components/mat-palette/mat-palette.html","270b79a811996d850fd7483056bdb575"],["/bower_components/mat-paper/mat-paper-behavior.html","227b6fb0a4f397197c7b922b0627a9b7"],["/bower_components/mat-paper/mat-paper-styles.html","21e04bb29a55feec6a1a90df8b74a9a5"],["/bower_components/mat-pressed-behavior/mat-pressed-behavior.html","77b83a28d6fa1ac9bcb05b870614267e"],["/bower_components/mat-pressed-behavior/mat-pressed-ink-behavior.html","e0633f2ac72a7f699ee89d0b7a3e8f98"],["/bower_components/mat-pressed-behavior/mat-pressed-ink-styles.html","01ec14433f8b79e9141b8e9a83657249"],["/bower_components/mat-pressed-behavior/mat-pressed-paper-behavior.html","aeef76aa00c103e95f312a1856b5b701"],["/bower_components/mat-pressed-behavior/mat-pressed-paper-styles.html","8861a01e1cf96527a9fc3b9089d5f065"],["/bower_components/mat-pressed-behavior/mat-pressed-styles.html","9cef7a094cb511cf75e0ef96556cb7b4"],["/bower_components/mat-shadow/mat-shadow.html","3e8d497affaeea0d7c89e23947bfcf68"],["/bower_components/mat-typography/mat-typography.html","951cca6caab81e58ae20e066c4b9b0ac"],["/bower_components/paper-behaviors/paper-inky-focus-behavior.html","51a1c5ccd2aae4c1a0258680dcb3e1ea"],["/bower_components/paper-behaviors/paper-ripple-behavior.html","b6ee8dd59ffb46ca57e81311abd2eca0"],["/bower_components/paper-drawer-panel/paper-drawer-panel.html","fb78c193d694a9ff0d9fc76dc5d5763a"],["/bower_components/paper-header-panel/paper-header-panel.html","bd966e2674c837eff7a8a9f240ee3b29"],["/bower_components/paper-icon-button/paper-icon-button.html","4a5cbc3fe046e2c070d4bf34ec7463d6"],["/bower_components/paper-item/paper-item-behavior.html","82636a7562fd8b0be5b15646ee461588"],["/bower_components/paper-item/paper-item-shared-styles.html","389eedfc65ee58b1f0d67281d0bad1a1"],["/bower_components/paper-item/paper-item.html","5099885c3bd34e04df7796d48851c4a4"],["/bower_components/paper-menu/paper-menu-shared-styles.html","d284d59303c2383edf6c626dd679302d"],["/bower_components/paper-menu/paper-menu.html","3d9cf400d7ee8753ab6d0cb6358bb711"],["/bower_components/paper-ripple/paper-ripple.html","30fa6456055a5725c6492f8e5a364f39"],["/bower_components/paper-styles/color.html","c53abb41659bf242d420a7f93b977e91"],["/bower_components/paper-styles/default-theme.html","25d95202be2ff5b60f651924e66abed2"],["/bower_components/paper-styles/typography.html","3f95c68bcd0bd4710f3469c4900533d6"],["/bower_components/paper-toolbar/paper-toolbar.html","e54bc7361f1e80997c80621b908dafdd"],["/bower_components/polymer/polymer-micro.html","ecf1ad808ec62a7adcec68e28cf3ffad"],["/bower_components/polymer/polymer-mini.html","e48d322a1d599c9db40523f050fbef23"],["/bower_components/polymer/polymer.html","837764153a0347c0e906b48d554941a0"],["/bower_components/webcomponentsjs/webcomponents-lite.js","9dc13c1fee8c627a241d629d0ea8fd7b"],["/bower_components/xp-anchor-behavior/xp-anchor-behavior.html","2749a5724cb39d5d0b8376106d0096ad"],["/bower_components/xp-anchor-behavior/xp-anchor-styles.html","79fc8638d4ad350b693be734e60ab3fe"],["/bower_components/xp-array-behavior/xp-array-behavior.html","124f7eeb081fefa2d07bf2fd25983d35"],["/bower_components/xp-breadcrumb-behavior/xp-breadcrumb-behavior.html","faab3452bc1bdf258ca62847a86702cf"],["/bower_components/xp-breadcrumb-behavior/xp-breadcrumb-step-behavior.html","a8b458b0a645e1bf764c0412eb9bd767"],["/bower_components/xp-breadcrumb-behavior/xp-breadcrumb-step-styles.html","0def3f5f924e99e6991c459a660e1ebb"],["/bower_components/xp-breadcrumb-behavior/xp-breadcrumb-styles.html","8706c132a9790cfaf61f79a8eb5c4946"],["/bower_components/xp-finder-behavior/xp-finder-behavior.html","5d8bdd9e10a18377a15f0fb6af013b67"],["/bower_components/xp-focused-behavior/xp-focused-behavior.html","fdba3179a92879631dce4a7608031f0e"],["/bower_components/xp-focused-behavior/xp-focused-styles.html","a778ae69855a2a61b2f8e11b03b0cbad"],["/bower_components/xp-icon-behavior/xp-icon-behavior.html","6c709413eebb98a27fdb1384d43eece7"],["/bower_components/xp-icon-behavior/xp-icon-styles.html","577ab61547c3e96a639800cc24d86a10"],["/bower_components/xp-iconset/xp-iconset-finder.html","7046b504acef1fd4660e5b8c1f0f5237"],["/bower_components/xp-iconset/xp-iconset.html","68716545fbba9eb1c7417a05dea32587"],["/bower_components/xp-list-behavior/xp-list-behavior.html","dfe869c754cf438a33321e3d668f9082"],["/bower_components/xp-list-behavior/xp-list-styles.html","c3b1a34150406dd2f81ce53c5bd68eb9"],["/bower_components/xp-master-behavior/xp-master-behavior.html","1584c50f452513e47aef32150a41f747"],["/bower_components/xp-media-query/xp-media-query.html","b27c49e6a6561718c9502b18ef8af7d2"],["/bower_components/xp-overlay/xp-overlay-injector.html","efe8208e23fb932986fcfb10b097fe75"],["/bower_components/xp-pressed-behavior/xp-pressed-behavior.html","1d16d3b98279d053b73d2f4ae808541e"],["/bower_components/xp-pressed-behavior/xp-pressed-styles.html","d8dccec05b2c07277a0e24956c21ff8e"],["/bower_components/xp-refirer-behavior/xp-refirer-behavior.html","aedf8f6bb0508877b0d3b6c804ac1445"],["/bower_components/xp-selector/xp-selector-behavior.html","e59ae4b5de4aa23376f5a33c5eea9602"],["/bower_components/xp-selector/xp-selector-multi-behavior.html","6ea62d7a96147e82154cd4bf1cd52645"],["/bower_components/xp-slave-behavior/xp-slave-behavior.html","2ccb48c6e1900a72a481e9b20219c96d"],["/bower_components/xp-targeter-behavior/xp-targeter-behavior.html","14d4870a888f96c40eaf2722a489a0c8"],["/index.html","67a50bef2ce8a540c5c8fd083c827ce1"],["/src/web-app/web-app.html","3e54f6030297667299d98d274478bce8"]];
/* eslint-enable quotes, comma-spacing */
var CacheNamePrefix = 'sw-precache-v1--' + (self.registration ? self.registration.scope : '') + '-';


var IgnoreUrlParametersMatching = [/^utm_/];



var addDirectoryIndex = function (originalUrl, index) {
    var url = new URL(originalUrl);
    if (url.pathname.slice(-1) === '/') {
      url.pathname += index;
    }
    return url.toString();
  };

var getCacheBustedUrl = function (url, param) {
    param = param || Date.now();

    var urlWithCacheBusting = new URL(url);
    urlWithCacheBusting.search += (urlWithCacheBusting.search ? '&' : '') +
      'sw-precache=' + param;

    return urlWithCacheBusting.toString();
  };

var isPathWhitelisted = function (whitelist, absoluteUrlString) {
    // If the whitelist is empty, then consider all URLs to be whitelisted.
    if (whitelist.length === 0) {
      return true;
    }

    // Otherwise compare each path regex to the path of the URL passed in.
    var path = (new URL(absoluteUrlString)).pathname;
    return whitelist.some(function(whitelistedPathRegex) {
      return path.match(whitelistedPathRegex);
    });
  };

var populateCurrentCacheNames = function (precacheConfig,
    cacheNamePrefix, baseUrl) {
    var absoluteUrlToCacheName = {};
    var currentCacheNamesToAbsoluteUrl = {};

    precacheConfig.forEach(function(cacheOption) {
      var absoluteUrl = new URL(cacheOption[0], baseUrl).toString();
      var cacheName = cacheNamePrefix + absoluteUrl + '-' + cacheOption[1];
      currentCacheNamesToAbsoluteUrl[cacheName] = absoluteUrl;
      absoluteUrlToCacheName[absoluteUrl] = cacheName;
    });

    return {
      absoluteUrlToCacheName: absoluteUrlToCacheName,
      currentCacheNamesToAbsoluteUrl: currentCacheNamesToAbsoluteUrl
    };
  };

var stripIgnoredUrlParameters = function (originalUrl,
    ignoreUrlParametersMatching) {
    var url = new URL(originalUrl);

    url.search = url.search.slice(1) // Exclude initial '?'
      .split('&') // Split into an array of 'key=value' strings
      .map(function(kv) {
        return kv.split('='); // Split each 'key=value' string into a [key, value] array
      })
      .filter(function(kv) {
        return ignoreUrlParametersMatching.every(function(ignoredRegex) {
          return !ignoredRegex.test(kv[0]); // Return true iff the key doesn't match any of the regexes.
        });
      })
      .map(function(kv) {
        return kv.join('='); // Join each [key, value] array into a 'key=value' string
      })
      .join('&'); // Join the array of 'key=value' strings into a string with '&' in between each

    return url.toString();
  };


var mappings = populateCurrentCacheNames(PrecacheConfig, CacheNamePrefix, self.location);
var AbsoluteUrlToCacheName = mappings.absoluteUrlToCacheName;
var CurrentCacheNamesToAbsoluteUrl = mappings.currentCacheNamesToAbsoluteUrl;

function deleteAllCaches() {
  return caches.keys().then(function(cacheNames) {
    return Promise.all(
      cacheNames.map(function(cacheName) {
        return caches.delete(cacheName);
      })
    );
  });
}

self.addEventListener('install', function(event) {
  event.waitUntil(
    // Take a look at each of the cache names we expect for this version.
    Promise.all(Object.keys(CurrentCacheNamesToAbsoluteUrl).map(function(cacheName) {
      return caches.open(cacheName).then(function(cache) {
        // Get a list of all the entries in the specific named cache.
        // For caches that are already populated for a given version of a
        // resource, there should be 1 entry.
        return cache.keys().then(function(keys) {
          // If there are 0 entries, either because this is a brand new version
          // of a resource or because the install step was interrupted the
          // last time it ran, then we need to populate the cache.
          if (keys.length === 0) {
            // Use the last bit of the cache name, which contains the hash,
            // as the cache-busting parameter.
            // See https://github.com/GoogleChrome/sw-precache/issues/100
            var cacheBustParam = cacheName.split('-').pop();
            var urlWithCacheBusting = getCacheBustedUrl(
              CurrentCacheNamesToAbsoluteUrl[cacheName], cacheBustParam);

            var request = new Request(urlWithCacheBusting,
              {credentials: 'same-origin'});
            return fetch(request).then(function(response) {
              if (response.ok) {
                return cache.put(CurrentCacheNamesToAbsoluteUrl[cacheName],
                  response);
              }

              console.error('Request for %s returned a response status %d, ' +
                'so not attempting to cache it.',
                urlWithCacheBusting, response.status);
              // Get rid of the empty cache if we can't add a successful response to it.
              return caches.delete(cacheName);
            });
          }
        });
      });
    })).then(function() {
      return caches.keys().then(function(allCacheNames) {
        return Promise.all(allCacheNames.filter(function(cacheName) {
          return cacheName.indexOf(CacheNamePrefix) === 0 &&
            !(cacheName in CurrentCacheNamesToAbsoluteUrl);
          }).map(function(cacheName) {
            return caches.delete(cacheName);
          })
        );
      });
    }).then(function() {
      if (typeof self.skipWaiting === 'function') {
        // Force the SW to transition from installing -> active state
        self.skipWaiting();
      }
    })
  );
});

if (self.clients && (typeof self.clients.claim === 'function')) {
  self.addEventListener('activate', function(event) {
    event.waitUntil(self.clients.claim());
  });
}

self.addEventListener('message', function(event) {
  if (event.data.command === 'delete_all') {
    console.log('About to delete all caches...');
    deleteAllCaches().then(function() {
      console.log('Caches deleted.');
      event.ports[0].postMessage({
        error: null
      });
    }).catch(function(error) {
      console.log('Caches not deleted:', error);
      event.ports[0].postMessage({
        error: error
      });
    });
  }
});


self.addEventListener('fetch', function(event) {
  if (event.request.method === 'GET') {
    var urlWithoutIgnoredParameters = stripIgnoredUrlParameters(event.request.url,
      IgnoreUrlParametersMatching);

    var cacheName = AbsoluteUrlToCacheName[urlWithoutIgnoredParameters];
    var directoryIndex = 'index.html';
    if (!cacheName && directoryIndex) {
      urlWithoutIgnoredParameters = addDirectoryIndex(urlWithoutIgnoredParameters, directoryIndex);
      cacheName = AbsoluteUrlToCacheName[urlWithoutIgnoredParameters];
    }

    var navigateFallback = '';
    // Ideally, this would check for event.request.mode === 'navigate', but that is not widely
    // supported yet:
    // https://code.google.com/p/chromium/issues/detail?id=540967
    // https://bugzilla.mozilla.org/show_bug.cgi?id=1209081
    if (!cacheName && navigateFallback && event.request.headers.has('accept') &&
        event.request.headers.get('accept').includes('text/html') &&
        /* eslint-disable quotes, comma-spacing */
        isPathWhitelisted([], event.request.url)) {
        /* eslint-enable quotes, comma-spacing */
      var navigateFallbackUrl = new URL(navigateFallback, self.location);
      cacheName = AbsoluteUrlToCacheName[navigateFallbackUrl.toString()];
    }

    if (cacheName) {
      event.respondWith(
        // Rely on the fact that each cache we manage should only have one entry, and return that.
        caches.open(cacheName).then(function(cache) {
          return cache.keys().then(function(keys) {
            return cache.match(keys[0]).then(function(response) {
              if (response) {
                return response;
              }
              // If for some reason the response was deleted from the cache,
              // raise and exception and fall back to the fetch() triggered in the catch().
              throw Error('The cache ' + cacheName + ' is empty.');
            });
          });
        }).catch(function(e) {
          console.warn('Couldn\'t serve response for "%s" from cache: %O', event.request.url, e);
          return fetch(event.request);
        })
      );
    }
  }
});




