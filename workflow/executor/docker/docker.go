package docker

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/kubework/work/errors"
	"github.com/kubework/work/util"
	"github.com/kubework/work/util/file"
	"github.com/kubework/work/workflow/common"
	execcommon "github.com/kubework/work/workflow/executor/common"
)

type DockerExecutor struct{}

func NewDockerExecutor() (*DockerExecutor, error) {
	log.Infof("Creating a docker executor")
	return &DockerExecutor{}, nil
}

func (d *DockerExecutor) GetFileContents(containerID string, sourcePath string) (string, error) {
	// Uses docker cp command to return contents of the file
	// NOTE: docker cp CONTAINER:SRC_PATH DEST_PATH|- streams the contents of the resource
	// as a tar archive to STDOUT if using - as DEST_PATH. Thus, we need to extract the
	// content from the tar archive and output into stdout. In this way, we do not need to
	// create and copy the content into a file from the wait container.
	dockerCpCmd := fmt.Sprintf("docker cp -a %s:%s - | tar -ax -O", containerID, sourcePath)
	cmd := exec.Command("sh", "-c", dockerCpCmd)
	log.Info(cmd.Args)
	out, err := cmd.Output()
	if err != nil {
		if exErr, ok := err.(*exec.ExitError); ok {
			log.Errorf("`%s` stderr:\n%s", cmd.Args, string(exErr.Stderr))
		}
		return "", errors.InternalWrapError(err)
	}
	return string(out), nil
}

func (d *DockerExecutor) CopyFile(containerID string, sourcePath string, destPath string) error {
	log.Infof("Archiving %s:%s to %s", containerID, sourcePath, destPath)
	dockerCpCmd := fmt.Sprintf("docker cp -a %s:%s - | gzip > %s", containerID, sourcePath, destPath)
	err := common.RunCommand("sh", "-c", dockerCpCmd)
	if err != nil {
		return err
	}
	copiedFile, err := os.Open(destPath)
	if err != nil {
		return err
	}
	defer util.Close(copiedFile)
	gzipReader, err := gzip.NewReader(copiedFile)
	if err != nil {
		return err
	}
	if !file.ExistsInTar(sourcePath, tar.NewReader(gzipReader)) {
		errMsg := fmt.Sprintf("path %s does not exist in archive %s", sourcePath, destPath)
		log.Warn(errMsg)
		return errors.Errorf(errors.CodeNotFound, errMsg)
	}
	log.Infof("Archiving completed")
	return nil
}

type cmdCloser struct {
	io.Reader
	cmd *exec.Cmd
}

func (c *cmdCloser) Close() error {
	err := c.cmd.Wait()
	if err != nil {
		return errors.InternalWrapError(err)
	}
	return nil
}

func (d *DockerExecutor) GetOutputStream(containerID string, combinedOutput bool) (io.ReadCloser, error) {
	cmd := exec.Command("docker", "logs", containerID)
	log.Info(cmd.Args)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, errors.InternalWrapError(err)
	}

	if !combinedOutput {
		err = cmd.Start()
		if err != nil {
			return nil, errors.InternalWrapError(err)
		}
		return stdout, nil
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, errors.InternalWrapError(err)
	}

	err = cmd.Start()
	if err != nil {
		return nil, errors.InternalWrapError(err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)
	reader, writer := io.Pipe()
	go func() {
		defer wg.Done()
		_, _ = io.Copy(writer, stdout)
	}()
	go func() {
		defer wg.Done()
		_, _ = io.Copy(writer, stderr)
	}()

	go func() {
		defer writer.Close()
		wg.Wait()
	}()

	return &cmdCloser{Reader: reader, cmd: cmd}, nil
}

func (d *DockerExecutor) WaitInit() error {
	return nil
}

// Wait for the container to complete
func (d *DockerExecutor) Wait(containerID string) error {
	return common.RunCommand("docker", "wait", containerID)
}

// killContainers kills a list of containerIDs first with a SIGTERM then with a SIGKILL after a grace period
func (d *DockerExecutor) Kill(containerIDs []string) error {
	killArgs := append([]string{"kill", "--signal", "TERM"}, containerIDs...)
	// docker kill will return with an error if a container has terminated already, which is not an error in this case.
	// We therefore ignore any error. docker wait that follows will re-raise any other error with the container.
	err := common.RunCommand("docker", killArgs...)
	if err != nil {
		log.Warningf("Ignored error from 'docker kill --signal TERM': %s", err)
	}
	waitArgs := append([]string{"wait"}, containerIDs...)
	waitCmd := exec.Command("docker", waitArgs...)
	log.Info(waitCmd.Args)
	if err := waitCmd.Start(); err != nil {
		return errors.InternalWrapError(err)
	}
	timer := time.AfterFunc(execcommon.KillGracePeriod*time.Second, func() {
		log.Infof("Timed out (%ds) for containers to terminate gracefully. Killing forcefully", execcommon.KillGracePeriod)
		forceKillArgs := append([]string{"kill", "--signal", "KILL"}, containerIDs...)
		forceKillCmd := exec.Command("docker", forceKillArgs...)
		log.Info(forceKillCmd.Args)
		// same as kill case above, we ignore any error
		err = forceKillCmd.Run()
		if err != nil {
			log.Warningf("Ignored error from 'docker kill --signal KILL': %s", err)
		}
	})
	err = waitCmd.Wait()
	_ = timer.Stop()
	if err != nil {
		return errors.InternalWrapError(err)
	}
	log.Infof("Containers %s killed successfully", containerIDs)
	return nil
}
