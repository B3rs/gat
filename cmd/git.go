package cmd

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/coreos/go-semver/semver"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func getLatestTagFromGit(repo *git.Repository) (*object.Commit, *plumbing.Reference, error) {

	tagRefs, err := repo.Tags()
	if err != nil {
		return nil, nil, err
	}

	var latestTagCommit *object.Commit
	var latestTagRef *plumbing.Reference
	err = tagRefs.ForEach(func(tagRef *plumbing.Reference) error {
		revision := plumbing.Revision(tagRef.Name().String())
		tagCommitHash, err := repo.ResolveRevision(revision)
		if err != nil {
			return err
		}

		commit, err := repo.CommitObject(*tagCommitHash)
		if err != nil {
			return err
		}

		if latestTagCommit == nil {
			latestTagCommit = commit
			latestTagRef = tagRef
		}

		if commit.Committer.When.After(latestTagCommit.Committer.When) {
			latestTagCommit = commit
			latestTagRef = tagRef
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return latestTagCommit, latestTagRef, nil

}

func isHead(repo *git.Repository, commit *object.Commit) (bool, error) {
	head, err := repo.Head()
	if err != nil {
		return false, err
	}

	return head.Hash() == commit.Hash, nil
}

func bumpVersion(v, action string) string {
	prefix := ""
	if strings.HasPrefix(v, "v") {
		prefix = "v"
		v = strings.TrimPrefix(v, "v")
	}
	version := semver.New(v)

	switch action {
	case "major":
		version.BumpMajor()
	case "minor":
		version.BumpMinor()
	case "patch":
		version.BumpPatch()
	}

	return prefix + version.String()
}

func tag(repo *git.Repository, version string) (*plumbing.Reference, error) {
	head, err := repo.Head()
	if err != nil {
		return nil, err
	}

	tagger, err := createTagger()
	if err != nil {
		return nil, err
	}
	return repo.CreateTag(version, head.Hash(), &git.CreateTagOptions{Message: version, Tagger: tagger})
}

func push(repo *git.Repository, ref *plumbing.Reference, remote, sshFile string, sshFilePassword string) error {
	auth, err := publicKey(sshFile, sshFilePassword)
	if err != nil {
		return err
	}
	po := &git.PushOptions{
		Auth:       auth,
		RemoteName: remote,
		Progress:   os.Stdout,
		RefSpecs: []config.RefSpec{
			config.RefSpec(ref.Name() + ":" + ref.Name()),
		},
	}

	return repo.Push(po)
}

func publicKey(filePath string, filePwd string) (*ssh.PublicKeys, error) {
	sshKey, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	publicKey, err := ssh.NewPublicKeys("git", []byte(sshKey), filePwd)
	if err != nil {
		return nil, err
	}
	return publicKey, err
}

func createTagger() (*object.Signature, error) {
	name, err := getUserName()
	if err != nil {
		return nil, err
	}
	email, err := getUserEmail()
	if err != nil {
		return nil, err
	}
	return &object.Signature{
		Name:  name,
		Email: email,
		When:  time.Now(),
	}, nil
}

func getUserName() (string, error) {
	cmd := exec.Command("git", "config", "user.name")
	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(stdout), nil
}

func getUserEmail() (string, error) {
	cmd := exec.Command("git", "config", "user.email")
	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(stdout), nil
}
