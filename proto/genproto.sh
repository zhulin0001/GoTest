#!/bin/sh
protoc -o account.pb account.proto
protoc -o cmd.pb cmd.proto