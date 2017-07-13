// Copyright 2015 The Polyverse Corporation. All rights reserved
// Polyverse proprietary code

// This package handles the swizzling decisions and keeps track of
// the use of swizzling for a given *.go file being compiled

package swizzle

import (
	"fmt"
	"math/rand"
	"time"
	"os"
	"strconv"
)

type SwzlCode	int
const (
	SWZL_BTS = iota
	SWZL_NIAR
	SWZL_REG
)

type swzlInfo_s struct {
	name	string	// Name of this type of swizzling
	envName	string	// Name of environment variable used to control this swizzling type
	prob	int	// Value is an unsigned integer of range [0,100] representing
			// the percentage probability that a particular swizzle occurs
	total	int	// Stats Tracking: Total number of swizzle decisions
	count	int	// Stats Tracking: Number of 'true' swizzle decisions
}

var	swzlInfo	map[SwzlCode]*swzlInfo_s	// Parameters/stats for each type of swizzle
var	swzlDebug	int;				// Value >0 increases debug verbosity

func Debug() int { return swzlDebug }

func init() {
	rand.Seed(time.Now().UTC().UnixNano())


	swzlInfo = map[SwzlCode]*swzlInfo_s {		// Parameters/stats for each type of swizzle
	    SWZL_BTS:	&swzlInfo_s{"BTS", "SWIZZLE_BTS",  0,0,0},
	    SWZL_NIAR:	&swzlInfo_s{"NIAR","SWIZZLE_NIAR", 0,0,0},
	    SWZL_REG:	&swzlInfo_s{"REG", "SWIZZLE_REG",  0,0,0},
	}

	var	tval	int
	var	err	error

	/*
	 * Process Swizzle-related environment variables
	 */

	 // Look for SWIZZLE_DEBUG
	if tvalStr, present := os.LookupEnv("SWIZZLE_DEBUG"); present {
		tval, err = strconv.Atoi(tvalStr)
		if len(tvalStr) <= 0 || err != nil || tval < 0 {
			// Environment value must be a valid integer in range [0,Inf)]
			panic("Environment variable \"SWIZZLE_DEBUG\" must an integer >= 0\n")
		}
		swzlDebug = tval
	}

	// Look for other SWIZZLE_* values defined in swzlInfo
	for key,info := range swzlInfo  {
		if tvalStr, present := os.LookupEnv(info.envName); present {
			tval, err = strconv.Atoi(tvalStr)
			if len(tvalStr) <= 0 || err != nil || tval < 0 || tval > 100 {
				// Environment value must be a valid integer in range [0,100]
				msg := fmt.Sprintf("Environment variable %s must contain an integer in the range [0,100]", info.envName)
				panic(msg)
			}
			swzlInfo[key].prob = tval
		}
	}
}

// Generates a boolean which is true 'n' percent of the time
func Decision(c SwzlCode) bool {
	const swzlRange int = 100	// Gives us an integer random variable in range [0,100)
	decision := rand.Intn(swzlRange) < swzlInfo[c].prob;
	swzlInfo[c].total++
	if decision { swzlInfo[c].count++ }
	if swzlDebug > 2 { fmt.Fprintf(os.Stderr, "Swizzle Decision (%s): %v\n", swzlInfo[c].name, decision) }
	return decision;
}

// Returns a string containing current swizzling statistics
func Stats() string {
	var statStr string
	for _,info := range swzlInfo  {
		statStr += fmt.Sprintf("%s=%d/%d ", info.name, info.count, info.total);
	}
	return statStr
}


