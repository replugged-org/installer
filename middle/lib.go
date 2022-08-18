package middle

import (
	"os"
	"path/filepath"
)

func ChownR(path string, uid int, gid int) error {
  return filepath.Walk(path, func(name string, info os.FileInfo, err error) error {
    if err == nil {
      err = os.Chown(name, uid, gid)
    }
    return err
  })
}
