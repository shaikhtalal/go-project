package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var (
	moduleName string
)

const (
	fileName = "gqlgen.yml"
)

// graphqlCmd represents the graphql command
var graphqlCmd = &cobra.Command{
	Use:   "gql",
	Short: "generate graphql resolver, model and schema resolver",
	RunE:  runGQL,
}

func runGQL(cmd *cobra.Command, args []string) error {
	moduleName = viper.GetString("module")

	err := validateModule()
	if err != nil {
		return err
	}

	_, err = os.Stat(fileName)
	if os.IsNotExist(err) {
		err = generateGenGql()
		if err != nil {
			return err
		}
	}

	genGQLs := GenGQL{
		Schema: []string{fmt.Sprintf("%s/%s/%s", "pkg", moduleName, "*.graphqls")},
		Exec: Exec{
			Filename: fmt.Sprintf("%s/%s/%s/%s/%s", "pkg", moduleName, "graph", "generated", "generated.go"),
			Package:  "generated",
		},
		Model: Model{
			Filename: fmt.Sprintf("%s/%s/%s/%s/%s", "pkg", moduleName, "graph", "model", "models_gen.go"),
			Package:  "model",
		},
		Resolver: Resolver{
			Layout:  "follow-schema",
			Dir:     fmt.Sprintf("%s/%s/%s", "pkg", moduleName, "graph"),
			Package: moduleName,
		},
	}

	b, err := marshalYaml(genGQLs)
	if err != nil {
		return err
	}

	if b != nil {
		path, err := filepath.Abs(fileName)
		if err != nil {
			return errors.Wrap(err, "no absolute path found")
		}

		ioutil.WriteFile(path, b, 0644)
	}
	return nil
}

func init() {
	graphqlCmd.Flags().StringVarP(&moduleName, "module", "m", "", "path to bind the package")

	RootCmd.AddCommand(graphqlCmd)
	viper.BindPFlags(graphqlCmd.Flags())
}

type GenGQL struct {
	Schema   []string `yaml:"schema"`
	Exec     Exec     `yaml:"exec"`
	Model    Model    `yaml:"model"`
	Resolver Resolver `yaml:"resolver"`
}

type Exec struct {
	Filename string `yaml:"filename"`
	Package  string `yaml:"package"`
}

type Model struct {
	Filename string `yaml:"filename"`
	Package  string `yaml:"package"`
}

type Resolver struct {
	Layout  string `yaml:"layout"`
	Dir     string `yaml:"dir"`
	Package string `yaml:"package"`
}

func getRootPath() string {
	rootPath, _ := os.Getwd()
	return rootPath
}

func validateModule() error {
	_, err := os.Stat(fmt.Sprintf("%s/%s/%s", getRootPath(), "pkg", moduleName))
	if os.IsNotExist(err) {
		return errors.Wrap(err, "no module found, please create first")
	}
	return nil
}

func generateGenGql() error {
	path := getRootPath()

	w, err := os.Create(fmt.Sprintf("%s/%s", path, "gqlgen.yml"))
	if err != nil {
		w.Close()
		return errors.Wrap(err, "failed to create file")
	}
	defer w.Close()
	return nil
}

func marshalYaml(genGQLs GenGQL) ([]byte, error) {
	b, err := yaml.Marshal(genGQLs)
	if err != nil {
		return nil, errors.Wrap(err, "yaml marshalling failed")
	}
	return b, nil
}
