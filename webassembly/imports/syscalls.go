package imports

import (
	"log"
)

func sys_ftruncate64(_ int32, _ int64) int32 {
	log.Fatal("Called into __sys_ftruncate64")
	return 0
}

func sys_unlinkat(_ int32, _ int32, _ int32) int32 {
	log.Fatal("Called into __sys_unlinkat")
	return 0
}

func sys_rmdir(_ int32) int32 {
	log.Fatal("Called into __sys_rmdir")
	return 0
}
