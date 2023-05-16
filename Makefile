BUILD_DIR := build
MODELS_DIR := models
INCLUDE_PATH := $(abspath ./whisper.cpp)
LIBRARY_PATH := $(abspath ./whisper.cpp)

all: clean whisper modtidy speech-to-text go-model-download 

whisper: mkdir
	@echo Build whisper
	@${MAKE} -C ./whisper.cpp libwhisper.a

speech-to-text:
	@C_INCLUDE_PATH=${INCLUDE_PATH} LIBRARY_PATH=${LIBRARY_PATH} go build ${BUILD_FLAGS} -o ${BUILD_DIR}/$(notdir $@) ./main.go

go-model-download: mkdir whisper modtidy
	@echo Build exe 'go-model-download' 
	@C_INCLUDE_PATH=${INCLUDE_PATH} LIBRARY_PATH=${LIBRARY_PATH} go build ${BUILD_FLAGS} -o ${BUILD_DIR}/$(notdir $@) ./models/go-model-download

download: mkdir go-model-download
	@echo Download base model 'ggml-base.en.bin' 
	@${BUILD_DIR}/go-model-download -out models ggml-base.en.bin

test: model-base whisper modtidy
	@C_INCLUDE_PATH=${INCLUDE_PATH} LIBRARY_PATH=${LIBRARY_PATH} go test -v .

mkdir:
	@echo Mkdir ${BUILD_DIR}
	@install -d ${BUILD_DIR}
	@echo Mkdir ${MODELS_DIR}
	@install -d ${MODELS_DIR}

modtidy:
	@go mod tidy

clean: 
	@echo Clean
	@rm -fr $(BUILD_DIR)
	@go clean
