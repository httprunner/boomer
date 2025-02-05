/*
 * Copyright 2020 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package data

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

// hrpPath is .hrp directory under the user directory.
var hrpPath string

//go:embed x509/*
var x509Dir embed.FS

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	hrpPath = filepath.Join(home, ".hrp")
	_ = EnsureFolderExists(filepath.Join(hrpPath, "x509"))
}

// Path returns the absolute path the given relative file or directory path
func Path(rel string) (destPath string) {
	destPath = rel
	if !filepath.IsAbs(rel) {
		destPath = filepath.Join(hrpPath, rel)
	}
	if !IsFilePathExists(destPath) {
		content, err := x509Dir.ReadFile(rel)
		if err != nil {
			return
		}

		err = os.WriteFile(destPath, content, 0o644)
		if err != nil {
			return
		}
	}
	return
}

func EnsureFolderExists(folderPath string) error {
	if !IsPathExists(folderPath) {
		err := CreateFolder(folderPath)
		return err
	} else if IsFilePathExists(folderPath) {
		return fmt.Errorf("path %v should be directory", folderPath)
	}
	return nil
}

// IsPathExists returns true if path exists, whether path is file or dir
func IsPathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// IsFilePathExists returns true if path exists and path is file
func IsFilePathExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		// path not exists
		return false
	}

	// path exists
	if info.IsDir() {
		// path is dir, not file
		return false
	}
	return true
}

func CreateFolder(folderPath string) error {
	log.Info().Str("path", folderPath).Msg("create folder")
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		log.Error().Err(err).Msg("create folder failed")
		return err
	}
	return nil
}
