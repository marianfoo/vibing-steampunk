# Changelog

All notable changes to this project will be documented in this file.
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Bug Fixes

- Add missing `items` to DebuggerGetVariables array schema (#24) (#25)([`9a7eebe`](https://github.com/oisee/vibing-steampunk/commit/9a7eebe9e31fe11abc03d0cfd799cb4dd7ee907b))
- Auto-retry on 401 Unauthorized after idle timeout (#35)([`d73460a`](https://github.com/oisee/vibing-steampunk/commit/d73460ade7035903fe638b4caf0500c64ef2a776))


### Features

- Add GetDependencyZIP function and tests for dependency retrieval (#60)([`5317105`](https://github.com/oisee/vibing-steampunk/commit/531710515c939cf6e3dbe8d67a5a79e4a07e033a))


### Tests

- Add unit test for DebuggerGetVariables tool schema validation([`d737629`](https://github.com/oisee/vibing-steampunk/commit/d737629639c26940e2fb702d54bd7f264c11e3b4))



## [2.27.0] - 2026-03-01
### Documentation

- Add article plan for "VSP Two Months Later"([`e894394`](https://github.com/oisee/vibing-steampunk/commit/e894394beac1c823281dc5342bd965a454402623))
- Add architecture diagrams (Mermaid) to docs/([`7505d84`](https://github.com/oisee/vibing-steampunk/commit/7505d84334c0082a975db82d89ea9a8b4e162aed))
- Add CLI coding agents guide (8 agents, 4 languages)([`24b6334`](https://github.com/oisee/vibing-steampunk/commit/24b6334d0770848c125f48ab9af11f0c434f3f83))
- Overhaul --help and fix stale tool counts (19/45 → 81/122)([`c4bf5e8`](https://github.com/oisee/vibing-steampunk/commit/c4bf5e824c2c4bbf332711d2e36487c9b815dc03))
- Add reviewer guide with 8 hands-on tasks([`b6d025f`](https://github.com/oisee/vibing-steampunk/commit/b6d025faba4a6a090f01f610e4a3b05458ddda75))
- Link reviewer guide from README top([`4866e07`](https://github.com/oisee/vibing-steampunk/commit/4866e0760ef675ae015ff6ed2b713b4cbbf3911d))
- Fix stale tool counts across CLAUDE.md and README (81/122)([`53a3fcb`](https://github.com/oisee/vibing-steampunk/commit/53a3fcbcc9d733c2d1dba72ff5bd76126cc51117))


### Features

- Iterative activation with package filtering + 100 stars article([`8d2c343`](https://github.com/oisee/vibing-steampunk/commit/8d2c343e50f79f48663418568deade412337cd03))


### Refactoring

- Refactoring zcl_vsp_tadir_move=>

zadt_c;_tadir_move=> zcl_vsp_tadir_move([`a94248c`](https://github.com/oisee/vibing-steampunk/commit/a94248cbda8f64f5e282af20d789535b7b2fbfba))



## [2.26.0] - 2026-02-04
### Bug Fixes

- PackageExists fails for local packages with $ in name([`83e8626`](https://github.com/oisee/vibing-steampunk/commit/83e86269f56eb5a3d6983385de3ff5276083d31e))


### Documentation

- Add transportable packages configuration to README([`029498a`](https://github.com/oisee/vibing-steampunk/commit/029498ad803159456bed07cec4773576a2262aab))
- Add RAP OData service creation guide to README([`a3dac1d`](https://github.com/oisee/vibing-steampunk/commit/a3dac1d245265d5404413da4fdaa0a88a4faf926))


### Tests

- Add namespace integration tests for Issue #18([`1ddbb22`](https://github.com/oisee/vibing-steampunk/commit/1ddbb22a980532a1b934f6adbde55e6be7ce42ac))
- Fix namespace and integration test issues([`78e5671`](https://github.com/oisee/vibing-steampunk/commit/78e5671ad7e088bbe4bf10468948bb6abb0fc555))



## [2.25.0] - 2026-02-03
### Bug Fixes

- Namespace URL encoding for all ADT operations([`59b4b90`](https://github.com/oisee/vibing-steampunk/commit/59b4b9061497d86fb6e599e5b37382edee865a1e))


### Features

- Allow transportable package creation with --enable-transports([`e483537`](https://github.com/oisee/vibing-steampunk/commit/e483537958dfd7243abfbce8be37214d0abe8ac2))
- CreatePackage software_component + viper env var fix([`c18309b`](https://github.com/oisee/vibing-steampunk/commit/c18309b0b9e14d90cd65e00eb2f77595a0d0f7cd))



## [2.24.0] - 2026-02-03
### Features

- V2.23.0 - GitExport to disk, GetAbapHelp via WebSocket([`ddf5c22`](https://github.com/oisee/vibing-steampunk/commit/ddf5c22f84ebdd9fbcfc5dcf771989487106af7f))
- V2.24.0 - Transportable Edits Safety Feature([`3a9b0b0`](https://github.com/oisee/vibing-steampunk/commit/3a9b0b0bea276e7ca9ae556a55cc710fd5a44831))



## [2.23.0] - 2026-02-02
### Chores

- Add sessions/ to gitignore([`13964b8`](https://github.com/oisee/vibing-steampunk/commit/13964b842a14ad665f823ff7e87609c7fdc4e32d))
- Disable tool aliases by default - reduce tool bloat([`4359cdc`](https://github.com/oisee/vibing-steampunk/commit/4359cdc9879675dc928d7020204cb5f405ad169c))


### Features

- Add granular tool visibility control via .vsp.json([`f8fd717`](https://github.com/oisee/vibing-steampunk/commit/f8fd717c0acbd62590aec602e88efc618be13d77))
- Add GetAbapHelp tool for ABAP keyword documentation (#10)([`434ed5e`](https://github.com/oisee/vibing-steampunk/commit/434ed5e83240cf52f3be334930c5b8602071c0cf))
- Add Level 2 GetAbapHelp - real docs from SAP system via ZADT_VSP([`b78803d`](https://github.com/oisee/vibing-steampunk/commit/b78803d339f76b2d3b92de4276cabcec106dc30a))
- GitExport saves ZIP to disk, GetAbapHelp uses amdpWSClient([`7c01351`](https://github.com/oisee/vibing-steampunk/commit/7c01351a783ca7588424a65c2fa64e2c21bce794))


### Sync

- Update embedded ABAP from SAP system([`a3468f7`](https://github.com/oisee/vibing-steampunk/commit/a3468f7dd454bd86091875f384fbc47f605e0360))



## [2.22.0] - 2026-02-01
### Bug Fixes

- Transport API 406 error and EditSource transport support([`c726bfe`](https://github.com/oisee/vibing-steampunk/commit/c726bfeb08d43357622853a4fa7d34d58a01469b))
- Honor HTTP_PROXY/HTTPS_PROXY environment variables (#13)([`a1af66f`](https://github.com/oisee/vibing-steampunk/commit/a1af66f83ad050a0799442c75645861c9a5ba680))


### Documentation

- Replace real terminal ID with random example([`b09f678`](https://github.com/oisee/vibing-steampunk/commit/b09f678a72eb158edcfa6eb6474abbb661a123ca))
- Update README for v2.22.0 release([`10c07f9`](https://github.com/oisee/vibing-steampunk/commit/10c07f99dd3e50714d63122287fff1d43ccf34d4))
- Add ABAP Help tool design report([`90278ec`](https://github.com/oisee/vibing-steampunk/commit/90278ec7ed92ada482b2475ce4b7c22ecae31bd4))


### Features

- Add MoveObject tool and refactor WebSocket code([`2d3d40c`](https://github.com/oisee/vibing-steampunk/commit/2d3d40cb472d4f0193f62870a5fcd172b35380cf))
- Add SAP_TERMINAL_ID config for SAP GUI breakpoint sharing([`677e7ce`](https://github.com/oisee/vibing-steampunk/commit/677e7cee84d456f5eb2b6009a4c47d9afcd7af31))


### Refactoring

- RunReport to use background jobs with spool output([`96d0117`](https://github.com/oisee/vibing-steampunk/commit/96d011709356aa7cb024aa502a7f9a607a2ed8ca))



## [2.21.0] - 2026-01-06
### Bug Fixes

- WebSocket reconnection check in report handlers([`52e17c9`](https://github.com/oisee/vibing-steampunk/commit/52e17c9d654607271bc923a47c863fff830ef0dd))
- Improve error handling in GetSystemInfo and CSRF fetch([`b9fb06b`](https://github.com/oisee/vibing-steampunk/commit/b9fb06b444a86c0057d26083d79176cee98a08eb))


### Features

- Add function module support to ImportFromFile([`c7997c0`](https://github.com/oisee/vibing-steampunk/commit/c7997c07105f1a35ac45e2fa1967bac56479762f))
- Add method-aware breakpoints with include resolution([`54417f6`](https://github.com/oisee/vibing-steampunk/commit/54417f6e9cdb06052332f81d0475aadbd83ea31f))
- Method-level source operations for GetSource, EditSource, WriteSource([`1fa5065`](https://github.com/oisee/vibing-steampunk/commit/1fa5065390f191fe1eeb4183d0a491c468082186))



## [2.20.0] - 2026-01-06
### Bug Fixes

- WebSocket client parameter order & mcp-to-vsp password sync([`29abb0c`](https://github.com/oisee/vibing-steampunk/commit/29abb0ce7e564720e165d528428e0618273750e5))
- Add .abapgit.xml to GitExport ZIP output([`93dc5ef`](https://github.com/oisee/vibing-steampunk/commit/93dc5ef05426d6ebdfbb1e96a5301711e0b08327))
- Use FULL folder logic for multi-package exports([`dafd1f5`](https://github.com/oisee/vibing-steampunk/commit/dafd1f52c6f4d55f92742a4b48d839fafdbdea6c))


### Chores

- Improve .gitignore for vsp config files([`3d3320a`](https://github.com/oisee/vibing-steampunk/commit/3d3320a8603b7e5ea9791bfd2fbff46159ea46f4))
- Use .vsp*.json wildcard pattern in .gitignore([`6e8c351`](https://github.com/oisee/vibing-steampunk/commit/6e8c351283ee21e704d81d4fb517aaee1c883c55))


### Documentation

- Document CLI mode and system profiles in README([`fcff096`](https://github.com/oisee/vibing-steampunk/commit/fcff096a5d33c45871738036add9bed7480aefc7))
- Update README for v2.20.0 CLI Mode release([`40210f9`](https://github.com/oisee/vibing-steampunk/commit/40210f9f0d931ad25d847e1ab07930a5de6f1524))


### Features

- Make sync-embedded for exporting ZADT_VSP from SAP([`ab47d27`](https://github.com/oisee/vibing-steampunk/commit/ab47d273b6c033e6cad98cc986eba877f4fc5f1b))
- CLI subcommands with system profiles([`cdab42c`](https://github.com/oisee/vibing-steampunk/commit/cdab42cb961d7bde5156e8a4e764daf5a94e20c8))
- Vsp config init/show commands([`bf90c25`](https://github.com/oisee/vibing-steampunk/commit/bf90c25b983caa4a7879112c887e65f7412467d1))
- Vsp config mcp-to-vsp and vsp-to-mcp commands([`717cd9a`](https://github.com/oisee/vibing-steampunk/commit/717cd9adb8909707c68a28cea1a1f8b954cd539c))
- Cookie authentication support in CLI system profiles([`d83080b`](https://github.com/oisee/vibing-steampunk/commit/d83080bbd466ad71cb97f1baae0b9b7f85049002))


### Refactoring

- Rename .vsp-systems.json to .vsp.json([`0447d00`](https://github.com/oisee/vibing-steampunk/commit/0447d00a2ce153befe99cc51ca14a41b7f0bae17))



## [2.19.1] - 2026-01-06
### Bug Fixes

- WebSocket TLS for self-signed certificates (#1)([`181f523`](https://github.com/oisee/vibing-steampunk/commit/181f52365c057a9aeb1c9184cf94ee4d34373b0e))


### Documentation

- Update roadmap for v2.19.1([`80de939`](https://github.com/oisee/vibing-steampunk/commit/80de939548c4f077d07abd450170fc6293a188c1))


### Features

- Tool aliases and heading texts support([`d29549a`](https://github.com/oisee/vibing-steampunk/commit/d29549a8ef29806639b9561d50ae1972435735e1))



## [2.19.0] - 2026-01-05
### Bug Fixes

- GetSystemInfo uses SQL fallback for reliability([`3c454a6`](https://github.com/oisee/vibing-steampunk/commit/3c454a6a3fd3d9f9e08e30aa9cdc49eebf2d24ef))


### Documentation

- DAP integration vision for ABAP debugging (Report 004)([`216e6cb`](https://github.com/oisee/vibing-steampunk/commit/216e6cbffe64aca56fa95401a500d22080328345))
- Add prioritized roadmap (Report 005)([`e71c491`](https://github.com/oisee/vibing-steampunk/commit/e71c491045639afd473ccd74eb09a0f1648206d1))
- Update README What's New for v2.17-v2.18([`3700518`](https://github.com/oisee/vibing-steampunk/commit/3700518a790dd7b846559fedc26c746a296202ac))
- Add missing research reports([`f912706`](https://github.com/oisee/vibing-steampunk/commit/f912706033d6d89f9c27ea7bdc298a46454b9cc8))
- DAP implementation plan (Report 2026-01-03-001)([`a73bc7b`](https://github.com/oisee/vibing-steampunk/commit/a73bc7b947ef59915831cf6e6e70a8ae36070624))
- Update README and roadmap for v2.19.0 debug CLI([`590e82e`](https://github.com/oisee/vibing-steampunk/commit/590e82edbdd6e3545004020f216b5039eda8d0b0))
- Update roadmap with v2.19.0 quick wins([`d5b0323`](https://github.com/oisee/vibing-steampunk/commit/d5b0323e8deae9cf8f3ae79d1bb31eb7023b27b1))
- V2.19.0 release documentation([`cf0f716`](https://github.com/oisee/vibing-steampunk/commit/cf0f71675f4f8e9e6d6f9156c2e8ce932e444e1b))


### Features

- Interactive CLI debugger (vsp debug)([`f1358e9`](https://github.com/oisee/vibing-steampunk/commit/f1358e9773e4b3f07ae32287126a5ceb3786cc94))
- Quick wins - GetMessages, ListDumps, ActivatePackage, X group([`2706797`](https://github.com/oisee/vibing-steampunk/commit/27067971ef521c7337257d0d534570f812f65be4))
- CreateTable tool + GetMessages fix([`a71ec42`](https://github.com/oisee/vibing-steampunk/commit/a71ec427e0548afdc572d78887aaae5eefa822e3))
- CompareSource, CloneObject, GetClassInfo tools([`8550435`](https://github.com/oisee/vibing-steampunk/commit/8550435b6bb82f0e9822cbed3772791788daa800))
- RunReportAsync and GetAsyncResult for background execution([`56dc11a`](https://github.com/oisee/vibing-steampunk/commit/56dc11af633cec85d13ddee46c2b149708c375b5))



## [2.18.0] - 2026-01-02
### Features

- WebSocket-based debugger tools via ZADT_VSP([`c3a3780`](https://github.com/oisee/vibing-steampunk/commit/c3a3780006c80c8d380d52ed3cfe41b60d25684e))
- Consolidate $ZADT_VSP package + lock cleanup fix([`5e4530a`](https://github.com/oisee/vibing-steampunk/commit/5e4530a4f3ea6f88acb3bb7e132078c531c1c4a5))
- Report execution tools + packageExists fix([`3df8955`](https://github.com/oisee/vibing-steampunk/commit/3df8955f110fd870ef24c98c7681865cbb6a0baf))



## [2.17.1] - 2025-12-24
### Bug Fixes

- Install tools upsert - proper package/object existence checks([`4505237`](https://github.com/oisee/vibing-steampunk/commit/450523755f3f9ad47151b1d0887e3d0bc4ee5d38))


### Documentation

- Cumulative progress article + ZADT_VSP self-deployment design([`b7686f8`](https://github.com/oisee/vibing-steampunk/commit/b7686f8bffa39975060b7712976ffb240b14c84c))


### Features

- InstallZADTVSP tool for one-command deployment([`1ee4962`](https://github.com/oisee/vibing-steampunk/commit/1ee496222403301e7db6615158d96b362c20aa07))
- InstallAbapGit tool + dependency embedding architecture([`a3f1fa0`](https://github.com/oisee/vibing-steampunk/commit/a3f1fa09960c7f554be5a9f919474d6690636bc5))



## [2.16.0] - 2025-12-23
### Documentation

- VSP v2.15 Possibilities Unlocked article([`6957372`](https://github.com/oisee/vibing-steampunk/commit/695737262ed04fdda0ba7671be53cfd1d01946f8))


### Features

- AbapGit WebSocket integration (Git domain)([`a73d2a6`](https://github.com/oisee/vibing-steampunk/commit/a73d2a6c9a9e797413a77c6ce61e2c4a1a5dfa45))
- Complete abapGit WebSocket integration (v2.16.0)([`78e2c6d`](https://github.com/oisee/vibing-steampunk/commit/78e2c6d16733a01cce29e2c7b4a7641bd1aba389))



## [2.15.1] - 2025-12-22
### Bug Fixes

- Correct unit test count 216 → 244([`c931533`](https://github.com/oisee/vibing-steampunk/commit/c93153344683b579f061766d9d5cbef557e79966))


### Documentation

- Close HTTP debugger topic - WebSocket ZADT_VSP is the solution([`89ea6bb`](https://github.com/oisee/vibing-steampunk/commit/89ea6bb275d868eabba08e14ee0312792d1f54e6))



## [2.15.0] - 2025-12-21
### Documentation

- Add RCA, Replay & Test Extraction section to README([`b7bf940`](https://github.com/oisee/vibing-steampunk/commit/b7bf940a6a26897e8d7fa8bf11e469191a5f9109))
- Test Extraction Implications - paradigm shift analysis([`81416ea`](https://github.com/oisee/vibing-steampunk/commit/81416ea3f4d52656318a6e0a441847da90f98fe0))
- Add Implications Analysis link to README([`90948b9`](https://github.com/oisee/vibing-steampunk/commit/90948b9e8c36a80d1a704a51b5539f88c91297f2))


### Features

- Variable History Recording (Phase 5.2)([`29e192d`](https://github.com/oisee/vibing-steampunk/commit/29e192d4c4510cd0b66204495547cae38da28888))
- Extended breakpoint types + Watchpoint Scripting (Phase 5.4)([`3dd20cd`](https://github.com/oisee/vibing-steampunk/commit/3dd20cd7b506264808dcec50ec649e6ee6351298))
- Force Replay - State Injection (Phase 5.5) - THE KILLER FEATURE([`70fb43f`](https://github.com/oisee/vibing-steampunk/commit/70fb43fe85da3d46759b40ef44321701a044a63d))
- Phase 5 TAS-Style Debugging Complete (v2.15.0)([`19405b2`](https://github.com/oisee/vibing-steampunk/commit/19405b2a4a13210f8809748d263f80f0524e4a61))



## [2.14.0] - 2025-12-21
### Documentation

- Add LinkedIn follow-up article on RCA tools([`fcb3b80`](https://github.com/oisee/vibing-steampunk/commit/fcb3b80a11b8ebae9abb4af05f7320806200cde2))
- TAS-style debugging and scripting vision document([`c9bd1f4`](https://github.com/oisee/vibing-steampunk/commit/c9bd1f4553666375305d1fe0847b1a67982cf41b))
- Test extraction and isolated replay design([`c57eee6`](https://github.com/oisee/vibing-steampunk/commit/c57eee6f6d87c1431bc754b7194062d6e1c4501a))
- Add VISION.md and ROADMAP.md, update README([`56f4660`](https://github.com/oisee/vibing-steampunk/commit/56f4660174f62363c53d20fb1c0917f5c80e7a9b))
- Phase 5 implementation plan (TAS-style debugging)([`a0f558b`](https://github.com/oisee/vibing-steampunk/commit/a0f558be48bcf6b244d71e0cba17f465ee6b494b))
- Force Replay - live state injection design([`cfaa542`](https://github.com/oisee/vibing-steampunk/commit/cfaa5425bca09621bb4574797dfdc912fe275081))


### Features

- Lua scripting integration (Phase 5.1)([`0e5c5c2`](https://github.com/oisee/vibing-steampunk/commit/0e5c5c2681fcca270d21a476139a387dfd73461a))



## [2.13.0] - 2025-12-21
### Bug Fixes

- External debugger breakpoint XML format & unit test parsing([`296b8f3`](https://github.com/oisee/vibing-steampunk/commit/296b8f31530810440db43eeb5609527bc9ec156c))
- GetDumps Accept header & add WebSocket debugging ADR([`2eb4a5e`](https://github.com/oisee/vibing-steampunk/commit/2eb4a5efd27241c866bc7a8c6234fa2f6471b7d5))


### Chores

- Prepare release v2.1.0([`bc9e768`](https://github.com/oisee/vibing-steampunk/commit/bc9e76871e7def766cf8da79e84f09d86b8a4fd6))


### Documentation

- External debugger investigation & batch API improvements([`2ee1d6d`](https://github.com/oisee/vibing-steampunk/commit/2ee1d6df9054473bd54f5f64649037771a68236f))
- AI-assisted RCA & ANST integration vision([`d04b1bb`](https://github.com/oisee/vibing-steampunk/commit/d04b1bb6a3ff306e313efafb9b3bed1c0654ff07))
- Mark RCA/ANST and WebSocket ADRs as proposals, not vanilla ADT([`c35fefe`](https://github.com/oisee/vibing-steampunk/commit/c35fefeaaade299ff9f4bdbd0c559dee023ed5d3))
- ADR-003 ZADT-VSP unified APC handler design([`0ccc60a`](https://github.com/oisee/vibing-steampunk/commit/0ccc60a7e9259e8f9daaaf6e238e6a7293aebf62))


### Features

- ZADT-VSP APC handler with RFC domain (ABAP)([`67e0024`](https://github.com/oisee/vibing-steampunk/commit/67e0024c750c4d6eae89c74067a7e5f8b0d16150))
- ZADT_VSP APC WebSocket handler - RFC domain operational([`c9109be`](https://github.com/oisee/vibing-steampunk/commit/c9109be2feb84a5bae21155e954997c4470dadfd))
- WebSocket RFC Handler (ZADT_VSP) with embedded ABAP source([`d36b1d6`](https://github.com/oisee/vibing-steampunk/commit/d36b1d6197154f38c97d33411c9ea3635f54e479))
- Add debug domain to WebSocket handler (ZADT_VSP)([`307d231`](https://github.com/oisee/vibing-steampunk/commit/307d23194918472feed5006c1d7340310a3c1d53))
- Full WebSocket debugging with TPDAPI integration (v2.0.0)([`fa4ada8`](https://github.com/oisee/vibing-steampunk/commit/fa4ada8b49c3ea504bb824abfa49ebab8a335b86))
- TPDAPI breakpoint integration verified working (v2.0.1)([`64050c6`](https://github.com/oisee/vibing-steampunk/commit/64050c600b2a793f2082ca25b7b8b35a75f9afd3))
- Add call graph traversal and RCA tools([`d8e3742`](https://github.com/oisee/vibing-steampunk/commit/d8e3742e3544c665b4c70386647a3fa12d3c5140))



## [2.12.6] - 2025-12-10
### Features

- EditSource support for class includes (testclasses, locals)([`3782380`](https://github.com/oisee/vibing-steampunk/commit/3782380101b3ba2edc155896c97ee580e40c786d))



## [2.12.5] - 2025-12-09
### Bug Fixes

- Normalize line endings in EditSource (CRLF → LF)([`fafbccf`](https://github.com/oisee/vibing-steampunk/commit/fafbccf304283dd44a698e26c987a3d8bd6214d7))


### Documentation

- V2.12.5 release notes - EditSource line ending fix([`6bd37bb`](https://github.com/oisee/vibing-steampunk/commit/6bd37bbb04dbde2a701fd43c0d49d91401cb85d4))



## [2.12.4] - 2025-12-09
### Features

- V2.12.4 - Feature Detection & Safety Network([`0d5693d`](https://github.com/oisee/vibing-steampunk/commit/0d5693d279e31e4f85c29d88584aa2b4300d9b04))



## [2.12.3] - 2025-12-08
### Bug Fixes

- Properly detect 404 in DeployFromFile for class includes([`d489743`](https://github.com/oisee/vibing-steampunk/commit/d489743dd965741466251447f47f54883c69f9d1))


### Features

- Auto-reconnect on SAP session timeout([`610bfeb`](https://github.com/oisee/vibing-steampunk/commit/610bfeb36e7680cbe977beee78707fc7dd634cd7))



## [2.12.2] - 2025-12-08
### Bug Fixes

- Extract class name from filename for class includes([`85fb919`](https://github.com/oisee/vibing-steampunk/commit/85fb919e58b12a00d875b6d592c4891c373b3169))



## [2.12.1] - 2025-12-07
### Documentation

- Add tagline and description to README header([`8023e5b`](https://github.com/oisee/vibing-steampunk/commit/8023e5b0a6a5c5c52650c4f9fbd8bd550420177f))
- Fix tagline - ADT↔MCP bridge, link to OData↔MCP([`0a2a1e6`](https://github.com/oisee/vibing-steampunk/commit/0a2a1e6d8e979d34d5563c322f2e9182e71e0be4))
- Add DSL mention to tagline([`b120f46`](https://github.com/oisee/vibing-steampunk/commit/b120f46e4cc377d8fda24743d87e0d1602fe92d5))


### Features

- Add CreatePackage tool to focused mode([`7452c48`](https://github.com/oisee/vibing-steampunk/commit/7452c484151fbfb3f57ca8d1dc79a7790ffb471b))



## [2.12.0] - 2025-12-07
### Features

- **amdp:** Enhance breakpoint functionality and testing([`76ca83b`](https://github.com/oisee/vibing-steampunk/commit/76ca83b539c1824f86b22f64abb29c6d5d78406e))
- V2.12.0 - abapGit-compatible format & batch operations([`c731e2e`](https://github.com/oisee/vibing-steampunk/commit/c731e2e8a13670bc0cc318a328d8b618978c8f0f))



## [2.11.0] - 2025-12-05

## [2.10.1] - 2025-12-05

## [2.10.0] - 2025-12-05

## [2.9.0] - 2025-12-05

## [2.8.0] - 2025-12-05

## [2.7.0] - 2025-12-05

## [2.6.0] - 2025-12-05

## [2.5.0] - 2025-12-05

## [2.4.0] - 2025-12-04

## [2.3.0] - 2025-12-04

## [2.2.0] - 2025-12-04

## [2.1.0] - 2025-12-04

## [2.0.0-before-rename] - 2025-12-04

## [2.0.0] - 2025-12-04

## [1.6.0] - 2025-12-04

## [1.5.0] - 2025-12-03
### Features

- Enhance tool descriptions with usage examples and workflows([`c52bd4f`](https://github.com/oisee/vibing-steampunk/commit/c52bd4fe2d4d0027281a8e89d3afbdf7555d272a))


### Milestone

- Add search/grep tools and enhance EditSource([`eeea7d9`](https://github.com/oisee/vibing-steampunk/commit/eeea7d973b0b31a2c9b165ce35a96174e2e123a5))



## [1.4.1] - 2025-12-03
### Bug Fixes

- Add missing SaveToFile and RenameObject MCP tool registrations([`67a5f1a`](https://github.com/oisee/vibing-steampunk/commit/67a5f1a061a0863cfff132f158039f93ac05cd4d))



## [1.4.0] - 2025-12-02
### Features

- Add file-based deployment tools solving token limit problem([`dc6b541`](https://github.com/oisee/vibing-steampunk/commit/dc6b541ae7e133169bb6fa741c38a0f63c787d43))


### Refactoring

- Simplify MCP tool interface (41→39 tools)([`688772b`](https://github.com/oisee/vibing-steampunk/commit/688772b0f26b92a8be420db588306ebe7264c779))



## [1.3.0] - 2025-12-02
### Features

- Add comprehensive research report on ABAP debugging and tracing capabilities([`0a1bb1e`](https://github.com/oisee/vibing-steampunk/commit/0a1bb1ef3d633e11598dce065a80f69fb662a4e6))
- Add roadmap section with ongoing and planned features for debugging and analysis tools([`b6c08db`](https://github.com/oisee/vibing-steampunk/commit/b6c08db98cdccbba75b4c3bbc4252224c514ab24))



## [1.2.0] - 2025-12-02

## [1.1.0] - 2025-12-02
### Features

- **adt:** Implement workflows for writing and creating ABAP programs and classes([`cdf3f98`](https://github.com/oisee/vibing-steampunk/commit/cdf3f98d401f2d571b93742c9e3755cd6027d9a7))



## [1.0.0] - 2025-12-02


