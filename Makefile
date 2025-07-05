.PHONY: package clean

# Default target - clean first, then package
all: clean package

# Package target - builds the Go binary and runs dxt pack
package:
	go build
	dxt pack
	# Bug for permissions tracked here:
	# https://github.com/anthropics/dxt/pull/14
	# # Until bug is fixed, need to run:L
	# chmod +x "/Users/bittelc/Library/Application Support/Claude/Claude Extensions/local.dxt.cole-bittel.spotify-mcp-server/spotify-mcp-server"
	cp spotify-mcp-server.dxt '/Users/bittelc/Library/Application Support/Claude/Claude Extensions/local.dxt.cole-bittel.spotify-mcp-server/spotify-mcp-server.dxt'
	cp spotify-mcp-server '/Users/bittelc/Library/Application Support/Claude/Claude Extensions/local.dxt.cole-bittel.spotify-mcp-server/spotify-mcp-server'
	cp -rf * '/Users/bittelc/Library/Application Support/Claude/Claude Extensions/local.dxt.cole-bittel.spotify-mcp-server/'
	chmod +x "/Users/bittelc/Library/Application Support/Claude/Claude Extensions/local.dxt.cole-bittel.spotify-mcp-server/spotify-mcp-server"
	chmod +x "/Users/bittelc/Library/Application Support/Claude/Claude Extensions/local.dxt.cole-bittel.spotify-mcp-server/spotify-mcp-server.dxt"

# Optional clean target to remove built artifacts
clean:
	rm -f spotify-mcp-server
	rm -f spotify-mcp-server.dxt
