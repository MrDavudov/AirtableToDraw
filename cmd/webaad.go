package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	
)

type config struct {
	hostURL         string
	username        string
	password        string
	ncFolderPath    string
	ncFileName      string
	localFolderPath string
	localFileName   string
	download        bool
	upload          bool
}

func Nextcloud() {
	c, err := getConfig()
	red := color.New(color.FgRed)
	if err != nil {
		red.Println(err)
		return
	}

	if c.upload == c.download {
		red.Println("Следует выбрать только загрузку или загрузку.")
		return
	}
	localFilePath := filepath.Join(c.localFolderPath, c.localFileName)
	nc := Config{HostURL: c.hostURL, Username: c.username, Password: c.password}
	currentDir, err := os.Getwd()
	if err != nil {
		red.Println(err)
		return
	}
	tempFolderName := filepath.Join(currentDir, "temp", generateRandomName())
	// Создать временный каталог
	if err := os.MkdirAll(tempFolderName, os.ModePerm); err != nil {
		red.Println(err)
		return
	}
	defer os.RemoveAll(tempFolderName)
	if c.download {
		color.Green("Downloading nextcloud file %s from nextcloud path %s to local directory %s\n", c.ncFileName, c.ncFolderPath, localFilePath)
		if err := nc.DownloadFile(c.ncFolderPath, c.ncFileName, tempFolderName); err != nil {
			red.Println(err)
			return
		}
		if err := os.Rename(filepath.Join(tempFolderName, c.ncFileName), localFilePath); err != nil {
			red.Println(err)
			return
		}
		color.Green("Download successful")
	} else if c.upload {
		ncPath := filepath.Join(c.ncFolderPath, c.ncFileName)
		color.Green("Uploading local file %s from local path %s to nextcloud directory %s\n", c.localFileName, c.localFolderPath, ncPath)
		tempFilePath := filepath.Join(tempFolderName, c.ncFileName)
		copyFile(localFilePath, tempFilePath)
		nc.UploadFile(c.ncFolderPath, c.ncFileName, tempFilePath)
		color.Green("Upload successful")
	}
}

// getConfig читает параметры из файла env и параметры cmd и возвращает конфигурацию в структуре
func getConfig() (config, error) {
	v := viper.New()
	// Support reading file & folder paths for nextcloud and local from flags
	pflag.String("nc_folder_path", "", "nextcloud folder path")
	v.BindPFlag("nc_folder_path", pflag.Lookup("nc_folder_path"))
	pflag.String("nc_file_name", "", "nextcloud file name")
	v.BindPFlag("nc_file_name", pflag.Lookup("nc_file_name"))
	pflag.String("local_folder_path", "", "local folder path")
	v.BindPFlag("local_folder_path", pflag.Lookup("local_folder_path"))
	pflag.String("local_file_name", "", "local file name")
	v.BindPFlag("local_file_name", pflag.Lookup("local_file_name"))

	pflag.Bool("d", false, "download from nextcloud server")
	v.BindPFlag("d", pflag.Lookup("d"))
	pflag.Bool("u", false, "upload to nextcloud server")
	v.BindPFlag("u", pflag.Lookup("u"))

	pflag.Parse()
	v.SetConfigName("config")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	if err := v.ReadInConfig(); err != nil {
		return config{}, errors.Wrap(err, "Cannot read env config, env file should be named config.env")
	}
	v.AutomaticEnv()

	// viper cannot unmarshal when using AutomaticEnv
	var c config = config{
		hostURL:         v.GetString("host_url"),
		username:        v.GetString("username"),
		password:        v.GetString("password"),
		ncFolderPath:    v.GetString("nc_folder_path"),
		ncFileName:      v.GetString("nc_file_name"),
		localFolderPath: v.GetString("local_folder_path"),
		localFileName:   v.GetString("local_file_name"),
		download:        v.GetBool("d"),
		upload:          v.GetBool("u"),
	}
	return c, nil
}

func generateRandomName() string {
	n := 5
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	s := fmt.Sprintf("%X", b)
	return s
}

func copyFile(src, dest string) error {
	// Open original file
	original, err := os.Open(src)
	if err != nil {
		return err
	}
	defer original.Close()

	// Create new file
	new, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer new.Close()

	// This will copy
	_, err = io.Copy(new, original)
	if err != nil {
		return err
	}
	return nil
}

// Config содержит необходимые данные для подключения к экземпляру nextcloud.
type Config struct {
	HostURL  string
	Username string
	Password string
}

// UploadFile загружает файл в указанный remoteFolderPath и fileName из текущего каталога
func (c *Config) UploadFile(remoteFolderPath, remoteFileName, localFilePath string) error {
	dat, err := ioutil.ReadFile(localFilePath)
	if err != nil {
		return err
	}

	uploadURL := fmt.Sprintf("%s/remote.php/dav/files/%s/%s/%s", c.HostURL, c.Username, remoteFolderPath, remoteFileName)
	client := &http.Client{}
	req, err := http.NewRequest("PUT", uploadURL, bytes.NewBuffer(dat))
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.Username, c.Password)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// DownloadFile загружает файл из заданного remoteFolderPath и remoteFileName в текущий каталог
func (c *Config) DownloadFile(remoteFolderPath, remoteFileName, localDownloadPath string) error {
	downloadURL := fmt.Sprintf("%s/remote.php/dav/%s/%s", c.HostURL, remoteFolderPath, remoteFileName)
	client := &http.Client{}
	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.Username, c.Password)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	downloadPath := filepath.Join(localDownloadPath, remoteFileName)

	out, err := os.Create(downloadPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	return err
}