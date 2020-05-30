package main

import (
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pb "github.com/vesari/grpccalc/grpccalc"
	"github.com/vesari/grpccalc/server"
	"google.golang.org/grpc"
)

type fakeWriter struct {
	written []byte
}

func (fw *fakeWriter) Write(p []byte) (int, error) {
	fw.written = append(fw.written, p...)
	return len(p), nil
}

func TestAdd(t *testing.T) {
	startServer(t)

	t.Log("Starting tests")

	t.Run("1 and 2", func(t *testing.T) {
		logger, fw := setUpTest(t, []string{"add", "1", "2"})

		err := realMain(logger)
		require.NoError(t, err)

		result := getResult(fw)
		assert.Equal(t, "3", result)
	})
}

func TestMultiplyF(t *testing.T) {
	startServer(t)

	t.Log("Starting tests")

	t.Run("1.3 and 2.1", func(t *testing.T) {
		logger, fw := setUpTest(t, []string{"multiplyF", "1.3", "2.1"})

		err := realMain(logger)
		require.NoError(t, err)

		result := getResult(fw)

		assert.Equal(t, "2.73", result)
	})
}

func setUpTest(t *testing.T, args []string) (*log.Logger, *fakeWriter) {
	origArgs := os.Args
	t.Cleanup(func() {
		os.Args = origArgs
	})

	os.Args = append([]string{"client"}, args...)
	fw := &fakeWriter{
		written: []byte{},
	}
	return log.New(fw, "", log.LstdFlags), fw
}

func getResult(fw *fakeWriter) string {
	written := string(fw.written)
	lines := strings.Split(written, "\n")

	reResult := regexp.MustCompile(`Result: (.+)$`)
	result := ""
	for _, l := range lines {
		matches := reResult.FindStringSubmatch(l)
		if matches == nil {
			continue
		}

		result = matches[1]
	}

	return result
}

func startServer(t *testing.T) {
	lis, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	elems := strings.Split(lis.Addr().String(), ":")
	pStr := elems[len(elems)-1]
	p, err := strconv.Atoi(pStr)
	require.NoError(t, err)
	s := grpc.NewServer()
	pb.RegisterCalcServer(s, &server.Server{})

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Printf("Server failed: %s", err)
		}
	}()

	origPort := os.Getenv("PORT")
	os.Setenv("PORT", strconv.Itoa(p))

	t.Cleanup(func() {
		t.Log("Stopping server")
		s.Stop()

		os.Setenv("PORT", origPort)
	})
}
