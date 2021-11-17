package exporter

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

type fileExporter struct {
	host     string
	operates map[string]Operate
}

func NewFileExporter(host string) Exporter {
	return &fileExporter{
		host:     host,
		operates: make(map[string]Operate),
	}
}

func (e *fileExporter) Register(name string, operate Operate) {
	e.operates[name] = operate
}

func (e *fileExporter) Run() error {
	engine := gin.Default()
	for k, o := range e.operates {
		path, err := o()
		if err != nil {
			return err
		}

		p, ok := path.(string)
		if !ok {
			return errors.New("operator not return string")
		}
		fmt.Println(p)
		engine.Static(k, p)
	}

	// Listen and serve on 0.0.0.0:8080
	return engine.Run(e.host)
}
