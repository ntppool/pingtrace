package traceroute

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net/netip"
	"os/exec"
)

type Traceroute struct {
	ip       netip.Addr
	sourceIP netip.Addr
	ctx      context.Context
	pipe     io.ReadCloser
	// ch   chan TraceRouteLine
	trp *TracerouteParser
}

func New(ip, sourceIP netip.Addr) (*Traceroute, error) {
	if !ip.IsValid() {
		return nil, fmt.Errorf("invalid IP")
	}
	return &Traceroute{
		ip:       ip,
		sourceIP: sourceIP,
		// ch: make(chan TraceRouteLine),
		trp: NewTracerouteParser(),
	}, nil
}

func Run(ctx context.Context, ip, sourceIP netip.Addr) (*TraceRouteResult, error) {
	tr, err := New(ip, sourceIP)
	if err != nil {
		return nil, err
	}
	err = tr.Start(ctx)
	if err != nil {
		return nil, err
	}
	lines, err := tr.ReadAll()
	if err != nil {
		return nil, err
	}
	return &TraceRouteResult{
		TargetIP: ip,
		Lines:    lines,
	}, nil

}

func (tr *Traceroute) Start(ctx context.Context) error {

	tr.ctx = ctx

	args := []string{"-q", "2", "-w", "3", "-n"}

	if tr.sourceIP.IsValid() {
		args = append(args, "-s", tr.sourceIP.String())
	}

	args = append(args, tr.ip.String())

	cmd := exec.CommandContext(ctx, "traceroute", args...)
	// cmd := exec.Command("./slowly.sh", "5")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	tr.pipe = stdout

	go tr.runner(cmd)

	return nil
}

func (tr *Traceroute) ReadAll() ([]*TraceRouteLine, error) {
	lines := make([]*TraceRouteLine, 0)
	for {
		trl, err := tr.Read()

		if trl != nil {
			lines = append(lines, trl)
		}

		if trl == nil || err != nil {
			return lines, err
		}
	}
}

func (tr *Traceroute) Read() (*TraceRouteLine, error) {
	select {
	case trl, ok := <-tr.trp.out:
		if !ok {
			return nil, nil
		}
		return &trl, nil
	case <-tr.ctx.Done():
		return nil, nil
	}
}

func (tr *Traceroute) runner(cmd *exec.Cmd) {

	r := bufio.NewReader(tr.pipe)
	defer tr.trp.Close()

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading from traceroute pipe: ", err)
			}
			break
		}

		tr.trp.Add(line)
		if err != nil {
			log.Printf("Could not parse '%s': %s", line, err)
			continue
		}
	}

	cmdRV := cmd.Wait()
	if cmdRV != nil {
		err := cmdRV.Error()
		if err != "signal: killed" {
			log.Printf("Error finishing command: %s", cmdRV.Error())
		}
		tr.trp.out <- TraceRouteLine{Err: fmt.Errorf("traceroute error: %s", cmdRV.Error())}
	}

}
