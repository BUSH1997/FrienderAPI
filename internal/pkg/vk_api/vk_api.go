package vk_api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

type VKApi struct {
	AccessToken string
	AlbumId     string
	GroupId     string
	Version     string
}

type UploadPhotoParamType string

const (
	Default          UploadPhotoParamType = "default"
	WithUploadServer UploadPhotoParamType = "with_upload_server"
)

type UploadPhotoParam struct {
	Type         UploadPhotoParamType
	Token        string
	AlbumId      string
	GroupId      string
	UploadServer string
}

const (
	file_for_upload = "C:\\Users\\ruduk\\Desktop\\События\\FrienderAPI\\123211.png"
	vk_api_url      = "https://api.vk.com/method"
)

func (vk *VKApi) UploadPhoto(file *multipart.FileHeader, param UploadPhotoParam) (string, error) {
	var err error
	uriServerUpload := ""

	if param.Type == WithUploadServer {
		uriServerUpload = param.UploadServer
	} else {
		uriServerUpload, err = vk.GetUploadServer(param)
		if err != nil {
			fmt.Println(err)
			return "", err
		}
	}

	jsonResp, err := vk.UploadPhotoOnUriServer(file, uriServerUpload)
	if err != nil {
		return "", err
	}

	jsonMap := jsonResp.(map[string]interface{})
	photos_list := jsonMap["photos_list"].(string)
	server := jsonMap["server"].(float64)
	stringFlaot := fmt.Sprintf("%v", server)
	hash := jsonMap["hash"].(string)
	idPhoto, err := vk.SaveFile(photos_list, stringFlaot, hash)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return idPhoto, nil
}

func (vk VKApi) UploadPhotoOnUriServer(file *multipart.FileHeader, uriServerUpload string) (interface{}, error) {
	b := new(bytes.Buffer)
	w := multipart.NewWriter(b)
	field, err := w.CreateFormFile("file1", file_for_upload)
	if err != nil {
		fmt.Println(err)
	}

	fileOpend, err := file.Open()
	if err != nil {
		fmt.Println(err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(fileOpend)

	field.Write(buf.Bytes())
	if err != nil {
		fmt.Println(err)
	}
	w.Close()
	resp, err := http.Post(uriServerUpload, w.FormDataContentType(), b)
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	var jsonResp interface{}
	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		return "", err
	}

	return jsonResp, nil
}

func (vk VKApi) UploadPhotosOnUriServer(files []*multipart.FileHeader, uriUploadServer string) ([]interface{}, error) {
	b := new(bytes.Buffer)
	w := multipart.NewWriter(b)
	fmt.Println(files)
	var res []interface{}

	for i := 0; i < len(files); i++ {
		field, err := w.CreateFormFile("file"+strconv.Itoa(i+1), files[i].Filename)
		if err != nil {
			fmt.Println(err)
		}

		fileOpend, err := files[i].Open()
		if err != nil {
			fmt.Println(err)
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(fileOpend)

		field.Write(buf.Bytes())
		if err != nil {
			fmt.Println(err)
		}

		if i != 0 && (i%10 == 4 || i%10 == 9) || i == len(files)-1 {
			fmt.Println(w)
			w.Close()

			resp, err := http.Post(uriUploadServer, w.FormDataContentType(), b)
			if err != nil {
				fmt.Println(err)
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				return res, err
			}

			var jsonResp interface{}
			err = json.Unmarshal(body, &jsonResp)
			if err != nil {
				return res, err
			}

			res = append(res, jsonResp)

			if i != len(files)-1 {
				b = new(bytes.Buffer)
				w = multipart.NewWriter(b)
			}
		}
	}
	fmt.Println(res)

	return res, nil
}

func (vk VKApi) GetUploadServer(param UploadPhotoParam) (string, error) {
	uri := fmt.Sprintf("%s/%s?access_token=%s&album_id=%s&group_id=%s&v=%s", vk_api_url,
		"photos.getUploadServer", vk.AccessToken, vk.AlbumId, vk.GroupId, vk.Version)

	resp, err := http.Get(uri)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	var jsonResp interface{}
	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		return "", err
	}

	jsonMap := jsonResp.(map[string]interface{})
	fieldResponse := jsonMap["response"].(map[string]interface{})
	if fieldResponse["upload_url"] != nil {
		return fieldResponse["upload_url"].(string), nil
	}

	return "", errors.New("Error getUploadServer")
}

func (vk VKApi) ReadFile() (bytes.Buffer, error) {
	f, err := os.Open(file_for_upload)
	if err != nil {
		fmt.Println(err)
		return bytes.Buffer{}, err
	}
	defer f.Close()

	wr := bytes.Buffer{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		wr.WriteString(sc.Text())
	}

	return wr, nil
}

func (vk VKApi) SaveFile(photosList string, server string, hash string) (string, error) {
	uri := fmt.Sprintf("%s/%s?access_token=%s&album_id=%s&group_id=%s&v=%s&photos_list=%s&server=%s&hash=%s", vk_api_url,
		"photos.save", vk.AccessToken, vk.AlbumId, vk.GroupId, vk.Version, photosList, server, hash)

	resp, err := http.Post(uri, "application/json", nil)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var jsonResp interface{}
	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		return "", err
	}

	jsonMap := jsonResp.(map[string]interface{})
	fieldResponse := jsonMap["response"].([]interface{})
	fieldResponse1 := fieldResponse[0].(map[string]interface{})
	id := fieldResponse1["id"].(float64)
	ownerId := fieldResponse1["owner_id"].(float64)
	stringId := fmt.Sprintf("%.0f_%.0f", ownerId, id)

	return stringId, nil
}
