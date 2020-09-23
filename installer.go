package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// InstallerOptions are used to create an installer instance.
type InstallerOptions struct {
	// dlURLFmt is the URL to the file to be downloaded.
	// It uses %s to inject the version into the URL.
	dlURLFmt string

	// dlDir is the absolute or relative path from the runtime to the download directory.
	dlDir string

	// fileFmt is the format the should be used for downloaded files.
	// It uses %s to inject the version.
	fileFmt string

	// install is called with the path to the dowloaded version file.
	// In this function you can perform the installation steps for the given version file.
	install func(path string)

	// noCache determines whether or not downloaded files should be reused.
	noCache bool
}

// Installer provides methods for downloading and installing a tool at a given URL.
// It is reccommended to call Create(InstallerOptions{}) or InstallerOptions.Create()
// instead of creating an instance manually.
type Installer struct {
	opts InstallerOptions

	// dlDir is the absolute path to the download directory.
	dlDir string

	// dlPathFmt is the absolute path to downloaded files.
	// It uses %s to inject the version into the filename.
	dlPathFmt string
}

// Create creates a new installer instance using the InstallerOptions.
func (o *InstallerOptions) Create() Installer {
	var dlDir string
	if filepath.IsAbs(o.dlDir) {
		dlDir = o.dlDir
	} else {
		dlDir = filepath.Join(RunDir(), o.dlDir)
	}

	return Installer{
		*o,
		dlDir,
		filepath.Join(dlDir, o.fileFmt),
	}
}

// Create creates a new installer instance using the InstallerOptions.
func Create(options InstallerOptions) Installer {
	return options.Create()
}

func (i *Installer) url(version string) string {
	return fmt.Sprintf(i.opts.dlURLFmt, version)
}

func (i *Installer) dlPath(version string) string {
	return fmt.Sprintf(i.dlPathFmt, version)
}

// Download downloads a specified version.
func (i *Installer) Download(version string) (string, error) {
	if !i.opts.noCache {
		fmt.Printf("Previously downloaded?...")
		dlPath := i.dlPath(version)
		info, err := os.Stat(dlPath)
		if err == nil && !info.IsDir() {
			fmt.Printf(" yes\n")
			return dlPath, nil
		}
		fmt.Printf(" no\n")
	}

	fmt.Printf("Download directory exists?...")
	dlDir := i.dlDir
	_, err := os.Stat(dlDir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf(" no\n")
			fmt.Printf("Creating download directory...")
			err := os.MkdirAll(dlDir, os.ModePerm)
			if err != nil {
				fmt.Printf(" error\n")
				return "", err
			}
			fmt.Printf(" done\n")
		} else {
			fmt.Printf(" error\n")
			return "", err
		}
	} else {
		fmt.Printf(" yes\n")
	}

	fmt.Printf("Downloading file...")
	resp, err := http.Get(i.url(version))
	if err != nil {
		fmt.Printf(" error\n")
		return "", err
	}
	fmt.Printf(" done\n")
	defer resp.Body.Close()

	fmt.Printf("Creating destination file...")
	dlPath := i.dlPath(version)
	file, err := os.Create(dlPath)
	if err != nil {
		fmt.Printf(" error\n")
		return "", err
	}
	fmt.Printf(" done\n")
	defer file.Close()

	fmt.Printf("Copying downloaded file...")
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Printf(" error\n")
		return "", err
	}
	fmt.Printf(" done\n")
	return dlPath, err
}

// Install installs a specified version.
func (i *Installer) Install(version string) error {
	path, err := i.Download(version)
	if err != nil {
		return err
	}

	i.opts.install(path)

	return nil
}
