// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package integrations

import (
	"net/http"
	"testing"

	"code.gitea.io/gitea/models"
	api "code.gitea.io/gitea/modules/structs"

	"github.com/stretchr/testify/assert"
)

func compareCommitFiles(t *testing.T, expect []string, files []*api.CommitAffectedFiles) {
	var actual []string
	for i := range files {
		actual = append(actual, files[i].Filename)
	}
	assert.ElementsMatch(t, expect, actual)
}

func TestAPIReposGitCommits(t *testing.T) {
	defer prepareTestEnv(t)()
	user := models.AssertExistsAndLoadBean(t, &models.User{ID: 2}).(*models.User)
	// Login as User2.
	session := loginUser(t, user.Name)
	token := getTokenForLoggedInUser(t, session)

	// check invalid requests
	req := NewRequestf(t, "GET", "/api/v1/repos/%s/repo1/git/commits/12345?token="+token, user.Name)
	session.MakeRequest(t, req, http.StatusNotFound)

	req = NewRequestf(t, "GET", "/api/v1/repos/%s/repo1/git/commits/..?token="+token, user.Name)
	session.MakeRequest(t, req, http.StatusUnprocessableEntity)

	req = NewRequestf(t, "GET", "/api/v1/repos/%s/repo1/git/commits/branch-not-exist?token="+token, user.Name)
	session.MakeRequest(t, req, http.StatusNotFound)

	for _, ref := range [...]string{
		"master", // Branch
		"v1.1",   // Tag
		"65f1",   // short sha
		"65f1bf27bc3bf70f64657658635e66094edbcb4d", // full sha
	} {
		req := NewRequestf(t, "GET", "/api/v1/repos/%s/repo1/git/commits/%s?token="+token, user.Name, ref)
		session.MakeRequest(t, req, http.StatusOK)
	}
}

func TestAPIReposGitCommitList(t *testing.T) {
	defer prepareTestEnv(t)()
	user := models.AssertExistsAndLoadBean(t, &models.User{ID: 2}).(*models.User)
	// Login as User2.
	session := loginUser(t, user.Name)
	token := getTokenForLoggedInUser(t, session)

	// Test getting commits (Page 1)
	req := NewRequestf(t, "GET", "/api/v1/repos/%s/repo16/commits?token="+token, user.Name)
	resp := session.MakeRequest(t, req, http.StatusOK)

	var apiData []api.Commit
	DecodeJSON(t, resp, &apiData)

	assert.Len(t, apiData, 3)
	assert.EqualValues(t, "69554a64c1e6030f051e5c3f94bfbd773cd6a324", apiData[0].CommitMeta.SHA)
	compareCommitFiles(t, []string{"readme.md"}, apiData[0].Files)
	assert.EqualValues(t, "27566bd5738fc8b4e3fef3c5e72cce608537bd95", apiData[1].CommitMeta.SHA)
	compareCommitFiles(t, []string{"readme.md"}, apiData[1].Files)
	assert.EqualValues(t, "5099b81332712fe655e34e8dd63574f503f61811", apiData[2].CommitMeta.SHA)
	compareCommitFiles(t, []string{"readme.md"}, apiData[2].Files)
}

func TestAPIReposGitCommitListPage2Empty(t *testing.T) {
	defer prepareTestEnv(t)()
	user := models.AssertExistsAndLoadBean(t, &models.User{ID: 2}).(*models.User)
	// Login as User2.
	session := loginUser(t, user.Name)
	token := getTokenForLoggedInUser(t, session)

	// Test getting commits (Page=2)
	req := NewRequestf(t, "GET", "/api/v1/repos/%s/repo16/commits?token="+token+"&page=2", user.Name)
	resp := session.MakeRequest(t, req, http.StatusOK)

	var apiData []api.Commit
	DecodeJSON(t, resp, &apiData)

	assert.Len(t, apiData, 0)
}

func TestAPIReposGitCommitListDifferentBranch(t *testing.T) {
	defer prepareTestEnv(t)()
	user := models.AssertExistsAndLoadBean(t, &models.User{ID: 2}).(*models.User)
	// Login as User2.
	session := loginUser(t, user.Name)
	token := getTokenForLoggedInUser(t, session)

	// Test getting commits (Page=1, Branch=good-sign)
	req := NewRequestf(t, "GET", "/api/v1/repos/%s/repo16/commits?token="+token+"&sha=good-sign", user.Name)
	resp := session.MakeRequest(t, req, http.StatusOK)

	var apiData []api.Commit
	DecodeJSON(t, resp, &apiData)

	assert.Len(t, apiData, 1)
	assert.Equal(t, "f27c2b2b03dcab38beaf89b0ab4ff61f6de63441", apiData[0].CommitMeta.SHA)
	compareCommitFiles(t, []string{"readme.md"}, apiData[0].Files)
}
