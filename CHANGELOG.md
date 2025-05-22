# Changelog

## [2.1.0](https://github.com/Snipa22/go-tari-grpc-lib/compare/v2.0.2...v2.1.0) (2025-05-22)


### Features

* **cmd:** add walletPaymentSender ([dd21911](https://github.com/Snipa22/go-tari-grpc-lib/commit/dd21911e3d7fc5c0618ab237c880eb641d0499fe))
* **cmd:** add walletUTXOCreator ([53c543f](https://github.com/Snipa22/go-tari-grpc-lib/commit/53c543f2d3916174a767e2cae6bfec17f63de25b))
* **nodeGRPC:** add getheaderbyhash ([85b247f](https://github.com/Snipa22/go-tari-grpc-lib/commit/85b247f9d797dae6feb7087acc08d7a416d96536))


### Bug Fixes

* **blockWinners:** parse flags ([7a107e7](https://github.com/Snipa22/go-tari-grpc-lib/commit/7a107e77ccda14ecf99cc22cd639a9eeb3f91703))
* **walletUTXOCreator:** parse flags ([dcef6ae](https://github.com/Snipa22/go-tari-grpc-lib/commit/dcef6ae23f16f3d8d57e47de4c746946bd1540cd))

## [2.0.2](https://github.com/Snipa22/go-tari-grpc-lib/compare/v2.0.1...v2.0.2) (2025-05-14)


### Bug Fixes

* **submodule:** bump imports ([e45d637](https://github.com/Snipa22/go-tari-grpc-lib/commit/e45d63705673288a95152dc0bdfe7c98b9936c1f))

## [2.0.1](https://github.com/Snipa22/go-tari-grpc-lib/compare/v2.0.0...v2.0.1) (2025-05-14)


### Bug Fixes

* **submodule:** bump to v2 ([b6da62e](https://github.com/Snipa22/go-tari-grpc-lib/commit/b6da62e111b4489faf3d25df2bc21f14cbfec516))

## [2.0.0](https://github.com/Snipa22/go-tari-grpc-lib/compare/v1.0.0...v2.0.0) (2025-05-14)


### ⚠ BREAKING CHANGES

* **wallet:** add wallet transactions

### Features

* **baseNodeGrpc:** adds GetNetworkState ([0976586](https://github.com/Snipa22/go-tari-grpc-lib/commit/0976586742263e566b6246cf637316a333da88ee))
* **baseNodeGrpc:** adds GetNewBlock ([eee0898](https://github.com/Snipa22/go-tari-grpc-lib/commit/eee0898570e2fa2a39539fcdeb97a03279b62480))
* **blockFindCalculator:** add basic calculator ([bd8602d](https://github.com/Snipa22/go-tari-grpc-lib/commit/bd8602d276840a565e696e6ebbe03586a25fe587))
* **blockWinners:** add a cli flag for depth ([8f2b3e2](https://github.com/Snipa22/go-tari-grpc-lib/commit/8f2b3e2a880049c42937dd7b537ec81a18a9a668))
* **blockWinners:** add blockwinners debug tool ([0b13a88](https://github.com/Snipa22/go-tari-grpc-lib/commit/0b13a8847806757c76191685908a85c62d0ae7a2))
* **init:** initial commit moving over the grpc code ([b99892d](https://github.com/Snipa22/go-tari-grpc-lib/commit/b99892d383183fa264bb1f7c9bda0dc5bc9fedb0))
* **wallet:** add wallet transaction grabber ([8763448](https://github.com/Snipa22/go-tari-grpc-lib/commit/87634489a0d02f5ad4964695e9428f73c68b44e6))
* **wallet:** add wallet transactions ([31c32d5](https://github.com/Snipa22/go-tari-grpc-lib/commit/31c32d5825b4d0a878cb940f06ce30e55fa0185b))


### Bug Fixes

* **baseNodeGrpc:** GetBlockWithCoinbases GRPC upstream ([acc1592](https://github.com/Snipa22/go-tari-grpc-lib/commit/acc1592d3c0ea6140a48b9106cace1b9b5a26801))
* **blockWinners:** add 0 mode to scan entire chain ([67fe80e](https://github.com/Snipa22/go-tari-grpc-lib/commit/67fe80ee0a3c57495f23ace06e94ff7b1fa9daa6))
* **blockWinners:** clarify unknown pools ([58082ec](https://github.com/Snipa22/go-tari-grpc-lib/commit/58082ec6ad4d99c5bdaac823e61e5a32c314a93d))
* **blockWinners:** handle hidden pagination ([cc42d81](https://github.com/Snipa22/go-tari-grpc-lib/commit/cc42d8168a368c3049520ff0f4c77a571d9a2ceb))
* **blockWinners:** Improve the print data ([1df6dfb](https://github.com/Snipa22/go-tari-grpc-lib/commit/1df6dfb64810d95f88ba8412c28d5c0eb98295d4))
* **blockWinners:** process rx and printable chars only ([60befab](https://github.com/Snipa22/go-tari-grpc-lib/commit/60befab93f4e11b04ad543736862d11ca2cf7216))
* **module:** correct module naming for github ([0282b1b](https://github.com/Snipa22/go-tari-grpc-lib/commit/0282b1bc7ddbe000136d440e3655de0e6a41cf2f))
* **proto:** correct go_package ([d2c87fb](https://github.com/Snipa22/go-tari-grpc-lib/commit/d2c87fbc919a760c65702b7ed63dc728ec71151c))

## [1.0.0](https://github.com/Snipa22/go-tari-grpc-lib/compare/v0.0.1...v1.0.0) (2025-05-14)


### ⚠ BREAKING CHANGES

* **wallet:** add wallet transactions

### Features

* **baseNodeGrpc:** adds GetNetworkState ([0976586](https://github.com/Snipa22/go-tari-grpc-lib/commit/0976586742263e566b6246cf637316a333da88ee))
* **baseNodeGrpc:** adds GetNewBlock ([eee0898](https://github.com/Snipa22/go-tari-grpc-lib/commit/eee0898570e2fa2a39539fcdeb97a03279b62480))
* **blockFindCalculator:** add basic calculator ([bd8602d](https://github.com/Snipa22/go-tari-grpc-lib/commit/bd8602d276840a565e696e6ebbe03586a25fe587))
* **blockWinners:** add a cli flag for depth ([8f2b3e2](https://github.com/Snipa22/go-tari-grpc-lib/commit/8f2b3e2a880049c42937dd7b537ec81a18a9a668))
* **blockWinners:** add blockwinners debug tool ([0b13a88](https://github.com/Snipa22/go-tari-grpc-lib/commit/0b13a8847806757c76191685908a85c62d0ae7a2))
* **wallet:** add wallet transaction grabber ([8763448](https://github.com/Snipa22/go-tari-grpc-lib/commit/87634489a0d02f5ad4964695e9428f73c68b44e6))
* **wallet:** add wallet transactions ([31c32d5](https://github.com/Snipa22/go-tari-grpc-lib/commit/31c32d5825b4d0a878cb940f06ce30e55fa0185b))


### Bug Fixes

* **baseNodeGrpc:** GetBlockWithCoinbases GRPC upstream ([acc1592](https://github.com/Snipa22/go-tari-grpc-lib/commit/acc1592d3c0ea6140a48b9106cace1b9b5a26801))
* **blockWinners:** add 0 mode to scan entire chain ([67fe80e](https://github.com/Snipa22/go-tari-grpc-lib/commit/67fe80ee0a3c57495f23ace06e94ff7b1fa9daa6))
* **blockWinners:** clarify unknown pools ([58082ec](https://github.com/Snipa22/go-tari-grpc-lib/commit/58082ec6ad4d99c5bdaac823e61e5a32c314a93d))
* **blockWinners:** handle hidden pagination ([cc42d81](https://github.com/Snipa22/go-tari-grpc-lib/commit/cc42d8168a368c3049520ff0f4c77a571d9a2ceb))
* **blockWinners:** Improve the print data ([1df6dfb](https://github.com/Snipa22/go-tari-grpc-lib/commit/1df6dfb64810d95f88ba8412c28d5c0eb98295d4))
* **blockWinners:** process rx and printable chars only ([60befab](https://github.com/Snipa22/go-tari-grpc-lib/commit/60befab93f4e11b04ad543736862d11ca2cf7216))
