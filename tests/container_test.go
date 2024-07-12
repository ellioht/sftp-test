package tests

import (
	"context"
	"github.com/ellioht/sftptest"
	"testing"
	"time"
)

func Test_CreateContainer(t *testing.T) {
	ctx := context.Background()

	cfg := &sftptest.Config{
		ImageName: sftptest.AtmozSftpImage,
		MountDir:  "upload",
	}

	ctr, err := sftptest.NewContainer(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer ctr.Close()

	t.Logf("Container running on port %s", ctr.Port)

	time.Sleep(5 * time.Second)
}
