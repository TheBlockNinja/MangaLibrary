package driver

import (
	"MangaLibrary/src/internal/converter"
	"MangaLibrary/src/internal/dao"
	"MangaLibrary/src/internal/dto"
	"MangaLibrary/src/internal/jobs"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/TheBlockNinja/WebParser"
)

func WebComponentsToComponent(wcs []*dto.WebComponent) *Component {
	sort.Slice(wcs, func(i, j int) bool {
		return wcs[i].ID < wcs[j].ID
	})
	output := &Component{}
	offSet := wcs[0].ID - 1
	for _, wc := range wcs {
		component := WebComponentToComponent(wc)
		if output.Name == "" {
			output = component
			continue
		}
		if wc.Parent == 0 {
			if output.Siblings == nil {
				output.Siblings = []*Component{}
			}
			output.Siblings = append(output.Siblings, component)
			fmt.Printf("adding sibling to %d %s\n", component.ID, component.Name)
		} else {
			fmt.Printf("looking for parent %d\n", wc.Parent)
			comp := output.FindID(wc.Parent + offSet)
			if comp == nil {
				fmt.Printf("looking for parent failed....%d\n", wc.Parent)
				continue
			}
			if comp.Children == nil {
				comp.Children = []*Component{}
			}
			comp.Children = append(comp.Children, component)

			fmt.Printf("adding child to %d %s\n", comp.ID, comp.Name)
		}
	}
	return output
}

func WebComponentToComponent(wc *dto.WebComponent) *Component {
	var linkAttributes []string
	meta := map[string]*MetaData{}
	ele := &Elements{
		Data: map[string][]*WebParser.HTMLData{},
	}
	if wc.LinkAttributes != "" {
		err := json.Unmarshal([]byte(wc.LinkAttributes), &linkAttributes)
		if err != nil {
			fmt.Printf("Link Attribute Err %s\n", err.Error())
		}
	}
	if wc.MetaData != "" {
		err := json.Unmarshal([]byte(wc.MetaData), &meta)
		if err != nil {
			fmt.Printf("MetaData Err %s\n", err.Error())

		}
	}
	if wc.ElementData != "" {
		_ = json.Unmarshal([]byte(wc.ElementData), ele)
	}
	component := &Component{
		ID:             wc.ID,
		Name:           wc.Name,
		Tag:            wc.Tag,
		Attribute:      wc.Attribute,
		Value:          wc.Value,
		IsLink:         wc.IsLink,
		IsDownload:     wc.IsDownload,
		LinkAttributes: linkAttributes,
		ElementData:    ele,
		Children:       []*Component{},
		Siblings:       []*Component{},
		Delay:          wc.Delay,
		MetaData:       meta,
		Reverse:        wc.Reverse,
	}
	return component
}

type Elements struct {
	Data     map[string][]*WebParser.HTMLData `json:"data"`
	MetaData map[string][]map[string]*MetaData
}

func (e *Elements) LoadData(Parser *WebParser.Parser, c *Component) error {
	if _, found := e.Data[Parser.URL]; !found {
		if c.Tag == "" {
			e.Data[Parser.URL] = Parser.Html.FindAttribute(c.Attribute, c.Value)
		} else if c.Tag != "" && c.Attribute != "" {
			e.Data[Parser.URL] = Parser.Html.Find(c.Tag, c.Attribute, c.Value)
		} else {
			e.Data[Parser.URL] = Parser.Html.FindTag(c.Tag)
		}
		if c.Reverse {
			e.Data[Parser.URL] = WebParser.Reverse(e.Data[Parser.URL])
		}
	}
	if len(e.Data[Parser.URL]) == 0 {
		return fmt.Errorf("failed to find data for %s %s %s\n", c.Tag, c.Attribute, c.Value)
	}
	return nil
}
func (e *Elements) GetLength(Parser *WebParser.Parser) int {
	return len(e.Data[Parser.URL])
}
func (e *Elements) GetLinks(Parser *WebParser.Parser, c *Component) ([]string, error) {
	var links []string
	for _, sites := range e.Data[Parser.URL] {
		link, err := sites.GetLink(c.LinkAttributes, Parser)
		if err != nil {
			return links, err
		}
		links = append(links, link)
	}
	return links, nil
}

func (e *Elements) HasURL(url string) bool {
	if v, found := e.Data[url]; found {
		if len(v) > 0 {
			return true
		}
	}
	return false
}

