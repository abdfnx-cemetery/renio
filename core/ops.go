package core

const (
  // FileSystem ops
  FSWrite     = 1
  FSRead      = 2
  FSExists    = 3
  FSDirExists = 4
  FSCwd       = 5
  FSStat      = 6
  FSRemove    = 7
  FSMkdir     = 9
  FSWalk      = 14
  // console ops
  Log          = 10
  // env ops
  Env          = 11
  // plugin ops
  Plugin       = 15
  // fetch ops
  Fetch        = 20
  // serve ops
  Serve        = 25
)
