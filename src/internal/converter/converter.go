package converter

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

func DirToPNG(dir string) error {
	files, err := ioutil.ReadDir(dir)
	sortedFiles := TimeSlice(files)
	if err != nil {
		return err
	}
	sort.Sort(sortedFiles)
	for _, f := range files {
		if !f.IsDir() {
			full := dir + "/" + f.Name()
			err = ToPNG(full)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func DirToJPG(dir string) error {
	files, err := ioutil.ReadDir(dir)
	sortedFiles := TimeSlice(files)
	if err != nil {
		return err
	}
	sort.Sort(sortedFiles)
	for _, f := range files {
		if !f.IsDir() {
			full := dir + "/" + f.Name()
			if strings.Contains(full, "jpg") {
				continue
			}
			err = ToJPG(full)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func ToPNG(filename string) error {
	if strings.Contains(filename, "jpeg") || strings.Contains(filename, "jpg") {
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		//imageData, imType,err := image.Decode(file)
		//if err != nil{
		//	return err
		//}
		data, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}
		//print(imType)
		data, err = ToPng(data)
		if err != nil {
			return err
		}
		baseFileNames := strings.Split(filename, ".")
		baseFileName := strings.Join(baseFileNames[:len(baseFileNames)-1], ".")
		f, err := os.Create(baseFileName + ".png")
		if err != nil {
			return err
		}
		defer f.Close()
		byt, err := f.Write(data)
		if err != nil {
			return err
		}
		fmt.Printf("wrote %d to %s \n", byt, baseFileName+".png")
		//err = os.WriteFile(baseFileName+".png",data,os.ModePerm)
		//if err != nil{
		//	return err
		//}
		_ = file.Close()
		err = os.Remove(filename)
		return err
	}
	return nil
}
func ToJPG(filename string) error {
	if strings.Contains(filename, "png") {
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		//imageData, imType,err := image.Decode(file)
		//if err != nil{
		//	return err
		//}
		data, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}
		//print(imType)
		data, err = ToJpg(data)
		if err != nil {
			return err
		}
		baseFileNames := strings.Split(filename, ".")
		baseFileName := strings.Join(baseFileNames[:len(baseFileNames)-1], ".")
		f, err := os.Create(baseFileName + ".jpg")
		if err != nil {
			return err
		}
		defer f.Close()
		byt, err := f.Write(data)
		if err != nil {
			return err
		}
		fmt.Printf("wrote %d to %s \n", byt, baseFileName+".jpg")
		//err = os.WriteFile(baseFileName+".png",data,os.ModePerm)
		//if err != nil{
		//	return err
		//}
		_ = file.Close()
		err = os.Remove(filename)
		return err
	}
	return nil
}

func ToPng(imageBytes []byte) ([]byte, error) {
	contentType := http.DetectContentType(imageBytes)

	switch contentType {
	case "image/png":
	case "image/jpeg":
		img, err := jpeg.Decode(bytes.NewReader(imageBytes))
		if err != nil {
			return nil, errors.Wrap(err, "unable to decode jpeg")
		}

		buf := new(bytes.Buffer)
		if err := png.Encode(buf, img); err != nil {
			return nil, errors.Wrap(err, "unable to encode png")
		}

		return buf.Bytes(), nil
	}

	return imageBytes, nil //fmt.Errorf("unable to convert %#v to png", contentType)
}

func ToJpg(imageBytes []byte) ([]byte, error) {
	contentType := http.DetectContentType(imageBytes)

	switch contentType {
	case "image/jpeg":
	case "image/png":
		img, err := png.Decode(bytes.NewReader(imageBytes))
		if err != nil {
			return nil, errors.Wrap(err, "unable to decode png")
		}

		buf := new(bytes.Buffer)

		if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 100}); err != nil {
			return nil, errors.Wrap(err, "unable to encode jpeg")
		}

		return buf.Bytes(), nil
	}

	return imageBytes, nil //fmt.Errorf("unable to convert %#v to jpeg", contentType)
}
