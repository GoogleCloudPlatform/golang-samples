// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package display is a driver for Grove LCD RGB backlight display.
//
// More information about the display is at http://wiki.seeed.cc/Grove-LCD_RGB_Backlight/.
// This package is ported from a reference Python package, see https://goo.gl/VI59c8.
package display

import (
	"fmt"
	"time"

	"golang.org/x/exp/io/i2c"
	"golang.org/x/exp/io/i2c/driver"
)

const (
	returnHome     = 0x02
	entryModeSet   = 0x04
	displayControl = 0x08

	functionSet = 0x20

	entryRight          = 0x00
	entryLeft           = 0x02
	entryShiftIncrement = 0x01
	entryShiftDecrement = 0x00

	displayOn  = 0x04
	displayOff = 0x00
	cursorOn   = 0x02
	cursorOff  = 0x00
	blinkOn    = 0x01
	blinkOff   = 0x00

	displayMove = 0x08
	cursorMove  = 0x00
	moveRight   = 0x04
	moveLeft    = 0x00

	mode8bit = 0x10
	mode4bit = 0x00

	line2 = 0x08
	line1 = 0x00

	dots10 = 0x04
	dots8  = 0x00

	// addresses
	lcdAddr = 0x3e
	rgbAddr = 0x62

	lcdFn  = 0x08
	lcdTxt = 0x40
)

// Device represents an Grove LCD RGB Backlight device.
type Device struct {
	lcd *i2c.Device
	rgb *i2c.Device
}

// Open opens a connection the the RGB backlight display.
// Once display is no longer in-use, it should be closed by Close.
func Open(o driver.Opener) (*Device, error) {
	lcd, err := i2c.Open(o, lcdAddr)
	if err != nil {
		return nil, fmt.Errorf("cannot open LCD device: %v", err)
	}
	rgb, err := i2c.Open(o, rgbAddr)
	if err != nil {
		return nil, fmt.Errorf("cannot open RGB device: %v", err)
	}

	// two lines, regular 10 dots font.
	if err := lcd.Write([]byte{lcdFn, functionSet | displayOn | line2 | dots10}); err != nil {
		return nil, err
	}
	// direction: left to right
	if err := lcd.Write([]byte{lcdFn, displayControl | displayOn}); err != nil {
		return nil, err
	}
	// display on
	if err := lcd.Write([]byte{lcdFn, entryModeSet | entryLeft | entryShiftDecrement}); err != nil {
		return nil, err
	}

	return &Device{lcd: lcd, rgb: rgb}, nil
}

// SetText clears the screen and prints the given text on the display.
func (d *Device) SetText(text string) error {
	if err := d.Clear(); err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)

	// num lines = 2
	if err := d.lcd.Write([]byte{lcdFn, 0x20 | 0x08}); err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)

	// return home
	if err := d.lcd.Write([]byte{lcdFn, returnHome}); err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)

	var row int
	var col int
	for _, c := range text {
		// If current col is larger than the width of display
		// or there is a new line, break the line.
		if c == '\n' || col == 16 {
			col = 0
			row++
			if row == 2 {
				return nil
			}
			// new line command
			if err := d.lcd.Write([]byte{lcdFn, 0xc0}); err != nil {
				return err
			}
			if c == '\n' {
				continue
			}
		}
		if err := d.lcd.Write([]byte{lcdTxt, byte(c)}); err != nil {
			return err
		}
		col++
	}
	return nil
}

// Clear clears the screen.
func (d *Device) Clear() error {
	return d.lcd.Write([]byte{0x80, 0x01}) // clear display
}

// SetRGB sets the backlight to the given color.
func (d *Device) SetRGB(r, g, b int) error {
	cmds := [][]byte{
		{0, 0},
		{1, 0},
		{0x08, 0xaa},
		{4, byte(r)},
		{3, byte(g)},
		{2, byte(g)},
	}
	for _, cmd := range cmds {
		if err := d.rgb.Write(cmd); err != nil {
			return err
		}
	}
	return nil
}

// Close closes connection to the device.
func (d *Device) Close() error {
	if err := d.rgb.Close(); err != nil {
		return err
	}
	return d.lcd.Close()
}
