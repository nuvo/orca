package yamlutils

import (
	"fmt"
	"io/ioutil"
	"log"

	genutils "orca/pkg/utils/general"

	yaml "gopkg.in/yaml.v2"
)

type ChartSpec struct {
	Name         string
	Version      string
	Dependencies []string
}

func ChartsYamlToStruct(file string) []ChartSpec {
	var charts []ChartSpec

	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalln(err)
	}

	var v map[string][]map[string]interface{}
	err = yaml.Unmarshal(data, &v)
	if err != nil {
		log.Fatalln(err)
	}

	for _, chart := range v["charts"] {

		var c ChartSpec

		c.Name = chart["name"].(string)
		c.Version = chart["version"].(string)

		if chart["depends_on"] != nil {
			for _, dep := range chart["depends_on"].([]interface{}) {
				depStr := dep.(string)
				c.Dependencies = append(c.Dependencies, depStr)
			}
		}
		charts = append(charts, c)
	}

	return charts
}

func (c ChartSpec) Print() {
	fmt.Println("name: " + c.Name)
	fmt.Println("name: " + c.Version)
	for _, dep := range c.Dependencies {
		fmt.Println("depends_on: " + dep)
	}
}

func RemoveChartFromDependencies(charts []ChartSpec, name string) []ChartSpec {

	var outCharts []ChartSpec

	for _, dependant := range charts {
		if genutils.Contains(dependant.Dependencies, name) {

			index := -1
			for i, elem := range dependant.Dependencies {
				if elem == name {
					index = i
				}
			}
			if index == -1 {
				panic("Could not find element in dependencies")
			}

			dependant.Dependencies[index] = dependant.Dependencies[len(dependant.Dependencies)-1]
			dependant.Dependencies[len(dependant.Dependencies)-1] = ""
			dependant.Dependencies = dependant.Dependencies[:len(dependant.Dependencies)-1]
		}
		outCharts = append(outCharts, dependant)
	}

	return outCharts
}
