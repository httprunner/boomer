package main

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"

	"github.com/httprunner/boomer"
	hrpboomer "github.com/httprunner/boomer/hrp"
	"github.com/httprunner/httprunner/v4/hrp"
)

// boomCmd represents the boom command
var boomCmd = &cobra.Command{
	Use:   "hrpboom",
	Short: "run load test with boomer",
	Long:  `run yaml/json testcase files for load test`,
	Example: `  $ hrpboom demo.json	# run specified json testcase file
  $ hrpboom demo.yaml	# run specified yaml testcase file
  $ hrpboom examples/	# run testcases in specified folder`,
	Args: cobra.MinimumNArgs(0),
	PreRun: func(cmd *cobra.Command, args []string) {
		boomer.SetUlimit(10240) // ulimit -n 10240
		// disable info logs for load testing
		log.Info().Msg("Set global log level to WARN for load testing")
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var paths []hrp.ITestCase
		for _, arg := range args {
			path := hrp.TestCasePath(arg)
			paths = append(paths, &path)
		}

		// if set profile, the priority is higher than the other commands
		if boomArgs.profile != "" {
			err := hrp.LoadFileObject(boomArgs.profile, &boomArgs.Profile)
			if err != nil {
				log.Error().Err(err).Msg("failed to load profile")
				return err
			}
		}

		// init boomer
		var hrpBoomer *hrpboomer.HRPBoomer
		if boomArgs.master {
			hrpBoomer = hrpboomer.NewMasterBoomer(boomArgs.masterBindHost, boomArgs.masterBindPort)
		} else if boomArgs.worker {
			hrpBoomer = hrpboomer.NewWorkerBoomer(boomArgs.masterHost, boomArgs.masterPort)
		} else {
			hrpBoomer = hrpboomer.NewStandaloneBoomer(boomArgs.SpawnCount, boomArgs.SpawnRate)
		}
		hrpBoomer.SetProfile(&boomArgs.Profile)
		ctx := hrpBoomer.EnableGracefulQuit(context.Background())

		// run boomer
		switch hrpBoomer.GetMode() {
		case "master":
			hrpBoomer.SetTestCasesPath(args)
			if boomArgs.autoStart {
				hrpBoomer.SetAutoStart()
				hrpBoomer.SetExpectWorkers(boomArgs.expectWorkers, boomArgs.expectWorkersMaxWait)
				hrpBoomer.SetSpawnCount(boomArgs.SpawnCount)
				hrpBoomer.SetSpawnRate(boomArgs.SpawnRate)
				hrpBoomer.SetRunTime(boomArgs.RunTime)
			}
			if boomArgs.autoStart {
				hrpBoomer.InitBoomer()
			} else {
				go hrpBoomer.StartServer(ctx, boomArgs.masterHttpAddress)
			}
			go hrpBoomer.PollTestCases(ctx)
			hrpBoomer.RunMaster()
		case "worker":
			if boomArgs.ignoreQuit {
				hrpBoomer.SetIgnoreQuit()
			}
			go hrpBoomer.PollTasks(ctx)
			hrpBoomer.RunWorker()
		case "standalone":
			hrpBoomer.InitBoomer()
			hrpBoomer.Run(paths...)
		}
		return nil
	},
}

type BoomArgs struct {
	boomer.Profile
	profile              string
	master               bool
	worker               bool
	ignoreQuit           bool
	masterHost           string
	masterPort           int
	masterBindHost       string
	masterBindPort       int
	masterHttpAddress    string
	autoStart            bool
	expectWorkers        int
	expectWorkersMaxWait int
}

var boomArgs BoomArgs

func init() {
	boomCmd.Flags().Int64Var(&boomArgs.MaxRPS, "max-rps", 0, "Max RPS that boomer can generate, disabled by default.")
	boomCmd.Flags().StringVar(&boomArgs.RequestIncreaseRate, "request-increase-rate", "-1", "Request increase rate, disabled by default.")
	boomCmd.Flags().Int64Var(&boomArgs.SpawnCount, "spawn-count", 1, "The number of users to spawn for load testing")
	boomCmd.Flags().Float64Var(&boomArgs.SpawnRate, "spawn-rate", 1, "The rate for spawning users")
	boomCmd.Flags().Int64Var(&boomArgs.RunTime, "run-time", 0, "Stop after the specified amount of time(s), Only used  --autostart. Defaults to run forever.")
	boomCmd.Flags().Int64Var(&boomArgs.LoopCount, "loop-count", -1, "The specify running cycles for load testing")
	boomCmd.Flags().StringVar(&boomArgs.MemoryProfile, "mem-profile", "", "Enable memory profiling.")
	boomCmd.Flags().DurationVar(&boomArgs.MemoryProfileDuration, "mem-profile-duration", 30*time.Second, "Memory profile duration.")
	boomCmd.Flags().StringVar(&boomArgs.CPUProfile, "cpu-profile", "", "Enable CPU profiling.")
	boomCmd.Flags().DurationVar(&boomArgs.CPUProfileDuration, "cpu-profile-duration", 30*time.Second, "CPU profile duration.")
	boomCmd.Flags().StringVar(&boomArgs.PrometheusPushgatewayURL, "prometheus-gateway", "", "Prometheus Pushgateway url.")
	boomCmd.Flags().BoolVar(&boomArgs.DisableConsoleOutput, "disable-console-output", false, "Disable console output.")
	boomCmd.Flags().BoolVar(&boomArgs.DisableCompression, "disable-compression", false, "Disable compression")
	boomCmd.Flags().BoolVar(&boomArgs.DisableKeepalive, "disable-keepalive", false, "Disable keepalive")
	boomCmd.Flags().StringVar(&boomArgs.profile, "profile", "", "profile for load testing")
	boomCmd.Flags().BoolVar(&boomArgs.master, "master", false, "master of distributed testing")
	boomCmd.Flags().StringVar(&boomArgs.masterBindHost, "master-bind-host", "127.0.0.1", "Interfaces (hostname, ip) that hrp master should bind to. Only used when running with --master. Defaults to * (all available interfaces).")
	boomCmd.Flags().IntVar(&boomArgs.masterBindPort, "master-bind-port", 5557, "Port that hrp master should bind to. Only used when running with --master. Defaults to 5557.")
	boomCmd.Flags().StringVar(&boomArgs.masterHttpAddress, "master-http-address", ":9771", "Interfaces (ip:port) that hrp master should control by user. Only used when running with --master. Defaults to *:9771.")
	boomCmd.Flags().BoolVar(&boomArgs.worker, "worker", false, "worker of distributed testing")
	boomCmd.Flags().BoolVar(&boomArgs.ignoreQuit, "ignore-quit", false, "ignores quit from master (only when --worker is used)")
	boomCmd.Flags().StringVar(&boomArgs.masterHost, "master-host", "127.0.0.1", "Host or IP address of hrp master for distributed load testing.")
	boomCmd.Flags().IntVar(&boomArgs.masterPort, "master-port", 5557, "The port to connect to that is used by the hrp master for distributed load testing.")
	boomCmd.Flags().BoolVar(&boomArgs.autoStart, "auto-start", false, "Starts the test immediately. Use --spawn-count and --spawn-rate to control user count and increase rate")
	boomCmd.Flags().IntVar(&boomArgs.expectWorkers, "expect-workers", 1, "How many workers master should expect to connect before starting the test (only when --autostart is used)")
	boomCmd.Flags().IntVar(&boomArgs.expectWorkersMaxWait, "expect-workers-max-wait", 120, "How many workers master should expect to connect before starting the test (only when --autostart is used")
}
