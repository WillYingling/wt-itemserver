CC=go build
SRC=*.go
BIN_DIR=bin
SERVER_BIN=webserver

all: $(BIN_DIR)/$(SERVER_BIN)

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

$(BIN_DIR)/$(SERVER_BIN): $(BIN_DIR) $(SRC)
	$(CC) -o $(BIN_DIR)/$(SERVER_BIN)

clean:
	rm -rf $(BIN_DIR)
	rm go.sum
