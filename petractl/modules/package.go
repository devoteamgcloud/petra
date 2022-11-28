package modules

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/devoteamgcloud/petra/petractl/storage"
	yaml "gopkg.in/yaml.v3"
)

const (
	petraConfigFile = "petra-config.yaml"
)

func PackageModules(workingDir string, recursive bool, b *storage.GCSBackend) error {
	var err error
	fmt.Printf("workingdir: %v and recursive %v\n", workingDir, recursive)
	if recursive {
		err = filepath.Walk(workingDir, func(path string, fi os.FileInfo, err error) error {
			if fi.Name() != petraConfigFile {
				fmt.Println(fi.Name())
				return nil
			}
			return processModule(path, b)
		})
		if err != nil {
			return fmt.Errorf("error: %v", err)
		}
	} else {
		fmt.Printf("filepath: %v\n", filepath.Join(workingDir, petraConfigFile))
		err = processModule(filepath.Join(workingDir, petraConfigFile), backend)
	}
	return err
}

func processModule(path string, b *storage.GCSBackend) error {
	var err error
	// Retrieve Module Specs
	mod := &module{}

	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, mod)

	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	moduleRoot := filepath.Dir(path)

	fmt.Println("Before Archive")
	// Create tgz archive
	archiveBuffer, err := archiveModule(moduleRoot)
	if err != nil {
		return err
	}

	// Upload archive to bucket
	modulePath := modPath(*mod)
	downloadURL, err := b.UploadModule(modulePath, archiveBuffer)
	if err != nil {
		return err
	}

	fmt.Printf("module successfully uploaded at : %s", downloadURL)

	return err
}

func archiveModule(root string) (io.Reader, error) {
	buf := new(bytes.Buffer)
	// ensure the src actually exists before trying to tar it
	if _, err := os.Stat(root); err != nil {
		return buf, fmt.Errorf("unable to tar files - %v", err.Error())
	}

	fmt.Println(root)

	gw := gzip.NewWriter(buf)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	err := filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// return on non-regular files
		if !fi.Mode().IsRegular() {
			return nil
		}

		fmt.Println(fi.Name())
		// create a new dir/file header
		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		// update the name to correctly reflect the desired destination when untaring
		header.Name = strings.TrimPrefix(strings.Replace(path, root, "", -1), string(filepath.Separator))

		fmt.Println("before write header")
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		fmt.Println("Before os open")
		data, err := os.Open(path)
		if err != nil {
			return err
		}

		fmt.Println("Before io copy")
		if _, err := io.Copy(tw, data); err != nil {
			return err
		}

		// manually close here after each file operation; deferring would cause each file close
		// to wait until all operations have completed.
		data.Close()

		return nil
	})

	fmt.Println(err)

	return buf, err
}
