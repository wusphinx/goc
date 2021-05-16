package cmd

import (
	"bytes"
	"go/format"
	"text/template"

	"github.com/sirupsen/logrus"
)

func formatCode(src []byte) ([]byte, error) {
	fs, err := format.Source(src)
	if err != nil {
		logrus.Errorf("format Source err:%v", err)
		return nil, err
	}

	return fs, nil
}

func RenderTemple(param map[string]string, tmpl string) ([]byte, error) {
	Tmpl := template.Must(template.New("tmpl").Parse(tmpl))
	var tmpBuf []byte
	buf := bytes.NewBuffer(tmpBuf)

	if err := Tmpl.Execute(buf, param); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

var gitignoreTmpl = `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with go test -c
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
vendor/
.vscode/
*.log
`

var golangciTmpl = `
linters:
  disable-all: true
  enable:
    - govet
    - golint
    - goimports
    - bodyclose
    - deadcode
    - misspell
    - errcheck
run:
  skip-dirs:
    - vendor
    - tests
  tests: true
  timeout: 1m
`

var mainTmpl = `
package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	r.Run()
}
`

var moduleTmpl = `
module {{.ModuleName}}

go {{.GoVersion}}
`

var buildTmpl = `#!/bin/bash
set -ex

#进入本地代码目录，也是代码编译目录，不需要更改
cd ${APP_COMPILE_PATH}

# 如果需要公司私有git仓库的代码，需要将下面一行注释去掉
git config --global url."https://registry.code.tuya-inc.top".insteadOf "https://gitlab.com"

# golang编译命令, 可以根据自己实际情况修改编译参数， 但是 -o .penglai/${APP_NAME} 后输出二进制执行文件不能改。
GOOS=linux go build -a -o .penglai/${APP_NAME}
`
var preCommitHookTmpl = `#!/usr/bin/env bash
# This file modified from k8s
# https://github.com/kubernetes/kubernetes/blob/master/hooks/pre-commit
# Now It's removed, The Reason is https://github.com/kubernetes/community/issues/729
# The PR is https://github.com/kubernetes/kubernetes/pull/47673

# How to use this hook?
# ln -fs ../../hooks/pre-commit .git/hooks/pre-commit
# In case hook is not executable
# chmod +x .git/hooks/pre-commit

readonly reset=$(tput sgr0)
readonly red=$(tput bold; tput setaf 1)
readonly green=$(tput bold; tput setaf 2)

# readonly goword=./tools/bin/goword

exit_code=0

echo -ne "Checking for files that need gofmt... "
files_need_gofmt=()
files=($(git diff --cached --name-only --diff-filter ACM | grep "\.go$" | grep -v -e "^_vendor"))
for file in "${files[@]}"; do
    # Check for files that fail gofmt.
    diff="$(git show ":${file}" | gofmt -s -d 2>&1)"
    if [[ -n "$diff" ]]; then
        files_need_gofmt+=("${file}")
    fi
done

if [[ "${#files_need_gofmt[@]}" -ne 0 ]]; then
    echo "${red}ERROR!"
    echo "Some files have not been gofmt'd. To fix these errors, "
    echo "copy and paste the following:"
    echo "  gofmt -s -w ${files_need_gofmt[@]}"
    exit_code=1
else
    echo "${green}OK"
fi
echo "${reset}"

# mpaas go mod check
echo -ne "Checking go mod... "
files_go_mod=()
files=("go.mod" "go.sum")
go mod tidy

for file in "${files[@]}"; do
    diff="$(git diff "${file}")"
    if [[ -n "$diff" ]]; then
        files_go_mod+=("${file}")
    fi
done

if [[ "${#files_go_mod[@]}" -ne 0 ]]; then
    echo "${red}ERROR!"
    echo "Some modules may missing or unused."
    exit_code=1
else
    echo "${green}OK"
fi
echo "${reset}"

if [[ "${exit_code}" != 0 ]]; then
    echo "${red}Aborting commit${reset}"
fi
exit ${exit_code}
`
