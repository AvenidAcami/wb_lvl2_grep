package cmd

import (
	"fmt"
	"wb_lvl2_grep/internal/grep"

	"github.com/spf13/cobra"
)

var (
	stringCountAfter              int
	stringCountBefore             int
	stringCountBeforeAndAfter     int
	printCount                    bool
	ignoreRegister                bool
	invertFilter                  bool
	sampleIsNotARegularExpression bool
	printStringNumberBeforeString bool
)

var grepCmd = &cobra.Command{
	Use:  "grep",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pattern := args[0]

		resChan, err := grep.FilterRows(pattern, grep.Options{
			StringCountAfter:              stringCountAfter,
			StringCountBefore:             stringCountBefore,
			StringCountBeforeAndAfter:     stringCountBeforeAndAfter,
			PrintCount:                    printCount,
			IgnoreRegister:                ignoreRegister,
			InvertFilter:                  invertFilter,
			SampleIsNotARegularExpression: sampleIsNotARegularExpression,
			PrintStringNumberBeforeString: printStringNumberBeforeString,
		})

		if err != nil {
			fmt.Println("error in regular expression")
			return err
		}

		for line := range resChan {
			fmt.Println(line)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(grepCmd)
	grepCmd.Flags().IntVarP(&stringCountAfter, "strings_after", "A", 0, "после каждой найденной строки дополнительно вывести N строк после неё")
	grepCmd.Flags().IntVarP(&stringCountBefore, "strings_before", "B", 0, "вывести N строк до каждой найденной строки")
	grepCmd.Flags().IntVarP(&stringCountBeforeAndAfter, "strings_before_and_after", "C", 0, "вывести N строк контекста вокруг найденной строки")
	grepCmd.Flags().BoolVarP(&printCount, "print_count", "c", false, "вывести только то количество строк, которые совпадают с шаблоном")
	grepCmd.Flags().BoolVarP(&ignoreRegister, "ignore_register", "i", false, "игнорировать регистр")
	grepCmd.Flags().BoolVarP(&invertFilter, "invert", "v", false, "инвертировать фильтр: выводить строки, не содержащие шаблон")
	grepCmd.Flags().BoolVarP(&sampleIsNotARegularExpression, "sample_is_a_substring", "F", false, "воспринимать шаблон как фиксированную строку, а не регулярное выражение")
	grepCmd.Flags().BoolVarP(&printStringNumberBeforeString, "print_string_number", "n", false, "выводить номер строки перед каждой найденной строкой")
}
