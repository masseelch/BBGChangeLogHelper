package changelog

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

func OpenRepository() (*git.Repository, error) {
	// Path to cache directory.
	cd, err := os.UserCacheDir()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not create cache dir: %s", err))
	}
	cd = filepath.Join(cd, "bbg-changelog-helper")

	// Make sure the cache directory exists.
	if err := os.MkdirAll(cd, os.ModePerm); err != nil {
		return nil, errors.New(fmt.Sprintf("Could not create cache dir: %s", err))
	}

	// Open the repository.
	r, err := git.PlainOpen(cd)
	if err != nil {
		if err != git.ErrRepositoryNotExists {
			return nil, errors.New(fmt.Sprintf("Could not open repository: %s", err))
		}

		// Path exits but repository is empty. Clone it.
		r, err = git.PlainClone(cd, true, &git.CloneOptions{
			URL:        "https://github.com/iElden/BetterBalancedGame",
			NoCheckout: true,
		})
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Could not open repository: %s", err))
		}
	}

	// Make sure we have fetched the latest changes.
	if err := r.Fetch(&git.FetchOptions{}); err != nil && err != git.NoErrAlreadyUpToDate {
		return nil, errors.New(fmt.Sprintf("Could not fetch updates from remote: %s", err))
	}

	return r, nil
}

func RetrieveCommit(r *git.Repository,  tagOrHash string) (*object.Commit, error) {
	var sc *object.Commit
	if t, err := r.Tag(tagOrHash); err == nil {
		// Tag was found. Retrieve the commit.
		if sc, err = r.CommitObject(t.Hash()); err != nil {
			return nil, errors.New(fmt.Sprintf("Commit not found: %s", err))
		}
	} else {
		// Tag was not found. Try to use given value as hash.
		if sc, err = r.CommitObject(plumbing.NewHash(viper.GetString(tagOrHash))); err != nil {
			return nil, errors.New(fmt.Sprintf("Commit not found: %s", err))
		}
	}

	return sc, nil
}
