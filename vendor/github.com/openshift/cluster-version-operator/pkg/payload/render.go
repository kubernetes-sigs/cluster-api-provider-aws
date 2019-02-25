package payload

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"

	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
)

// Render renders all the manifests from /manifests to outputDir.
func Render(outputDir, releaseImage string) error {
	var (
		manifestsDir  = filepath.Join(DefaultPayloadDir, CVOManifestDir)
		oManifestsDir = filepath.Join(outputDir, "manifests")
		bootstrapDir  = "/bootstrap"
		oBootstrapDir = filepath.Join(outputDir, "bootstrap")

		renderConfig = manifestRenderConfig{ReleaseImage: releaseImage}
	)

	tasks := []struct {
		idir      string
		odir      string
		skipFiles sets.String
	}{{
		idir:      manifestsDir,
		odir:      oManifestsDir,
		skipFiles: sets.NewString("image-references"),
	}, {
		idir:      bootstrapDir,
		odir:      oBootstrapDir,
		skipFiles: sets.NewString(),
	}}
	var errs []error
	for _, task := range tasks {
		if err := renderDir(renderConfig, task.idir, task.odir, task.skipFiles); err != nil {
			errs = append(errs, err)
		}
	}

	agg := utilerrors.NewAggregate(errs)
	if agg != nil {
		return fmt.Errorf("error rendering manifests: %v", agg.Error())
	}
	return nil
}

func renderDir(renderConfig manifestRenderConfig, idir, odir string, skipFiles sets.String) error {
	if err := os.MkdirAll(odir, 0666); err != nil {
		return err
	}
	files, err := ioutil.ReadDir(idir)
	if err != nil {
		return err
	}
	var errs []error
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if skipFiles.Has(file.Name()) {
			continue
		}

		ipath := filepath.Join(idir, file.Name())
		iraw, err := ioutil.ReadFile(ipath)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		rraw, err := renderManifest(renderConfig, iraw)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		opath := filepath.Join(odir, file.Name())
		if err := ioutil.WriteFile(opath, rraw, 0666); err != nil {
			errs = append(errs, err)
			continue
		}
	}

	agg := utilerrors.NewAggregate(errs)
	if agg != nil {
		return fmt.Errorf("error rendering manifests: %v", agg.Error())
	}
	return nil
}

type manifestRenderConfig struct {
	ReleaseImage string
}

// renderManifest Executes go text template from `manifestBytes` with `config`.
func renderManifest(config manifestRenderConfig, manifestBytes []byte) ([]byte, error) {
	tmpl, err := template.New("manifest").Parse(string(manifestBytes))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse manifest")
	}

	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, config); err != nil {
		return nil, errors.Wrapf(err, "failed to execute template")
	}

	return buf.Bytes(), nil
}
