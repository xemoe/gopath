package config

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"runtime"

	homedir "github.com/mitchellh/go-homedir"
)

func reloadProfile() {
}

var re = regexp.MustCompile(`export\s+GOPATH=(.+)`)

const PROFILE = ".profile"
const BASH_PROFILE = ".bash_profile"
const BASHRC = ".bashrc"

const SOURCE_COMMAND = "source"

type linux struct {
	gopath  string
	profile string
}

var linuxProfile *linux

func init() {

	if runtime.GOOS == "linux" {
		linuxProfile = &linux{}
		profile()
		PATH = linuxProfile
	}

}

func profile() {

	var profiles = []string{PROFILE, BASH_PROFILE, BASHRC}
	homepath, _ := homedir.Dir()

	for _, pf := range profiles {

		profilePath := path.Join(homepath, pf)
		_, err := os.Stat(profilePath)
		if err == nil {
			if linuxProfile.gopath == "" {
				linuxProfile.profile = profilePath
			}

			gopath := getGoPath(profilePath)

			if gopath != "" {
				linuxProfile.profile = profilePath
				linuxProfile.gopath = gopath
			}
		}

	}

}

func getGoPath(profile string) (gopath string) {
	f, _ := os.Open(profile)
	defer f.Close()
	data, _ := ioutil.ReadAll(bufio.NewReader(f))

	result := re.FindAllStringSubmatch(string(data), 1)
	for _, value := range result {
		for _, s := range value {
			gopath = s
		}
		return
	}

	return
}

func setGoPath(profile, gopath string) (result string, err error) {
	f, err := os.Open(profile)
	if err != nil {
		return
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	f.Close()
	result = re.ReplaceAllString(string(data), fmt.Sprintf(`export GOPATH=%#v`, gopath))

	ioutil.WriteFile(profile, []byte(result), 0644)

	return
}

func (l *linux) GetPath() (path string, err error) {

	path = l.gopath
	return
}

func (l *linux) ChangePath(path string) (err error) {
	linuxProfile.gopath = path

	_, err = setGoPath(linuxProfile.profile, linuxProfile.gopath)

	if err != nil {
		return
	}

	reloadProfile()

	return
}
