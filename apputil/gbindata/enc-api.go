package gbindata

// Reference: https://github.com/a-urth/go-bindata

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"github.com/shawnwyckoff/gpkg/container/gstring"
	"github.com/shawnwyckoff/gpkg/sys/gfs"
	"os"
)

var templateBegin = `
// DO NOT EDIT BY HAND
//
// Generated by bindata

package %s

var %s = "0x`

var templateEnd = `"`

// packageName: 生成的包名
// varName: 生成的包含文件数据的变量名
func Enc(binary_filename, encoded_go_filename, packageName, varName string) error {
	if !gfs.FileExits(binary_filename) {
		return errors.Errorf("file %s not exists", binary_filename)
	}
	if !gstring.EndWith(encoded_go_filename, ".go") {
		return errors.Errorf("filename %s doesn't end with .go", encoded_go_filename)
	}

	os.Remove(encoded_go_filename)

	src_f, err := os.Open(binary_filename)
	if err != nil {
		return err
	}
	defer src_f.Close()

	dst_f, err := os.Create(encoded_go_filename)
	if err != nil {
		return err
	}
	defer dst_f.Close()

	dst_f.WriteString(fmt.Sprintf(templateBegin, packageName, varName))
	scan := bufio.NewScanner(src_f)
	scan.Split(bufio.ScanBytes)
	for scan.Scan() {
		dst_f.WriteString(hex.EncodeToString(scan.Bytes()))
		if scan.Err() != nil {
			return err
		}
	}
	dst_f.WriteString(templateEnd)

	return nil
}
