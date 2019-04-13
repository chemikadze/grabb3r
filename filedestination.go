package grabb3r

import (
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"syscall"
)

type fileDestination struct {
	rootDir string
}

func NewFileDestination(rootDir string) SolutionDestination {
	return &fileDestination{rootDir}
}

func (dst *fileDestination) Initialize() error {
	info, err := os.Stat(dst.rootDir)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(dst.rootDir, os.ModeDir|0755)
		return err
	} else if err != nil {
		return err
	}
	if !info.IsDir() {
		return &os.PathError{"stat", dst.rootDir, syscall.ENOTDIR}
	}
	return nil
}

func normalizeString(str string) string {
	re := regexp.MustCompile("[.,$%&`\"'! ]")
	return re.ReplaceAllString(str, "_")
}

func (dst *fileDestination) solutionPath(solution SolutionDesc) string {
	fileName := normalizeString(solution.ProblemName())
	fileExt := extensionForLanguage(solution.Language())
	return path.Join(dst.rootDir, fmt.Sprintf("%s.%s", fileName, fileExt))
}

func extensionForLanguage(lang Language) string {
	switch lang {
	case Language("python"):
		return "py"
	case Language("go"):
		return "go"
	default:
		return fmt.Sprintf("%s.txt", normalizeString(string(lang)))
	}
}

func (dst *fileDestination) SaveSolution(solution Solution) error {
	fileName := dst.solutionPath(solution.Desc())
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	_, err = file.WriteString(solution.Code())
	if err != nil {
		_ = file.Close()
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}
	err = os.Chtimes(fileName, solution.Desc().SubmittedTime(), solution.Desc().SubmittedTime())
	if err != nil {
		return err
	}
	log.Printf("Saved solution file: %s", fileName)
	return err
}

func (dst *fileDestination) ContainsSolution(solution SolutionDesc) (bool, error) {
	info, err := os.Stat(dst.solutionPath(solution))
	if err != nil && os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	} else if info.IsDir() {
		return false, &os.PathError{"stat", dst.rootDir, syscall.EISDIR}
	} else {
		return true, nil
	}
}
