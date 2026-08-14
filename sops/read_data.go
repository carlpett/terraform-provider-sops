package sops

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/cmd/sops/common"
	"github.com/getsops/sops/v3/cmd/sops/formats"
	"github.com/getsops/sops/v3/config"
	"github.com/getsops/sops/v3/decrypt"
	"gopkg.in/yaml.v3"

	"github.com/carlpett/terraform-provider-sops/sops/internal/dotenv"
	"github.com/carlpett/terraform-provider-sops/sops/internal/ini"
)

// sopsData holds the decrypted contents of a sops file along with the metadata
// surfaced from the encrypted document.
type sopsData struct {
	data         map[string]string
	raw          string
	lastModified time.Time
}

func readData(content []byte, format string) (sopsData, error) {
	cleartext, err := decrypt.Data(content, format)
	if userErr, ok := err.(sops.UserError); ok {
		err = userErr
	}
	if err != nil {
		return sopsData{}, fmt.Errorf("Error decrypting sops file: %w", err)
	}

	lastModified, err := readLastModified(content, format)
	if err != nil {
		return sopsData{}, fmt.Errorf("Error reading sops metadata: %w", err)
	}

	var data map[string]interface{}
	switch format {
	case "json":
		err = json.Unmarshal(cleartext, &data)
	case "yaml":
		err = yaml.Unmarshal(cleartext, &data)
	case "dotenv":
		err = dotenv.Unmarshal(cleartext, &data)
	case "ini":
		err = ini.Unmarshal(cleartext, &data)
	}
	if err != nil {
		return sopsData{}, fmt.Errorf("Error parsing decrypted data: %w", err)
	}

	return sopsData{
		data:         flatten(data),
		raw:          string(cleartext),
		lastModified: lastModified,
	}, nil
}

// readLastModified parses the sops metadata of an encrypted document (see
// loadMetadata) and returns the lastmodified timestamp recorded there. It
// never decrypts the document, so it needs no access to the data key and
// never reads any secret values.
func readLastModified(content []byte, format string) (time.Time, error) {
	metadata, err := loadMetadata(content, format)
	if err != nil {
		return time.Time{}, err
	}
	if metadata.LastModified.IsZero() {
		return time.Time{}, errors.New("missing lastmodified timestamp")
	}
	return metadata.LastModified, nil
}

// loadMetadata parses the sops metadata (e.g. lastmodified) from the encrypted
// document without decrypting it. Loading only parses the document structure and
// does not require access to the data key, so it is cheap relative to decryption.
//
// NOTE: this deliberately reaches into sops's `common`, `formats` and `config`
// packages, which — unlike the `decrypt` package this file otherwise relies on —
// are not part of sops's documented stable API (see the decrypt package doc:
// "It is the only package in SOPS with a stable API."). The stable API does not
// expose the decrypted tree's metadata, so there is currently no alternative.
// This means the document is parsed twice per read (once here, once inside
// decrypt.Data), and a future sops bump could change these internals. If sops
// ever exposes metadata through the stable API, prefer that and drop these imports.
func loadMetadata(content []byte, format string) (sops.Metadata, error) {
	store := common.StoreForFormat(formats.FormatFromString(format), config.NewStoresConfig())
	tree, err := store.LoadEncryptedFile(content)
	if err != nil {
		return sops.Metadata{}, err
	}
	return tree.Metadata, nil
}
