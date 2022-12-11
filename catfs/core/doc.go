// Layout of the key/value store:
//
// objects/<NODE_HASH>                   => NODE_METADATA
// tree/<FULL_NODE_PATH>                 => NODE_HASH
// index/<CMT_INDEX>                     => COMMIT_HASH
// inode/<INODE>                         => NODE_HASH
// moves/<INODE>                         => MOVE_INFO
// moves/overlay/<INODE>                 => MOVE_INFO
//
// stage/objects/<NODE_HASH>             => NODE_METADATA
// stage/tree/<FULL_NODE_PATH>           => NODE_HASH
// stage/STATUS                          => COMMIT_METADATA
// stage/moves/<INODE>                   => MOVE_INFO
// stage/moves/overlay/<INODE>           => MOVE_INFO
//
// stats/max-inode                       => UINT64
// refs/<REFNAME>                        => NODE_HASH
//
// Defined by caller:
//
// metadata/                             => BYTES (Caller defined data)
//
// NODE is either a Commit, a Directory or a File.
// FULL_NODE_PATH may contain slashes and in case of directories,
// it will contain a trailing slash.
//
// The following refs are defined by the system:
// HEAD -> Points to the latest finished commit, or nil.
// CURR -> Points to the staging commit.

package core
