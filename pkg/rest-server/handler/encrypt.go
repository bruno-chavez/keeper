package handler

import (
	"io"
	"log"
	"net/http"
)

const wrongContentTypeMessage = "wrong content type, please use 'text/plain'"
const parsingErrorMessage = "error parsing content, body does not has the appropriate format"
const invalidContextMessage = "invalid context query parameter, expected comma-separated key:value pairs (e.g. context=purpose:develop-freeze)"
const encryptionErrorMessage = "error encrypting content"

func (h Handler) Encrypt(w http.ResponseWriter, request *http.Request) {
	defer func() {
		err := request.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	if request.Header.Get("Content-Type") != "text/plain" {
		err := writeResponse(w, http.StatusBadRequest, wrongContentTypeMessage)
		if err != nil {
			log.Println(err)
		}
		return
	}

	encCtx, err := parseEncryptionContext(request.URL.Query().Get("context"))
	if err != nil {
		log.Printf("error parsing encryption context: %v\n", err)
		err := writeResponse(w, http.StatusBadRequest, invalidContextMessage)
		if err != nil {
			log.Println(err)
		}
		return
	}

	data, err := io.ReadAll(request.Body)
	if err != nil {
		log.Printf("error parsing request body: %v\n", err)
		err := writeResponse(w, http.StatusBadRequest, parsingErrorMessage)
		if err != nil {
			log.Println(err)
		}
		return
	}

	encryptedData, err := h.encryptor.Encrypt(request.Context(), string(data), encCtx)
	if err != nil {
		log.Printf("error encrypting content: %v\n", err)
		err := writeResponse(w, http.StatusBadRequest, encryptionErrorMessage)
		if err != nil {
			log.Println(err)
		}
		return
	}

	err = writeResponse(w, http.StatusOK, encryptedData)
	if err != nil {
		log.Println(err)
	}
}
