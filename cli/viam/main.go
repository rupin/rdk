// Package main is the CLI command itself.
package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	rdkcli "go.viam.com/rdk/cli"
)

func main() {
	app := &cli.App{
		Name:            "viam",
		Usage:           "interact with your Viam robots",
		HideHelpCommand: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:   "base-url",
				Hidden: true,
				Value:  "https://app.viam.com:443",
				Usage:  "base URL of app",
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "load configuration from `FILE`",
			},
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"vvv"},
				Usage:   "enable debug logging",
			},
		},
		Commands: []*cli.Command{
			{
				Name: "login",
				// NOTE(benjirewis): maintain `auth` as an alias for backward compatibility.
				Aliases:         []string{"auth"},
				Usage:           "login to app.viam.com",
				HideHelpCommand: true,
				Action:          rdkcli.LoginAction,
				Subcommands: []*cli.Command{
					{
						Name:   "print-access-token",
						Usage:  "print the access token associated with current credentials",
						Action: rdkcli.PrintAccessTokenAction,
					},
				},
			},
			{
				Name:   "logout",
				Usage:  "logout from current session",
				Action: rdkcli.LogoutAction,
			},
			{
				Name:   "whoami",
				Usage:  "get currently logged-in user",
				Action: rdkcli.WhoAmIAction,
			},
			{
				Name:            "organizations",
				Usage:           "work with organizations",
				HideHelpCommand: true,
				Subcommands: []*cli.Command{
					{
						Name:   "list",
						Usage:  "list organizations for the current user",
						Action: rdkcli.ListOrganizationsAction,
					},
				},
			},
			{
				Name:            "locations",
				Usage:           "work with locations",
				HideHelpCommand: true,
				Subcommands: []*cli.Command{
					{
						Name:      "list",
						Usage:     "list locations for the current user",
						ArgsUsage: "[organization]",
						Action:    rdkcli.ListLocationsAction,
					},
				},
			},
			{
				Name:            "data",
				Usage:           "work with data",
				HideHelpCommand: true,
				Subcommands: []*cli.Command{
					{
						Name:  "export",
						Usage: "download data from Viam cloud",
						UsageText: fmt.Sprintf("viam data export <%s> <%s> [other options]",
							rdkcli.DataFlagDestination, rdkcli.DataFlagDataType),
						Flags: []cli.Flag{
							&cli.PathFlag{
								Name:     rdkcli.DataFlagDestination,
								Required: true,
								Usage:    "output directory for downloaded data",
							},
							&cli.StringFlag{
								Name:     rdkcli.DataFlagDataType,
								Required: true,
								Usage:    "data type to be downloaded: either binary or tabular",
							},
							&cli.StringSliceFlag{
								Name:  rdkcli.DataFlagOrgIDs,
								Usage: "orgs filter",
							},
							&cli.StringSliceFlag{
								Name:  rdkcli.DataFlagLocationIDs,
								Usage: "locations filter",
							},
							&cli.StringFlag{
								Name:  rdkcli.DataFlagRobotID,
								Usage: "robot-id filter",
							},
							&cli.StringFlag{
								Name:  rdkcli.DataFlagPartID,
								Usage: "part id filter",
							},
							&cli.StringFlag{
								Name:  rdkcli.DataFlagRobotName,
								Usage: "robot name filter",
							},
							&cli.StringFlag{
								Name:  rdkcli.DataFlagPartName,
								Usage: "part name filter",
							},
							&cli.StringFlag{
								Name:  rdkcli.DataFlagComponentType,
								Usage: "component type filter",
							},
							&cli.StringFlag{
								Name:  rdkcli.DataFlagComponentName,
								Usage: "component name filter",
							},
							&cli.StringFlag{
								Name:  rdkcli.DataFlagMethod,
								Usage: "method filter",
							},
							&cli.StringSliceFlag{
								Name:  rdkcli.DataFlagMimeTypes,
								Usage: "mime types filter",
							},
							&cli.UintFlag{
								Name:        rdkcli.DataFlagParallelDownloads,
								Usage:       "number of download requests to make in parallel",
								DefaultText: "10",
							},
							&cli.StringFlag{
								Name:  rdkcli.DataFlagStart,
								Usage: "ISO-8601 timestamp indicating the start of the interval filter",
							},
							&cli.StringFlag{
								Name:  rdkcli.DataFlagEnd,
								Usage: "ISO-8601 timestamp indicating the end of the interval filter",
							},
							&cli.StringSliceFlag{
								Name: rdkcli.DataFlagTags,
								Usage: "tags filter. " +
									"accepts tagged for all tagged data, untagged for all untagged data, or a list of tags for all data matching any of the tags",
							},
							&cli.StringSliceFlag{
								Name: rdkcli.DataFlagBboxLabels,
								Usage: "bbox labels filter. " +
									"accepts string labels corresponding to bounding boxes within images",
							},
						},
						Action: rdkcli.DataExportAction,
					},
					{
						Name:      "delete",
						Usage:     "delete data from Viam cloud",
						UsageText: fmt.Sprintf("viam data delete <%s> [other options]", rdkcli.DataFlagDataType),
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     rdkcli.DataFlagDataType,
								Required: true,
								Usage:    "data type to be deleted: either binary or tabular",
							},
							&cli.StringSliceFlag{
								Name:  rdkcli.DataFlagOrgIDs,
								Usage: "orgs filter",
							},
							&cli.StringSliceFlag{
								Name:  rdkcli.DataFlagLocationIDs,
								Usage: "locations filter",
							},
							&cli.StringFlag{
								Name:  rdkcli.DataFlagRobotID,
								Usage: "robot id filter",
							},
							&cli.StringFlag{
								Name:  rdkcli.DataFlagPartID,
								Usage: "part id filter",
							},
							&cli.StringFlag{
								Name:  rdkcli.DataFlagRobotName,
								Usage: "robot name filter",
							},
							&cli.StringFlag{
								Name:  rdkcli.DataFlagPartName,
								Usage: "part name filter",
							},
							&cli.StringFlag{
								Name:  rdkcli.DataFlagComponentType,
								Usage: "component type filter",
							},
							&cli.StringFlag{
								Name:  rdkcli.DataFlagComponentName,
								Usage: "component name filter",
							},
							&cli.StringFlag{
								Name:  rdkcli.DataFlagMethod,
								Usage: "method filter",
							},
							&cli.StringSliceFlag{
								Name:  rdkcli.DataFlagMimeTypes,
								Usage: "mime types filter",
							},
							&cli.StringFlag{
								Name:  rdkcli.DataFlagStart,
								Usage: "ISO-8601 timestamp indicating the start of the interval filter",
							},
							&cli.StringFlag{
								Name:  rdkcli.DataFlagEnd,
								Usage: "ISO-8601 timestamp indicating the end of the interval filter",
							},
						},
						Action: rdkcli.DataDeleteAction,
					},
				},
			},
			{
				Name:            "robots",
				Usage:           "work with robots",
				HideHelpCommand: true,
				Subcommands: []*cli.Command{
					{
						Name:  "list",
						Usage: "list robots in an organization and location",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:        "organization",
								DefaultText: "first organization alphabetically",
							},
							&cli.StringFlag{
								Name:        "location",
								DefaultText: "first location alphabetically",
							},
						},
						Action: rdkcli.ListRobotsAction,
					},
				},
			},
			{
				Name:            "robot",
				Usage:           "work with a robot",
				HideHelpCommand: true,
				Subcommands: []*cli.Command{
					{
						Name:      "status",
						Usage:     "display robot status",
						UsageText: "viam robot status <robot> [other options]",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:        "organization",
								DefaultText: "first organization alphabetically",
							},
							&cli.StringFlag{
								Name:        "location",
								DefaultText: "first location alphabetically",
							},
							&cli.StringFlag{
								Name:     "robot",
								Required: true,
							},
						},
						Action: rdkcli.RobotStatusAction,
					},
					{
						Name:      "logs",
						Usage:     "display robot logs",
						UsageText: "viam robot logs <robot> [other options]",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:        "organization",
								DefaultText: "first organization alphabetically",
							},
							&cli.StringFlag{
								Name:        "location",
								DefaultText: "first location alphabetically",
							},
							&cli.StringFlag{
								Name:     "robot",
								Required: true,
							},
							&cli.BoolFlag{
								Name:  "errors",
								Usage: "show only errors",
							},
						},
						Action: rdkcli.RobotLogsAction,
					},
					{
						Name:            "part",
						Usage:           "work with a robot part",
						HideHelpCommand: true,
						Subcommands: []*cli.Command{
							{
								Name:      "status",
								Usage:     "display part status",
								UsageText: "viam robot part status <robot> <part> [other options]",
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:        "organization",
										DefaultText: "first organization alphabetically",
									},
									&cli.StringFlag{
										Name:        "location",
										DefaultText: "first location alphabetically",
									},
									&cli.StringFlag{
										Name:     "robot",
										Required: true,
									},
									&cli.StringFlag{
										Name:     "part",
										Required: true,
									},
								},
								Action: rdkcli.RobotPartStatusAction,
							},
							{
								Name:      "logs",
								Usage:     "display part logs",
								UsageText: "viam robot part logs <robot> <part> [other options]",
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:        "organization",
										DefaultText: "first organization alphabetically",
									},
									&cli.StringFlag{
										Name:        "location",
										DefaultText: "first location alphabetically",
									},
									&cli.StringFlag{
										Name:     "robot",
										Required: true,
									},
									&cli.StringFlag{
										Name:     "part",
										Required: true,
									},
									&cli.BoolFlag{
										Name:  "errors",
										Usage: "show only errors",
									},
									&cli.BoolFlag{
										Name:    "tail",
										Aliases: []string{"f"},
										Usage:   "follow logs",
									},
								},
								Action: rdkcli.RobotPartLogsAction,
							},
							{
								Name:      "run",
								Usage:     "run a command on a robot part",
								UsageText: "viam robot part run <organization> <location> <robot> <part> [other options] <service.method>",
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:     "organization",
										Required: true,
									},
									&cli.StringFlag{
										Name:     "location",
										Required: true,
									},
									&cli.StringFlag{
										Name:     "robot",
										Required: true,
									},
									&cli.StringFlag{
										Name:     "part",
										Required: true,
									},
									&cli.StringFlag{
										Name:    "data",
										Aliases: []string{"d"},
									},
									&cli.DurationFlag{
										Name:    "stream",
										Aliases: []string{"s"},
									},
								},
								Action: rdkcli.RobotPartRunAction,
							},
							{
								Name:        "shell",
								Usage:       "start a shell on a robot part",
								Description: `In order to use the shell command, the robot must have a valid shell type service.`,
								UsageText:   "viam robot part shell <organization> <location> <robot> <part>",
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:     "organization",
										Required: true,
									},
									&cli.StringFlag{
										Name:     "location",
										Required: true,
									},
									&cli.StringFlag{
										Name:     "robot",
										Required: true,
									},
									&cli.StringFlag{
										Name:     "part",
										Required: true,
									},
								},
								Action: rdkcli.RobotPartShellAction,
							},
						},
					},
				},
			},
			{
				Name:            "module",
				Usage:           "manage your modules in Viam's registry",
				HideHelpCommand: true,
				Subcommands: []*cli.Command{
					{
						Name:  "create",
						Usage: "create & register a module on app.viam.com",
						Description: `Creates a module in app.viam.com to simplify code deployment.
Ex: 'viam module create --name my-great-module --org-id <my org id>'
Will create the module and a corresponding meta.json file in the current directory.

If your org has set a namespace in app.viam.com then your module name will be 'my-namespace:my-great-module' and
you won't have to pass a namespace or org-id in future commands. Otherwise there will be no namespace
and you will have to provide the org-id to future cli commands. You cannot make your module public until you claim an org-id.

After creation, use 'viam module update' to push your new module to app.viam.com.`,
						UsageText: "viam module create <name> [other options]",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "name",
								Usage:    "name of your module (cannot be changed once set)",
								Required: true,
							},
							&cli.StringFlag{
								Name:  "public-namespace",
								Usage: "the public namespace where the module will reside (alternative way of specifying the org id)",
							},
							&cli.StringFlag{
								Name:  "org-id",
								Usage: "id of the organization that will host the module",
							},
						},
						Action: rdkcli.CreateModuleAction,
					},
					{
						Name:  "update",
						Usage: "update a module's metadata on app.viam.com",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:        "module",
								Usage:       "path to meta.json",
								DefaultText: "./meta.json",
								TakesFile:   true,
							},
							&cli.StringFlag{
								Name:  "public-namespace",
								Usage: "the public namespace where the module resides (alternative way of specifying the org id)",
							},
							&cli.StringFlag{
								Name:  "org-id",
								Usage: "id of the organization that hosts the module",
							},
						},
						Action: rdkcli.UpdateModuleAction,
					},
					{
						Name:  "upload",
						Usage: "upload a new version of your module",
						Description: `Upload an archive containing your module's file(s) for a specified platform

Example for linux/amd64:
tar -czf packaged-module.tar.gz my-binary   # the meta.json entrypoint is relative to the root of the archive, so it should be "./my-binary"
viam module upload --version "0.1.0" --platform "linux/amd64" packaged-module.tar.gz
                        `,
						UsageText: "viam module upload <version> <platform> [other options] <packaged-module.tar.gz>",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:        "module",
								Usage:       "path to meta.json",
								DefaultText: "./meta.json",
								TakesFile:   true,
							},
							&cli.StringFlag{
								Name:  "public-namespace",
								Usage: "the public namespace where the module resides (alternative way of specifying the org id)",
							},
							&cli.StringFlag{
								Name:  "org-id",
								Usage: "id of the organization that hosts the module",
							},
							&cli.StringFlag{
								Name:  "name",
								Usage: "name of the module (used if you don't have a meta.json)",
							},
							&cli.StringFlag{
								Name:     "version",
								Usage:    "version of the module to upload (semver2.0) ex: \"0.1.0\"",
								Required: true,
							},
							&cli.StringFlag{
								Name: "platform",
								Usage: `platform of the binary you are uploading. Must be one of:
                        linux/amd64
                        linux/arm64
                        darwin/amd64 (for intel macs)
                        darwin/arm64 (for non-intel macs)`,
								Required: true,
							},
						},
						Action: rdkcli.UploadModuleAction,
					},
				},
			},
			{
				Name:   "version",
				Usage:  "print version info for this program",
				Action: rdkcli.VersionAction,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		rdkcli.Errorf(app.ErrWriter, err.Error())
	}
}
