# Changelog

## [2.4.0](https://github.com/Snipa22/go-tari-grpc-lib/compare/v2.3.0...v2.4.0) (2025-06-20)


### Features

* **cmd:** add auto utxo gen ([1e26c75](https://github.com/Snipa22/go-tari-grpc-lib/commit/1e26c7584e66c2292f40f920104285f73c46c79a))
* **cmd:** add auto wallet rebooter ([b13b882](https://github.com/Snipa22/go-tari-grpc-lib/commit/b13b8827d9ea35a897a63b57544724a532c1aa18))
* **cmd:** add get node identity bin ([266120f](https://github.com/Snipa22/go-tari-grpc-lib/commit/266120f31ed082428428ece152e68c7353f625a6))
* **cmd:** add node auto reboot ([00fa216](https://github.com/Snipa22/go-tari-grpc-lib/commit/00fa2162aba51af5025ed93ea6b7bb0ef6877660))
* **cmd:** add wallet base node setter ([830e792](https://github.com/Snipa22/go-tari-grpc-lib/commit/830e7928be025cf518e4848e9dfc46637b6a78ac))
* **cmd:** add wallet validator ([d75530c](https://github.com/Snipa22/go-tari-grpc-lib/commit/d75530c02725c693c3b64fceaf190d4f20fc24fa))
* **cmd:** add walletSweeper ([e2e5538](https://github.com/Snipa22/go-tari-grpc-lib/commit/e2e553802eccb344c4db01a06ddb3d938105e140))
* **getNetworkStats:** add root reward ([673aab4](https://github.com/Snipa22/go-tari-grpc-lib/commit/673aab49d81bf35f9abef5a6ba191785169729a8))
* **walletGRPC:** add address grabber ([fa8f1b5](https://github.com/Snipa22/go-tari-grpc-lib/commit/fa8f1b5d26c5f7e929a9b076ce4778677b3db58f))
* **walletSweeper:** safety checks ([00f4b29](https://github.com/Snipa22/go-tari-grpc-lib/commit/00f4b299f1eebb6c00a628a3bf9959dafe1e0811))


### Bug Fixes

* **blockWinners:** handle short ids ([5bd8687](https://github.com/Snipa22/go-tari-grpc-lib/commit/5bd8687b4bf75925ebaf2fc7975ecdf0110123b6))
* **grpc:** allow larger response ([8be3791](https://github.com/Snipa22/go-tari-grpc-lib/commit/8be37918767f5676482ebc5782db78ce7ed2c933))
* **smartutxo:** add state prints ([83cbd82](https://github.com/Snipa22/go-tari-grpc-lib/commit/83cbd82a05f211f7ed4a0aa50acafa382ee2d0a8))
* **smartutxo:** inf-loop + run at startup ([ddc8276](https://github.com/Snipa22/go-tari-grpc-lib/commit/ddc82767ae427567d684b3ef24bbc519d23263d5))
* **statsExporter:** Add caching to the system. ([d8e5eeb](https://github.com/Snipa22/go-tari-grpc-lib/commit/d8e5eeb9bea088b24d48d8ba287bddb214283e0f))
* **walletbooter:** add debug logging ([810318d](https://github.com/Snipa22/go-tari-grpc-lib/commit/810318d5fa835968f23397479dba37d963a846bd))
* **walletSweeper:** get addresses in base58 ([eb509fc](https://github.com/Snipa22/go-tari-grpc-lib/commit/eb509fc43c18073abed0020d130f3741a91aa747))

## [2.3.0](https://github.com/Snipa22/go-tari-grpc-lib/compare/v2.2.1...v2.3.0) (2025-05-27)


### Features

* **blockWinners:** scan for coinbases ([3b167b5](https://github.com/Snipa22/go-tari-grpc-lib/commit/3b167b59252dc613725ebd79cf082a44e6eae8a3))
* **statsExporter:** add new stats exporter ([a465a3a](https://github.com/Snipa22/go-tari-grpc-lib/commit/a465a3a0f0244fb21ad659ebc44d0d86b33473ff))
* **waller:** add balance exporter ([93c9c45](https://github.com/Snipa22/go-tari-grpc-lib/commit/93c9c45a8f596cf44049e515516975b1c899b710))
* **walletGRPC:** add get txinfo by id ([e336f99](https://github.com/Snipa22/go-tari-grpc-lib/commit/e336f99e4acfe6b1d39e8f35d24b5844ba284068))


### Bug Fixes

* **blockWinners:** break to exit ([ca60e63](https://github.com/Snipa22/go-tari-grpc-lib/commit/ca60e63a0887c6e85afcd93bd3249ac2d7116745))
* **blockWinners:** correct the bucket ([d82586b](https://github.com/Snipa22/go-tari-grpc-lib/commit/d82586b05ef7e8eeb9cd62e33c1055cfaebe8c5a))
* **blockWinners:** make sure the append bucket is right ([d670627](https://github.com/Snipa22/go-tari-grpc-lib/commit/d670627b5274d2598e39e115cb3a61ed4df5216b))
* **blockWinners:** pickup the pool shim ([e830763](https://github.com/Snipa22/go-tari-grpc-lib/commit/e830763192f4a79cf3ee51255a5ece4a2195893e))
* **blockWinners:** update for RXT ([b2e0973](https://github.com/Snipa22/go-tari-grpc-lib/commit/b2e097388d0919426d0376bf894c4103dfe476aa))
* **getNetworkStats:** rxt bt goes into rxt diff ([ca8b549](https://github.com/Snipa22/go-tari-grpc-lib/commit/ca8b5499b46276f3f6dd9ddff7646e82d9c0d68b))
* **statsExporter:** export sha3x to the right var ([0d2a587](https://github.com/Snipa22/go-tari-grpc-lib/commit/0d2a5870efc7e9d92259c285b329211b41a5482f))

## [2.2.1](https://github.com/Snipa22/go-tari-grpc-lib/compare/v2.2.0...v2.2.1) (2025-05-22)


### Bug Fixes

* **clients:** remove flag ([69050ea](https://github.com/Snipa22/go-tari-grpc-lib/commit/69050ead4190800f3dc1b83e8bf404a9d63b3cca))

## [2.2.0](https://github.com/Snipa22/go-tari-grpc-lib/compare/v2.1.0...v2.2.0) (2025-05-22)


### Features

* **generated:** new tari protos ([1f23fae](https://github.com/Snipa22/go-tari-grpc-lib/commit/1f23fae4906f8cc9f70c8e1d2019b0831f5e9156))

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
