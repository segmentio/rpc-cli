
v1.0.1 / 2018-09-24
===================

  * fix(build): add ldflags to set version in binary

v1.0.0 / 2018-09-24
===================

  * fix(cmd): add '-legacy' to avoid conflicting with a newer cli
  * feat: add User-Agent header (#12)
  * dev: suffix -legacy for package and republish
  * Use govendor; add Makefile to produce pkg
  * feat: use json number in decoder to support large numbers better (#10)
  * print application-level errors when using HTTP (#9)
  * ocd
  * rpc: detect stdin, fixes #4
  * internal/rpc: handle empty inputs
  * lazy license
  * ocd
  * Initial commit
