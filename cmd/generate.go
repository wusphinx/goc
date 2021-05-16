package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/sirupsen/logrus"
)

func init() {
	formatter := runtime.Formatter{ChildFormatter: &logrus.TextFormatter{
		FullTimestamp: true,
	}}
	formatter.Line = true
	logrus.SetFormatter(&formatter)
	logrus.SetLevel(logrus.WarnLevel)
}

var dirList = []string{"docs", "configs", "logs", "hooks"}

func Generator() error {
	dir, err := os.Getwd()
	if err != nil {
		logrus.Fatal(err)
	}

	if strings.HasPrefix(appName, ".") {
		parentDir = dir
		ss := strings.Split(dir, "/")
		appName = ss[len(ss)-1]
	} else {
		parentDir = filepath.Join(dir, appName)
		if err := os.MkdirAll(parentDir, os.ModePerm); err != nil {
			logrus.Fatal(err)
		}
	}

	if err := os.Chdir(parentDir); err != nil {
		logrus.Fatal(err)
	}

	genHandler()

	return nil
}

func genHandler() error {
	genDirs(dirList)
	genGitignore()
	genPreCommitHooks()
	getGolangci()
	genMod()
	genMain()

	return nil
}

func genPreCommitHooks() error {
	fp := filepath.Join(parentDir, "hooks", "pre-commit")
	return genFile(fp, preCommitHookTmpl, nil)

}

func genMod() error {
	param := map[string]string{
		"ModuleName": moduleName,
		"GoVersion":  goVersion,
	}
	return genFile(filepath.Join(parentDir, "go.mod"), moduleTmpl, param)
}

func genFile(filePath, tmpl string, param map[string]string) error {
	fp, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer fp.Close()

	data, err := RenderTemple(param, tmpl)
	if err != nil {
		logrus.Fatal(err)
	}

	if _, err := fp.Write(data); err != nil {
		logrus.Fatal(err)
	}

	return nil
}

func genGitignore() error {
	return genFile(filepath.Join(parentDir, ".gitignore"), gitignoreTmpl, nil)
}

func genDirs(ds []string) (err error) {
	defer func() {
		if err != nil {
			logrus.Fatal(err)
		}
	}()

	for _, d := range ds {
		dir := filepath.Join(parentDir, d)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			logrus.Fatal(err)
		}
	}

	return nil
}

func getGolangci() (err error) {
	return genFile(filepath.Join(parentDir, ".golangci.yml"), golangciTmpl, nil)
}

func genMain() {
	filePath := filepath.Join(parentDir, "main.go")
	fp, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer fp.Close()

	data, err := RenderTemple(nil, mainTmpl)
	if err != nil {
		logrus.Fatal(err)
	}

	code, err := formatCode(data)
	if err != nil {
		logrus.Fatal(err)
	}

	if _, err := fp.Write(code); err != nil {
		logrus.Fatal(err)
	}

	output, err := exec.Command("go", "mod", "tidy").CombinedOutput()
	if err != nil {
		logrus.Fatalf("cmd--err:%v", err)
	}

	logrus.Infof("go mod tidy %s", output)
}
