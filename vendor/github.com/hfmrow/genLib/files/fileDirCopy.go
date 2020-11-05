// fileDirCopy.go

/*
	Copyright ©2020 H.F.M - Files/Dirs copy library v1.0 part of H.F.M genLib
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php

	-Copy files and / or directory contents and store source and destination filenames
	 so you can use them later.
	-Preserve file/dir permissions ans allow choosing target owner.
	-Extension/Files mask may include or exclude files.

	[notice:] mkDirAll methode is a clone of original golang command os.MkDirAll(), with a
	simple modification that permit to record each created directory and assign wanted
	permissions on folder recently created.

*/

package files

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"syscall"
)

// FileDirCopyStruct: Used to copy folder and content to another place.
// Preserve file/dir permissions and allow choosing target owner.
type FileDirCopyStruct struct {
	Src,
	Dest string

	FilesListStrings,
	Masks []string
	// true: keep only the files that match the extension masks
	// false: keep only files that do NOT match the extension masks
	IncludeMasks,
	RemoveBeforeCopy,
	SortResult bool

	currFilesList,
	FilesList,
	FilesListNotCopied []fileInfo

	// Fill these variables to chown destination files/dir.
	OwnerUid,
	OwnerGid int
}

// fileInfo: structure to hold single file
type fileInfo struct {
	fInfo os.FileInfo
	FNameSrc,
	FNameDst string
	Err error
}

// FileDirCopyStructNew:
func FileDirCopyStructNew() (fsi *FileDirCopyStruct) {
	fsi = new(FileDirCopyStruct)
	fsi.OwnerUid, fsi.OwnerGid = -1, -1
	return
}

// GetFileInfos: Retrieve os.FileInfo for each file starting at "fsi.Src"
func (fsi *FileDirCopyStruct) GetFilenames() (err error) {
	var ok bool
	var stat os.FileInfo
	var tmpDest string

	// Clear files entries
	fsi.currFilesList = fsi.currFilesList[:0]

	var checkExt = func(path string, info os.FileInfo) (err error) {
		if info != nil {
			// Combine destination filename
			tmpDest, err = filepath.Rel(fsi.Src, path)
			tmpDest = filepath.Join(fsi.Dest, tmpDest)
			// Check for .ext exclusion.
			if len(fsi.Masks) != 0 {
				for _, mask := range fsi.Masks {
					if ok, err = filepath.Match(mask, filepath.Base(path)); err != nil {
						return err
					}
					if ok {
						if fsi.IncludeMasks {
							fsi.currFilesList = append(fsi.currFilesList,
								fileInfo{fInfo: info, FNameSrc: path, FNameDst: tmpDest})
						}
						break
					}
				}
				if !ok && !fsi.IncludeMasks {
					fsi.currFilesList = append(fsi.currFilesList,
						fileInfo{fInfo: info, FNameSrc: path, FNameDst: tmpDest})
				}
			} else {
				fsi.currFilesList = append(fsi.currFilesList,
					fileInfo{fInfo: info, FNameSrc: path, FNameDst: tmpDest})
			}
		}
		return err
	}

	stat, err = os.Lstat(fsi.Src)
	if err == nil {
		switch {
		case stat.IsDir():
			if err = filepath.Walk(fsi.Src, func(path string, info os.FileInfo, err error) error {
				return checkExt(path, info)
			}); err != nil {
				return
			}
		default:
			if err = checkExt(fsi.Src, stat); err != nil {
				return
			}
		}
	}
	// Store filenames
	for _, fi := range fsi.currFilesList {
		switch {
		case fi.fInfo.IsDir():

		default:
			fsi.FilesList = append(fsi.FilesList, fi)
			fsi.FilesListStrings = append(fsi.FilesListStrings, fi.FNameSrc)
		}
	}
	return
}

