package vcs

/*
This package implements floo's sync algorithm "powder" (Floo Powder)

The sync algorithm tries to handle the following special cases:
 - Propagate moves (most of them, at least)
 - Propagate deletes (configurable)
 - Also sync empty directories.

Terminology:
 - Destination (short "dst") is used to reference our own storage.
 - Source (short: "src") is used to reference the remote storage.

The sync algorithm can be roughly divided in 4 stages:
 - Stage 1: "Move Marking":
   Iterate over all ghosts in the tree and check if they were either moved
   (has sibling) or removed (has no sibling). In case of directories, the
   second mapping stage is already executed.

- Stage 2: "Mapping":
   Finding pairs of files that possibly adding, merging or conflict handling.
   Equal files will already be sorted out at this point. Every already
   visited node in the remote linker will be marked. The mapping algorithm
   starts at the root node and uses the attributes of the merkle trees
   (same hash = same content) to skip over same parts.

- Stage 3: "Resolving":
   For each file a decision needs to be made. This decision defines the next step
   and can be one of the following.

   - The file was added on the remote, we should add it too -> Add them.
   - The file was removed on the remote, we might want to also delete it.
   - The file was only moved on the remote node, we might want to move it also.
   - The file has compatible changes on the both sides. -> Merge them.
   - The file has incompatible changes on both sides -> Do conflict resolution.
   - The nodes have differing types (directory vs files). Report them.

- Stage 4: "Handling"
   Only at this stage "sync" and "diff" differ.
   Sync will take the files from Stage 3 and add/remove/merge files.
   Diff will create a report out of those files and also includes files that
   are simply missing on the source side (but do not need to be removed).

 Everything except Stage 4 is read-only. If a user wants to only show the diff
 between two linkers, he just prints what would be done instead of actually doing it.
 This makes the diff and sync implementation share most of the code.
*/
