package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/coreos/go-semver/semver"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/plumbing/transport"
	gogitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"golang.org/x/crypto/ssh"
)

func getLatestVersionFromGit(path string) (string, error) {
	tags, err := getGitTags(".")
	if err != nil {
		return "", err
	}

	prefix, versions := makeVersions(tags)
	latestVersion := getLatestVersion(versions)
	return prefix + latestVersion.String(), nil

}

func getLatestVersion(versions []*semver.Version) *semver.Version {
	if len(versions) == 0 {
		versions = append(versions, semver.New("0.0.0"))
	}

	semver.Sort(versions)

	lastVersion := versions[len(versions)-1]
	lastVersionCopy := *lastVersion
	return &lastVersionCopy
}

func makeVersions(tagrefs storer.ReferenceIter) (string, []*semver.Version) {
	versionPrefix := ""
	versions := []*semver.Version{}

	err := tagrefs.ForEach(func(t *plumbing.Reference) error {

		tagName := string(t.Name().Short())

		if strings.HasPrefix(tagName, "v") {
			versionPrefix = "v"
			tagName = strings.TrimPrefix(tagName, "v")
		}

		v, err := semver.NewVersion(tagName)
		if err != nil {
			fmt.Println("found malformed tag: ", tagName)
			return nil
		}
		versions = append(versions, v)
		return nil
	})
	if err != nil {
		return "", nil
	}

	return versionPrefix, versions
}

func copyVersion(v *semver.Version) *semver.Version {
	return semver.New(v.String())
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

func getGitTags(path string) (storer.ReferenceIter, error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	tagrefs, err := r.Tags()
	if err != nil {
		return nil, err
	}
	return tagrefs, nil
}

func printAction(lastVersion, newVersion string) {
	fmt.Printf("%s => %s\n", lastVersion, newVersion)
}

func tagAndPush(path, remote, sshFile, version string) error {
	r, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	head, err := r.Head()
	if err != nil {
		return err
	}

	fmt.Printf("tagging repo version %s...\n", version)
	ref, err := r.CreateTag(version, head.Hash(), &git.CreateTagOptions{Message: version})
	if err != nil {
		return err
	}

	fmt.Printf("pushing %s to %s...\n", ref.Name(), remote)

	if err := r.Push(&git.PushOptions{
		Auth:       getSSHKeyAuth(sshFile),
		RemoteName: remote,
		RefSpecs: []config.RefSpec{
			config.RefSpec(ref.Name() + ":" + ref.Name()),
		},
	}); err != nil {
		return err
	}

	fmt.Println("done")
	return nil
}

func getSSHKeyAuth(privateSshKeyFile string) transport.AuthMethod {
	var auth transport.AuthMethod
	sshKey, _ := ioutil.ReadFile(privateSshKeyFile)
	signer, _ := ssh.ParsePrivateKey([]byte(sshKey))
	auth = &gogitssh.PublicKeys{User: "git", Signer: signer}
	return auth
}
