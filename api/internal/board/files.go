package board

import (
	"path"
	"strings"
)

const MaxAssetSize = 5 << 20

var assetTypes = map[string]string{
	".png":  "image/png",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".gif":  "image/gif",
	".webp": "image/webp",
	".svg":  "image/svg+xml",
}

func WantedFile(name string, size int64) bool {
	ext := strings.ToLower(path.Ext(name))
	if ext == ".md" {
		return true
	}
	_, isAsset := assetTypes[ext]
	return isAsset && size <= MaxAssetSize
}

func AssetType(name string) (string, bool) {
	t, ok := assetTypes[strings.ToLower(path.Ext(name))]
	return t, ok
}
