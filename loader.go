package main

// Loads message packages to avoid cycle if loaded inside of network package

import (
	_ "github.com/brenfwd/gocraft/network/messages/serverbound"
)