// CopyFiles: Copy file and create directory if not exists;
// Given parameters override internal structure (Src, Dest) variables.
func (fsi *FileDirCopyStruct) CopyFiles(scrAndDestFilename ...string) (err error) {
	var ok bool
	var stat os.FileInfo
	var tmpDest string

	if len(scrAndDestFilename) == 2 {
		fsi.Src = scrAndDestFilename[0]
		fsi.Dest = scrAndDestFilename[1]
	}

	// Clear files entries
	fsi.currFilesList = fsi.currFilesList[:0]

	var checkExt = func(path string, info os.FileInfo) (err error) {
		if info != nil {
			// Combine destination filename
			tmpDest, err = filepath.Rel(fsi.Src, path)
			tmpDest = filepath.Join(fsi.Dest, tmpDest)
			// Check for .ext exclusion.
			if len(fsi.Masks) != 0 {
				for _, mask := range fsi.Masks {
					if ok, err = filepath.Match(mask, filepath.Base(path)); err != nil {
						return err
					}
					if ok {
						if fsi.IncludeMasks {
							fsi.currFilesList = append(fsi.currFilesList,
								fileInfo{fInfo: info, FNameSrc: path, FNameDst: tmpDest})
						}
						break
					}
				}
				if !ok && !fsi.IncludeMasks {
					fsi.currFilesList = append(fsi.currFilesList,
						fileInfo{fInfo: info, FNameSrc: path, FNameDst: tmpDest})
				}
			} else {
				fsi.currFilesList = append(fsi.currFilesList,
					fileInfo{fInfo: info, FNameSrc: path, FNameDst: tmpDest})
			}
		}
		return err
	}

	stat, err = os.Lstat(fsi.Src)
	if err == nil {
		switch {
		case stat.IsDir():
			if err = filepath.Walk(fsi.Src, func(path string, info os.FileInfo, err error) error {
				return checkExt(path, info)
			}); err != nil {
				return
			}
		default:
			if err = checkExt(fsi.Src, stat); err != nil {
				return
			}
		}
	}
	return fsi.copy()
}

var OS = OsPermsStructNew()

// copy: with permission préservation and specific owner if requested
func (fsi *FileDirCopyStruct) copy() (err error) {
	var permsDir = os.ModePerm & (OS.USER_W | OS.ALL_RX) // 0755
	var inBytes []byte

	for _, fi := range fsi.currFilesList {
		switch {
		case fi.fInfo.IsDir():
			// Build directory tree
			err = fsi.MkdirAll(fi.FNameDst, fi.fInfo.Mode())
		default:
			// Check if the destination file's directory. exists and create if not.
			tmpDestDir := filepath.Dir(fi.FNameDst)
			if _, err = os.Stat(tmpDestDir); os.IsNotExist(err) {
				err = fsi.MkdirAll(tmpDestDir, permsDir)
			}
			// Read & Copy file
			if inBytes, err = ioutil.ReadFile(fi.FNameSrc); err == nil {
				if _, err = os.Stat(fi.FNameDst); !os.IsNotExist(err) && fsi.RemoveBeforeCopy {
					if err = os.RemoveAll(fi.FNameDst); err != nil {
						fmt.Printf("No way to delete existing file before creating it :%s\n", err.Error())
					}
				}
				if err = ioutil.WriteFile(fi.FNameDst, inBytes, fi.fInfo.Mode()); err != nil {
					fi.Err = err
					fsi.FilesListNotCopied = append(fsi.FilesListNotCopied, fi)
				} else {
					fsi.FilesList = append(fsi.FilesList, fi)
					fsi.FilesListStrings = append(fsi.FilesListStrings, fi.FNameDst)
				}
			}
		}
	}

	// Change files/dirs owner if required
	if fsi.OwnerUid+fsi.OwnerGid > -1 {
		for _, fi := range fsi.currFilesList {
			if err = os.Chown(fi.FNameDst, fsi.OwnerUid, fsi.OwnerGid); err != nil && !os.IsNotExist(err) {
				return
			}
		}
	}

	// Sort string preserving order ascendant
	if fsi.SortResult {
		sort.SliceStable(fsi.FilesList, func(i, j int) bool {
			return fsi.FilesList[i].FNameSrc < fsi.FilesList[j].FNameSrc
		})
	}
	return
}

