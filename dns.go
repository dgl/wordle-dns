package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/miekg/dns"
)

const (
	SOA = `wd.ip.wtf SOA wd-ns.ip.wtf dns-admin.oo.fail 42 42 42 42 42`
	NS  = `wd.ip.wtf NS wd-ns.ip.wtf`
	TXT = ` TXT ""`
)

var (
	flagName   = flag.String("name", "DNS name this server is accessible on", "ns.example.com")
	flagRname  = flag.String("rname", "DNS name for administrator (email with '@' replaced with '.')", "")
	flagListen = flag.String("listen", "DNS listen address", ":5300")
)

// MustNewRR is a shortcut to dns.NewRR that panics on error.
func MustNewRR(s string) dns.RR {
	r, err := dns.NewRR(s)
	if err != nil {
		panic(err)
	}
	return r
}

var wordles = map[string]string{
	"0": "zeros",

	"1":  "amaze",
	"2":  "magic",
	"3":  "house",
	"4":  "cloud",
	"5":  "nodes",
	"6":  "paint",
	"7":  "disco",
	"8":  "hedge",
	"9":  "poems",
	"10": "flair",
	"11": "batch",
	"12": "stout",
	"13": "candy",
	"14": "larva",
	"15": "maple",
	"16": "ladle",
	"17": "bacon",
	"18": "snake",
	"19": "chart",
	"20": "oasis",
	"21": "space",
	"22": "basic",
	"23": "slope",
	"24": "great",
	"25": "first",
	"26": "rhino",
	"27": "screw",
	"28": "force",
	"29": "water",
	"30": "hyper",
	"31": "cakes",

	"example": "names",
}

func dnsServe() {
	dns.HandleFunc("wd.ip.wtf", func(w dns.ResponseWriter, r *dns.Msg) {
		m := new(dns.Msg)
		m.SetReply(r)
		m.Authoritative = true

		rname := strings.Split(r.Question[0].Name, ".")
		name := strings.TrimSuffix(r.Question[0].Name, ".")
		day := fmt.Sprintf("%d", 1+((12 + time.Now().UTC().Day()) % 31))

		if len(rname) == 4 && strings.ToLower(name) == "wd.ip.wtf" && r.Question[0].Qtype == dns.TypeNS {
			m.Answer = []dns.RR{MustNewRR(NS)}
		} else if len(rname) == 4 && strings.ToLower(name) == "wd.ip.wtf" && r.Question[0].Qtype == dns.TypeSOA {
			m.Answer = []dns.RR{MustNewRR(SOA)}
		} else if len(rname[0]) != 5 && len(rname) <= 5 {
			rr := MustNewRR(name + TXT)
			txt := rr.(*dns.TXT)
			txt.Txt[0] = `Welcome to Wordle over DNS! Today's puzzle is #` + day + `: <guess>.` + day + `.wd.ip.wtf`
			rr2 := MustNewRR(name + TXT)
			txt2 := rr2.(*dns.TXT)
			txt2.Txt = []string{
				`This shell function makes it easier to play`,
				`wd() { dig +short txt $1.` + day + `.wd.ip.wtf | perl -pe's/\\([0-9]{1,3})/chr$1/eg'; }`}
			m.Answer = []dns.RR{rr, rr2}
		} else if len(rname[0]) == 5 && len(rname) <= 5 {
			rr := MustNewRR(name + " CNAME " + rname[0] + "." + day + ".wd.ip.wtf.")
			m.Answer = []dns.RR{rr}
		} else if word, ok := wordles[rname[1]]; ok {
			result := wordle(strings.ToLower(rname[0]), word)

			if result == "" || !checkGuessValid(rname[0]) {
				result = fmt.Sprintf("Your guess must be a valid word and %d letters.", len(word))
			}

			rr := MustNewRR(name + TXT)
			txt := rr.(*dns.TXT)
			txt.Txt[0] = result
			m.Answer = []dns.RR{rr}
			if result == "🟩🟩🟩🟩🟩" {
				rr := MustNewRR(name + TXT)
				txt := rr.(*dns.TXT)
				txt.Txt = []string{
					"🥳🎉🥳🎉🪅",
				}
				m.Answer = append(m.Answer, rr)
			}
		} else {
			m.Rcode = dns.RcodeNameError
			m.Ns = []dns.RR{MustNewRR(SOA)}
		}
		w.WriteMsg(m)
	})

	go func() {
		srv := &dns.Server{Addr: *flagListen, Net: "udp"}
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Failed to set udp listener %s\n", err.Error())
		}
	}()

	go func() {
		srv := &dns.Server{Addr: *flagListen, Net: "tcp"}
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Failed to set tcp listener %s\n", err.Error())
		}
	}()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	log.Fatalf("Signal (%v) received, stopping\n", s)
}

func main() {
	flag.Parse()
	dictLoad(5)
	dnsServe()
}
