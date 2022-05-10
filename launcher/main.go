package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"net/url"
	"omega_launcher/embed_binary"
	"omega_launcher/utils"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pterm/pterm"
)

var STOARGE_REPO = "http://124.222.6.29:6000"

type BotConfig struct {
	Code   string `json:"租赁服号"`
	Passwd string `json:"租赁服密码"`
}

func main() {
	targetHash := GetRemoteOmegaHash()
	currentHash := GetCurrentOmegaHash()
	// fmt.Println(targetHash)
	// fmt.Println(currentHash)
	if targetHash == currentHash {
		pterm.Success.Println("太好了，你的 omega 已经是最新的了!")
	} else {
		pterm.Warning.Println("我们将为你下载最新 omega, 请保持耐心...")
		DownloadOmega()
	}
	if err := os.Chdir(GetCurrentDir()); err != nil {
		panic(err)
	}
	if utils.IsDir(path.Join(GetCurrentDir(), "omega_storage")) {
		CQHttpEnablerHelper()
	}
	StartOmegaHelper()
}

//go:embed config.yml
var defaultConfigBytes []byte

//go:embed 组件-群服互通-1.json
var defaultQGroupLinkConfigByte []byte

func CQHttpEnablerHelper() {
	pterm.Info.Printf("要启用群服互通吗 要请输入 y 不要请输入 n ")
	accept := utils.GetInputYN()
	if !accept {
		return
	}
	if err := utils.WriteFileData(GetCqHttpExec(), GetCqHttpBinary()); err != nil {
		panic(err)
	}
	configFile := path.Join(GetCurrentDir(), "config.yml")
	omegaConfigFile := path.Join(GetCurrentDir(), "omega_storage", "配置", "组件-群服互通-1.json")
	if !utils.IsFile(configFile) {
		pterm.Info.Printf("请输入QQ账号: ")
		Code := utils.GetValidInput()
		pterm.Info.Printf("请输入QQ密码（想扫码登录则留空）: ")
		Passwd := utils.GetInput()
		if Passwd == "" {
			Passwd = "''"
		}
		defaultConfigStr := string(defaultConfigBytes)
		cfgStr := strings.ReplaceAll(defaultConfigStr, "[QQ账号]", Code)
		cfgStr = strings.ReplaceAll(cfgStr, "[QQ密码]", Passwd)
		utils.WriteFileData(configFile, []byte(cfgStr))
		pterm.Info.Printf("请输入想链接的群号: ")
		GroupCode := utils.GetValidInput()
		groupCfgStr := strings.ReplaceAll(string(defaultQGroupLinkConfigByte), "[群号]", GroupCode)
		utils.WriteFileData(omegaConfigFile, []byte(groupCfgStr))
	}
	pterm.Warning.Println("将使用 " + configFile + " 的配置进行 QQ 登录，您可以自行修改这份文件")
	pterm.Warning.Println("将使用 " + omegaConfigFile + " 的配置进行群服互通，您可以自行修改这份文件")
	RunCQHttp()
}

func WaitConnect() {
	for {
		u := url.URL{Scheme: "ws", Host: "127.0.0.1:6700"}
		var err error
		_, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			time.Sleep(1)
			continue
		} else {
			return
		}
	}
}

func RunCQHttp() {
	cmd := exec.Command(GetCqHttpExec())
	cqHttpOut, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	go func() {
		reader := bufio.NewReader(cqHttpOut)
		for {
			readString, err := reader.ReadString('\n')
			if err != nil || err == io.EOF {
				fmt.Println("reader exit")
				return
			}
			fmt.Print("[CQHTTP] " + readString)
		}
	}()
	go func() {
		err = cmd.Start()
		if err != nil {
			fmt.Println(err)
		}
		err = cmd.Wait()
		if err != nil {
			fmt.Println(err)
		}
	}()
	WaitConnect()
	pterm.Success.Println("CQ-Http已经成功启动了！")
}

func StartOmegaHelper() {
	pterm.Success.Println("开始配置Omega登录")
	botConfig := &BotConfig{}
	reconfigFlag := true
	if err := utils.GetJsonData(path.Join(GetCurrentDir(), "服务器登录配置.json"), botConfig); err == nil && botConfig.Code != "" {
		pwd := " 密码为空"
		if botConfig.Passwd != "" {
			pwd = " 密码为: " + botConfig.Passwd
		}
		pterm.Info.Println("租赁服账号为: " + botConfig.Code + pwd)
		pterm.Info.Printf("接受这个登录配置请输入 y ,需要修改请输入 n: ")
		accept := utils.GetInputYN()
		if accept {
			reconfigFlag = false
		}
	}
	if reconfigFlag {
		pterm.Info.Printf("请输入租赁服账号: ")
		botConfig.Code = utils.GetValidInput()
		pterm.Info.Printf("请输入租赁服密码（没有则留空）: ")
		botConfig.Passwd = utils.GetInput()
	}
	if err := utils.WriteJsonData(path.Join(GetCurrentDir(), "服务器登录配置.json"), botConfig); err != nil {
		pterm.Error.Println("无法记录租赁服配置，不过可能不是什么大问题")
	}
	RunOmega(botConfig)
}