// mkdirAll: creates a directory named path, along with any
// necessary parents, and returns nil, or else returns an error.
// The permission bits perm (before umask) are used for all
// directories that MkdirAll creates. If path is already a
// directory, MkdirAll does nothing and returns nil.
// This is a modified version of original golang function
// this one store created directories to be handled later.
func (fsi *FileDirCopyStruct) MkdirAll(path string, perm os.FileMode) error {
	// Fast path: if we can tell whether path is a directory or file,
	// stop with success or error.
	dir, err := os.Stat(path)
	if err == nil {
		if dir.IsDir() {
			return nil
		}
		return &os.PathError{"mkdir", path, syscall.ENOTDIR}
	}

	// Slow path: make sure parent exists and then call Mkdir for path.
	i := len(path)
	for i > 0 && os.IsPathSeparator(path[i-1]) { // Skip trailing path separator.
		i--
	}

	j := i
	for j > 0 && !os.IsPathSeparator(path[j-1]) { // Scan backward over element.
		j--
	}

	if j > 1 {
		// Create parent.
		err = fsi.MkdirAll(path[:j-1], perm)
		if err != nil {
			return err
		}
	}

	// Parent now exists; invoke Mkdir and use its result.
	err = os.Mkdir(path, perm)

	tmpFI := fileInfo{
		fInfo:    dir,
		FNameDst: path,
		Err:      err,
	}
	if err != nil {
		// Handle arguments like "foo/." by
		// double-checking that directory doesn't exist.
		dir, err1 := os.Lstat(path)
		if err1 == nil && dir.IsDir() {
			return nil
		}
		// Add to Error history
		fsi.FilesListNotCopied = append(fsi.FilesListNotCopied, tmpFI)
		return err
	}
	// Add to Copied history
	fsi.FilesList = append(fsi.FilesList, tmpFI)
	fsi.FilesListStrings = append(fsi.FilesListStrings, tmpFI.FNameDst)
	fsi.currFilesList = append(fsi.currFilesList, tmpFI) // To apply chown.
	return nil
}

// /**************************************************************

// 		golang copy directory recursively (MIT) license.

// 	Hiromu OCHIAI
// 	otiai10
// 	https://github.com/otiai10

// 		Retrieved from: https://github.com/otiai10/copy

// **************************************************************/

const (
	// tmpPermissionForDirectory makes the destination directory writable,
	// so that stuff can be copied recursively even if any original directory is NOT writable.
	// See https://github.com/otiai10/copy/pull/9 for more information.
	tmpPermissionForDirectory = os.FileMode(0755)
)

// Copy copies src to dest, doesn't matter if src is a directory or a file
func FileDirCopy(src, dest string) error {
	info, err := os.Lstat(src)
	if err != nil {
		return err
	}
	return copy(src, dest, info)
}

// copy dispatches copy-funcs according to the mode.
// Because this "copy" could be called recursively,
// "info" MUST be given here, NOT nil.
func copy(src, dest string, info os.FileInfo) error {
	if info.Mode()&os.ModeSymlink != 0 {
		return lcopy(src, dest, info)
	}
	if info.IsDir() {
		return dcopy(src, dest, info)
	}
	return fcopy(src, dest, info)
}

// fcopy is for just a file,
// with considering existence of parent directory
// and file permission.
func fcopy(src, dest string, info os.FileInfo) error {

	if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	if err = os.Chmod(f.Name(), info.Mode()); err != nil {
		return err
	}

	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	_, err = io.Copy(f, s)
	return err
}

// dcopy is for a directory,
// with scanning contents inside the directory
// and pass everything to "copy" recursively.
func dcopy(srcdir, destdir string, info os.FileInfo) error {

	originalMode := info.Mode()

	// Make dest dir with 0755 so that everything writable.
	if err := os.MkdirAll(destdir, tmpPermissionForDirectory); err != nil {
		return err
	}
	// Recover dir mode with original one.
	defer os.Chmod(destdir, originalMode)

	contents, err := ioutil.ReadDir(srcdir)
	if err != nil {
		return err
	}

	for _, content := range contents {
		cs, cd := filepath.Join(srcdir, content.Name()), filepath.Join(destdir, content.Name())
		if err := copy(cs, cd, content); err != nil {
			// If any error, exit immediately
			return err
		}
	}

	return nil
}

// lcopy is for a symlink,
// with just creating a new symlink by replicating src symlink.
func lcopy(src, dest string, info os.FileInfo) error {
	src, err := os.Readlink(src)
	if err != nil {
		return err
	}
	return os.Symlink(src, dest)
}
