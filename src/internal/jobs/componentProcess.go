package jobs

import (
	"MangaLibrary/src/internal/dao"
	"MangaLibrary/src/internal/dto"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/TheBlockNinja/WebParser"
)

type Component struct {
	ID             int                  `json:"id"`
	Name           string               `db:"name"  json:"name"`
	Tag            string               `db:"tag"  json:"tag"`
	Attribute      string               `db:"attribute"  json:"attribute"`
	Value          string               `db:"value"  json:"value"`
	IsLink         bool                 `db:"is_link"  json:"is_link"`
	IsDownload     bool                 `db:"is_download"  json:"is_download"`
	LinkAttributes []string             `db:"link_attributes"  json:"link_attributes"`
	ElementData    *Elements            `db:"element_data"  json:"element_data"`
	Children       []*Component         `json:"children"`
	Siblings       []*Component         `json:"siblings"`
	Delay          int                  `db:"delay"  json:"delay"`
	MetaData       map[string]*MetaData `db:"meta_data"  json:"meta_data"`
	Reverse        bool                 `db:"reverse" json:"reverse"`
}

type MetaData struct {
	Attribute string `json:"attribute"`
	Data      string `json:"data"`
	Flatten   bool   `json:"flatten"`
}

func (c *Component) FindID(id int) *Component {
	if c.ID == id {
		return c
	}
	for _, child := range c.Children {
		tmp := child.FindID(id)
		if tmp != nil {
			return tmp
		}
	}
	for _, sibling := range c.Siblings {
		tmp := sibling.FindID(id)
		if tmp != nil {
			return tmp
		}
	}
	return nil
}

func (c *Component) LoadSiteData(parser *WebParser.Parser, job *Job, newURL string, jobDAO *dao.JobDAO) error {
	var err error
	/*
		If Data already exists
	*/

	/*
		Get Data
	*/
	if job.ForceStop {
		return fmt.Errorf("force stop")
	}
	didLoad := true
	if !c.ElementData.HasURL(parser.URL) {
		if newURL != "" {
			parser, err = LoadParserData(parser, newURL)
			if err != nil {
				return err
			}
		} else {
			parser, err = LoadParserData(parser, job.Ctx.URL)
			if err != nil {
				return err
			}
		}
		time.Sleep(time.Duration(c.Delay) * time.Second)
	} else {
		didLoad = false
	}
	err = c.ElementData.LoadData(parser, c)
	if err != nil {
		return err
	}
	err = c.ElementData.GetMetaData(parser, c)
	if err != nil {
		return err
	}
	if !didLoad || job.TotalProgress == 0 {
		job.TotalProgress += c.ElementData.GetLength(parser)
		fmt.Printf("JOB(%d/%d)\n", job.CurrentProgress, job.TotalProgress)

		job.UpdateDB(jobDAO)

	}

	if len(c.Siblings) > 0 {
		for _, siblings := range c.Siblings {
			if job.ForceStop {
				return fmt.Errorf("force stop")
			}
			err = siblings.LoadSiteData(parser, job, "", jobDAO)
			if err != nil {
				return err
			}
		}
		metadata := c.GetAllMetaData()
		var value string
		value = FindInMetaData(metadata, "name")
		if value != "" {
			job.Name = value
			parser.Logger.Info("updating job name..")
			job.UpdateDB(jobDAO)
		}
	}
	if len(c.Children) > 0 {
		for _, child := range c.Children {
			if job.ForceStop {
				return fmt.Errorf("force stop")
			}
			if child.IsLink {
				/*
					Seperate links to only get one at a time
					Then add the meta data from the c.element to it
				*/
				for _, sites := range c.ElementData.Data[parser.URL] {
					link, err := sites.GetLink(c.LinkAttributes, parser)
					if err != nil {
						return err
					}
					childParser := &WebParser.Parser{}
					if !c.ElementData.HasURL(link) {
						childParser.URL = link
						err = child.LoadSiteData(childParser, job, link, jobDAO)
						if err != nil {
							return err
						}
						if child.ElementData.MetaData[link] == nil {
							child.ElementData.MetaData[link] = []map[string]*MetaData{}
						}
						parentMetaData := map[string]*MetaData{} //c.MetaData
						for k, v := range c.MetaData {
							parentMetaData[k] = &MetaData{
								Attribute: v.Attribute,
								Data:      "",
								Flatten:   v.Flatten,
							}
							data := *sites
							if v.Attribute == "text" {
								if v.Flatten {
									parentMetaData[k].Data += FlattenText(&data)
								} else {
									parentMetaData[k].Data += data.TextData
								}

							} else {
								if v.Flatten {
									data = *(&data).Flatten()
								}
								parentMetaData[k].Data += data.Attributes[v.Attribute]
							}
						}
						// todo create a copy of the parent meta data and add it to the child
						//c.ElementData.Data[parser.URL][i].TextData
						child.ElementData.MetaData[link] = append(child.ElementData.MetaData[link], parentMetaData)
						job.CurrentProgress += 1

						job.UpdateTime()
						job.UpdateDB(jobDAO)

					}

				}
			} else {
				/*
					Add support for child object without link data
				*/
				//childParser := &WebParser.Parser{}
				//for i, d := range c.ElementData.Data[parser.URL] {
				//	childParser.Html = d
				//	childParser.URL = parser.URL + fmt.Sprintf("-%d", i)
				//	err = child.LoadSiteData(childParser, job, parser.URL+fmt.Sprintf("-%d", i))
				//	if err != nil {
				//		return err
				//	}
				//}

			}
		}

	} else {
		if !didLoad || job.TotalProgress == 0 {
			job.CurrentProgress += c.ElementData.GetLength(parser)
		}
	}

	return nil
}

