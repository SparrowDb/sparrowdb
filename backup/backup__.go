package backup

/*
var (
	searchDir = "../data/teste3"
	targetDir = "../data/snapshot"
)

func extractPath(path string) string {
	return strings.Replace(path, searchDir, "", -1)
}

func getFiles() []string {
	var fileList []string

	filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, extractPath(path))
		return nil
	})

	return fileList
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	return fileInfo.IsDir(), err
}

func exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err != nil, err
}

func createHardLink(src, dst string) {
	err := os.Link(src, dst)
	if err != nil {
		log.Fatalln(err)
	}
}

func rename(src, dst string) {
	err := os.Rename(src, dst)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func checkFolderAndCreateHardLink(dirs []string) {
	targetFolder := filepath.Join(targetDir)
	if ok, _ := util.Exists(targetFolder); ok == false {
		util.CreateDir(targetFolder)
	}

	for _, file := range dirs {
		sourcePath := filepath.Join(searchDir, file)
		targetPath := filepath.Join(targetDir, extractPath(sourcePath))

		if ok, _ := isDirectory(sourcePath); ok == true {
			if fok, _ := exists(targetDir); fok == false {
				util.CreateDir(targetPath)
			}
		} else {
			if fok, _ := exists(targetDir); fok == false {
				createHardLink(sourcePath, targetPath)
			}
		}
	}

}

func do() {
	fmt.Println("dsadasdasadA")

	oldTargetDir := targetDir
	targetDir = targetDir + "_tmp"

	var fileList []string
	fileList = append(fileList, getFiles()...)
	checkFolderAndCreateHardLink(fileList)

	rename(targetDir, oldTargetDir)
}
*/
