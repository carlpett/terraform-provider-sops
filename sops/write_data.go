package sops

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	sopslib "github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/aes"
	"github.com/getsops/sops/v3/cmd/sops/common"
	sopsformats "github.com/getsops/sops/v3/cmd/sops/formats"
	sopsconfig "github.com/getsops/sops/v3/config"
	"github.com/getsops/sops/v3/keyservice"
	"github.com/getsops/sops/v3/version"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// writeEncryptedFileFromDynamic marshals value (from dynamicToGo), encrypts it with sops and
// writes it to destinationFile.
func writeEncryptedFileFromDynamic(destinationFile types.String, outputType types.String, value any) error {
	destinationFileValue := destinationFile.ValueString()

	format, err := formatForFile(destinationFileValue, outputType)
	if err != nil {
		return err
	}

	content, err := marshalDynamicValue(format, value)
	if err != nil {
		return newSummaryError("Error marshalling data", err)
	}

	return encryptAndWrite(destinationFileValue, format, content)
}

func encryptAndWrite(destinationFileValue string, format string, content []byte) error {
	absPath, err := filepath.Abs(destinationFileValue)
	if err != nil {
		return newSummaryError("Error resolving file path", err)
	}

	store := common.StoreForFormat(sopsformats.FormatFromString(format), sopsconfig.NewStoresConfig())

	branches, err := store.LoadPlainFile(content)
	if err != nil {
		return newSummaryError("Error parsing content", err)
	}
	if len(branches) < 1 {
		return newSummaryError("Error parsing content", fmt.Errorf("content must contain at least one document"))
	}

	metadata, err := sopsMetadataForFile(absPath, store)
	if err != nil {
		return newSummaryError("Error determining encryption keys", err)
	}

	tree := sopslib.Tree{
		Branches: branches,
		Metadata: metadata,
		FilePath: absPath,
	}

	dataKey, errs := tree.GenerateDataKeyWithKeyServices([]keyservice.KeyServiceClient{keyservice.NewLocalClient()})
	if len(errs) > 0 {
		return newSummaryError("Error generating data key", fmt.Errorf("%v", errs))
	}

	if err := common.EncryptTree(common.EncryptTreeOpts{
		DataKey: dataKey,
		Tree:    &tree,
		Cipher:  aes.NewCipher(),
	}); err != nil {
		return newSummaryError("Error encrypting content", err)
	}

	encrypted, err := store.EmitEncryptedFile(tree)
	if err != nil {
		return newSummaryError("Error marshalling encrypted content", err)
	}

	if err := os.WriteFile(destinationFileValue, encrypted, 0644); err != nil {
		return newSummaryError("Error writing file", err)
	}

	return nil
}

func formatForFile(destinationFileValue string, outputType types.String) (string, error) {
	var format string
	if !outputType.IsNull() {
		format = outputType.ValueString()
	} else {
		switch ext := path.Ext(destinationFileValue); ext {
		case ".json":
			format = "json"
		case ".yaml", ".yml":
			format = "yaml"
		default:
			return "", newSummaryError("Unknown file type", fmt.Errorf("don't know how to encode file with extension %s, set output_type to \"json\" or \"yaml\"", ext))
		}
	}

	switch format {
	case "json", "yaml":
		return format, nil
	default:
		return "", newSummaryError("Invalid output type", fmt.Errorf("data_wo only supports output_type %q or %q, got %q", "json", "yaml", format))
	}
}

// sopsMetadataForFile mirrors how `sops encrypt` derives keys when none are passed explicitly.
func sopsMetadataForFile(absPath string, store common.Store) (sopslib.Metadata, error) {
	configResult, err := sopsconfig.LookupConfigFile(absPath)
	if err != nil {
		return sopslib.Metadata{}, fmt.Errorf("no .sops.yaml config file found for %s: %w", absPath, err)
	}

	conf, err := sopsconfig.LoadCreationRuleForFile(configResult.Path, absPath, nil)
	if err != nil {
		return sopslib.Metadata{}, err
	}
	if conf == nil {
		return sopslib.Metadata{}, fmt.Errorf("%s has no creation rules", configResult.Path)
	}

	isSingleValueStore := false
	if svs, ok := store.(sopslib.SingleValueStore); ok {
		isSingleValueStore = svs.IsSingleValueStore()
	}

	unencryptedSuffix := conf.UnencryptedSuffix
	encryptedSuffix := conf.EncryptedSuffix
	unencryptedRegex := conf.UnencryptedRegex
	encryptedRegex := conf.EncryptedRegex
	unencryptedCommentRegex := conf.UnencryptedCommentRegex
	encryptedCommentRegex := conf.EncryptedCommentRegex
	macOnlyEncrypted := conf.MACOnlyEncrypted

	if isSingleValueStore {
		// A single-value store has one key, which should always be encrypted.
		unencryptedSuffix = ""
		encryptedSuffix = ""
		unencryptedRegex = ""
		encryptedRegex = ""
		unencryptedCommentRegex = ""
		encryptedCommentRegex = ""
		macOnlyEncrypted = false
	} else if unencryptedSuffix == "" && encryptedSuffix == "" && unencryptedRegex == "" &&
		encryptedRegex == "" && unencryptedCommentRegex == "" && encryptedCommentRegex == "" {
		unencryptedSuffix = sopslib.DefaultUnencryptedSuffix
	}

	return sopslib.Metadata{
		KeyGroups:               conf.KeyGroups,
		ShamirThreshold:         conf.ShamirThreshold,
		UnencryptedSuffix:       unencryptedSuffix,
		EncryptedSuffix:         encryptedSuffix,
		UnencryptedRegex:        unencryptedRegex,
		EncryptedRegex:          encryptedRegex,
		UnencryptedCommentRegex: unencryptedCommentRegex,
		EncryptedCommentRegex:   encryptedCommentRegex,
		MACOnlyEncrypted:        macOnlyEncrypted,
		Version:                 version.Version,
	}, nil
}
