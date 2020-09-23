package main

import (
	"testing"
)

const (
	DlURLFmt = "https://update.code.visualstudio.com/%s/linux-x64/stable"
	FileFmt  = "code-%s.tar.gz"
)

func dl(o InstallerOptions, version string, t *testing.T) {
	i := o.Create()

	_, err := i.Download(version)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDownloadNoDlDir(t *testing.T) {
	o := InstallerOptions{
		dlURLFmt: DlURLFmt,
		fileFmt:  FileFmt,
	}

	dl(o, "1.49", t)
}

func TestDownloadRelDir(t *testing.T) {
	o := InstallerOptions{
		dlURLFmt: DlURLFmt,
		fileFmt:  FileFmt,
		dlDir:    "./cache/code",
	}

	dl(o, "1.48", t)
}
