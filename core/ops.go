package core

const (
  // FileSystem ops
  FS_Write     = 1
  FS_Read      = 2
  FS_Exists    = 3
  FS_DirExists = 4
  FS_Cwd       = 5
  FS_Stat      = 6
  FS_Remove    = 7
  FS_Mkdir     = 9
  FS_Walk      = 14
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
