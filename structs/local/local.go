package local
import (
	"os"
	"path/filepath"

	"github.com/kazoup/platform/structs"
	"github.com/kazoup/platform/structs/categories"
)

// LocalFile model
type LocalFile struct {
	Path string
	Info os.FileInfo
}


func NewKazoupFileFromLocal(lf *LocalFile ) *structs.KazoupFile {

	original:= struct{
		Mode os.FileMode `json:"mode"`
		}{
			Mode:lf.Info.Mode(),
		}
	return &structs.KazoupFile{
		ID:       structs.GetMD5Hash(lf.Path),
		Name:     lf.Info.Name(),
		URL:      "/local" + lf.Path,
		Modified: lf.Info.ModTime(),
		Size:     lf.Info.Size(),
		IsDir:    lf.Info.IsDir(),
		Category: categories.GetDocType(filepath.Ext(lf.Info.Name())),
		Depth:    structs.UrlDepth(lf.Path),
		Original : original,
	}
}

