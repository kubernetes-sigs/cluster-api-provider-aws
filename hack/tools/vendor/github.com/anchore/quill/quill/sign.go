package quill

import (
	"fmt"
	"os"
	"path"

	blacktopMacho "github.com/blacktop/go-macho"

	macholibre "github.com/anchore/go-macholibre"
	"github.com/anchore/quill/internal/bus"
	"github.com/anchore/quill/internal/log"
	"github.com/anchore/quill/quill/event"
	"github.com/anchore/quill/quill/macho"
	"github.com/anchore/quill/quill/pki"
	"github.com/anchore/quill/quill/pki/load"
	"github.com/anchore/quill/quill/sign"
)

type SigningConfig struct {
	SigningMaterial pki.SigningMaterial
	Identity        string
	Path            string
}

func NewSigningConfigFromPEMs(binaryPath, certificate, privateKey, password string, failWithoutFullChain bool) (*SigningConfig, error) {
	var signingMaterial pki.SigningMaterial
	if certificate != "" {
		sm, err := pki.NewSigningMaterialFromPEMs(certificate, privateKey, password, failWithoutFullChain)
		if err != nil {
			return nil, err
		}

		signingMaterial = *sm
	}

	return &SigningConfig{
		Path:            binaryPath,
		Identity:        path.Base(binaryPath),
		SigningMaterial: signingMaterial,
	}, nil
}

func NewSigningConfigFromP12(binaryPath string, p12Content load.P12Contents, failWithoutFullChain bool) (*SigningConfig, error) {
	signingMaterial, err := pki.NewSigningMaterialFromP12(p12Content, failWithoutFullChain)
	if err != nil {
		return nil, err
	}

	return &SigningConfig{
		Path:            binaryPath,
		Identity:        path.Base(binaryPath),
		SigningMaterial: *signingMaterial,
	}, nil
}

func (c *SigningConfig) WithIdentity(id string) *SigningConfig {
	if id != "" {
		c.Identity = id
	}
	return c
}

func (c *SigningConfig) WithTimestampServer(url string) *SigningConfig {
	c.SigningMaterial.TimestampServer = url
	return c
}

