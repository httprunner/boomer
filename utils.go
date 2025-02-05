package boomer

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
)

func round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

// genMD5 returns the md5 hash of strings.
func genMD5(slice ...string) string {
	h := md5.New()
	for _, v := range slice {
		io.WriteString(h, v)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

// startMemoryProfile starts memory profiling and save the results in file.
func startMemoryProfile(file string, duration time.Duration) (err error) {
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	log.Info().Dur("duration", duration).Msg("Start memory profiling")
	time.AfterFunc(duration, func() {
		err := pprof.WriteHeapProfile(f)
		if err != nil {
			log.Error().Err(err).Msg("failed to write memory profile")
		}
		f.Close()
		log.Info().Dur("duration", duration).Msg("Stop memory profiling")
	})
	return nil
}

// startCPUProfile starts cpu profiling and save the results in file.
func startCPUProfile(file string, duration time.Duration) (err error) {
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	log.Info().Dur("duration", duration).Msg("Start CPU profiling")
	err = pprof.StartCPUProfile(f)
	if err != nil {
		f.Close()
		return err
	}

	time.AfterFunc(duration, func() {
		pprof.StopCPUProfile()
		f.Close()
		log.Info().Dur("duration", duration).Msg("Stop CPU profiling")
	})
	return nil
}

// generate a random nodeID like locust does, using the same algorithm.
func getNodeID() (nodeID string) {
	hostname, _ := os.Hostname()
	id := strings.Replace(uuid.NewV4().String(), "-", "", -1)
	nodeID = fmt.Sprintf("%s_%s", hostname, id)
	return
}

// GetCurrentPidCPUUsage get current pid CPU usage
func GetCurrentPidCPUUsage() float64 {
	currentPid := os.Getpid()
	p, err := process.NewProcess(int32(currentPid))
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("failed to get CPU percent\n"))
		return 0.0
	}
	percent, err := p.CPUPercent()
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("failed to get CPU percent\n"))
		return 0.0
	}
	return percent
}

// GetCurrentPidCPUPercent get the percentage of current pid cpu used
func GetCurrentPidCPUPercent() float64 {
	currentPid := os.Getpid()
	p, err := process.NewProcess(int32(currentPid))
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("failed to get CPU percent\n"))
		return 0.0
	}
	percent, err := p.Percent(time.Second)
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("failed to get CPU percent\n"))
		return 0.0
	}
	return percent
}

// GetCurrentCPUPercent get the percentage of current cpu used
func GetCurrentCPUPercent() float64 {
	percent, _ := cpu.Percent(time.Second, false)
	return percent[0]
}

// GetCurrentMemoryPercent get the percentage of current memory used
func GetCurrentMemoryPercent() float64 {
	memInfo, _ := mem.VirtualMemory()
	return memInfo.UsedPercent
}

// GetCurrentPidMemoryUsage get current Memory usage
func GetCurrentPidMemoryUsage() float64 {
	currentPid := os.Getpid()
	p, err := process.NewProcess(int32(currentPid))
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("failed to get CPU percent\n"))
		return 0.0
	}
	percent, err := p.MemoryPercent()
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("failed to get CPU percent\n"))
		return 0.0
	}
	return float64(percent)
}

func Float32ToByte(v float32) []byte {
	bits := math.Float32bits(v)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}

func ByteToFloat32(v []byte) float32 {
	bits := binary.LittleEndian.Uint32(v)
	return math.Float32frombits(bits)
}

func Float64ToByte(v float64) []byte {
	bits := math.Float64bits(v)
	bts := make([]byte, 8)
	binary.LittleEndian.PutUint64(bts, bits)
	return bts
}

func ByteToFloat64(v []byte) float64 {
	bits := binary.LittleEndian.Uint64(v)
	return math.Float64frombits(bits)
}

func Int64ToBytes(n int64) []byte {
	bytesBuf := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytesBuf, binary.BigEndian, n)
	return bytesBuf.Bytes()
}

func BytesToInt64(bys []byte) (data int64) {
	byteBuff := bytes.NewBuffer(bys)
	_ = binary.Read(byteBuff, binary.BigEndian, &data)
	return
}

func SplitInteger(m, n int) (ints []int) {
	quotient := m / n
	remainder := m % n
	if remainder >= 0 {
		for i := 0; i < n-remainder; i++ {
			ints = append(ints, quotient)
		}
		for i := 0; i < remainder; i++ {
			ints = append(ints, quotient+1)
		}
		return
	} else if remainder < 0 {
		for i := 0; i < -remainder; i++ {
			ints = append(ints, quotient-1)
		}
		for i := 0; i < n+remainder; i++ {
			ints = append(ints, quotient)
		}
	}
	return
}

func Bytes2File(data []byte, filename string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0o755)
	defer file.Close()
	if err != nil {
		log.Error().Err(err).Msg("failed to generate file")
	}
	count, err := file.Write(data)
	if err != nil {
		return err
	}
	log.Info().Msg(fmt.Sprintf("write file %s len: %d \n", filename, count))
	return nil
}

func Dump2JSON(data interface{}, path string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		log.Error().Err(err).Msg("convert absolute path failed")
		return err
	}
	log.Info().Str("path", path).Msg("dump data to json")

	// init json encoder
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "    ")

	err = encoder.Encode(data)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, buffer.Bytes(), 0o644)
	if err != nil {
		log.Error().Err(err).Msg("dump json path failed")
		return err
	}
	return nil
}
