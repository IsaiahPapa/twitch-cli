// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package util

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
)

// RandomUserID generates a random user ID from 1->100,000,000 for use in mock events
func RandomUserID() string {
	uid, err := rand.Int(rand.Reader, big.NewInt(1*100*100*100*100))
	if err != nil {
		log.Fatal(err.Error())
	}
	return uid.String()
}

// RandomGUID generates a random GUID for use with creating IDs in the local store and for mock events
func RandomGUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return uuid
}

// RandomClientID generates a fake client ID of length 30
func RandomClientID() string {
	b := make([]byte, 30)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:30]
}

// RandomViewerCount generates a fake viewercount between 0->100,000
func RandomViewerCount() int64 {
	viewer, err := rand.Int(rand.Reader, big.NewInt(10*100*100))
	if err != nil {
		log.Fatal(err.Error())
	}
	return viewer.Int64()
}

// RandomInt generates a random integer between 0->max
func RandomInt(max int64) int64 {
	someInt, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		log.Fatal(err.Error())
	}

	return someInt.Int64()
}

// RandomType generates a fake type; Either bits, subscription, or other, in roughly even distribution
func RandomType() string {
	someInt, err := rand.Int(rand.Reader, big.NewInt(1*10*100*100*100))
	if err != nil {
		log.Fatal(err.Error())
	}
	if (someInt.Int64() % 3) == 0 {
		return "bits"
	} else if (someInt.Int64() % 3) == 1 {
		return "other"
	} else {
		return "subscription"
	}
}

type RGBColor struct {
	Red   int64
	Green int64
	Blue  int64
}

func RandomColorInRgb() RGBColor {
	Red := RandomInt(255)
	Green := RandomInt(255)
	blue := RandomInt(255)
	c := RGBColor{Red, Green, blue}
	return c
}

func RandomColorInHex() string {
	color := RandomColorInRgb()
	hex := "#" + getHex(color.Red) + getHex(color.Green) + getHex(color.Blue)
	return hex
}

func getHex(num int64) string {
	hex := fmt.Sprintf("%x", num)
	if len(hex) == 1 {
		hex = "0" + hex
	}
	return hex
}