func Sign(cfg SigningConfig) error {
	f, err := os.Open(cfg.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	if macholibre.IsUniversalMachoBinary(f) {
		return signMultiarchBinary(cfg)
	}

	mon := bus.PublishTask(
		event.Title{
			Default:      "Sign binary",
			WhileRunning: "Signing binary",
			OnSuccess:    "Signed binary",
		},
		cfg.Path,
		-1,
	)

	err = signSingleBinary(cfg)
	if err != nil {
		mon.Err = err
	} else {
		mon.SetCompleted()
	}
	return err
}

//nolint:funlen
func signMultiarchBinary(cfg SigningConfig) error {
	log.WithFields("binary", cfg.Path).Info("signing multi-arch binary")

	f, err := os.Open(cfg.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	dir, err := os.MkdirTemp("", "quill-extract-"+path.Base(cfg.Path))
	if err != nil {
		return fmt.Errorf("unable to create temp directory to extract multi-arch binary: %w", err)
	}
	defer os.RemoveAll(dir)

	extractMon := bus.PublishTask(
		event.Title{
			Default:      "Extract universal binary",
			WhileRunning: "Extracting universal binary",
			OnSuccess:    "Extracted universal binary",
		},
		cfg.Path,
		-1,
	)

	extractedFiles, err := macholibre.Extract(f, dir)
	if err != nil {
		extractMon.Err = err
		return fmt.Errorf("unable to extract multi-arch binary: %w", err)
	}

	extractMon.Stage.Current = fmt.Sprintf("%d nested binaries", len(extractedFiles))

	extractMon.SetCompleted()

	log.WithFields("binary", cfg.Path, "arches", len(extractedFiles)).Trace("discovered nested binaries within multi-arch binary")

	var cfgs []SigningConfig
	for _, ef := range extractedFiles {
		c := cfg
		c.Path = ef.Path
		cfgs = append(cfgs, c)
	}

	signMon := bus.PublishTask(
		event.Title{
			Default:      "Sign binaries",
			WhileRunning: "Signing binaries",
			OnSuccess:    "Signed binaries",
		},
		cfg.Path,
		len(cfgs),
	)

	defer signMon.SetCompleted()

	for _, c := range cfgs {
		signMon.Stage.Current = path.Base(c.Path)
		if err := signSingleBinary(c); err != nil {
			signMon.Err = err
			return err
		}
		signMon.N++
	}

	signMon.Stage.Current = ""

	var paths []string
	for _, c := range cfgs {
		paths = append(paths, c.Path)
	}

	log.WithFields("binary", cfg.Path, "arches", len(cfgs)).Info("packaging signed binaries into single multi-arch binary")

	packMon := bus.PublishTask(
		event.Title{
			Default:      "Repack universal binary",
			WhileRunning: "Repacking universal binary",
			OnSuccess:    "Repacked universal binary",
		},
		cfg.Path,
		-1,
	)

	defer packMon.SetCompleted()

	if err := macholibre.Package(cfg.Path, paths...); err != nil {
		packMon.Err = err
		return err
	}

	return nil
}

func signSingleBinary(cfg SigningConfig) error {
	log.WithFields("binary", cfg.Path).Info("signing binary")

	m, err := macho.NewFile(cfg.Path)
	if err != nil {
		return err
	}

	// check there already isn't a LcCodeSignature loader already (if there is, bail)
	if m.HasCodeSigningCmd() {
		log.Debug("binary already signed, removing signature...")
		if err := m.RemoveSigningContent(); err != nil {
			return fmt.Errorf("unable to remove existing code signature: %+v", err)
		}
	}

	if cfg.SigningMaterial.Signer == nil {
		bus.Notify("Warning: performed ad-hoc sign, which means that anyone can alter the binary contents without you knowing (there is no cryptographic signature)")
		log.Warnf("only ad-hoc signing, which means that anyone can alter the binary contents without you knowing (there is no cryptographic signature)")
	}

	// (patch) add empty LcCodeSignature loader (offset and size references are not set)
	if err = m.AddEmptyCodeSigningCmd(); err != nil {
		return err
	}

	// first pass: add the signed data with the dummy loader
	log.Debugf("estimating signing material size")
	superBlobSize, sbBytes, err := sign.GenerateSigningSuperBlob(cfg.Identity, m, cfg.SigningMaterial, 0)
	if err != nil {
		return fmt.Errorf("failed to add signing data on pass=1: %w", err)
	}

	// (patch) make certain offset and size references to the superblob are finalized in the binary
	log.Debugf("patching binary with updated superblob offsets")
	if err = sign.UpdateSuperBlobOffsetReferences(m, uint64(len(sbBytes))); err != nil {
		return nil
	}

	// second pass: now that all of the sizing is right, let's do it again with the final contents (replacing the hashes and signature)
	log.Debug("creating signature for binary")
	_, sbBytes, err = sign.GenerateSigningSuperBlob(cfg.Identity, m, cfg.SigningMaterial, superBlobSize)
	if err != nil {
		return fmt.Errorf("failed to add signing data on pass=2: %w", err)
	}

	// (patch) append the superblob to the __LINKEDIT section
	log.Debugf("patching binary with signature")

	codeSigningCmd, _, err := m.CodeSigningCmd()
	if err != nil {
		return err
	}

	if err = m.Patch(sbBytes, len(sbBytes), uint64(codeSigningCmd.DataOffset)); err != nil {
		return fmt.Errorf("failed to patch super blob onto macho binary: %w", err)
	}

	return nil
}

func IsSigned(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	if macholibre.IsUniversalMachoBinary(f) {
		log.WithFields("binary", path).Trace("binary is a universal binary")
		mf, err := blacktopMacho.NewFatFile(f)
		if mf == nil || err != nil {
			return false, fmt.Errorf("failed to parse universal macho binary: %w", err)
		}
		defer mf.Close()

		success := true
		for _, arch := range mf.Arches {
			sig := arch.CodeSignature()
			if sig == nil {
				log.WithFields("binary", path, "arch", arch.String()).Trace("no code signature block found")

				return false, nil
			}
			log.WithFields("length", len(sig.CMSSignature), "arch", arch.String()).Trace("CMS signature found")

			success = success && len(sig.CMSSignature) > 0
		}

		return success, nil
	}

	log.WithFields("binary", path).Trace("binary is for a single architecture")

	mf, err := blacktopMacho.NewFile(f)
	if mf == nil || err != nil {
		return false, fmt.Errorf("failed to parse macho binary: %w", err)
	}

	defer mf.Close()

	sig := mf.CodeSignature()
	if sig == nil {
		log.WithFields("binary", path).Trace("no code signature block found")
		return false, nil
	}

	log.WithFields("length", len(sig.CMSSignature)).Trace("CMS signature found")

	return len(sig.CMSSignature) > 0, nil
}
