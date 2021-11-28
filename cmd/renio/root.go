package renio

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/abdfnx/renio/core/options"
	"github.com/abdfnx/renio/config"
	_package "github.com/abdfnx/renio/package"
	"github.com/abdfnx/renio/tools"
	"github.com/abdfnx/renio/cmd/factory"

	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

var homeDir, _ = homedir.Dir()

var installDir = path.Join(homeDir, "./.renio/")

// Renio functions expected to be passed into cmd
type Renio struct {
	Run    func(opt options.Options)
	Dev    func(og string, opt options.Options)
	Bundle func(file string, minify bool, config *config.Config) string
}

// Execute start the CLI
func Execute(renio Renio, f *factory.Factory) *cobra.Command {
	tools.CheckDotRenioDir()

	// Load renio mod.toml
	config, err := config.ConfigLoad()
	tools.Check(err)

	color.NoColor = config.Options.NoColor

	var fsFlag bool
	var netFlag bool
	var minifyFlag bool
	var envFlag bool
	const desc = `🦏 A secure, lightweight, and fast runtime for JavaScript and TypeScript.`

	// Root command
	var rootCmd = &cobra.Command{
		Use:   "renio <subcommand> [flags]",
		Short:  desc,
		Long: desc,
		SilenceErrors: true,
		Example: heredoc.Doc(`
			renio run <file> --net
			renio test --fs
			renio bundle <file> --minify
			renio dev <file> --env
		`),
		Annotations: map[string]string{
			"help:tellus": heredoc.Doc(`
				Open an issue at https://github.com/abdfnx/renio/issues
			`),
		},
	}

	rootCmd.SetOut(f.IOStreams.Out)
	rootCmd.SetErr(f.IOStreams.ErrOut)

	cs := f.IOStreams.ColorScheme()

	helpHelper := func(command *cobra.Command, args []string) {
		rootHelpFunc(cs, command, args)
	}

	rootCmd.PersistentFlags().Bool("help", false, "Help for renio")
	rootCmd.SetHelpFunc(helpHelper)
	rootCmd.SetUsageFunc(rootUsageFunc)
	rootCmd.SetFlagErrorFunc(rootFlagErrorFunc)

	// Run sub-command
	var runCmd = &cobra.Command{
		Use:   "run [file]",
		Short: "Run a JavaScript or TypeScript source file",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) >= 0 {
				env := options.Environment{
					NoColor:  config.Options.NoColor,
					Args:     args[1:],
					RunTests: false,
				}
				
				// check if the argument start with "http" and download it
				if strings.HasPrefix(args[0], "http") {
					// fmt.Println(color.Cyan("Download ") + args[0])

					c := color.New(color.FgCyan)

					c.Print("Download ")
					fmt.Println(args[0])
					
					tools.Download(args[0])

					b := renio.Bundle(tools.FileX(args[0]), true, config)

					opt := options.Options{
						SourceFile: path.Join(installDir, "temp.js"),
						Source:     b,
						Perms:      &options.Perms{fsFlag, netFlag, envFlag},
						Env:        env,
					}

					renio.Run(opt)
				} else {
					bundle := renio.Bundle(args[0], true, config)

					opt := options.Options{
						SourceFile: args[0],
						Source:     bundle,
						Perms:      &options.Perms{fsFlag, netFlag, envFlag},
						Env:        env,
					}

					renio.Run(opt)
				}
			}
		},
	}

	// --fs, --net, --env flags
	/*
		--fs: allow filesystem access
		--net: allow network access
		--env: allow environment variables
	*/

	runCmd.Flags().BoolVar(&fsFlag, "fs", false, "Allow file system access")
	runCmd.Flags().BoolVar(&netFlag, "net", false, "Allow net access")
	runCmd.Flags().BoolVar(&envFlag, "env", false, "Allow Environment Variables access")

	// dev sub-command to run script in development mode
	var devCmd = &cobra.Command{
		Use:   "dev [file]",
		Short: "Run a script in development mode.",
		Long:  `Run a script in development mode. It enables type-checking using the inbuilt TypeScript compiler 📦.`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) >= 0 {
				bundle := renio.Bundle(args[0], true, config)
				env := options.Environment{
					NoColor:  config.Options.NoColor,
					Args:     args[1:],
					RunTests: false,
				}

				opt := options.Options{
					SourceFile: args[0],
					Source:     bundle,
					Perms: &options.Perms{
						Fs:  true,
						Env: true,
						Net: true,
					},
					Env: env,
				}

				og, _ := ioutil.ReadFile(args[0])
				renio.Dev(string(og), opt)
			}
		},
	}

	// bundle sub-command to bundle a source file
	var bundleCmd = &cobra.Command{
		Use:   "bundle [file]",
		Short: "Bundle your script to a single JavaScript file",
		Long:  `Bundle your script to a single JavaScript file. It utilizes esbuild for super fast bundling.`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) >= 0 {
				out := renio.Bundle(args[0], minifyFlag, config)
				fmt.Println(out)
			}
		},
	}

	// --minify flag for bundling
	bundleCmd.Flags().BoolVarP(&minifyFlag, "minify", "m", false, "Minify the output bundle")

	// package sub-command for trigger the package
	var packageCmd = &cobra.Command{
		Use:   "package [file]",
		Short: "Package your script to a standalone executable.",
		Long:  `Package your script to a standalone executable.`,
		Aliases: []string{"pkg"},
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) >= 0 {
				_package.PkgSource(args[0])
			}
		},
	}

	// test sub-command to run test files
	var testCmd = &cobra.Command{
		Use:   "test",
		Short: "Run tests for your Renio scripts.",
		Long:  `Run tests for your Renio scripts. All files matching *_test.js | *_test.ts | *.test.js | *.test.ts are run.`,
		Run: func(cmd *cobra.Command, args []string) {
			env := options.Environment{
				NoColor:  config.Options.NoColor,
				Args:     args,
				RunTests: true,
			}

			opt := options.Options{
				Perms: &options.Perms{fsFlag, netFlag, envFlag},
				Env:   env,
			}

			tests := CollectAllTests()

			for _, test := range tests {
				opt.SourceFile = test
				bundle := renio.Bundle(test, true, config)
				opt.Source = bundle
				renio.Run(opt)
			}
		},
	}

	// --fs, --net, --env perms
	testCmd.Flags().BoolVar(&fsFlag, "fs", false, "Allow file system access")
	testCmd.Flags().BoolVar(&netFlag, "net", false, "Allow net access")
	testCmd.Flags().BoolVar(&envFlag, "env", false, "Allow Environment Variables access")

	// Add sub-commands to root command
	rootCmd.AddCommand(bundleCmd, runCmd, packageCmd, devCmd, testCmd)

	return rootCmd
}

// replace a string if it is windows
func isWindows(toChange *string, replacement string) {
	if runtime.GOOS == "windows" {
		*toChange = replacement
	}
}

func shebang(loc string) string {
	exec := `
#!/bin/sh
renio "run" "%s" "$@"`
	// windows
	isWindows(&exec, `@renio "run" "%s" %*`)

	return fmt.Sprintf(exec, loc)
}

// match test files
func matchedFiles(name string) bool {
	matchedJS, err := filepath.Match("*_test.js", name)
	matchedTS, err := filepath.Match("*_test.ts", name)
	matchedJSTest, err := filepath.Match("*.test.js", name)
	matchedTSTest, err := filepath.Match("*.test.ts", name)

	if err != nil {
		log.Fatal(err)
	}

	return (matchedJS || matchedTS || matchedTSTest || matchedJSTest)
}

// CollectAllTests files
func CollectAllTests() []string {
	var testFiles []string

	e := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err == nil {
			if err != nil {
				return nil
			}

			if matchedFiles(info.Name()) {
				testFiles = append(testFiles, path)
			}
		}

		return nil
	})

	if e != nil {
		log.Fatal(e)
	}

	return testFiles
}
