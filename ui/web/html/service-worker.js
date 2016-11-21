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
var PrecacheConfig = [["/bower_components/expandjs/dist/expandjs.min.js","3aab86173354c175f716c9ab588dfa9f"],["/bower_components/expandjs/expandjs.html","d6d82e14c9996cb71df9d61107a9ddfd"],["/bower_components/filesize/lib/filesize.min.js","6e608c936328907be4e042e0fb5ac354"],["/bower_components/font-roboto/roboto.html","09500fd5adfad056ff5aa05e2aae0ec5"],["/bower_components/iron-a11y-announcer/iron-a11y-announcer.html","a3bd031e39dde38cb8e619f670ee50f7"],["/bower_components/iron-a11y-keys-behavior/iron-a11y-keys-behavior.html","b9a8e766d0ab03a5d13e275754ec3d54"],["/bower_components/iron-ajax/iron-ajax.html","d606b330d7bd040660a53a5cda7f8acf"],["/bower_components/iron-ajax/iron-request.html","c2d289c4b20653353cff315cf247a45e"],["/bower_components/iron-behaviors/iron-button-state.html","6565a80d1af09299c1201f8286849c3b"],["/bower_components/iron-behaviors/iron-control-state.html","1c12ee539b1dbbd0957ae26b3549cc13"],["/bower_components/iron-checked-element-behavior/iron-checked-element-behavior.html","6fd1055c2c04382401dc910a0db569c6"],["/bower_components/iron-dropdown/iron-dropdown-scroll-manager.html","70904f32a519b07ec427d1a9a0c71528"],["/bower_components/iron-dropdown/iron-dropdown.html","60a45cf71d0893d16b8c29d063dec959"],["/bower_components/iron-fit-behavior/iron-fit-behavior.html","8d3799ca2f619ed4f31261bb03284671"],["/bower_components/iron-flex-layout/iron-flex-layout-classes.html","90471c0acb830c41b01e02a2507bed3c"],["/bower_components/iron-flex-layout/iron-flex-layout.html","3987521c615734e4fe403f9acecfea54"],["/bower_components/iron-form-element-behavior/iron-form-element-behavior.html","a64177311979fc6a6aae454cb85ea2be"],["/bower_components/iron-icon/iron-icon.html","23fe3af4b80a767dc9ec5e2e0ac5ab42"],["/bower_components/iron-icons/av-icons.html","b69fba5107077e3c4448351591a7cef5"],["/bower_components/iron-icons/device-icons.html","207f23207025327d038aba8ea236f82f"],["/bower_components/iron-icons/editor-icons.html","0c73cde432a2d5b140bc031ac459e8b1"],["/bower_components/iron-icons/image-icons.html","30ef0224c9cf6acd66c506818396ccf7"],["/bower_components/iron-icons/iron-icons.html","c8f9154ae89b94e658e4a52eee690a16"],["/bower_components/iron-icons/social-icons.html","7c0d7482ea9c4ff9b2b76dac1198d9a9"],["/bower_components/iron-iconset-svg/iron-iconset-svg.html","8fb45b1b4668dae069f5efb5004c2af4"],["/bower_components/iron-image/iron-image.html","bde5467573acb26172d1049ab2fca607"],["/bower_components/iron-input/iron-input.html","3e393eda6c241be2817ce0acc512bcf6"],["/bower_components/iron-list/iron-list.html","0186f9177a96a750e08a0b2bf3363b5a"],["/bower_components/iron-media-query/iron-media-query.html","7436f9608ebd2d31e4b346921651f84b"],["/bower_components/iron-menu-behavior/iron-menu-behavior.html","8c4fc9ccbb28f3bf68c621ebc3859fb7"],["/bower_components/iron-menu-behavior/iron-menubar-behavior.html","a0cc6674fc6d9d7ba0b68ff680b4e56b"],["/bower_components/iron-meta/iron-meta.html","dd4ef14e09c5771e70292d70963f6718"],["/bower_components/iron-overlay-behavior/iron-overlay-backdrop.html","35013b4b97041ed6b63cf95dbb9fbcb4"],["/bower_components/iron-overlay-behavior/iron-overlay-behavior.html","3d22bc0675145cda7ca42168bce673ab"],["/bower_components/iron-overlay-behavior/iron-overlay-manager.html","51191fe7de42d8a6a31503eab8ada83b"],["/bower_components/iron-pages/iron-pages.html","5872a2ad58225c94b14ddea747f67cbd"],["/bower_components/iron-range-behavior/iron-range-behavior.html","34f5b83882b6b6c5cfad7543caab080e"],["/bower_components/iron-resizable-behavior/iron-resizable-behavior.html","e93449ccd4312e4e30060c87bd52ed25"],["/bower_components/iron-scroll-target-behavior/iron-scroll-target-behavior.html","0185cbe8d7139c9e92af3a9af67feb17"],["/bower_components/iron-selector/iron-multi-selectable.html","46d6620acd7bad986d81097d9ca91692"],["/bower_components/iron-selector/iron-selectable.html","65b04f3f5f1b551d91a82b36f916f6b6"],["/bower_components/iron-selector/iron-selection.html","83545b7d1eae4020594969e6b9790b65"],["/bower_components/iron-selector/iron-selector.html","4d2657550768bec0788eba5190cddc66"],["/bower_components/iron-validatable-behavior/iron-validatable-behavior.html","02bf0434cc1a0d09e18413dea91dcea1"],["/bower_components/jsOAuth/dist/jsOAuth-1.3.7.min.js","324d718651bd7f7ecacacc49dbc7dfba"],["/bower_components/jshashes/hashes.min.js","6efbdb3858e6739a871cdd15a9337c98"],["/bower_components/l2t-paper-slider/l2t-paper-slider.html","5bb270d4ad6611b5172b6bb01d907a44"],["/bower_components/mat-avatar/mat-avatar.html","0331ed3e9c7176a87efa70c10c02f21c"],["/bower_components/mat-breadcrumb/mat-breadcrumb-step.html","97ecc7ce3dddfbf34c2c3272cc53d569"],["/bower_components/mat-breadcrumb/mat-breadcrumb.html","b2059c2464c6c3b72dfa000753d1f856"],["/bower_components/mat-button/mat-button.html","f3665cba3c320c67d25fe383bfa40a12"],["/bower_components/mat-dialog/mat-dialog.html","102e7a1a9a9bd880640ec8dddab6aac0"],["/bower_components/mat-divider/mat-divider.html","e552c268b66e3d15c5f3b7b780f06d6e"],["/bower_components/mat-icon-button/mat-icon-button.html","36b5adc8a9c44c8e9b4c41de314748d3"],["/bower_components/mat-icon/mat-icon.html","717bb03f6435b9d6c106736605bd4354"],["/bower_components/mat-icons/action-icons.html","7f1ebf5ae2d9601d8adc3547db3c81a3"],["/bower_components/mat-icons/hardware-icons.html","7bb0ac6f9874f4734df4f213107e3c88"],["/bower_components/mat-icons/image-icons.html","636f309ad243efd39db3ace314fd4496"],["/bower_components/mat-icons/mat-icons.html","3aba3e418f9b4fb586c77c68c3552727"],["/bower_components/mat-ink/mat-ink-behavior.html","1984ec8e24c17c86d3430f6f32c8fe52"],["/bower_components/mat-ink/mat-ink-styles.html","de69ee22c289ffc5d3d6d57d102e0c97"],["/bower_components/mat-ink/mat-ink.html","1b77c0e54d7e9953afa109e6f026370d"],["/bower_components/mat-input-behavior/mat-input-behavior.html","f3ec70a80049e51e9c993e4ba4ba1def"],["/bower_components/mat-input-behavior/mat-input-styles.html","fce0ed6c6b41fd78e13fe3aa0f3c45b7"],["/bower_components/mat-item/mat-item.html","789541d272b077a5c5d32c9ff84353bd"],["/bower_components/mat-list/mat-list.html","8aada473cbb203a50747a118617dbe44"],["/bower_components/mat-menu/mat-menu.html","6a9d5b35a71d637d163f6728ace27312"],["/bower_components/mat-option/mat-option.html","c1d4349b2e14b638ed2523452c3dc39e"],["/bower_components/mat-palette/mat-palette.html","270b79a811996d850fd7483056bdb575"],["/bower_components/mat-paper/mat-paper-behavior.html","227b6fb0a4f397197c7b922b0627a9b7"],["/bower_components/mat-paper/mat-paper-styles.html","e994cd5a2d18333ae93a05c01fca3723"],["/bower_components/mat-pressed-behavior/mat-pressed-behavior.html","77b83a28d6fa1ac9bcb05b870614267e"],["/bower_components/mat-pressed-behavior/mat-pressed-ink-behavior.html","e0633f2ac72a7f699ee89d0b7a3e8f98"],["/bower_components/mat-pressed-behavior/mat-pressed-ink-styles.html","01ec14433f8b79e9141b8e9a83657249"],["/bower_components/mat-pressed-behavior/mat-pressed-paper-behavior.html","aeef76aa00c103e95f312a1856b5b701"],["/bower_components/mat-pressed-behavior/mat-pressed-paper-styles.html","8861a01e1cf96527a9fc3b9089d5f065"],["/bower_components/mat-pressed-behavior/mat-pressed-styles.html","17dd9cdc440eac985cb1a0d165444167"],["/bower_components/mat-raised-button/mat-raised-button.html","eae407581c646ff93cc58155ff2f005d"],["/bower_components/mat-ripple/mat-ripple.html","a018dc59de0466ed13f589a198b995d6"],["/bower_components/mat-shadow/mat-shadow.html","6c002abe5979634e11588da1e793c8b7"],["/bower_components/mat-snackbar/mat-snackbar.html","d3e454c50439240d80b7605479393ff0"],["/bower_components/mat-text-field/mat-text-field.html","170b9cff7053dd2eb1585e2ccd57e4e8"],["/bower_components/mat-toast/mat-toast.html","f8cf052c6961659196b054e0b775625d"],["/bower_components/mat-typography/mat-typography.html","9b1677a3d77d2c6d9eb26ed0e65638a7"],["/bower_components/neon-animation/animations/fade-in-animation.html","b814c818dbcffe2bb50563e1406497df"],["/bower_components/neon-animation/animations/fade-out-animation.html","44988226230af0e6d92f0988fc8688e2"],["/bower_components/neon-animation/animations/opaque-animation.html","8e2f63bbc648796f3ed96834a5553d07"],["/bower_components/neon-animation/neon-animatable-behavior.html","270f52231303cae4cb8e3fadb5a805c1"],["/bower_components/neon-animation/neon-animated-pages.html","8bb61cb8467f755163cec87e954425fc"],["/bower_components/neon-animation/neon-animation-behavior.html","eb1cdd9cd9d780a811fd25e882ba1f8e"],["/bower_components/neon-animation/neon-animation-runner-behavior.html","782cac67e6cb5661d36fb32d9129ff03"],["/bower_components/neon-animation/web-animations.html","b310811179297697d51eac3659df7776"],["/bower_components/paper-behaviors/paper-button-behavior.html","edddd3f97cf3ea944f3a48b4154939e8"],["/bower_components/paper-behaviors/paper-checked-element-behavior.html","59702db25efd385b161ad862b8027819"],["/bower_components/paper-behaviors/paper-inky-focus-behavior.html","51a1c5ccd2aae4c1a0258680dcb3e1ea"],["/bower_components/paper-behaviors/paper-ripple-behavior.html","b6ee8dd59ffb46ca57e81311abd2eca0"],["/bower_components/paper-button/paper-button.html","babcc1f776831ddfab647db31d97ee05"],["/bower_components/paper-drawer-panel/paper-drawer-panel.html","62b6d6beb252fdb633d7958df807bef6"],["/bower_components/paper-fab/paper-fab.html","b3ad1553b61ad0996a0371b94efe2164"],["/bower_components/paper-header-panel/paper-header-panel.html","b883923e580ff87f152f354358fc324b"],["/bower_components/paper-icon-button/paper-icon-button.html","2a75de00f858ae1e894ab21344464787"],["/bower_components/paper-input/paper-input-addon-behavior.html","de7b482dc1fb01847efba9016db16206"],["/bower_components/paper-input/paper-input-behavior.html","3960579058d3ba0a74ae7b67b78f53c2"],["/bower_components/paper-input/paper-input-char-counter.html","94c2003f281325951e3bf5b927a326bb"],["/bower_components/paper-input/paper-input-container.html","e3c61b8a6e35b134c99c09381ef48067"],["/bower_components/paper-input/paper-input-error.html","b90f3a86d797f1bdbbb4d158aeae06ab"],["/bower_components/paper-input/paper-input.html","3385511052db3467ca6ec155faa619ad"],["/bower_components/paper-item/paper-icon-item.html","17d1540072712073af1a84ae9b0ba06a"],["/bower_components/paper-item/paper-item-behavior.html","82636a7562fd8b0be5b15646ee461588"],["/bower_components/paper-item/paper-item-shared-styles.html","31466267014182098267f1b9303f656e"],["/bower_components/paper-item/paper-item.html","e614731572c425b144aa8a9da24ee3ea"],["/bower_components/paper-material/paper-material-shared-styles.html","d0eeeb696e55702f3a38ef1ad0058f59"],["/bower_components/paper-material/paper-material.html","47301784c93c3d9989abfbab68ec9859"],["/bower_components/paper-menu-button/paper-menu-button-animations.html","a6d6ed67a145ca00d4487c40c4b06273"],["/bower_components/paper-menu-button/paper-menu-button.html","328142b62de76860bca169b4a2d12342"],["/bower_components/paper-menu/paper-menu-shared-styles.html","9b2ae6e8b26011a37194ea3b4bdd043d"],["/bower_components/paper-menu/paper-menu.html","5270d7b4b603d9fdfcfdb079c750a3cd"],["/bower_components/paper-progress/paper-progress.html","f936509bf1fa2e404afc2a23413c943c"],["/bower_components/paper-ripple/paper-ripple.html","e22bc21b61184cb28125d16f9d80fb59"],["/bower_components/paper-spinner/paper-spinner-behavior.html","82e814c4460e8803f6f57cc457658adf"],["/bower_components/paper-spinner/paper-spinner-styles.html","a2122d2c0f3ac98e6911160d8886d31a"],["/bower_components/paper-spinner/paper-spinner.html","940e64bbde54dad918d0f5c0e247a8bd"],["/bower_components/paper-styles/color.html","430305db311431da78c0a6e1236f9ebe"],["/bower_components/paper-styles/default-theme.html","c910188e898624eb890898418de20ee0"],["/bower_components/paper-styles/shadow.html","7fd97f2613eb356e1bb901e37c3e8980"],["/bower_components/paper-styles/typography.html","bdd7f31bb85f116ce97061c4135b74c9"],["/bower_components/paper-tabs/paper-tab.html","395fdc6be051eb7218b1c77a94eff726"],["/bower_components/paper-tabs/paper-tabs-icons.html","9fb57777c667562392afe684d85ddbe2"],["/bower_components/paper-tabs/paper-tabs.html","2bf908cedd6ff6e67c18dbf337787bcc"],["/bower_components/paper-toast/paper-toast.html","f64d10724104f3751cae8b764aef56ff"],["/bower_components/paper-toggle-button/paper-toggle-button.html","a05d11b38ff158b663c307d9253564f4"],["/bower_components/paper-toolbar/paper-toolbar.html","ff99e4e6d522685e7e4d04f290e8ac9b"],["/bower_components/paper-tooltip/paper-tooltip.html","7b8150e196989eeba7c3f28727bb2f01"],["/bower_components/polymer/polymer-micro.html","7739e37db5581472b91925e5fa9bde55"],["/bower_components/polymer/polymer-mini.html","9e9dfb157eae29a59c98343288d4d120"],["/bower_components/polymer/polymer.html","2dce719d53b7ea549067d3d21a2b2c29"],["/bower_components/promise-polyfill/Promise.js","5afb14fd81497aca81bf25929d65b02d"],["/bower_components/promise-polyfill/promise-polyfill-lite.html","06470312beff013fc54695e0bdda2cb3"],["/bower_components/qwest/qwest.min.js","ac3ef2d3309ca02a1a5f8613e8e645cd"],["/bower_components/wave-player/wave-player.html","3a438dfec7ff907445f728e81a6fbf39"],["/bower_components/wavesurfer.js/dist/wavesurfer.min.js","2d13dd8d65bd74bf65a279eee0cc68f5"],["/bower_components/web-animations-js/web-animations-next-lite.min.js","607006febe5750df440a89ccd00660d7"],["/bower_components/web-socket/web-socket.html","2e33d752d9b17947cf11bb673bdfc30d"],["/bower_components/webcomponentsjs/webcomponents-lite.js","9dc13c1fee8c627a241d629d0ea8fd7b"],["/bower_components/xp-anchor-behavior/xp-anchor-behavior.html","2749a5724cb39d5d0b8376106d0096ad"],["/bower_components/xp-anchor-behavior/xp-anchor-styles.html","09af3aa103035f83a6f6eeb91794b945"],["/bower_components/xp-array-behavior/xp-array-behavior.html","124f7eeb081fefa2d07bf2fd25983d35"],["/bower_components/xp-breadcrumb-behavior/xp-breadcrumb-behavior.html","faab3452bc1bdf258ca62847a86702cf"],["/bower_components/xp-breadcrumb-behavior/xp-breadcrumb-step-behavior.html","a8b458b0a645e1bf764c0412eb9bd767"],["/bower_components/xp-breadcrumb-behavior/xp-breadcrumb-step-styles.html","11ede5a3b1c26c799d8b9512d5d8ff8b"],["/bower_components/xp-breadcrumb-behavior/xp-breadcrumb-styles.html","48e1bd8665489f6a3b78eb8ba4a2e43d"],["/bower_components/xp-dialog-behavior/xp-dialog-behavior.html","9bab0ff1a618e3f4056d080f1bd7207c"],["/bower_components/xp-dialog-behavior/xp-dialog-styles.html","8b6c59ec41110c9615a35ba085282d26"],["/bower_components/xp-finder-behavior/xp-finder-behavior.html","5d8bdd9e10a18377a15f0fb6af013b67"],["/bower_components/xp-focused-behavior/xp-focused-behavior.html","fdba3179a92879631dce4a7608031f0e"],["/bower_components/xp-focused-behavior/xp-focused-styles.html","11f5e1dd998a5c1dd9b97356cb436d1f"],["/bower_components/xp-icon-behavior/xp-icon-behavior.html","6c709413eebb98a27fdb1384d43eece7"],["/bower_components/xp-icon-behavior/xp-icon-styles.html","b49dcf77c6883e743c2a1e4320662979"],["/bower_components/xp-iconset/xp-iconset-finder.html","7046b504acef1fd4660e5b8c1f0f5237"],["/bower_components/xp-iconset/xp-iconset.html","9202aff7ada1128bff5f08c6339cc7ee"],["/bower_components/xp-input-behavior/xp-input-behavior.html","43c8fe9caf2262ede28dee696d0a5029"],["/bower_components/xp-input-behavior/xp-input-styles.html","a2019243c7d65a37526ae57032afbc77"],["/bower_components/xp-list-behavior/xp-list-behavior.html","dfe869c754cf438a33321e3d668f9082"],["/bower_components/xp-list-behavior/xp-list-styles.html","46a091fb7f9da268afa939a0b13830b7"],["/bower_components/xp-master-behavior/xp-master-behavior.html","1584c50f452513e47aef32150a41f747"],["/bower_components/xp-media-query/xp-media-query.html","b27c49e6a6561718c9502b18ef8af7d2"],["/bower_components/xp-menu-behavior/xp-menu-behavior.html","9a7c60c062228e6c1401e35762658d7b"],["/bower_components/xp-menu-behavior/xp-menu-styles.html","9d18ad6e4920f89145c7307f5470b8dd"],["/bower_components/xp-overlay/xp-overlay-behavior.html","877cf34bd3647116d514bf1d816d0a05"],["/bower_components/xp-overlay/xp-overlay-injector.html","efe8208e23fb932986fcfb10b097fe75"],["/bower_components/xp-overlay/xp-overlay-styles.html","26d90365776db2390b179e6f0ef67036"],["/bower_components/xp-pressed-behavior/xp-pressed-behavior.html","1d16d3b98279d053b73d2f4ae808541e"],["/bower_components/xp-pressed-behavior/xp-pressed-styles.html","c3b4c0456e51e7b7f50ce5eeedb8493a"],["/bower_components/xp-refirer-behavior/xp-refirer-behavior.html","aedf8f6bb0508877b0d3b6c804ac1445"],["/bower_components/xp-selector/xp-selector-behavior.html","e59ae4b5de4aa23376f5a33c5eea9602"],["/bower_components/xp-selector/xp-selector-multi-behavior.html","6ea62d7a96147e82154cd4bf1cd52645"],["/bower_components/xp-slave-behavior/xp-slave-behavior.html","2ccb48c6e1900a72a481e9b20219c96d"],["/bower_components/xp-snackbar-behavior/xp-snackbar-behavior.html","17104aecc64b52208942185490c79653"],["/bower_components/xp-snackbar-behavior/xp-snackbar-styles.html","ef47facb3afe82c407f805ac4dbba30e"],["/bower_components/xp-targeter-behavior/xp-targeter-behavior.html","14d4870a888f96c40eaf2722a489a0c8"],["/bower_components/xp-toast-behavior/xp-toast-behavior.html","92ec377c5f8776634680c3ccb4c13448"],["/bower_components/xp-toast-behavior/xp-toast-styles.html","deaf2913176ece878783ca6bc525cc0f"],["/index-web.js","af66fa9bda1affcb395605823a8ca696"],["/index.html","8f45a338ee950c8918a5f08aca79efd0"],["/src/auth0/login-auth0.html","55458e25dddb6134379ef2e096a1272c"],["/src/behaviors/async-behavior.html","d26417e96fec0753ede12e75e76ea34c"],["/src/behaviors/common-behavior.html","b9b29ccc39c065346e95b2f76d34bca4"],["/src/behaviors/endpoints-behavior.html","4a3c9e15e27c625796fa1b9be9a83c3c"],["/src/behaviors/map-behavior.html","06aaddc4bd0c9231702521f672fbd851"],["/src/behaviors/nav-behavior.html","ca1e1037542293d5a66a833edae2dd86"],["/src/behaviors/util-behavior.html","9f9aaf49517494381d874413e676ecbd"],["/src/configuration/configuration-page.html","5bbb7e0ab124c4c7ea8c27e11630364e"],["/src/datasources/datasources-page.html","e9b99eeca7f25adb40a03c82409b09b8"],["/src/detail-view/detail-view.html","cf03f165de95fa77b1f2b557c455eb33"],["/src/empty-state/empty-state.html","07d1228e38f85b46b5fafbf3def7f1ba"],["/src/filter-converter/filter-converter.html","ec63ab1935a02a22d7226808e9146840"],["/src/kazoup-icons/kazoup-icons.html","917742a2941e80571572753e8da800c8"],["/src/onboarding/onboarding-page.html","4f8dda2e3d3b5832937de91eb87cd94d"],["/src/settings/settings-page.html","f22e33dbdcc3b0301a1387f7ace6f5c5"],["/src/sockets/check-connection.html","a5aed65052069ccf62ecf17a889b8ec1"],["/src/sockets/notify-messages.html","96a2053b7dbbe96ce7b871487dd1b458"],["/src/static/auth.js","141accb6b2f933866e1a275a34564e8c"],["/src/styles/theme.css","c5d7fd062c2b14421924d3439329e220"],["/src/styles/theme.html","b0d93a56efc455de05def6667d534b93"],["/src/third-party-actions/create-file.html","f6d554413b6dd76e4b94842e7f3b3eb2"],["/src/third-party-actions/delete-file.html","72a14a8f29601d47691036b81c4623ea"],["/src/third-party-actions/share-box-file.html","e2621001a91290e226186f1acd16583a"],["/src/third-party-actions/share-dropbox-file.html","15a2ed2c8822378b772998a315f8c937"],["/src/third-party-actions/share-google-drive-file.html","9d680a6578243172be3f4cfa5a498a75"],["/src/third-party-actions/share-one-drive-file.html","56b4f1b67b934f2d59fd732a282be7cc"],["/src/third-party-actions/share-slack-file.html","1693680ef5ab556138fd8b25339b7d8f"],["/src/time/time-ago.html","9c45d0f52e31e0c5d0b8fcc539a46c3e"],["/src/web-app/web-app.html","7407c1e34c8550f1dd7c8e782ebc86e2"]];
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




