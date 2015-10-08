package model

type BackupFile struct {
	Name       string
	FullPath   string
	Size       int64
	CreateTime int64
}

type BackupFiles []*BackupFile

func (b BackupFiles) Len() int {
	return len(b)
}

func (b BackupFiles) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b BackupFiles) Less(i, j int) bool {
	return b[i].CreateTime > b[j].CreateTime
}
