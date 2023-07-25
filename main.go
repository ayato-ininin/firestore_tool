package main

import (
	"bufio"
	"context"
	"firestore_tool/fs_pkg"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"cloud.google.com/go/firestore"
)

type Setting struct {
	credentialJsonPath string
	mode               string
	docPath            string
	updateStruct []firestore.Update
}

var setting Setting

func main() {
	flag.StringVar(&setting.credentialJsonPath, "credentialJsonPath", "", "エントリーデフォルトパス")
	flag.Parse()

	ctx := context.Background()
	client, err := fs_pkg.FirebaseInit(ctx, setting.credentialJsonPath)
	if err != nil {
		log.Fatalf("FirebaseInit error: %v", err)
	}
	defer client.Close()

	intro()

	// create a channel to indicat e when the program can quit
	doneChan := make(chan bool)

	// start a gorouutin to read user input and run program
	go readUserInput(os.Stdin, doneChan, ctx, client)

	// block until the done chan gets a value
	<-doneChan

	// close the channel
	close(doneChan)

	fmt.Println("done!")

}

func intro() {
	fmt.Println("Firebaseと接続しました。")
	fmt.Println("メソッドを入力してください。")
	prompt()
}

func prompt() {
	fmt.Print("-> ")
}

func readUserInput(in io.Reader, doneChan chan bool, ctx context.Context, client *firestore.Client) {
	scanner := bufio.NewScanner(in)

	for {
		scanner.Scan()
		mode := scanner.Text()
		// invertedIndexからscanner.Text()に対応するデータ取得
		if mode == "exit" || mode == "q" {
			doneChan <- true
			return
		}
		if mode != "delete" && mode != "update" {
			fmt.Println("Invalid mode entered, exiting.")
			doneChan <- true
			return
		}
		setting.mode = mode

		fmt.Println("対象のドキュメントパスを入力してください。")
		prompt()
		scanner.Scan()
		docPath := scanner.Text()
		setting.docPath = docPath

		if mode == "update" {
			fmt.Println("対象のドキュメントキーを入力してください。")
			prompt()
			scanner.Scan()
			key := scanner.Text()

			fmt.Println("保存する値を入力してください。")
			prompt()
			scanner.Scan()
			value := scanner.Text()
			setting.updateStruct = append(setting.updateStruct, firestore.Update{Path: key, Value: value})
		}

		// ユーザーからの入力を取得した後、データを削除または更新します。
		if setting.mode == "delete" {
			// データを削除する関数を呼び出します。
			fs_pkg.Delete(ctx, client, setting.docPath)
		} else if setting.mode == "update" {
			// データを更新する関数を呼び出します。
			fs_pkg.Update(ctx, client, setting.docPath, setting.updateStruct)
		}

		doneChan <- true
		return
	}
}
