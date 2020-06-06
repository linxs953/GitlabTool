package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

var CommandClient http.Client = http.Client{Timeout: time.Second * 3}

type Conf struct {
	TOKEN           string
	ISSUEPATH       string
	GROUPS          string
	ISSUEBASE       string
	CREATEISSUEAPI  string
	GETUSERAPI      string
	GETPROJECTIDAPI string
	CICD            string
	GETALLGROUPAPI  string
	GETPROJECTSAPI  string
	PROJECT         string
}

type GroupInfo struct {
	FullPath string `json:"full_path"`
}

type GroupsList struct {
	groups []GroupInfo `GroupInfo:"full_path"`
}

type ProjectInfo struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
}

type ProjectsList struct {
	prjos []ProjectInfo `GroupInfo:"name"`
}

type CreateIssueStruct struct {
	PrivateToken string `json:"private_token"`
	ID           string `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	AssigneeIds  []int  `json:"assignee_ids"`
}

type IssueInfo struct {
	IID int64 `json:"iid"`
}

type UserInfo struct {
	ID int64 `json:"id"`
}

type UserList struct {
	Users []UserInfo `UserInfo:"id"`
}

func ReadConfig(filename string, configType string) (Conf, error) {
	var err error
	rootPath := GetConfigPath()
	filePath := rootPath + filename + "." + configType
	config, errReturn := Conf{}, Conf{}
	if !Exists(filePath) {
		log.Print("File not found")
		return errReturn, err
	}
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if err != nil {
		log.Error().Err(err).Msg("Open File error")
		return errReturn, err
	}
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error().Err(err).Msg("Read File Content Failed")
		return errReturn, err
	}
	err = json.Unmarshal(fileBytes, &config)
	if err != nil {
		log.Error().Err(err).Msg("Parse Config Failed")
		return errReturn, err
	}
	return config, err
}

func WriteConfig(key string, value string, filename string) {
	configRootPath := GetConfigPath()
	if configRootPath == "" {
		log.Print("Get Config File Root Path Error")
		return
	}
	filePath := fmt.Sprintf("%s/%s", configRootPath, filename)
	config, err := ReadConfig("config", "json")
	if err != nil {
		log.Error().Err(err).Msg("获取配置失败")
		return
	}
	switch {
	case strings.ToUpper(key) == "TOKEN":
		config.TOKEN = value
	case strings.ToUpper(key) == "ISSUEPATH":
		config.ISSUEPATH = value
	case strings.ToUpper(key) == "GROUPS":
		config.GROUPS = value
	case strings.ToUpper(key) == "ISSUEBASE":
		config.ISSUEBASE = value
	case strings.ToUpper(key) == "CREATEISSUEAPI":
		config.CREATEISSUEAPI = value
	case strings.ToUpper(key) == "GETUSERAPI":
		config.GETUSERAPI = value
	case strings.ToUpper(key) == "GETPROJECTIDAPI":
		config.GETPROJECTIDAPI = value
	case strings.ToUpper(key) == "CICD":
		config.CICD = value
	case strings.ToUpper(key) == "GETALLGROUPAPI":
		config.GETALLGROUPAPI = value
	case strings.ToUpper(key) == "GETPROJECTSAPI":
		config.GETPROJECTSAPI = value
	case strings.ToUpper(key) == "PROJECT":
		config.PROJECT = value
	default:
		log.Printf("Can not support arg %s ", strings.ToUpper(key))
		return
	}
	fileContent, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		log.Error().Err(err).Msg("Struct to Byte error")
		return
	}
	err = ioutil.WriteFile(filePath, fileContent, 0644)
	if err != nil {
		log.Error().Err(err).Msg("Write to config error")
		return
	}
}

func Init() {
	homeEnv := os.Getenv("HOME")
	if homeEnv == "" {
		fmt.Println("Get $HOME ENV nil")
	}
	configDir := homeEnv + "/.mt"
	if !Exists(configDir) {
		err := os.Mkdir(configDir, os.ModePerm)
		if err != nil {
			log.Error().Err(err).Msg("创建目录失败")
			return
		}
		err = os.Mkdir(configDir+"/IssueTemplate", os.ModePerm)
		if err != nil {
			log.Error().Err(err).Msg("创建目录失败")
			return
		}
		configFilePath := homeEnv + "/.mt/config.json"
		if !Exists(configFilePath) {
			newJwtConfigJson()
		}
	}
}

func GetConfigPath() string {
	homeEnv := os.Getenv("HOME")
	if homeEnv == "" {
		fmt.Println("Get $HOME ENV nil")
		return ""
	}
	path := homeEnv + "/.mt/"
	return path
}

func newJwtConfigJson() {
	fileContent := `{
    "TOKEN": "",
    "ISSUEPATH": "/IssueTemplate/",
    "GROUPS": "",
    "ISSUEBASE": "",
    "CREATEISSUEAPI": "",
    "GETUSERAPI": "",
    "GETPROJECTIDAPI": "",
    "CICD": "",
    "GETALLGROUPAPI": "",
    "GETPROJECTSAPI": "",
    "PROJECT": ""
	}`
	path := GetConfigPath()
	if path == "" {
		return
	}
	path = path + "config.json"
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		log.Printf("Occur error %s", err.Error())
		return
	}
	_, err = file.Write([]byte(fileContent))
	if err != nil {
		log.Error().Err(err).Msg("写入文件出错")
	}
}

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func ListAllGroups(config Conf) {
	token, getAllGroupAPI := config.TOKEN, config.GETALLGROUPAPI
	if token == "" || getAllGroupAPI == "" {
		log.Print("get config key nil")
		return
	}
	getGroupsURL := fmt.Sprintf(getAllGroupAPI, token, 1000000)
	resp, err := CommandClient.Get(getGroupsURL)
	if err != nil {
		log.Error().Err(err).Msg("Invaild response")
		return
	}
	if resp.StatusCode != 200 {
		log.Print("Unexpected StatusCode")
		return
	}
	groupsList := GroupsList{}
	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Convert response to bytes error ")
		return
	}
	err = json.Unmarshal(respByte, &groupsList.groups)
	if err != nil {
		log.Error().Err(err).Msg("parse json error")
	}
	for _, value := range groupsList.groups {
		fmt.Println(value.FullPath)
	}
}

func ListAllProjects(config Conf, group string) {
	token, getAllProjectAPI := config.TOKEN, config.GETPROJECTSAPI
	if token == "" || getAllProjectAPI == "" {
		log.Print("Get config key nil")
		return
	}
	group = url.QueryEscape(group)
	getAllProjectURL := fmt.Sprintf(getAllProjectAPI, group, token)
	resp, err := CommandClient.Get(getAllProjectURL)
	if err != nil {
		log.Error().Err(err).Msg("Invaild Response")
		return
	}
	if resp.StatusCode != 200 {
		log.Printf("Unexpected Status Code %s", resp.StatusCode)
		return
	}
	projectList := ProjectsList{}
	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Parse Resp body error")
		return
	}
	err = json.Unmarshal(respByte, &projectList.prjos)
	if err != nil {
		log.Error().Err(err).Msg("Parse json error")
		return
	}
	for _, value := range projectList.prjos {
		fmt.Println(value.Name)
	}
}

func GetProjectID(config Conf) int64 {
	getProjectIDAPI, group, project, token := config.GETPROJECTIDAPI, config.GROUPS, config.PROJECT, config.TOKEN
	if getProjectIDAPI == "" || group == "" || project == "" || token == "" {
		log.Print("Get Config Key Error")
		return -1
	}
	projectPath := url.QueryEscape(group + "/" + project)
	getProjectIDURL := fmt.Sprintf(getProjectIDAPI, projectPath, token)
	resp, err := CommandClient.Get(getProjectIDURL)
	if err != nil {
		log.Error().Err(err).Msg("Invaild Response")
		return -1
	}
	resContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Get Resp Content Error")
		return -1
	}
	projectInfo := ProjectInfo{}
	err = json.Unmarshal(resContent, &projectInfo)
	if err != nil {
		log.Error().Err(err).Msg("Parse json error")
		return -1
	}
	return projectInfo.ID
}

func GetDescription(config Conf, filename string) string {
	configRootPath, dirname := GetConfigPath(), config.ISSUEPATH
	if configRootPath == "" || dirname == "" {
		log.Print("Get Config Key nil ")
		return ""
	}
	dirname = strings.Replace(dirname, "/", "", -1)
	markdownFilePath := fmt.Sprintf("%s%s/%s.md", configRootPath, dirname, filename)
	fileContent, err := openFile(markdownFilePath)
	if err != nil {
		log.Error().Err(err).Msg("Open File Error")
		return ""
	}
	return string(fileContent)
}

func openFile(filepath string) ([]byte, error) {
	fileByte, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Error().Err(err).Msg("Open File Error")
		return []byte(""), err
	}
	return fileByte, err
}

func GetAssigneeID(config Conf, username string) int64 {
	token, getUserIDAPI := config.TOKEN, config.GETUSERAPI
	if token == "" || getUserIDAPI == "" {
		log.Print("Get Config Key Error")
		return -1
	}
	getUserIDURL := fmt.Sprintf(getUserIDAPI, username, token)
	resp, err := CommandClient.Get(getUserIDURL)
	if err != nil {
		log.Error().Err(err).Msg("Invaild Response")
		return -1
	}
	if resp.StatusCode != 200 {
		log.Printf("Unexpected Status Code %s", resp.StatusCode)
		return -1
	}
	respContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Get Resp Content Error")
		return -1
	}
	userList := UserList{}
	err = json.Unmarshal(respContent, &userList.Users)
	if err != nil {
		log.Error().Err(err).Msg("Parse Json Error")
		return -1
	}
	return userList.Users[0].ID
}

func GetToken(config Conf) string {
	token := config.TOKEN
	if token == "" {
		log.Print("Get Config Key TOKEN nil")
		return ""
	}
	return token
}

func NewIssue(token string, projectId string, title string, description string, assigneeId string, config Conf) (string, error) {
	var err error
	var createIssueStruct CreateIssueStruct
	createIssueStruct.PrivateToken = token
	createIssueStruct.ID = projectId
	createIssueStruct.Title = title
	createIssueStruct.Description = description
	assigneeIds := make([]int, 1)
	num, err := strconv.Atoi(assigneeId)
	if err != nil {
		log.Error().Err(err).Str("assigneeId", assigneeId).Msg("createIssue Atoi error")
		return "", err
	}
	assigneeIds[0] = num
	createIssueStruct.AssigneeIds = assigneeIds

	data, err := json.Marshal(createIssueStruct)
	if err != nil {
		log.Error().Err(err).Str("assigneeId", assigneeId).Msg("createIssue Marshal error")
		return "", err
	}
	createIssueURL := config.CREATEISSUEAPI
	if createIssueURL == "" {
		log.Print("Get Config Key Error")
		return "", err
	}
	url := fmt.Sprintf(createIssueURL, projectId)
	res, err := CommandClient.Post(url, "application/json", strings.NewReader(string(data)))
	if err != nil {
		log.Error().Err(err).Str("assigneeId", assigneeId).Msg("createIssue Post error")
		return "", err
	}
	resByte, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error().Err(err).Msg("Get Resp content error")
		return "", err
	}
	issueInfo := IssueInfo{}
	err = json.Unmarshal(resByte, &issueInfo)
	if err != nil {
		log.Error().Err(err).Msg("Parse byte to json error")
		return "", err
	}
	if res != nil {
		defer res.Body.Close()
	}
	issueURLBase, group, project := config.ISSUEBASE, config.GROUPS, config.PROJECT
	if issueURLBase == "" || group == "" || project == "" {
		log.Print("Get Config Key Error")
		return "", err
	}
	issueURL := fmt.Sprintf(config.ISSUEBASE, config.GROUPS, config.PROJECT, strconv.FormatInt(issueInfo.IID, 10))
	openURL(issueURL)
	return strconv.FormatInt(issueInfo.IID, 10), err
}

func openURL(url string) {
	openCmd := exec.Command("open", url)
	err := openCmd.Run()
	if err != nil {
		log.Error().Err(err).Msg("Open URL error")
		return
	}
	log.Print("Already Open in Browser")
}
