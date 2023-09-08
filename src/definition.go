package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
)

// TODO delete unnecessary fields

type valueAmendment func(any) (any, error)

func millisToSeconds(v any) (any, error) {
	asFloat, err := strconv.ParseFloat(fmt.Sprintf("%v", v), 64)
	if err != nil {
		return v, fmt.Errorf("expected numeric value: %w", err)
	}
	return (time.Duration(asFloat) * time.Millisecond).Seconds(), nil
}

type metricDefinition struct {
	MetricName string
	SourceType metric.SourceType
	Amend      valueAmendment
}

func (m *metricDefinition) value(metricValue any) (any, error) {
	if m.Amend != nil {
		return m.Amend(metricValue)
	}
	return metricValue, nil
}

// HAProxyFrontendStats holds the metric definitions for a frontend
var HAProxyFrontendStats = map[string]metricDefinition{
	"pxname":        {MetricName: "frontend.proxyName", SourceType: metric.ATTRIBUTE},
	"svname":        {MetricName: "frontend.serviceName", SourceType: metric.ATTRIBUTE},
	"scur":          {MetricName: "frontend.currentSessions", SourceType: metric.GAUGE},
	"smax":          {MetricName: "frontend.maxSessions", SourceType: metric.GAUGE},
	"stot":          {MetricName: "frontend.sessionsPerSecond", SourceType: metric.PRATE},
	"bin":           {MetricName: "frontend.bytesInPerSecond", SourceType: metric.PRATE},
	"bout":          {MetricName: "frontend.bytesOutPerSecond", SourceType: metric.PRATE},
	"dreq":          {MetricName: "frontend.requestsDenied.securityConcernsPerSecond", SourceType: metric.PRATE},
	"dresp":         {MetricName: "frontend.responsesDenied.securityConcernsPerSecond", SourceType: metric.PRATE},
	"ereq":          {MetricName: "frontend.requestErrorsPerSecond", SourceType: metric.PRATE},
	"status":        {MetricName: "frontend.status", SourceType: metric.ATTRIBUTE},
	"type":          {MetricName: "frontend.type", SourceType: metric.GAUGE},
	"rate_max":      {MetricName: "frontend.maxSessionsPerSecond", SourceType: metric.GAUGE},
	"hrsp_1xx":      {MetricName: "frontend.http100ResponsesPerSecond", SourceType: metric.PRATE},
	"hrsp_2xx":      {MetricName: "frontend.http200ResponsesPerSecond", SourceType: metric.PRATE},
	"hrsp_3xx":      {MetricName: "frontend.http300ResponsesPerSecond", SourceType: metric.PRATE},
	"hrsp_4xx":      {MetricName: "frontend.http400ResponsesPerSecond", SourceType: metric.PRATE},
	"hrsp_5xx":      {MetricName: "frontend.http500ResponsesPerSecond", SourceType: metric.PRATE},
	"hrsp_other":    {MetricName: "frontend.httpOtherResponsesPerSecond", SourceType: metric.PRATE},
	"req_rate_max":  {MetricName: "frontend.httpRequests.maxPerSecond", SourceType: metric.PRATE},
	"req_tot":       {MetricName: "frontend.httpRequestsPerSecond", SourceType: metric.PRATE},
	"mode":          {MetricName: "frontend.mode", SourceType: metric.ATTRIBUTE},
	"conn_rate_max": {MetricName: "frontend.maxConnectionsPerSecond", SourceType: metric.GAUGE},
	"conn_tot":      {MetricName: "frontend.connectionsPerSecond", SourceType: metric.PRATE},
	"intercepted":   {MetricName: "frontend.interceptedRequestsPerSecond", SourceType: metric.PRATE},
	"dcon":          {MetricName: "frontend.requestsDenied.tcpRequestConnectionRulesPerSecond", SourceType: metric.PRATE},
	"dses":          {MetricName: "frontend.requestsDenied.tcpRequestSessionRulesPerSecond", SourceType: metric.PRATE},
}

