# boomer

This module is initially forked from [myzhan/boomer@v1.6.0] and made a lot of changes.

[myzhan/boomer@v1.6.0]: https://github.com/myzhan/boomer/tree/v1.6.0

```shell
$ go run ./hrp/cli -h
run yaml/json testcase files for load test

Usage:
  hrpboom [flags]

Examples:
  $ hrpboom demo.json   # run specified json testcase file
  $ hrpboom demo.yaml   # run specified yaml testcase file
  $ hrpboom examples/   # run testcases in specified folder

Flags:
      --auto-start                      Starts the test immediately. Use --spawn-count and --spawn-rate to control user count and increase rate
      --cpu-profile string              Enable CPU profiling.
      --cpu-profile-duration duration   CPU profile duration. (default 30s)
      --disable-compression             Disable compression
      --disable-console-output          Disable console output.
      --disable-keepalive               Disable keepalive
      --expect-workers int              How many workers master should expect to connect before starting the test (only when --autostart is used) (default 1)
      --expect-workers-max-wait int     How many workers master should expect to connect before starting the test (only when --autostart is used (default 120)
  -h, --help                            help for boom
      --ignore-quit                     ignores quit from master (only when --worker is used)
      --loop-count int                  The specify running cycles for load testing (default -1)
      --master                          master of distributed testing
      --master-bind-host string         Interfaces (hostname, ip) that hrp master should bind to. Only used when running with --master. Defaults to * (all available interfaces). (default "127.0.0.1")
      --master-bind-port int            Port that hrp master should bind to. Only used when running with --master. Defaults to 5557. (default 5557)
      --master-host string              Host or IP address of hrp master for distributed load testing. (default "127.0.0.1")
      --master-http-address string      Interfaces (ip:port) that hrp master should control by user. Only used when running with --master. Defaults to *:9771. (default ":9771")
      --master-port int                 The port to connect to that is used by the hrp master for distributed load testing. (default 5557)
      --max-rps int                     Max RPS that boomer can generate, disabled by default.
      --mem-profile string              Enable memory profiling.
      --mem-profile-duration duration   Memory profile duration. (default 30s)
      --profile string                  profile for load testing
      --prometheus-gateway string       Prometheus Pushgateway url.
      --request-increase-rate string    Request increase rate, disabled by default. (default "-1")
      --run-time int                    Stop after the specified amount of time(s), Only used  --autostart. Defaults to run forever.
      --spawn-count int                 The number of users to spawn for load testing (default 1)
      --spawn-rate float                The rate for spawning users (default 1)
      --worker                          worker of distributed testing
```