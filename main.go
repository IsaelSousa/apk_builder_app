package main

import (
	"fmt"
	"io"
	"io/fs"
	"strings"
	"time"

	"log"
	"os"
	"os/exec"

	"github.com/common-nighthawk/go-figure"
	"github.com/manifoldco/promptui"

	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
)

func main() {
	clearTerminal()
	startScreen(information())
	promptRunner()
}

func startScreen(version InfJson) {
	fmt.Println("===============================================================================")
	fmt.Println("===============================================================================")
	myFigure := figure.NewFigure("APK-Builder", "", true)
	myFigure.Print()
	fmt.Println("===============================================================================")
	fmt.Println("===============================================================================")
	fmt.Println("Hello dear! Here you can generate yours bundle of Android.")
	fmt.Println("Version: " + version.Version)
	fmt.Println("Author: " + version.Author)
	fmt.Println(" ")
}

func promptRunner() {
	menu := []string{"Generate APK (Release)", "Generate APK (Debug)", "Clean", "Exit"}

	prompt := promptui.Select{
		Label: "Select an option!",
		Items: menu,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	switch i {
	case 0:
		gradlewAssemble("apk-release", "assembleRelease")
	case 1:
		gradlewAssemble("apk-debug", "assembleDebug")
	case 2:
		gradlewClean()
	case 3:
		os.Exit(0)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
}

func clearTerminal() {
	fmt.Print("\033[H\033[2J")
}

func copyFile(source, destiny string) error {
	arquivoOrigem, err := os.Open(source)
	if err != nil {
		return err
	}
	defer arquivoOrigem.Close()

	arquivoDestino, err := os.Create(destiny)
	if err != nil {
		return err
	}
	defer arquivoDestino.Close()

	_, err = io.Copy(arquivoDestino, arquivoOrigem)
	if err != nil {
		return err
	}

	return arquivoDestino.Sync()
}

func countdown() {
	duration := 15
	for i := duration; i > 0; i-- {
		fmt.Printf("\rEnding... %2d seconds", i)
		time.Sleep(1 * time.Second)
	}
	clearTerminal()
}

func gradlewAssemble(apkName string, apkType string) {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	projectDir := pwd + "\\android"

	var gradlew string
	if runtime.GOOS == "windows" {
		gradlew = filepath.Join(projectDir, "gradlew.bat")
	} else {
		gradlew = filepath.Join(projectDir, "gradlew")
	}

	if _, err := os.Stat(gradlew); os.IsNotExist(err) {
		log.Fatalf("Gradle file not found at %s", gradlew)
	}

	clearTerminal()

	cmd := exec.Command(gradlew, apkType)
	cmd.Dir = projectDir

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Start Build Android...")
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Fail to execute gradlew: %v", err)
	}

	getApk := projectDir + "\\app\\build\\outputs\\apk\\release"

	err = filepath.Walk(getApk, func(path string, info fs.FileInfo, err error) error {

		if strings.Contains(path, ".apk") {
			createDir := os.MkdirAll(filepath.Join(pwd, "apk"), 0755)

			if createDir == nil {
				errCopyFile := copyFile(path, pwd+"\\apk\\"+apkName+".apk")
				if errCopyFile != nil {
					fmt.Println("Error to copy file!")
					fmt.Println(errCopyFile)
				}
			}

		}

		return err
	})

	if err != nil {
		fmt.Println(err)
	}

	countdown()
	startScreen(information())
	promptRunner()
}

func gradlewClean() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	projectDir := pwd + "\\android"

	var gradlew string
	if runtime.GOOS == "windows" {
		gradlew = filepath.Join(projectDir, "gradlew.bat")
	} else {
		gradlew = filepath.Join(projectDir, "gradlew")
	}

	if _, err := os.Stat(gradlew); os.IsNotExist(err) {
		log.Fatalf("Gradle file not found at %s", gradlew)
	}

	cmd := exec.Command(gradlew, "clean")
	cmd.Dir = projectDir

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Start Clean Android...")
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Fail to execute gradlew: %v", err)
	}
}

type InfJson struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Author  string `json:"author"`
}

func information() InfJson {
	var inf InfJson

	inf.Name = "APK Android Builder CLI"
	inf.Version = "1.0.1"
	inf.Author = "PajeDeath"

	return inf
}
