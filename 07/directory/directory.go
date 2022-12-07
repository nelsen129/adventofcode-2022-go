package directory

type Directory struct {
	name              string
	child_directories []*Directory
	parent_directory  *Directory
	size              int
	total_size        int
}

func (D *Directory) GetName() string {
	return D.name
}

func (D *Directory) AddSize(s int) {
	D.size += s
}

func (D *Directory) AddSubdirectory(name string) {
	cd := Directory{}
	cd.name = name
	cd.parent_directory = D
	D.child_directories = append(D.child_directories, &cd)
}

func (D *Directory) GetSubdirectoryFromName(name string) *Directory {
	for i := range D.child_directories {
		if D.child_directories[i].name == name {
			return D.child_directories[i]
		}
	}

	return D
}

func (D *Directory) GetParentDirectory() *Directory {
	return D.parent_directory
}

func (D *Directory) GetTotalSize() int {
	if D.total_size != 0 {
		return D.total_size
	}

	for i := range D.child_directories {
		D.total_size += D.child_directories[i].GetTotalSize()
	}

	D.total_size += D.size

	return D.total_size
}

func (D *Directory) GetTotalSizeSubdirectoriesLessThanEqualTo(size int) int {
	total_size := 0

	if D.GetTotalSize() <= size {
		total_size += D.total_size*2 - D.size
	} else {
		for i := range D.child_directories {
			total_size += D.child_directories[i].GetTotalSizeSubdirectoriesLessThanEqualTo(size)
		}
	}

	return total_size
}
