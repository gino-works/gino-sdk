package crds

import (
	"embed"
	"fmt"
	"io/fs"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	yaml "sigs.k8s.io/yaml"
)

// Assets is the embedded assets objects for the included crd yaml data.
//
//go:embed *.yaml
var Assets embed.FS

// LoadCRDs loads crds from the embed filesystem.
func LoadCRDs() ([]*apiextensionsv1.CustomResourceDefinition, error) {
	files, err := fs.Glob(Assets, "*.yaml")
	if err != nil {
		return nil, fmt.Errorf("globing crds: %w", err)
	}

	crds := make([]*apiextensionsv1.CustomResourceDefinition, 0, len(files))

	for _, file := range files {
		b, err := fs.ReadFile(Assets, file)
		if err != nil {
			return nil, fmt.Errorf("reading crd %q: %w", file, err)
		}

		var crd apiextensionsv1.CustomResourceDefinition

		err = yaml.Unmarshal(b, &crd)
		if err != nil {
			return nil, fmt.Errorf("unmarshaling crd %q: %w", file, err)
		}

		crds = append(crds, &crd)
	}

	return crds, nil
}
