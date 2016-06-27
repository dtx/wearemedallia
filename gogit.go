package vcs

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// GitVCSDriver is the VCS driver for Globe to manage template versions with Git
type GitVCSDriver struct {
	remote     string
	mainBranch string
}

func (g *GitVCSDriver) runCommand(cmdStr string) error {
	cmdSlice := strings.Split(cmdStr, " ")
	cmd := exec.Command(cmdSlice[0], cmdSlice[1:]...)
	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}
	fmt.Printf("Globe command: %q", out.String())
	return nil
}

// ResetCommitID resets the local repo to the provided commit ID
func (g *GitVCSDriver) ResetCommitID(id string) {
	command := fmt.Sprintf("git reset %s", id)
	g.runCommand(command)
}

// GetLatest is a hard reset to remote HEAD, so make sure that this is a different clone than the work you are developing on to avoid
// complications. The goal to encourage seperation between developers and user clones. Globe aims to move around a lot of commits
// depending on the version of the target cluster, hence it makes sense to operate it on a seperate clone that just tracks the master.
func (g *GitVCSDriver) GetLatest() {
	fetchall := "git fetch --all"
	resetToCommit := fmt.Sprintf("git reset --hard %s/%s", g.remote, g.mainBranch)
	g.runCommand(fetchall)
	g.runCommand(resetToCommit)
}

// ChangeBranch Git checkout a remote branch with --track option. This method should return err when there is no remote branch
// with the name specified
func (g *GitVCSDriver) ChangeBranch(name string) error {
	command := fmt.Sprintf("git checkout --track %s/%s", g.remote, name)
	err := g.runCommand(command)
	if err != nil {
		return errors.New("There was a problem with checking out the branch. Does it exist? Confirm and try again")
	}
	return nil
}

func main() {
	gg := GitVCSDriver{remote: "origin", mainBranch: "master"}
	gg.GetLatest()
}
