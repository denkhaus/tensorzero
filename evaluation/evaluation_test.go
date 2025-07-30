//go:build unit

package evaluation

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRunRequest(t *testing.T) {
	variants := map[string]string{"A": "variant1", "B": "variant2"}
	tags := map[string]string{"eval_group": "groupA"}
	projectName := "project_X"
	displayName := "Eval Run 1"

	req := RunRequest{
		Variants:    variants,
		Tags:        tags,
		ProjectName: &projectName,
		DisplayName: &displayName,
	}

	assert.Equal(t, variants, req.Variants)
	assert.Equal(t, tags, req.Tags)
	assert.Equal(t, projectName, *req.ProjectName)
	assert.Equal(t, displayName, *req.DisplayName)
}

func TestEpisodeRequest(t *testing.T) {
	runID := uuid.New()
	taskName := "task_Y"
	datapointName := "dp_Z"
	tags := map[string]string{"type": "episode"}

	req := EpisodeRequest{
		RunID:         runID,
		TaskName:      &taskName,
		DatapointName: &datapointName,
		Tags:          tags,
	}

	assert.Equal(t, runID, req.RunID)
	assert.Equal(t, taskName, *req.TaskName)
	assert.Equal(t, datapointName, *req.DatapointName)
	assert.Equal(t, tags, req.Tags)
}

func TestRunResponse(t *testing.T) {
	runID := uuid.New()
	response := RunResponse{RunID: runID}
	assert.Equal(t, runID, response.RunID)
}

func TestEpisodeResponse(t *testing.T) {
	episodeID := uuid.New()
	response := EpisodeResponse{EpisodeID: episodeID}
	assert.Equal(t, episodeID, response.EpisodeID)
}