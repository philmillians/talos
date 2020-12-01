// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package libretechallh3cch5

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/sys/unix"

	"github.com/talos-systems/talos/internal/app/machined/pkg/runtime"
	"github.com/talos-systems/talos/pkg/machinery/constants"
)

var (
	bin       = fmt.Sprintf("/usr/install/u-boot/%s/u-boot-sunxi-with-spl.bin", constants.BoardLibretechAllH3CCH5)
	off int64 = 1024 * 8
)

// LibretechAllH3CCH5 represents the Libre Computer ALL-H3-CC (Tritium).
//
// Reference: https://libre.computer/products/boards/all-h3-cc/
type LibretechAllH3CCH5 struct{}

// Name implenents the runtime.Board.
func (l *LibretechAllH3CCH5) Name() string {
	return constants.BoardLibretechAllH3CCH5
}

// Install implenents the runtime.Board.
func (l *LibretechAllH3CCH5) Install(disk string) (err error) {
	var f *os.File

	if f, err = os.OpenFile(disk, os.O_RDWR|unix.O_CLOEXEC, 0o666); err != nil {
		return err
	}
	// nolint: errcheck
	defer f.Close()

	var uboot []byte

	uboot, err = ioutil.ReadFile(bin)
	if err != nil {
		return err
	}

	log.Printf("writing %s at offset %d", bin, off)

	var n int

	n, err = f.WriteAt(uboot, off)
	if err != nil {
		return err
	}

	log.Printf("wrote %d bytes", n)

	// NB: In the case that the block device is a loopback device, we sync here
	// to esure that the file is written before the loopback device is
	// unmounted.
	err = f.Sync()
	if err != nil {
		return err
	}

	return nil
}

// PartitionOptions implenents the runtime.Board.
func (l *LibretechAllH3CCH5) PartitionOptions() *runtime.PartitionOptions {
	return &runtime.PartitionOptions{PartitionsOffset: 2048}
}