func RunOmega(cfg *BotConfig) {
	args := []string{"-M", "-O", "--no-update-check", "-c", cfg.Code}
	if cfg.Passwd != "" {
		args = append(args, "-p")
		args = append(args, cfg.Passwd)
	}
	readC := make(chan string)
	go func() {
		for {
			s := utils.GetInput()
			readC <- s
		}
	}()
	t := time.NewTicker(10 * time.Second)
	for {
		cmd := exec.Command(GetOmegaExecName(), args...)
		omega_out, err := cmd.StdoutPipe()
		if err != nil {
			panic(err)
		}
		omega_in, err := cmd.StdinPipe()
		if err != nil {
			panic(err)
		}
		pterm.Success.Println("如果Omega崩溃了，它会在最长 30 秒后自动重启")

		stopped := false
		go func() {
			reader := bufio.NewReader(omega_out)
			for {
				readString, err := reader.ReadString('\n')
				if err != nil || err == io.EOF {
					fmt.Println("reader exit")
					return
				}
				fmt.Print(readString)
			}
		}()

		go func() {
			for {
				s := <-readC
				if stopped {
					readC <- s
					return
				}
				omega_in.Write([]byte(s + "\n"))
			}
		}()

		err = cmd.Start()
		if err != nil {
			fmt.Println(err)
		}
		err = cmd.Wait()
		if err != nil {
			fmt.Println(err)
		}
		stopped = true
		pterm.Warning.Println("Omega将在最长 30 秒后自动重启")
		time.Sleep(10)
		<-t.C
	}
}

func GetCqHttpBinary() []byte {
	return embed_binary.GetCqHttpBinary()
}

func GetCurrentDir() string {
	pathExecutable, err := os.Executable()
	if err != nil {
		panic(err)
	}
	dirPathExecutable := filepath.Dir(pathExecutable)
	return dirPathExecutable
}

func GetOmegaExecName() string {
	omega := "fastbuilder"
	if GetPlantform() == embed_binary.WINDOWS_x86_64 {
		omega = "fastbuilder.exe"
	}
	omega = path.Join(GetCurrentDir(), omega)
	p, err := filepath.Abs(omega)
	if err != nil {
		panic(err)
	}
	return p
}

func GetCqHttpExec() string {
	cqhttp := "cqhttp"
	if GetPlantform() == embed_binary.WINDOWS_x86_64 {
		cqhttp = "cqhttp.exe"
	}
	cqhttp = path.Join(GetCurrentDir(), cqhttp)
	p, err := filepath.Abs(cqhttp)
	if err != nil {
		panic(err)
	}
	return p
}

func GetPlantform() string {
	return embed_binary.GetPlantform()
}

func GetRemoteOmegaHash() string {
	url := ""
	switch GetPlantform() {
	case embed_binary.WINDOWS_x86_64:
		url = STOARGE_REPO + "/fastbuilder-windows.hash"
	case embed_binary.Linux_x86_64:
		url = STOARGE_REPO + "/fastbuilder-linux.hash"
	case embed_binary.MACOS_x86_64:
		url = STOARGE_REPO + "/fastbuilder-macos.hash"
	default:
		panic("未知平台" + GetPlantform())
	}
	hashBytes := utils.DownloadSmallContent(url)
	return string(hashBytes)
}

func GetFileHash(fname string) string {
	if utils.IsFile(fname) {
		fileData, err := utils.GetFileData(fname)
		if err != nil {
			panic(err)
		}
		return utils.GetBinaryHash(fileData)
	}
	return ""
}

func GetCurrentOmegaHash() string {
	exec := GetOmegaExecName()
	return GetFileHash(exec)
}

func GetCQHttpHash() string {
	exec := GetCqHttpExec()
	return GetFileHash(exec)
}

func GetEmbeddedCQHttpHash() string {
	return utils.GetBinaryHash(GetCqHttpBinary())
}

func DownloadOmega() {
	exec := GetOmegaExecName()
	url := ""
	switch GetPlantform() {
	case embed_binary.WINDOWS_x86_64:
		url = STOARGE_REPO + "/fastbuilder-windows.exe"
	case embed_binary.Linux_x86_64:
		url = STOARGE_REPO + "/fastbuilder-linux"
	case embed_binary.MACOS_x86_64:
		url = STOARGE_REPO + "/fastbuilder-macos"
	default:
		panic("未知平台" + GetPlantform())
	}
	execBytes := utils.DownloadSmallContent(url)
	if err := utils.WriteFileData(exec, execBytes); err != nil {
		panic(err)
	}
}
