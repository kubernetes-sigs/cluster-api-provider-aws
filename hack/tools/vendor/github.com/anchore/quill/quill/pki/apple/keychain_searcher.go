package apple

import (
	"crypto/x509"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/anchore/quill/internal/log"
	"github.com/anchore/quill/quill/pki/certchain"
	"github.com/anchore/quill/quill/pki/load"
)

var _ certchain.Searcher = (*keychainSearcher)(nil)

type keychainSearcher struct {
	keychainPath string
}

func NewKeychainSearcher(keychainPath string) certchain.Searcher {
	return &keychainSearcher{
		keychainPath: keychainPath,
	}
}

func (k keychainSearcher) CertificatesByCN(commonName string) ([]*x509.Certificate, error) {
	contents, err := searchKeychain(commonName, k.keychainPath)
	if err != nil {
		return nil, fmt.Errorf("unable to search keychain: %w", err)
	}

	certs, err := load.CertificatesFromPEM([]byte(contents))
	if err != nil {
		return nil, fmt.Errorf("unable to load certificates from PEM on keychain: %w", err)
	}

	return certs, nil
}

func searchKeychain(certCNSearch, keychainPath string) (string, error) {
	if !isMacOS() {
		log.Warn("searching the keychain is only supported on macOS. Skipping...")
		return "", nil
	}

	contents, err := run("security", "find-certificate", "-a", "-c", certCNSearch, "-p", keychainPath)
	if err != nil {
		return "", err
	}
	return contents, nil
}

func isMacOS() bool {
	return strings.Contains(strings.ToLower(runtime.GOOS), "darwin")
}

func run(args ...string) (string, error) {
	baseCmd := args[0]
	cmdArgs := args[1:]

	log.Tracef("running command: %q", strings.Join(args, " "))

	cmd := exec.Command(baseCmd, cmdArgs...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
