.PHONY: build test build-skill clean

# Go 构建
BIN_DIR    := cli/bin
BIN_NAME   := bangumi
SKILL_DIR  := skill
SKILL_BIN  := $(SKILL_DIR)/bangumi

VERSION := $(shell cd cli && git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT  := $(shell cd cli && git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE    := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || echo "unknown")

LDFLAGS := -s -w \
	-X 'main.version=$(VERSION)' \
	-X 'main.commit=$(COMMIT)' \
	-X 'main.date=$(DATE)' \
	-X 'github.com/kasuganosora/bangumi.skill/cli/api.version=$(VERSION)'

# 构建 CLI 二进制到 cli/bin/
build:
	cd cli && go build -ldflags "$(LDFLAGS)" -o bin/bangumi .

# 运行测试
test:
	cd cli && go test ./... -count=1 -timeout=60s

# 构建 skill 包：复制二进制到 skill/ 目录
build-skill: build
	@mkdir -p $(SKILL_DIR)
	@cp $(BIN_DIR)/bangumi $(SKILL_DIR)/ 2>/dev/null || cp $(BIN_DIR)/bangumi.exe $(SKILL_DIR)/bangumi.exe 2>/dev/null || echo "Warning: binary not found"
	@echo "Skill package ready: $(SKILL_DIR)/"

# 清理
clean:
	rm -rf $(BIN_DIR)
	rm -f $(SKILL_DIR)/bangumi $(SKILL_DIR)/bangumi.exe
