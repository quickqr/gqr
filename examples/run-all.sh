#!/bin/bash

for f in $(ls *.go); do go run "$f"; done