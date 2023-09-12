package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strings"

	log15 "gopkg.in/inconshreveable/log15.v2"

	"git.softndit.com/collector/backend/config"
	"git.softndit.com/collector/backend/dal"
	dalpg "git.softndit.com/collector/backend/dal/pg"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/services"
	"github.com/urfave/cli"
)

const (
	separator = "=======8========="
)

var (
	logger   = log15.New("app", "collector")
	cliFlags = []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Value: "collectord.toml",
		},
		cli.BoolFlag{
			Name: "dryRun",
		},
		cli.StringFlag{
			Name:  "log",
			Value: "migration.log",
		},
	}
)

type context struct {
	DBM dal.TrManager

	SwiftFileStorage *services.SwiftFiler
	FsFileStorage    services.FileStorage
	LogFile          io.Writer
	MigratedPaths    map[string]struct{}

	dryRun bool
}

func migrateURL(oldURL url.URL) (*url.URL, error) {
	// oldURL.Scheme = "http"
	// oldURL.Host = "asd"
	return &oldURL, nil
}

func main() {
	app := cli.NewApp()
	app.Flags = cliFlags

	app.Action = func(c *cli.Context) {
		// log file
		logFile, err := os.OpenFile(c.String("log"), os.O_APPEND|os.O_RDWR, 0600)
		if err != nil {
			panic(err)
		}
		defer logFile.Close()

		processedFsPaths := make(map[string]struct{})
		// read readed
		scanner := bufio.NewScanner(logFile)
		for scanner.Scan() {
			text := scanner.Text()
			if len(text) > 0 {
				paths := strings.Split(text, separator)
				if len(paths) != 2 {
					panic(paths)
				}
				path := strings.Replace(paths[0], "\"", "", -1)
				processedFsPaths[path] = struct{}{}
			}

		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		configPath := c.String("config")

		mainConfig := &config.Config{}
		if err := mainConfig.ReadConfig(configPath); err != nil {
			panic(err)
		}

		instanceHolder, err := config.NewInstanceHolder(mainConfig)
		if err != nil {
			panic(err)
		}

		// DBM
		serverConfig := mainConfig.ServerConfig
		pgConfig, err := instanceHolder.GetPgConfig(serverConfig.DBM.PgRef)
		if err != nil {
			panic(err)
		}
		DBM, err := dalpg.NewManager(&dalpg.ManagerContext{
			PoolConfig: pgConfig,
		})
		if err != nil {
			panic(err)
		}

		// file storage
		fsCfg, err := instanceHolder.GetGenericFilerConfig("fs")
		if err != nil {
			panic(err)
		}
		// file storage
		swiftCfg, err := instanceHolder.GetGenericFilerConfig("swift")
		if err != nil {
			panic(err)
		}

		// fs builder
		fsBuilder, err := services.GetTomlFilerBuilder("fs")
		if err != nil {
			panic(err)
		}

		fsStorage, err := fsBuilder(&services.TomlFileStorageContext{
			Log:    logger.New("service", "fsFiler"),
			Config: fsCfg.Config,
		})
		if err != nil {
			panic(err)
		}

		// swift builder
		swiftBuilder, err := services.GetTomlFilerBuilder("swift")
		if err != nil {
			panic(err)
		}

		swiftStorage, err := swiftBuilder(&services.TomlFileStorageContext{
			Log:    logger.New("service", "swiftFiler"),
			Config: swiftCfg.Config,
		})
		if err != nil {
			panic(err)
		}
		swiftFiler, casted := swiftStorage.(*services.SwiftFiler)
		if !casted {
			panic("not swift")
		}

		// dryrun
		dryRun := c.Bool("dryRun") // || true

		perPage := 50

		context := &context{
			DBM:              DBM,
			dryRun:           dryRun,
			SwiftFileStorage: swiftFiler,
			FsFileStorage:    fsStorage,
			LogFile:          logFile,
			MigratedPaths:    processedFsPaths,
		}

	MainLoop:
		for paginator := dto.NewPagePaginator(0, int16(perPage)); ; paginator.NextPage() {
			medias, err := DBM.GetMediaByPage(dto.MediaTypeList, paginator)
			if err != nil {
				panic(err)
			}

			// TODO: bulk update
		Media:
			for i, media := range medias {
				fmt.Printf("process media with id: %d\n", media.ID)
				if false && (i > 4) {
					break MainLoop
				}
				used, err := DBM.IsMediaUsed(media.ID)
				if err != nil {
					panic(err)
				}
				if !used {
					fmt.Printf("media is not used and to delete: %d\n", media.ID)
					if !dryRun {
						err := DBM.DeleteMedias([]int64{media.ID})
						if err != nil {
							panic(err)
						}
					}
					continue Media
				}

				if err := migrateMediaURL(context, media); err != nil {
					panic(err)
				}
				if !dryRun {
					if err := DBM.UpdateMedia(media); err != nil {
						panic(err)
					}
				}
			}

			if mlen := len(medias); mlen == 0 || mlen < perPage {
				break
			}
		}

	}
	app.Run(os.Args)
}

func migrateMediaURL(context *context, media *dto.Media) error {
	fmt.Printf("migrate URL media: %d\n", media.ID)
	if photo := media.Photo; photo != nil {
		for _, variant := range photo.Variants {
			URLString := variant.URI
			URL, err := url.Parse(URLString)
			if err != nil {
				return err
			}

			newURL, err := migrateURL(*URL)
			if err != nil {
				return err
			}

			newURLString := newURL.String()
			variant.URI = newURLString

			if err := migrateMedia(context, URLString, newURLString); err != nil {
				return err
			}
			fmt.Printf("from: %s, to: %s\n", URLString, newURLString)
		}
	}

	if doc := media.Document; doc != nil {
		URLString := doc.URI
		URL, err := url.Parse(URLString)
		if err != nil {
			return err
		}

		newURL, err := migrateURL(*URL)
		if err != nil {
			return err
		}

		newURLString := newURL.String()
		doc.URI = newURLString

		if err := migrateMedia(context, URLString, newURLString); err != nil {
			return err
		}
		fmt.Printf("from: %s, to: %s\n", URLString, newURLString)
	}

	return nil
}

func migrateMedia(context *context, fsURL, swiftURL string) error {
	if _, found := context.MigratedPaths[fsURL]; found {
		fmt.Printf("path already processed: %s", fsURL)
		return nil
	}

	swiftLocation, err := context.SwiftFileStorage.GetFileLocation(swiftURL)
	if err != nil {
		return err
	}

	context.LogFile.Write([]byte(fmt.Sprintf("\"%s\"%s\"%s\"\n", fsURL, separator, swiftURL)))
	if context.SwiftFileStorage.IsExist(swiftLocation) {
		fmt.Printf("file %s; exists in swiftStorage\n", swiftURL)
		return nil
	}

	fsLocation, err := context.FsFileStorage.GetFileLocation(fsURL)
	if err != nil {
		fmt.Printf("111 cant find file: %s\n", fsURL)
		return nil
	}

	fileContent, err := context.FsFileStorage.GetFile(fsLocation)
	if err != nil {
		fmt.Printf("111 cant find file: %s\n", fsURL)
		return nil
	}

	fmt.Printf("fsURL: %s, swiftURL: %s\n", fsURL, swiftURL)
	fmt.Printf("fsLocation: %+v\n", fsLocation)

	if !context.dryRun {
		err := context.SwiftFileStorage.SaveByPath(swiftLocation.FullURL(), fileContent)
		if err != nil {
			return err
		}
	}

	return nil
}
