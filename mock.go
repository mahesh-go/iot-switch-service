package main

import (
	"math/rand"
)

/*
 * Mock function to generate random voltage and wattage
 */
func GetRandomSwitchValues() (int, int) {
	Voltage := []int{110, 120, 130, 220, 240}
	Wattage := []int{20, 25, 40, 60, 80, 110, 150}

	i := rand.Intn(4)
	j := rand.Intn(6)
	return Voltage[i], Wattage[j]
}
