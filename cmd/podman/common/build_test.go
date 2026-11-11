package common_test

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.podman.io/podman/v6/cmd/podman/common"
)

func TestParseBuildOptsManifestTagConflict(t *testing.T) {
	for _, tc := range []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{"only tag", []string{"--manifest", "mylist", "--tag", "mylist"}, true},
		{"additional tag", []string{"--tag", "myimage", "--manifest", "mylist", "--tag", "mylist"}, true},
		// --stage-labels without --save-stages makes ParseBuildOpts stop at a
		// later check, so a different error proves the names were accepted.
		{"different names", []string{"--tag", "myimage", "--manifest", "mylist", "--stage-labels"}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			buildOpts := common.BuildFlagsWrapper{}
			cmd := &cobra.Command{}
			common.DefineBuildFlags(cmd, &buildOpts, false)
			require.NoError(t, cmd.ParseFlags(tc.args))

			_, err := common.ParseBuildOpts(cmd, []string{t.TempDir()}, &buildOpts)
			if tc.wantErr {
				assert.ErrorContains(t, err, "the same name must not be specified for both '--tag' and '--manifest'")
				return
			}
			assert.ErrorContains(t, err, `"--stage-labels" requires "--save-stages"`)
		})
	}
}
