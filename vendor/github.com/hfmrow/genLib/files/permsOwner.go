package files

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
)

/*
Usage to Create any directories needed to put this file in them:

func DirTree(filename string){
    var dir_file_mode os.FileMode
	OS := OsPermsStructNew()

    dir_file_mode = os.ModeDir | (OS.USER_RWX | OS.ALL_R)
    os.MkdirAll(filename, dir_file_mode)
}

func TestPerms(){
	OS := OsPermsStructNew()

	permsFile := os.ModePerm & (OS.USER_RW | OS.GROUP_R | OS.OTH_R)     // 0644
	permsFileX := os.ModePerm & (OS.USER_RWX | OS.GROUP_RX | OS.OTH_RX) // 0755
	permsDir := os.ModePerm & (OS.USER_RWX | OS.GROUP_RX | OS.OTH_RX)   // 0755
	fmt.Printf("%#o, %#o, %#o, \n", permsFile, permsFileX, permsDir)
}

Memo:
   0 : - - - (aucun droit)
   1 : - - x (exécution)
   2 : - w - (écriture)
   3 : - w x (écriture et exécution)
   4 : r - - (lecture seule)
   5 : r - x (lecture et exécution)
   6 : r w - (lecture et écriture)
   7 : r w x (lecture, écriture et exécution)
*/

// Structure that hold file permissions
type osPerms struct {
	USER_R,
	USER_W,
	USER_X,
	USER_RW,
	USER_RX,
	USER_RWX,

	GROUP_R,
	GROUP_W,
	GROUP_X,
	GROUP_RW,
	GROUP_RX,
	GROUP_RWX,

	OTH_R,
	OTH_W,
	OTH_X,
	OTH_RW,
	OTH_RX,
	OTH_RWX,

	ALL_R,
	ALL_W,
	ALL_X,
	ALL_RW,
	ALL_RX,
	ALL_RWX,

	rEAD,
	wRITE,
	eX,
	uSER_SHIFT,
	gROUP_SHIFT,
	oTH_SHIFT os.FileMode
}

func OsPermsStructNew() (osp *osPerms) {
	osp = new(osPerms)

	osp.rEAD = 04
	osp.wRITE = 02
	osp.eX = 01
	osp.uSER_SHIFT = 6
	osp.gROUP_SHIFT = 3
	osp.oTH_SHIFT = 0

	osp.USER_R = osp.rEAD << osp.uSER_SHIFT
	osp.USER_W = osp.wRITE << osp.uSER_SHIFT
	osp.USER_X = osp.eX << osp.uSER_SHIFT
	osp.USER_RW = osp.USER_R | osp.USER_W
	osp.USER_RX = osp.USER_R | osp.USER_X
	osp.USER_RWX = osp.USER_RW | osp.USER_X

	osp.GROUP_R = osp.rEAD << osp.gROUP_SHIFT
	osp.GROUP_W = osp.wRITE << osp.gROUP_SHIFT
	osp.GROUP_X = osp.eX << osp.gROUP_SHIFT
	osp.GROUP_RW = osp.GROUP_R | osp.GROUP_W
	osp.GROUP_RX = osp.GROUP_R | osp.GROUP_X
	osp.GROUP_RWX = osp.GROUP_RW | osp.GROUP_X

	osp.OTH_R = osp.rEAD << osp.oTH_SHIFT
	osp.OTH_W = osp.wRITE << osp.oTH_SHIFT
	osp.OTH_X = osp.eX << osp.oTH_SHIFT
	osp.OTH_RW = osp.OTH_R | osp.OTH_W
	osp.OTH_RX = osp.OTH_R | osp.OTH_X
	osp.OTH_RWX = osp.OTH_RW | osp.OTH_X

	osp.ALL_R = osp.USER_R | osp.GROUP_R | osp.OTH_R
	osp.ALL_W = osp.USER_W | osp.GROUP_W | osp.OTH_W
	osp.ALL_X = osp.USER_X | osp.GROUP_X | osp.OTH_X
	osp.ALL_RW = osp.ALL_R | osp.ALL_W
	osp.ALL_RX = osp.ALL_R | osp.ALL_X
	osp.ALL_RWX = osp.ALL_RW | osp.ALL_X
	return
}

