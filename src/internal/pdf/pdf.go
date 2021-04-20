package pdf

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"sort"

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
func CreatePDFV2(filename string, outputFile string, refresh bool) error {
	fmt.Printf("pdffile:%s\n", outputFile)
	if !refresh {

		if _, err := os.Stat(outputFile); err == nil {
			return nil
		}
	}

	files, err := ioutil.ReadDir(filename)
	sortedFiles := TimeSlice(files)
	if err != nil {
		return err
	}
	sort.Sort(sortedFiles)
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	fmt.Printf("files %d\n", len(sortedFiles))
	for _, f := range sortedFiles {
		fName := filename + f.Name()
		if f.IsDir() {
			continue
		}
		fmt.Printf("file:%s\n", fName)
		pdf.AddPage()
		err = pdf.Image(fName, 0, 0, gopdf.PageSizeA4)
		if err != nil {
			fmt.Printf("ERROR %s\n", err.Error())
			continue
		}

	}
	return pdf.WritePdf(outputFile)
}

func GetFileSize(filename string) int {
	Myfile, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file!!!")
	}

	//We moved file pointer to end of file.
	size, err := Myfile.Seek(0, 2)
	if err != nil {
		fmt.Println(err)
	}
	Myfile.Close()
	return int(size)
}