func (c *Component) Download(parser *WebParser.Parser, basePath string, job *Job, jobDAO *dao.JobDAO, book *dto.Books) error {
	metadata := c.GetAllMetaData()
	var value string
	value = FindInMetaData(metadata, "name")
	if value != "" {
		book.Name = value
		d, err := json.Marshal(metadata)
		if err != nil {
			parser.Logger.Error("failed getting all meta data", zap.Error(err))
		}
		book.Metadata = string(d)
		book.Description = FindInMetaData(metadata, "description")
	}
	//if value == "" {
	//	value = FindInMetaData(metadata, "volume")
	//}
	//if value == "" {
	//	value += FindInMetaData(metadata, "chapter")
	//}
	basePath += "/" + value
	err := c.ElementData.Download(parser, c, basePath, job, jobDAO, book)
	if err != nil {
		fmt.Printf("error downloading...%s...basepath %s\n", err.Error(), basePath)
		return err
	}

	for _, child := range c.Children {
		err = child.Download(parser, basePath, job, jobDAO, book)
		if err != nil {
			return err
		}
	}
	for _, child := range c.Siblings {
		err = child.Download(parser, basePath, job, jobDAO, book)
		if err != nil {
			return err
		}
	}
	return nil
}

func FindInMetaData(metadata []map[string]*MetaData, value string) string {
	shortest := ""
	for _, d := range metadata {
		if v, found := d[value]; found {
			name := v.Data
			if name != "" && (len(name) < len(shortest) || shortest == "") {
				shortest = name
			}

		}
	}
	return shortest
}

func (c *Component) GetAllSiblingMetaData() []map[string]*MetaData {
	output := []map[string]*MetaData{}
	for _, sibling := range c.Siblings {
		output = append(output, sibling.ElementData.GetAllMetaData()...)
	}
	return output
}
func (c *Component) GetAllMetaData() []map[string]*MetaData {
	output := c.ElementData.GetAllMetaData()
	bytes, _ := json.Marshal(output)
	fmt.Printf("META DATA:%s\n", string(bytes))
	for _, sibling := range c.Siblings {
		output = append(output, sibling.GetAllMetaData()...)
	}
	for _, child := range c.Children {
		output = append(output, child.GetAllMetaData()...)
	}
	return output
}

func LoadParserData(parser *WebParser.Parser, currentURL string) (*WebParser.Parser, error) {
	if parser.Html == nil {
		parser.URL = ""
		err := parser.Get(currentURL)
		if err != nil {
			return parser, err
		}
		parser.URL = currentURL
	}
	return parser, nil
}