// HAProxyBackendStats holds the metric definitions for a backend
var HAProxyBackendStats = map[string]metricDefinition{
	"pxname":      {MetricName: "backend.proxyName", SourceType: metric.ATTRIBUTE},
	"qcur":        {MetricName: "backend.currentQueuedRequestsWithoutServer", SourceType: metric.GAUGE},
	"qmax":        {MetricName: "backend.maxQueuedRequestsWithoutServer", SourceType: metric.GAUGE},
	"scur":        {MetricName: "backend.currentSessions", SourceType: metric.GAUGE},
	"smax":        {MetricName: "backend.maxSessions", SourceType: metric.GAUGE},
	"stot":        {MetricName: "backend.sessionsPerSecond", SourceType: metric.PRATE},
	"bin":         {MetricName: "backend.bytesInPerSecond", SourceType: metric.PRATE},
	"bout":        {MetricName: "backend.bytesOutPerSecond", SourceType: metric.PRATE},
	"dreq":        {MetricName: "backend.requestsDenied.securityConcernsPerSecond", SourceType: metric.PRATE},
	"dresp":       {MetricName: "backend.responsesDenied.securityConcernsPerSecond", SourceType: metric.PRATE},
	"econ":        {MetricName: "backend.connectingRequestErrorsPerSecond", SourceType: metric.PRATE},
	"eresp":       {MetricName: "backend.responseErrorsPerSecond", SourceType: metric.PRATE},
	"wretr":       {MetricName: "backend.connectionRetriesPerSecond", SourceType: metric.PRATE},
	"wredis":      {MetricName: "backend.requestRedispatchPerSecond", SourceType: metric.PRATE},
	"status":      {MetricName: "backend.status", SourceType: metric.ATTRIBUTE},
	"weight":      {MetricName: "backend.totalWeight", SourceType: metric.GAUGE},
	"act":         {MetricName: "backend.activeServers", SourceType: metric.GAUGE},
	"bck":         {MetricName: "backend.backupServers", SourceType: metric.GAUGE},
	"chkdown":     {MetricName: "backend.upToDownTransitionsPerSecond", SourceType: metric.PRATE},
	"lastchg":     {MetricName: "backend.timeSinceLastUpDownTransitionInSeconds", SourceType: metric.GAUGE},
	"downtime":    {MetricName: "backend.downtimeInSeconds", SourceType: metric.GAUGE},
	"lbtot":       {MetricName: "backend.serverSelectedPerSecond", SourceType: metric.PRATE},
	"type":        {MetricName: "backend.type", SourceType: metric.GAUGE},
	"rate_max":    {MetricName: "backend.maxSessionsPerSecond", SourceType: metric.GAUGE},
	"hrsp_1xx":    {MetricName: "backend.http100ResponsesPerSecond", SourceType: metric.PRATE},
	"hrsp_2xx":    {MetricName: "backend.http200ResponsesPerSecond", SourceType: metric.PRATE},
	"hrsp_3xx":    {MetricName: "backend.http300ResponsesPerSecond", SourceType: metric.PRATE},
	"hrsp_4xx":    {MetricName: "backend.http400ResponsesPerSecond", SourceType: metric.PRATE},
	"hrsp_5xx":    {MetricName: "backend.http500ResponsesPerSecond", SourceType: metric.PRATE},
	"hrsp_other":  {MetricName: "backend.httpOtherResponsesPerSecond", SourceType: metric.PRATE},
	"req_tot":     {MetricName: "backend.httpRequestsPerSecond", SourceType: metric.PRATE},
	"cli_abrt":    {MetricName: "backend.dataTransfersAbortedByClientPerSecond", SourceType: metric.PRATE},
	"srv_abrt":    {MetricName: "backend.dataTransfersAbortedByServerPerSecond", SourceType: metric.PRATE},
	"comp_in":     {MetricName: "backend.httpResponseBytesFedToCompressorPerSecond", SourceType: metric.PRATE},
	"comp_out":    {MetricName: "backend.httpResponseBytesEmittedByCompressorPerSecond", SourceType: metric.PRATE},
	"comp_byp":    {MetricName: "backend.bytesThatBypassedCompressorPerSecond", SourceType: metric.PRATE},
	"comp_rsp":    {MetricName: "backend.httpResponsesCompressedPerSecond", SourceType: metric.PRATE},
	"lastsess":    {MetricName: "backend.timeSinceLastSessionAssignedInSeconds", SourceType: metric.GAUGE},
	"qtime":       {MetricName: "backend.averageQueueTimeInSeconds", SourceType: metric.GAUGE, Amend: millisToSeconds},
	"ctime":       {MetricName: "backend.averageConnectTimeInSeconds", SourceType: metric.GAUGE, Amend: millisToSeconds},
	"rtime":       {MetricName: "backend.averageResponseTimeInSeconds", SourceType: metric.GAUGE, Amend: millisToSeconds},
	"ttime":       {MetricName: "backend.averageTotalSessionTimeInSeconds", SourceType: metric.GAUGE, Amend: millisToSeconds},
	"cookie":      {MetricName: "backend.cookieName", SourceType: metric.ATTRIBUTE},
	"mode":        {MetricName: "backend.mode", SourceType: metric.ATTRIBUTE},
	"intercepted": {MetricName: "backend.interceptedRequestsPerSecond", SourceType: metric.PRATE},
}

