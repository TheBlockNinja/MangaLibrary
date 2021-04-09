package api

import (
	"MangaLibrary/src/internal/converter"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/TheBlockNinja/WebParser"
	"go.uber.org/zap"
)

type API struct {
	Parser       *WebParser.Parser
	SiteData     map[string]*WebComponentV2
	BasePath     string
	UsefulInfo   map[string]*WebComponentV2
	MetaData     map[string]interface{}
	UseSubPath   bool
	Name         *WebComponentV2
	NameURL      string
	SearchURL    string
	ProgressData *Progress
}
type Progress struct {
	Name          string    `json:"name"`
	URL           string    `json:"url"`
	Current       int       `json:"current"`
	Total         int       `json:"total"`
	Message       string    `json:"Message"`
	ForceStop     bool      `json:"force_stop"`
	RemainingTime string    `json:"remaining_time"`
	LastUpdate    time.Time `json:"last_update"`
}

func NewAPI(logger *zap.Logger) *API {
	parser := WebParser.NewParser(logger)
	return &API{
		Parser:     parser,
		SiteData:   map[string]*WebComponentV2{},
		BasePath:   "",
		UseSubPath: false,
		Name:       nil,
		SearchURL:  "",
		MetaData:   map[string]interface{}{},
		UsefulInfo: map[string]*WebComponentV2{},
	}
}
func (a *API) GetURL(url string, search map[string]string) string {
	output := url
	for k, v := range search {
		output = strings.ReplaceAll(output, fmt.Sprintf("{%s}", k), v)
	}
	return output
}

func (a API) Process(url string, urlData map[string]string, progress *Progress) error {
	a.Parser.Logger.Info(fmt.Sprintf("getting url %s", url))
	err := a.Parser.Get(url)
	if err != nil {
		return err
	}
	mainHTMLData := *a.Parser.Html
	fullName := "tmp"
	if a.NameURL != "" {
		nameUrl := a.GetURL(a.NameURL, urlData)
		a.Parser.Logger.Info(fmt.Sprintf("getting url name %s", nameUrl))
		fullName, err = a.GetNameWithURL(nameUrl)
		if err != nil {
			a.Parser.Logger.Error("failed getting name data", zap.Error(err))
			return err
		}
	} else {
		fullName, err = a.GetName()
		a.Parser.Logger.Info(fmt.Sprintf("name %s", fullName))
	}
	if progress != nil {
		progress.Name = fullName
		progress.LastUpdate = time.Now()
	}
	for k, v := range a.SiteData {
		a.Parser.Html = &mainHTMLData
		err = v.Get(a.Parser, fullName, a.UseSubPath, progress)
		if err != nil {
			return err
		}
		a.MetaData[k] = v.ElementLinks
	}
	if progress != nil {
		progress.Message = "Success"
	}

	return nil
}

func (a API) GetName() (string, error) {
	if a.Name != nil {
		err := a.Name.Get(a.Parser, "", false, nil)
		if err != nil {
			return "", err
		}
		space := regexp.MustCompile(`\s+`)
		fullName := a.BasePath + "/"
		for _, v := range a.Name.ElementLinks {
			if a.Name.OtherData != nil && len(a.Name.OtherData) > 0 {
				for k, _ := range a.Name.OtherData {
					name := space.ReplaceAllString(v[0].Attributes[k], "_")
					fullName += name
				}
			} else {
				name := space.ReplaceAllString(v[0].TextData, "_")
				fullName += name
			}

		}
		a.Parser.Logger.Info(fmt.Sprintf("library name: %s", fullName))
		return fullName, nil
	}
	return "", nil
}

