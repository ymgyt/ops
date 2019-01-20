package filesystem

func New() *FileSystem {
	return &FileSystem{}
}

type FileSystem struct {
}

func (fs *FileSystem) DiskUsageReporter() *DiskUsageReporter {
	return &DiskUsageReporter{}
}
