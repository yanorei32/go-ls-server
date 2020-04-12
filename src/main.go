package main

import (
	"path/filepath"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"gopkg.in/go-playground/validator.v9"
	yaml "gopkg.in/yaml.v3"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Configure struct {
	Directory	string		`yaml:"directory" validate:"required"`
	Port		int			`yaml:"port" validate:"required,gte=1,lte=65536"`
}

type Image struct {
	LastModified	int64	`json:"lastModified"`
	FileName		string	`json:"filename"`
}

func readConf(path string) (Configure, error) {
	var c Configure

	buf, err := ioutil.ReadFile(path)

	if err != nil {
		return c, err
	}

	if err := yaml.Unmarshal(buf, &c); err != nil {
		return c, err
	}

	v := validator.New()

	if err := v.Struct(&c); err != nil {
		return c, err
	}

	return c, nil
}

func main() {
	binpath, err := os.Executable()

	if err != nil {
		log.Fatal(err)
	}

	cfg, err := readConf(filepath.Join(
		filepath.Dir(binpath),
		"configure.yml",
	))

	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/", func(c echo.Context) error {
		var imgs []Image

		files, err := ioutil.ReadDir(cfg.Directory)

		if err != nil {
			return err
		}

		for _, f := range files {
			if f.IsDir() {
				continue
			}

			imgs = append(
				imgs,
				Image {
					LastModified: f.ModTime().UnixNano(),
					FileName: f.Name(),
				},
			)
		}

		return c.JSON(http.StatusOK, imgs)
	})

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(cfg.Port)))
}

