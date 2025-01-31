package gpt

import (
	"backend/conf"
	"encoding/xml"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"github.com/imroc/req"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestBody struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
}
type ResponseBody struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

type SlidesXML struct {
	XMLName  xml.Name     `xml:"slides"`
	Sections []SectionXML `xml:"section"`
}

type SectionXML struct {
	XMLName xml.Name `xml:"section"`
	Class   string   `xml:"class,attr"`
	Content []string `xml:"p"`
}

func ErrorScanner(gptResponse string, genXMLType interface{}) (string, error) {
	// 确保程序按照genXMLType的类型进行解析
	value := reflect.New(reflect.TypeOf(genXMLType)).Interface()
	err := xml.Unmarshal([]byte(gptResponse), value)
	if err != nil {
		return "", fmt.Errorf("genXMLType格式与gptResponse不匹配: %s", err.Error())
	}

	// 移除所有转义字符
	r := regexp.MustCompile(`\\.`)
	gptResponse = r.ReplaceAllString(gptResponse, "")

	// Remove all newline characters from the response
	gptResponse = strings.ReplaceAll(gptResponse, "\n", "")

	// 移除所有不必要空格，"<>"标签内内容除外
	inScope := false
	var builder strings.Builder

	for _, ch := range gptResponse {
		// If we encounter a '<', we are inside a tag
		if ch == '<' {
			inScope = true
			builder.WriteRune(ch)
			continue
		}
		// If we encounter a '>', we are no longer inside a tag
		if ch == '>' {
			inScope = false
			builder.WriteRune(ch)
			continue
		}
		// If we are outside a tag, remove unnecessary spaces
		if !inScope {
			if unicode.IsSpace(ch) {
				// Only add a single space between non-tag content (outside tags)
				if builder.Len() > 0 && builder.String()[builder.Len()-1] != ' ' {
					builder.WriteRune(' ')
				}
				continue
			}
		}
		// Otherwise, write the character (preserving spaces within tags)
		builder.WriteRune(ch)
	}

	return builder.String(), nil
}

func RequestGpt(prompt string, genXmlType interface{}) (string, error) {
	var apikey string

	// 当GetApiKey()得到合法的apikey时，才进行下一步
	for apikey == "" {
		apikey, _ = GetApiKey()
	}
	defer ReleaseApiKey(apikey)

	var body RequestBody
	body.Model = conf.GetGptModel()
	body.Messages = append(body.Messages, Message{
		Role:    "user",
		Content: prompt,
	})
	body.Temperature = 0.7
	fmt.Println("Prompt", prompt)

	req_url := conf.GetGptApiUrl()

	// 最多尝试2次
	retryCount := 0

	for retryCount < 2 {
		// 进行http请求
		resp, err := req.Post(req_url, req.BodyJSON(&body), req.Header(map[string]string{
			"Authorization": "Bearer " + apikey,
			"Content-Type":  "application/json",
		}))
		if err != nil {
			return "", err
		}
		if resp.Response().StatusCode != 200 {
			return "", fmt.Errorf("GPT请求失败，状态码为%d", resp.Response().StatusCode)
		}

		var res ResponseBody
		resp.ToJSON(&res)

		// 扫描错误
		result, err := ErrorScanner(res.Choices[0].Message.Content, genXmlType)
		if err != nil {
			retryCount++
			continue
		}

		// 成功
		fmt.Println("GPT response\n", res.Choices[0].Message.Content)
		fmt.Println("GPT response XML\n", result)
		return result, nil
	}

	// 3次尝试均失败
	return "", fmt.Errorf("all retries failed")
}

func GuideContentSection(outline string) (string, error) {
	template := conf.GetGuideSinglePromptTemplate()
	template = strings.ReplaceAll(template, "{{outline}}", outline)

	guide_slide, err := RequestGpt(template, SectionXML{}) // <section></section>
	if err != nil {
		return "", err
	}

	guide_slide = strings.ReplaceAll(guide_slide, "\n", "")
	return guide_slide, nil
}
