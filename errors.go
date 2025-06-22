// SPDX-FileCopyrightText: 2023 Luca <luca.tagliavini5@studio.unibo.it>
// SPDX-FileCopyrightText: 2023 Samuele Musiani <samu@teapot.ovh>
// SPDX-FileCopyrightText: 2024 Eyad Issa <eyadlorenzo@gmail.com>
//
// SPDX-License-Identifier: AGPL-3.0-only

package filenameslinter

import (
	"fmt"
)

type RegexMatchError struct {
	Regexp   string
	Filename string
	Path     string
}

func (e RegexMatchError) Error() string {
	return fmt.Sprintf("Filname %s doesn't match the regexp %s", e.Filename, e.Regexp)
}
