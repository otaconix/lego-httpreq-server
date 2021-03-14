package main

import (
	"encoding/json"
  "fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns"
)

const (
  PRESENT CallType = 0
  CLEANUP CallType = 1
)

type CallType int

func (c *CallType) Format() string {
  if *c == PRESENT {
    return "presenting"
  } else {
    return "cleaning up"
  }
}

type messageType struct {
	Domain  string `json:"domain"`
	Token   string `json:"token"`
	KeyAuth string `json:"keyAuth"`
}

type LegoHttpreqServer struct {
  DnsProvider challenge.Provider
}

func EnvOrDefault(name string, defaultValue string) (value string) {
  value, isPresent := os.LookupEnv(name)

  if !isPresent {
    value = defaultValue
  }

  return
}

func readMessage(r io.Reader) (messageType, error) {
  message := messageType {}
  err := json.NewDecoder(r).Decode(&message)

  return message, err
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(200)
}

func determineCallType(r *http.Request) (CallType, error) {
  if r.URL.Path == "/present" {
    return PRESENT, nil
  } else if r.URL.Path == "/cleanup" {
    return CLEANUP, nil
  }

  return 0, fmt.Errorf("Could not determine call type for path: %s", r.URL.Path)
}

func (s *LegoHttpreqServer) Handler(w http.ResponseWriter, r *http.Request) {
  callType, err := determineCallType(r)
  if err != nil {
    w.WriteHeader(400)
    log.Println(err)
  }

  message, err := readMessage(r.Body)

  if err != nil {
    w.WriteHeader(400)
    return
  }

  if callType == PRESENT {
    err = s.DnsProvider.Present(message.Domain, message.Token, message.KeyAuth)
  } else if callType == CLEANUP {
    err = s.DnsProvider.CleanUp(message.Domain, message.Token, message.KeyAuth)
  }

  if err != nil {
    w.WriteHeader(500)
    log.Println("There was an error", callType.Format(), err)
    return
  }

  log.Println("Successfully finished", callType.Format(), "for domain", message.Domain)

  w.WriteHeader(200)
}

func main() {
  providerName, providerNamePresent := os.LookupEnv("DNS_PROVIDER")
  if !providerNamePresent {
    log.Fatal("DNS_PROVIDER environment variable absent")
    return
  }

  provider, error := dns.NewDNSChallengeProviderByName(providerName)
  
  if error != nil {
    log.Fatal(error)
  }

  log.Println("Created provider with name:", providerName)

  server := LegoHttpreqServer {
    DnsProvider: provider,
  }

  http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/present", server.Handler)
	http.HandleFunc("/cleanup", server.Handler)

  listenAddress := EnvOrDefault("LISTEN_ADDRESS", ":8080")
  log.Println("Starting server on address:", listenAddress)

	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