func (a API) GetNameWithURL(url string) (string, error) {
	err := a.Parser.Get(url)
	if err != nil {
		return "", err
	}

	if a.Name != nil {
		err = a.Name.Get(a.Parser, "", false, nil)
		if err != nil {
			a.Parser.Logger.Error("failed getting name", zap.Error(err))
			return "", err
		}
		space := regexp.MustCompile(`\s+`)
		fullName := a.BasePath + "/"
		for _, v := range a.Name.ElementLinks {
			if a.Name.OtherData != nil && len(a.Name.OtherData) > 0 {
				for k, _ := range a.Name.OtherData {
					name := space.ReplaceAllString(v[0].Attributes[k], "_")
					fullName += name
				}
			} else {
				name := space.ReplaceAllString(v[0].TextData, "_")
				fullName += name
			}
		}
		a.Parser.Logger.Info(fmt.Sprintf("library name: %s", fullName))

		return fullName, nil
	}
	return "", nil
}

func (a API) CreatePDF(path string) ([]string, error) {
	output := []string{}
	pdfLocation := path + "/pdf"
	_ = os.MkdirAll(pdfLocation, os.ModePerm)
	dirs, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	sortedDirs := converter.TimeSlice(dirs)
	sort.Sort(sortedDirs)
	for _, directories := range sortedDirs {
		pdfFileName := pdfLocation + "/" + directories.Name() + ".pdf"
		if directories.IsDir() {
			err = converter.CreatePDF(path+"/"+directories.Name(), pdfFileName, true)
			output = append(output, pdfFileName)
			if err != nil {
				return nil, err
			}

		} else {
			name := strings.Split(path, "/")
			pdfFileName = pdfLocation + "/" + name[len(name)-1] + ".pdf"
			err = converter.CreatePDF(path, pdfFileName, true)
			output = append(output, pdfFileName)
			if err != nil {
				return nil, err
			}
			return output, nil
		}
	}

	return output, nil
}

func (a API) SaveMetaData() {

}

type WebComponentV2 struct {
	Tag            string `json:"tag"`
	Attribute      string
	Value          string
	IsLink         bool
	IsDownload     bool
	Name           string
	Reverse        bool
	linkAttributes []string
	Elements       []*WebParser.HTMLData
	ElementLinks   map[string][]*WebParser.HTMLData
	Children       []*WebComponentV2
	Delay          int
	OtherData      map[string]string
	Type           string
}

