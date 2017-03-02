package user

import (
	"io/ioutil"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/Unknwon/com"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/rai-project/auth"
	"github.com/rai-project/config"
)

var (
	DefaultProfilePath string
)

type Profile struct {
	Username  string `toml:"username"`
	AccessKey string `toml:"access_key"`
	SecretKey string `toml:"secret_key"`
}

func NewProfile(path string) (*Profile, error) {
	if path == "" {
		path = DefaultProfilePath
	}
	if !com.IsFile(path) {
		return nil, errors.Errorf("unable to locate %v. not such file or directory", path)
	}
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read profile")
	}
	profile := new(Profile)
	_, err = toml.Decode(string(buf), profile)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (p *Profile) Verify() bool {
	return auth.Verify(p.Username, p.AccessKey, p.SecretKey)
}

func init() {
	config.AfterInit(func() {
		homeDir, err := homedir.Dir()
		if err != nil {
			return
		}

		appName := config.App.Name

		// load ~/.rai_profile
		homeProfileFile := filepath.Join(homeDir, "."+appName+"_profile")
		if com.IsFile(homeProfileFile) {
			DefaultProfilePath = homeProfileFile
			return
		}

		// load ~/.rai_env
		homeEnvFile := filepath.Join(homeDir, "."+appName+"_env")
		if com.IsFile(homeEnvFile) {
			DefaultProfilePath = homeEnvFile
			return
		}

		// load ~/.rai.profile
		homeProfileFile = filepath.Join(homeDir, "."+appName+".profile")
		if com.IsFile(homeProfileFile) {
			DefaultProfilePath = homeProfileFile
			return
		}

		// load ~/.rai.env
		homeEnvFile = filepath.Join(homeDir, "."+appName+".env")
		if com.IsFile(homeEnvFile) {
			DefaultProfilePath = homeEnvFile
			return
		}

		DefaultProfilePath = appName + ".profile"
	})
}
