package fileservice

var ValidExtensions = map[string]struct{}{
	"pdf":  {},
	"docx": {},
	"doc":  {},

	"jpg":  {},
	"jpeg": {},
	"png":  {},

	"mp4": {},

	"xlsx": {},
	"xls":  {},
	"csv":  {},
}

var ValidCategories = map[string]struct{}{
	"word":  {},
	"pdf":   {},
	"image": {},
	"video": {},
	"excel": {},
	"audio": {},
	"zip":   {},
	"other": {},
}

var CategoryExtensions = map[string]string{
	"pdf":  "pdf",
	"docx": "word",
	"doc":  "word",

	"jpg":  "image",
	"jpeg": "image",
	"png":  "image",

	"mp4": "video",

	"xlsx": "excel",
	"xls":  "excel",
	"csv":  "excel",
}