func (w *WebComponentV2) Get(Parser *WebParser.Parser, path string, useSubPath bool, progress *Progress) error {
	if _, found := w.ElementLinks[Parser.URL]; !found {
		if Parser.Html == nil {
			if Parser.URL != "" {
				err := Parser.Get(Parser.URL)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("there is no html data in the parser")
			}
		}
	}
	if w.ElementLinks == nil {
		w.ElementLinks = map[string][]*WebParser.HTMLData{}
	}
	didFindElementLink := false
	if _, found := w.ElementLinks[Parser.URL]; !found {
		if w.Tag == "" {
			w.ElementLinks[Parser.URL] = Parser.Html.FindAttribute(w.Attribute, w.Value)
		} else if w.Tag != "" && w.Attribute != "" {
			w.ElementLinks[Parser.URL] = Parser.Html.Find(w.Tag, w.Attribute, w.Value)
		} else {
			w.ElementLinks[Parser.URL] = Parser.Html.FindTag(w.Tag)
		}
		if w.Reverse {
			w.ElementLinks[Parser.URL] = WebParser.Reverse(w.ElementLinks[Parser.URL])
		}
	} else {
		didFindElementLink = true
	}
	if len(w.ElementLinks[Parser.URL]) == 0 {
		return fmt.Errorf("failed to find data for %s %s %s", w.Tag, w.Attribute, w.Value)
	}
	Parser.Logger.Info(fmt.Sprintf("found %d elements", len(w.ElementLinks[Parser.URL])))
	if len(w.Children) > 0 && w.IsLink {
		for _, l := range w.ElementLinks {
			if progress != nil {
				progress.Total += len(l) * len(w.Children)
				if progress.ForceStop {
					return fmt.Errorf("force stop implemented")
				}
			} else {
				Parser.Logger.Error("progress is nil")
			}
			for _, d := range l {
				for i, _ := range w.Children {

					link, err := d.GetLink(w.linkAttributes, Parser)
					if err != nil {
						return err
					}
					if progress != nil {
						otherTime := 0
						if progress.RemainingTime != "" {
							tmp, err := time.Parse("2006-01-02 15:04:05", progress.RemainingTime)
							if err != nil {

							} else {
								otherTime = tmp.Second()
							}
						}

						updateTime := time.Now().Second() - progress.LastUpdate.Second()
						remaining := progress.Total - progress.Current
						willFinish := time.Now().Add(time.Second * time.Duration(((updateTime*remaining)+otherTime)/2))
						Parser.Logger.Info(fmt.Sprintf("Seconds since last update %d", updateTime))
						progress.RemainingTime = willFinish.Format("2006-01-02 15:04:05")
						progress.LastUpdate = time.Now()
						progress.Current += 1
						progress.URL = link
						progress.Message = "Downloading"
						if progress.ForceStop {
							progress.Message = "Force Stopped"
							return fmt.Errorf("force stop implemented")
						}
					}
					name := d.Attributes[w.Name]
					if name == "" {
						test := d.Flatten()
						space := regexp.MustCompile(`\s+`)
						name = space.ReplaceAllString(test.TextData, "_")
						Parser.Logger.Info(fmt.Sprintf("current element text: %s", name))

					}

					pathName := fmt.Sprintf("%s/%s", path, name)
					if useSubPath {
						pathName = path
					}

					childParser := WebParser.NewParser(Parser.Logger)
					if !didFindElementLink {
						time.Sleep(time.Duration(w.Delay) * time.Second)
					}
					err = childParser.Get(link)
					if err != nil {
						return err
					}
					childParser.URL = link

					err = w.Children[i].Get(childParser, pathName, useSubPath, progress)
					if err != nil {
						return err
					}
				}

			}
		}
	}
	if w.IsDownload {
		w.Download(path, Parser, Parser.URL)
	}
	return nil
}
func (w *WebComponentV2) Download(rootPath string, mainParser *WebParser.Parser, url string) {
	DownloadParser := WebParser.NewParser(mainParser.Logger)
	page_count := 0
	if v, found := w.ElementLinks[url]; found {
		for _, d := range v {
			link, err := d.GetLink(w.linkAttributes, mainParser)
			if err != nil {
				continue
			}
			urlParse := strings.Split(link, ".")
			ext := urlParse[len(urlParse)-1]
			name := d.Attributes[w.Name]
			if name == "" {
				tmp := urlParse[len(urlParse)-2]
				tmpAr := strings.Split(tmp, "/")
				name = tmpAr[len(tmpAr)-1]
			}
			page_count += 1
			filename := fmt.Sprintf("%s/%s.%s", rootPath, name, ext)
			filenamePNG := fmt.Sprintf("%s/%s.png", rootPath, name)
			filenameJPG := fmt.Sprintf("%s/%s.jpg", rootPath, name)
			mainParser.Logger.Info(fmt.Sprintf("download image %s : %s  %d out of %d", link, rootPath, page_count, len(v)))
			if _, err := os.Stat(filename); err == nil {
				_ = setFileModTime(filename)
				continue
			}
			if _, err := os.Stat(filenamePNG); err == nil {
				_ = setFileModTime(filenamePNG)
				continue
			}
			if _, err := os.Stat(filenameJPG); err == nil {
				_ = setFileModTime(filenameJPG)
				continue
			}
			_ = os.MkdirAll(rootPath, os.ModePerm)
			time.Sleep(time.Duration(w.Delay) * time.Second)
			_ = DownloadParser.Download(link, filename)
			if ext == "webp" {
				mainParser.Logger.Info(fmt.Sprintf("converting %s to png", filename))
				err = converter.ConvertWebPToPng(filename, filenamePNG)
				if err != nil {
					mainParser.Logger.Error("failed converting webp to png", zap.Error(err))
					continue
				}

				_ = os.Remove(filename)
			}
		}
	}

}
func setFileModTime(filename string) error {
	if _, err := os.Stat(filename); err == nil {
		currentTime := time.Now().Local()
		err := os.Chtimes(filename, currentTime, currentTime)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}
