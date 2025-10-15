package git

import (
	"errors"
	"strconv"
	"strings"
	"time"

	gitv5 "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

// RawCommit represents a single, unprocessed commit from the git history.
type RawCommit struct {
    Hash    string
    Message string
}

// GetLatestTag returns the most recent git tag from the repository.
func GetLatestTag(repoPath string) (string, error) {
    repo, err := gitv5.PlainOpen(repoPath)
    if err != nil {
        return "", err
    }
    tagsIter, err := repo.Tags()
    if err != nil {
        return "", err
    }
    var latest string
    var latestTime time.Time
    _ = tagsIter.ForEach(func(ref *plumbing.Reference) error {
        var commitHash plumbing.Hash
        // Try annotated tag
        if obj, err := repo.TagObject(ref.Hash()); err == nil {
            commitHash = obj.Target
        } else {
            // Lightweight tag points directly to the commit
            commitHash = ref.Hash()
        }
        if commit, err := repo.CommitObject(commitHash); err == nil {
            candidate := ref.Name().Short()
            commitTime := commit.Committer.When
            if commitTime.After(latestTime) || (commitTime.Equal(latestTime) && semverGreater(candidate, latest)) {
                latestTime = commitTime
                latest = candidate
            }
        }
        return nil
    })
    if latest == "" {
        return "", errors.New("no tags found")
    }
    return latest, nil
}

// semverGreater compares two tag strings like v1.2.3 and returns true if a > b.
func semverGreater(a, b string) bool {
    if a == "" {
        return false
    }
    if b == "" {
        return true
    }
    strip := func(s string) string {
        if strings.HasPrefix(s, "v") || strings.HasPrefix(s, "V") {
            return s[1:]
        }
        return s
    }
    ap := strings.Split(strip(a), ".")
    bp := strings.Split(strip(b), ".")
    for len(ap) < 3 {
        ap = append(ap, "0")
    }
    for len(bp) < 3 {
        bp = append(bp, "0")
    }
    for i := 0; i < 3; i++ {
        ai, _ := strconv.Atoi(ap[i])
        bi, _ := strconv.Atoi(bp[i])
        if ai > bi {
            return true
        }
        if ai < bi {
            return false
        }
    }
    return false
}

// GetCommitsSince reads the git log and returns all commits between the 'fromRef' and HEAD.
func GetCommitsSince(repoPath, fromRef string) ([]RawCommit, error) {
    repo, err := gitv5.PlainOpen(repoPath)
    if err != nil {
        return nil, err
    }

    // List all commits reachable from HEAD
    head, err := repo.Head()
    if err != nil {
        return nil, err
    }
    cIter, err := repo.Log(&gitv5.LogOptions{From: head.Hash()})
    if err != nil {
        return nil, err
    }
    var out []RawCommit
    stopAt := plumbing.ZeroHash
    if fromRef != "" {
        ref, err := repo.Tag(fromRef)
        if err == nil {
            obj, err := repo.TagObject(ref.Hash())
            if err == nil {
                // Annotated tag
                stopAt = obj.Target
            } else {
                // Lightweight tag
                stopAt = ref.Hash()
            }
        }
    }
    err = cIter.ForEach(func(c *object.Commit) error {
        if stopAt != plumbing.ZeroHash && c.Hash == stopAt {
            return storer.ErrStop
        }
        out = append(out, RawCommit{Hash: c.Hash.String(), Message: c.Message})
        return nil
    })
    if err != nil {
        return nil, err
    }
    return out, nil
}