func (e *Elements) GetMetaData(Parser *WebParser.Parser, c *Component) error {
	if e.MetaData == nil {
		e.MetaData = map[string][]map[string]*MetaData{}
	}
	for _, d := range e.Data[Parser.URL] {
		data, err := GetMetaData(d, c)
		if err != nil {
			return err
		}
		if len(data) > 0 {
			e.MetaData[Parser.URL] = append(e.MetaData[Parser.URL], data)
		}

	}
	return nil
}
func (e *Elements) GetAllMetaData() []map[string]*MetaData {
	output := []map[string]*MetaData{}
	for _, v := range e.MetaData {
		output = append(output, v...)
	}
	return output
}
func GetMetaData(htmlData *WebParser.HTMLData, c *Component) (map[string]*MetaData, error) {
	tmp := map[string]*MetaData{}

	for k, v := range c.MetaData {
		tmp[k] = &MetaData{
			Attribute: v.Attribute,
		}
		data := *htmlData
		if v.Attribute == "text" {
			if v.Flatten {
				tmp[k].Data += FlattenText(&data)
			} else {
				tmp[k].Data += data.TextData
			}

		} else {
			if v.Flatten {
				data = *(&data).Flatten()
			}
			tmp[k].Data += data.Attributes[v.Attribute]

		}
	}

	return tmp, nil
}
func FlattenText(h *WebParser.HTMLData) string {
	output := h.TextData
	for _, c := range h.Child {
		output += FlattenText(c)
	}
	return output
}
func (e *Elements) Download(parser *WebParser.Parser, c *Component, path string, job *jobs.Job, jobDAO *dao.JobDAO, book *dto.Books) error {
	if !c.IsDownload {
		return nil
	}
	job.TotalProgress += e.GetTotalElements()
	book.Pages = job.TotalProgress
	if jobDAO != nil {
		job.UpdateDB(jobDAO)
	}
	keys := make([]string, 0, len(e.Data))
	for k := range e.Data {
		keys = append(keys, k)
	}
	re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)

	sort.Slice(keys, func(i, j int) bool {
		submatchallA := re.FindAllString(keys[i], -1)
		submatchallB := re.FindAllString(keys[j], -1)

		numA, _ := strconv.ParseFloat(submatchallA[len(submatchallA)-1], 64)
		numB, _ := strconv.ParseFloat(submatchallB[len(submatchallB)-1], 64)
		return numA < numB
	})
	//sort.Strings(keys)

	for _, k := range keys {
		parser.URL = k
		err := e.DownloadElements(parser, c, path+"/", job, jobDAO, book)
		if err != nil {
			return err
		}
	}
	if jobDAO != nil && len(e.Data[parser.URL]) == 1 {
		book.JobId = job.Id
		booksDAO := dao.BooksDAO{DB: jobDAO.DB}
		err := booksDAO.NewBook(book)
		if err != nil {
			parser.Logger.Error("failed creating new book", zap.Error(err))
		}
	}
	return nil
}

func (e *Elements) GetTotalElements() int {
	total := 0
	for _, v := range e.Data {
		total += len(v)
	}
	return total
}

func (e *Elements) DownloadElements(parser *WebParser.Parser, c *Component, path string, job *jobs.Job, jobDAO *dao.JobDAO, book *dto.Books) error {
	DownloadParser := WebParser.NewParser(parser.Logger)
	page_count := 0
	if !c.IsDownload {
		return nil
	}
	NewPath := strings.ReplaceAll(path, "//", "/") + e.GetElementPath(parser)
	book.UserID = job.User
	book.FilePath = NewPath
	book.Chapter = e.GetElementData(parser, "chapter")
	book.Volume = e.GetElementData(parser, "volume")
	book.Description = e.GetElementData(parser, "description")

	if v, found := e.Data[parser.URL]; found {
		for _, d := range v {
			job.CurrentProgress += 1
			job.UpdateTime()
			if jobDAO != nil {
				job.UpdateDB(jobDAO)
			}
			link, err := d.GetLink(c.LinkAttributes, parser)
			if err != nil {
				continue
			}
			urlParse := strings.Split(link, ".")
			ext := urlParse[len(urlParse)-1]
			name := d.Attributes[c.Name]
			if name == "" {
				tmp := urlParse[len(urlParse)-2]
				tmpAr := strings.Split(tmp, "/")
				name = tmpAr[len(tmpAr)-1]
			}
			page_count += 1
			filename := fmt.Sprintf("%s%s.%s", NewPath, name, ext)
			filenamePNG := fmt.Sprintf("%s%s.png", NewPath, name)
			filenameJPG := fmt.Sprintf("%s%s.jpg", NewPath, name)
			parser.Logger.Debug(fmt.Sprintf("download image %s : %s  %d out of %d", link, NewPath, page_count, len(v)))
			if _, err := os.Stat(filename); err == nil {
				_ = setFileModTime(filename)
				if book.CoverImage == "" {
					book.CoverImage = filename
				}
				continue
			}
			if _, err := os.Stat(filenamePNG); err == nil {
				_ = setFileModTime(filenamePNG)
				if book.CoverImage == "" {
					book.CoverImage = filenamePNG
				}
				continue
			}
			if _, err := os.Stat(filenameJPG); err == nil {
				_ = setFileModTime(filenameJPG)
				if book.CoverImage == "" {
					book.CoverImage = filenameJPG
				}
				continue
			}
			if book.CoverImage == "" {
				book.CoverImage = filename
			}
			_ = os.MkdirAll(NewPath, os.ModePerm)
			time.Sleep(time.Duration(c.Delay) * time.Second)
			_ = DownloadParser.Download(link, filename)
			if ext == "webp" {
				parser.Logger.Debug(fmt.Sprintf("converting %s to png", filename))
				err = converter.ConvertWebPToPng(filename, filenamePNG)
				if err != nil {
					parser.Logger.Error("failed converting webp to png", zap.Error(err))
					continue
				}

				_ = os.Remove(filename)
			}

		}
	}
	if jobDAO != nil && len(e.Data[parser.URL]) > 1 {
		booksDAO := dao.BooksDAO{DB: jobDAO.DB}
		book.Name = "_"
		err := booksDAO.NewBook(book)
		if err != nil {
			parser.Logger.Error("failed creating new book", zap.Error(err))
		}

	}
	return nil
}
func (e *Elements) GetElementPath(parser *WebParser.Parser) string {
	if metaList, found := e.MetaData[parser.URL]; found {
		for _, m := range metaList {
			if v, found := m["volume"]; found {
				return strings.ReplaceAll(v.Data, "/", "") + "/"
			}
			if v, found := m["chapter"]; found {
				return strings.ReplaceAll(v.Data, "/", "") + "/"
			}
		}
	}
	return ""
}

func (e *Elements) GetElementData(parser *WebParser.Parser, key string) string {
	if metaList, found := e.MetaData[parser.URL]; found {
		for _, m := range metaList {
			if v, found := m[key]; found {
				return v.Data
			}
		}
	}
	return ""
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
