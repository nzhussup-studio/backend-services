# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

## [1.6.0](https://github.com/nzhussup-studio/backend-services/compare/v1.5.0...v1.6.0) (2026-11-10)


### Features

* account for redirect cases and add health endpoints where needed ([00c95cc](https://github.com/nzhussup-studio/backend-services/commit/00c95cc83d7b0ccd880b963e14b77c186d9e3215))
* add configuration endpoint for llm ([0d7d303](https://github.com/nzhussup-studio/backend-services/commit/0d7d303a8dd923d4fa660cb181c179ec8d71831f))
* add custom login page ([bc43050](https://github.com/nzhussup-studio/backend-services/commit/bc43050917955d94f23ecabd81b2612349a9e855))
* add github fetched descriptions into the base-service ([0459cb7](https://github.com/nzhussup-studio/backend-services/commit/0459cb76d7bda398896137d07b282626ccbe2e65))
* add optional issuer in certificates ([bc6c20e](https://github.com/nzhussup-studio/backend-services/commit/bc6c20ec8945a437acdc418238c72a3ede64ca89))
* remove .com references ([0134ee6](https://github.com/nzhussup-studio/backend-services/commit/0134ee65d68913b3c92669859da0d2e3a3d8a1ee))
* update default llm model ([158e519](https://github.com/nzhussup-studio/backend-services/commit/158e51905be29c554f5f8fba60b5400defd04ca7))


### Bug Fixes

* add another cors url for envoy to support 2 frontends simultaneously ([bb159f5](https://github.com/nzhussup-studio/backend-services/commit/bb159f5e1f8fabd0d55245822b3c9e76d7bfd70f))
* caching empty project desc responses ([a2d04e0](https://github.com/nzhussup-studio/backend-services/commit/a2d04e0167004f495e79e792f6412abbea9221d7))
* change default model ([6da9d2f](https://github.com/nzhussup-studio/backend-services/commit/6da9d2f128aa050c84a2d6ee35c02d4354752853))
* concurrency bug in image service ([2676d5f](https://github.com/nzhussup-studio/backend-services/commit/2676d5f708dc39ee935a2475f2ade561179b8bed))
* fix concurrency multipart error with introducing temporary save ([a03edfd](https://github.com/nzhussup-studio/backend-services/commit/a03edfdc2646be2b9e07c328a9323f9670ee82f8))
* keycloak dockerfile ([80a1666](https://github.com/nzhussup-studio/backend-services/commit/80a166609eef45b09c7210d0de20abb1b8b8d609))
* llm-service html validation ([4ffad98](https://github.com/nzhussup-studio/backend-services/commit/4ffad985637047e6e75d1d80796c73124443ef7a))
* replace KZ lang code to KK ([3a59fa6](https://github.com/nzhussup-studio/backend-services/commit/3a59fa601fd68e05cf1b256de888e198ad33ba0c))

## [1.5.0](https://github.com/nzhussup-studio/backend-services/compare/v1.4.0...v1.5.0) (2026-04-02)


### Features

* persist cv generator preferences ([749f871](https://github.com/nzhussup-studio/backend-services/commit/749f871e1f88b594d32f617e1110bc3c327ae249))

## [1.4.0](https://github.com/nzhussup-studio/backend-services/compare/v1.3.0...v1.4.0) (2026-03-27)


### Features

* add account-service ([353a067](https://github.com/nzhussup-studio/backend-services/commit/353a067176edf68b9c99a2177a0d8b64afe55b98))
* add account-service as a thin wrapper for keycloak delete service ([2759bbc](https://github.com/nzhussup-studio/backend-services/commit/2759bbc9e8a630b454bbfb88e89c5fdd57894e14))
* add security config into nginx ([5075dfe](https://github.com/nzhussup-studio/backend-services/commit/5075dfeab313563d435292b7de2447c68d6c6052))


### Bug Fixes

* cleanup unused envvars ([7951a7f](https://github.com/nzhussup-studio/backend-services/commit/7951a7f33cc5d5868f62761d44a5f0c6d786c036))
* fix nginx config ([c96f357](https://github.com/nzhussup-studio/backend-services/commit/c96f357b5dedfee1c6bcce612cce14cb89e698e8))
* fix status code for non admin in image service ([8cd9eef](https://github.com/nzhussup-studio/backend-services/commit/8cd9eefcf2f0767e4d90615315a743627eec6385))
* json ([1f7d407](https://github.com/nzhussup-studio/backend-services/commit/1f7d40768767288d92876ff23fc39f38b7e41ede))
* remove resolver from nginx ([f760e53](https://github.com/nzhussup-studio/backend-services/commit/f760e53afa0c3fe054b30b1493f68ad66860ab4a))

## [1.3.0](https://github.com/nzhussup-studio/backend-services/compare/v1.2.0...v1.3.0) (2026-03-20)


### Features

* fix dev mode and add custom template based keycloak ([53b314e](https://github.com/nzhussup-studio/backend-services/commit/53b314e6fe09e75e00cebf4a9c224116284d77e3))
* make keycloak only build and deploy ([4f8b390](https://github.com/nzhussup-studio/backend-services/commit/4f8b39058a098cd9383ff89e4ff18618269f9f30))


### Bug Fixes

* adjust example env ([f91577a](https://github.com/nzhussup-studio/backend-services/commit/f91577a75e8cd3808ba6827f8b39336a493a9622))
* cicd ([82ff8e4](https://github.com/nzhussup-studio/backend-services/commit/82ff8e49849e9584dead830cc73b6457457a6131))
* delete go cache ([f549766](https://github.com/nzhussup-studio/backend-services/commit/f5497666f19fa9dd113fd6e6c75f1710ef5a8a6f))
* delete go cache ([80a16eb](https://github.com/nzhussup-studio/backend-services/commit/80a16eb508e4a7fd066f0f49eafe0e959947b783))
* rename keycloak to be consistent with namings of external services ([ed0fb28](https://github.com/nzhussup-studio/backend-services/commit/ed0fb2815a98a2bdae13903d8ecec3e615709a9b))
* replace url ([18cc8d9](https://github.com/nzhussup-studio/backend-services/commit/18cc8d9dd24a0885bf1665097f8c0702ea3a0c18))

## [1.2.0](https://github.com/nzhussup-studio/backend-services/compare/v1.1.0...v1.2.0) (2026-03-19)


### Features

* add dev env ([29434f3](https://github.com/nzhussup-studio/backend-services/commit/29434f3aad6e20f043b45ad9cd26fe70131ae803))
* add release guide ([67804be](https://github.com/nzhussup-studio/backend-services/commit/67804be0a9f3ffe85194644cb9e5fd215cb1a4ee))


### Bug Fixes

* prepare for dev env ([25b7512](https://github.com/nzhussup-studio/backend-services/commit/25b75128b146587c18b41a26450b10558568b9e1))
* reference redis directly from opensource image ([56bb00e](https://github.com/nzhussup-studio/backend-services/commit/56bb00e657dbf35684c5843e0d2a6cc5899d56f3))

## 1.1.0 (2026-03-19)


### Features

* add openapi.yaml generator ([1720a18](https://github.com/nzhussup-studio/backend-services/commit/1720a1820a446b4c8be5f1123adb12ecefaf6353))
* add readme and make script executable ([a9135e9](https://github.com/nzhussup-studio/backend-services/commit/a9135e976f462505321294e84010f9b0e5df5d46))
* add release cicd ([67387c7](https://github.com/nzhussup-studio/backend-services/commit/67387c74ba42f0b4c2a3213e3ce825d21d5294ed))


### Bug Fixes

* add backend services from admin panel ([d10598b](https://github.com/nzhussup-studio/backend-services/commit/d10598bb6a7b9969ee6b701a1e03bb7c171f5a5a))
* cicd ([a1fb5a2](https://github.com/nzhussup-studio/backend-services/commit/a1fb5a202888a09064bba19cf2e172da45e475c7))
* do single platform build ([1b7cc14](https://github.com/nzhussup-studio/backend-services/commit/1b7cc14c4855c7fff7e2a41d2556f9f551edcb8b))
* fix cicd ([6c2b339](https://github.com/nzhussup-studio/backend-services/commit/6c2b339240166a825afbb6068742dbb3287e7e5e))
* fix cicd ([dc4d372](https://github.com/nzhussup-studio/backend-services/commit/dc4d372f5fb9c2872c133cbe18068357283e0b9e))
* unify cicd and release ([cb5fa2d](https://github.com/nzhussup-studio/backend-services/commit/cb5fa2dbc3ea7072418c30155736622877bc068f))
