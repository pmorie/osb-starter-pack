// Package pkg in the `osb-broker-lib` project provides building blocks to build an
// Open Service Broker.
//
// Main components of this library are:
//
// - broker/Interface: the interface containing the callbacks for the OSB
//   API actions
// - rest/APISurface: provides the glue between the OSB REST API and
//   your broker.Interface
// - server/Server: the HTTP server for the OSB and metrics APIs
package pkg
