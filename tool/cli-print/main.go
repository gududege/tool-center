package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:] // skip program name

	fmt.Println("╔══════════════════════════════════════════════╗")
	fmt.Println("║         CLI Print — Arguments Received        ║")
	fmt.Println("╚══════════════════════════════════════════════╝")
	fmt.Println()

	if len(args) == 0 {
		fmt.Println("  (no arguments)")
		os.Exit(0)
	}

	total := len(args)
	var positional []string
	var options []string
	switchCount := 0
	i := 0
	for i < len(args) {
		a := args[i]
		if strings.HasPrefix(a, "--") || (strings.HasPrefix(a, "-") && len(a) == 2 && a[1] >= 'a' && a[1] <= 'z') {
			// Could be switch or option — check if next arg exists and is not a flag
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				options = append(options, a, args[i+1])
				i += 2
				continue
			}
			// standalone flag -> switch
			fmt.Printf("  %-30s  (switch)\n", a)
			switchCount++
		} else {
			positional = append(positional, a)
		}
		i++
	}

	if len(positional) > 0 {
		fmt.Println("  Positional Arguments:")
		for idx, p := range positional {
			fmt.Printf("    [%d] %s\n", idx+1, p)
		}
		fmt.Println()
	}

	if len(options) > 0 {
		fmt.Println("  Options:")
		for j := 0; j < len(options); j += 2 {
			flag := options[j]
			val := options[j+1]
			fmt.Printf("    %-20s  %s\n", flag, val)
		}
		fmt.Println()
	}

	if switchCount > 0 {
		fmt.Printf("  Switches: %d enabled\n", switchCount)
		fmt.Println()
	}

	fmt.Printf("  ─────────────────────────────────────\n")
	fmt.Printf("  Total: %d argument(s)\n", total)
	fmt.Printf("         %d positional · %d option · %d switch\n",
		len(positional), len(options)/2, switchCount)
}
