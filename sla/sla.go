/*
Real-time Online/Offline Charging System (OCS) for Telecom & ISP environments
Copyright (C) ITsysCOM GmbH

*/

package sla

import (
	"time"

	"github.com/cgrates/cgrates/cores"
	"github.com/cgrates/cgrates/utils"
)

const (
	CustomerID   = "WeKall"
	EncryptKey   = "CGRateS.org!3010"
	CapsLimit    = 30
	CapsStrategy = utils.MetaBusy
	SlaSv1Info   = "SlaSv1.Info"
)

func init() {
	utils.ConcurrentReqsLimit = CapsLimit // sets the system wide limit to the customer one
	utils.ConcurrentReqsStrategy = CapsStrategy
}

// SLAInfo is retrieved via API call
type SlaInfo struct {
	CustomerID   string
	CapsLicensed int
	CapsCurrent  int
	//CapsAverage  int
	Time string
	Hash string
}

func NewSlaSv1(cncReqs *cores.Caps) *SlaSv1 {
	return &SlaSv1{cncReqs: cncReqs}
}

// SlaSv1 exports RPC
type SlaSv1 struct {
	cncReqs *cores.Caps
}

// Call implements rpcclient.ClientConnector interface for internal RPC
func (slav1 *SlaSv1) Call(serviceMethod string,
	args interface{}, reply interface{}) error {
	return utils.APIerRPCCall(slav1, serviceMethod, args, reply)
}

// Ping is used for health checks
func (slav1 *SlaSv1) Ping(ign struct{}, reply *string) (err error) {
	*reply = utils.Pong
	return
}

// Info returns the license information
func (slav1 *SlaSv1) Info(ign struct{}, info *SlaInfo) (err error) {
	used := slav1.cncReqs.Allocated()
	tm := time.Now().String()
	var hsh string
	if hsh, err = utils.ComputeHash(EncryptKey, tm); err != nil {
		return utils.NewErrServerError(err)
	}
	*info = SlaInfo{
		CustomerID:   CustomerID,
		CapsLicensed: CapsLimit,
		CapsCurrent:  used,
		Time:         tm,
		Hash:         hsh,
	}
	return
}

type ArgVerifyHash struct {
	Hash string
	Keys []string
}

// Info returns the license information
func (slav1 *SlaSv1) VerifyHash(arg ArgVerifyHash, match *bool) (err error) {
	keys := append([]string{EncryptKey}, arg.Keys...)
	*match = utils.VerifyHash(arg.Hash, keys...)
	return
}