// DispPerms: display right.
// i.e: DispPerms(OS.GROUP_RW | OS.USER_RW | OS.ALL_R) -> 664
func DispPerms(value int) {
	fmt.Printf("os.ModePerm & %#o\n", value)
}

// ChangeFileOwner: set the owner of the file to real user instead
// of root. Obtaining root rights differs when sudo or pkexec is
// used, this function checks which command was used and acts to
// perform the job correctly.
func ChangeFileOwner(filename string) (err error) {
	var uid, gid int

	if _, realUser, _, err := GetRootCurrRealUser(); err == nil {
		if uid, err = strconv.Atoi(realUser.Uid); err == nil {
			if gid, err = strconv.Atoi(realUser.Gid); err == nil {
				err = os.Chown(filename, uid, gid)
			}
		}
	}
	return err
}

// GetRootCurrRealUser: Retrieve informations about root state
// and current and real user. Handle different behavior between
// sudo and pkexec commands.
func GetRootCurrRealUser() (currentUser, realUser *user.User, root bool, err error) {
	realUser = new(user.User)
	if currentUser, err = user.Current(); err == nil {
		*(realUser) = *(currentUser)                                  // Copy data content
		if root = (currentUser.Uid + currentUser.Gid) == "00"; root { // root rights acquired ?
			// fmt.Println("pkexec used")
			if realUser, err = user.LookupId(os.Getenv("PKEXEC_UID")); err != nil {
				// fmt.Println("sudo used")
				if realUser, err = user.Lookup(os.Getenv("SUDO_USER")); err != nil {
					return
				}
			}
		}
	}
	return
}

// const (
// 	OS_READ        = 04
// 	OS_WRITE       = 02
// 	OS_EX          = 01
// 	OS_USER_SHIFT  = 6
// 	OS_GROUP_SHIFT = 3
// 	OS_OTH_SHIFT   = 0

// 	OS_USER_R   = OS_READ << OS_USER_SHIFT
// 	OS_USER_W   = OS_WRITE << OS_USER_SHIFT
// 	OS_USER_X   = OS_EX << OS_USER_SHIFT
// 	OS_USER_RW  = OS_USER_R | OS_USER_W
// 	OS_USER_RX  = OS_USER_R | OS_USER_X
// 	OS_USER_RWX = OS_USER_RW | OS_USER_X

// 	OS_GROUP_R   = OS_READ << OS_GROUP_SHIFT
// 	OS_GROUP_W   = OS_WRITE << OS_GROUP_SHIFT
// 	OS_GROUP_X   = OS_EX << OS_GROUP_SHIFT
// 	OS_GROUP_RW  = OS_GROUP_R | OS_GROUP_W
// 	OS_GROUP_RX  = OS_GROUP_R | OS_GROUP_X
// 	OS_GROUP_RWX = OS_GROUP_RW | OS_GROUP_X

// 	OS_OTH_R   = OS_READ << OS_OTH_SHIFT
// 	OS_OTH_W   = OS_WRITE << OS_OTH_SHIFT
// 	OS_OTH_X   = OS_EX << OS_OTH_SHIFT
// 	OS_OTH_RW  = OS_OTH_R | OS_OTH_W
// 	OS_OTH_RX  = OS_OTH_R | OS_OTH_X
// 	OS_OTH_RWX = OS_OTH_RW | OS_OTH_X

// 	OS_ALL_R   = OS_USER_R | OS_GROUP_R | OS_OTH_R
// 	OS_ALL_W   = OS_USER_W | OS_GROUP_W | OS_OTH_W
// 	OS_ALL_X   = OS_USER_X | OS_GROUP_X | OS_OTH_X
// 	OS_ALL_RW  = OS_ALL_R | OS_ALL_W
// 	OS_ALL_RX  = OS_ALL_R | OS_ALL_X
// 	OS_ALL_RWX = OS_ALL_RW | OS_ALL_X
// )
