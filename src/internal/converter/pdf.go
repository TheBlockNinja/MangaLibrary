package converter

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/signintech/gopdf"
)

type TimeSlice []fs.FileInfo

func (p TimeSlice) Len() int {
	return len(p)
}

func (p TimeSlice) Less(i, j int) bool {
	return p[i].ModTime().Before(p[j].ModTime())
}

func (p TimeSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func CreatePDF(filename string, outputFile string, refresh bool) error {
	if !refresh {
		if _, err := os.Stat(outputFile); err == nil {
			return nil
		}
	}
	err := DirToJPG(filename)
	if err != nil {
		return err
	}
	files, err := ioutil.ReadDir(filename)
	sortedFiles := TimeSlice(files)
	if err != nil {
		return err
	}
	sort.Sort(sortedFiles)
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	for _, f := range sortedFiles {
		fName := filename + "/" + f.Name()
		if strings.Contains(fName, "pdf") || f.IsDir() {
			continue
		}
		pdf.AddPage()
		err = pdf.Image(fName, 0, 0, gopdf.PageSizeA4)
		if err != nil {
			fmt.Printf("Error %s %s\n", fName, err.Error())
			_ = os.Remove(fName)
			return nil
		}

	}
	return pdf.WritePdf(outputFile)
}