// HAProxyServerStats holds the metric definitions for a server
var HAProxyServerStats = map[string]metricDefinition{
	"pxname":         {MetricName: "server.proxyName", SourceType: metric.ATTRIBUTE},
	"svname":         {MetricName: "server.serviceName", SourceType: metric.ATTRIBUTE},
	"qcur":           {MetricName: "server.currentQueuedRequestsWithoutServer", SourceType: metric.GAUGE},
	"qmax":           {MetricName: "server.maxQueuedRequestsWithoutServer", SourceType: metric.GAUGE},
	"scur":           {MetricName: "server.currentSessions", SourceType: metric.GAUGE},
	"smax":           {MetricName: "server.maxSessions", SourceType: metric.GAUGE},
	"stot":           {MetricName: "server.sessionsPerSecond", SourceType: metric.PRATE},
	"bin":            {MetricName: "server.bytesInPerSecond", SourceType: metric.PRATE},
	"bout":           {MetricName: "server.bytesOutPerSecond", SourceType: metric.PRATE},
	"dreq":           {MetricName: "server.requestsDenied.securityConcernsPerSecond", SourceType: metric.PRATE},
	"dresp":          {MetricName: "server.responsesDenied.securityConcernsPerSecond", SourceType: metric.PRATE},
	"econ":           {MetricName: "server.connectingRequestErrorsPerSecond", SourceType: metric.PRATE},
	"eresp":          {MetricName: "server.responseErrorsPerSecond", SourceType: metric.PRATE},
	"wretr":          {MetricName: "server.connectionRetriesPerSecond", SourceType: metric.PRATE},
	"wredis":         {MetricName: "server.requestRedispatchPerSecond", SourceType: metric.PRATE},
	"status":         {MetricName: "server.status", SourceType: metric.ATTRIBUTE},
	"weight":         {MetricName: "server.serverWeight", SourceType: metric.GAUGE},
	"act":            {MetricName: "server.isActive", SourceType: metric.GAUGE},
	"bck":            {MetricName: "server.isBackup", SourceType: metric.GAUGE},
	"chkfail":        {MetricName: "server.failedChecksPerSecond", SourceType: metric.PRATE},
	"chkdown":        {MetricName: "server.upToDownTransitionsPerSecond", SourceType: metric.PRATE},
	"lastchg":        {MetricName: "server.timeSinceLastUpDownTransitionInSeconds", SourceType: metric.GAUGE},
	"downtime":       {MetricName: "server.downtimeInSeconds", SourceType: metric.GAUGE},
	"sid":            {MetricName: "server.serverID", SourceType: metric.ATTRIBUTE},
	"throttle":       {MetricName: "server.throttlePercentage", SourceType: metric.GAUGE},
	"lbtot":          {MetricName: "server.serverSelectedPerSecond", SourceType: metric.PRATE},
	"type":           {MetricName: "server.type", SourceType: metric.GAUGE},
	"rate_max":       {MetricName: "server.maxSessionsPerSecond", SourceType: metric.GAUGE},
	"check_status":   {MetricName: "server.healthCheckStatus", SourceType: metric.ATTRIBUTE},
	"check_code":     {MetricName: "server.layerCode", SourceType: metric.ATTRIBUTE},
	"check_duration": {MetricName: "server.healthCheckDurationInMilliseconds", SourceType: metric.GAUGE},
	"hrsp_1xx":       {MetricName: "server.http100ResponsesPerSecond", SourceType: metric.PRATE},
	"hrsp_2xx":       {MetricName: "server.http200ResponsesPerSecond", SourceType: metric.PRATE},
	"hrsp_3xx":       {MetricName: "server.http300ResponsesPerSecond", SourceType: metric.PRATE},
	"hrsp_4xx":       {MetricName: "server.http400ResponsesPerSecond", SourceType: metric.PRATE},
	"hrsp_5xx":       {MetricName: "server.http500ResponsesPerSecond", SourceType: metric.PRATE},
	"hrsp_other":     {MetricName: "server.httpOtherResponsesPerSecond", SourceType: metric.PRATE},
	"hanafail":       {MetricName: "server.failedHealthCheckDetails", SourceType: metric.ATTRIBUTE},
	"cli_abrt":       {MetricName: "server.dataTransfersAbortedByClientPerSecond", SourceType: metric.PRATE},
	"srv_abrt":       {MetricName: "server.dataTransfersAbortedByServerPerSecond", SourceType: metric.PRATE},
	"lastsess":       {MetricName: "server.timeSinceLastSessionAssignedInSeconds", SourceType: metric.GAUGE},
	"last_chk":       {MetricName: "server.healthCheckContents", SourceType: metric.ATTRIBUTE},
	"last_agt":       {MetricName: "server.agentCheckContents", SourceType: metric.ATTRIBUTE},
	"qtime":          {MetricName: "server.averageQueueTimeInSeconds", SourceType: metric.GAUGE, Amend: millisToSeconds},
	"ctime":          {MetricName: "server.averageConnectTimeInSeconds", SourceType: metric.GAUGE, Amend: millisToSeconds},
	"rtime":          {MetricName: "server.averageResponseTimeInSeconds", SourceType: metric.GAUGE, Amend: millisToSeconds},
	"ttime":          {MetricName: "server.averageTotalSessionTimeInSeconds", SourceType: metric.GAUGE, Amend: millisToSeconds},
	"agent_status":   {MetricName: "server.agentStatus", SourceType: metric.ATTRIBUTE},
	"agent_duration": {MetricName: "server.agentDurationSeconds", SourceType: metric.GAUGE, Amend: millisToSeconds},
	"check_desc":     {MetricName: "server.checkStatusDescription", SourceType: metric.ATTRIBUTE},
	"agent_desc":     {MetricName: "server.agentStatusDescription", SourceType: metric.ATTRIBUTE},
	"cookie":         {MetricName: "server.cookieValue", SourceType: metric.ATTRIBUTE},
	"mode":           {MetricName: "server.mode", SourceType: metric.ATTRIBUTE},
}

// HAProxyInventory holds the list of inventory items to be collected from the stats request
var HAProxyInventory = map[string]struct{}{
	"slim":     {},
	"pid":      {},
	"sid":      {},
	"iid":      {},
	"rate_lim": {},
	"qmax":     {},
}
