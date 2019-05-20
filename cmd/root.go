package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/kitagry/genshijin"
	"github.com/spf13/cobra"
	"github.com/tcnksm/go-gitconfig"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

var message string

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "git-genshijin",
		Short: "ゲンシジン オマエ カワリ コミット スル",
		Run:   GenshiCommit,
	}

	rootCmd.Flags().StringVarP(&message, "message", "m", "", "コミット メッセージ")
	rootCmd.MarkFlagRequired("message")

	return rootCmd
}

func GenshiCommit(cmd *cobra.Command, args []string) {
	r, err := git.PlainOpen(".")
	if err != nil {
		fmt.Println(err)
		return
	}

	w, err := r.Worktree()
	if err != nil {
		fmt.Println(err)
		return
	}

	status, err := w.Status()
	if err != nil {
		fmt.Println(err)
		return
	}

	existAddedFile := false
	for _, stat := range status {
		if stat.Staging != git.Untracked && stat.Staging != git.Unmodified {
			existAddedFile = true
		}
	}

	if !existAddedFile {
		fmt.Println("オマエ コミット デキル ファイル ナイ")
		os.Exit(-1)
	}

	userName, _ := gitconfig.Username()
	email, _ := gitconfig.Email()

	commit, err := w.Commit(genshijin.ToGenshijin(message), &git.CommitOptions{
		Author: &object.Signature{
			Name:  userName,
			Email: email,
			When:  time.Now(),
		},
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	obj, err := r.CommitObject(commit)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Println(obj)
}
