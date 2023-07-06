/*
 * MIT License
 *
 * Copyright (c) 2022 Anton (stremovskyy) Stremovskyy <stremovskyy@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package manager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/meklis/gofondy/consts"
	"github.com/meklis/gofondy/models"
)

type v1Client struct {
	client  *http.Client
	options *ClientOptions
}

func newV1Client(client *http.Client, options *ClientOptions) *v1Client {
	return &v1Client{client: client, options: options}
}
func (m *v1Client) do(url consts.FondyURL, request *models.FondyRequestObject, credit bool, merchantAccount *models.MerchantAccount, reservationData *models.ReservationData) (*[]byte, error) {
	requestID := uuid.New().String()
	methodPost := "POST"

	if reservationData != nil {
		request.ReservationData = reservationData.Base64Encoded()
	}

	if m.options.IsDebug {
		log.Printf("[GO FONDY] Request ID: %v\n", requestID)
		log.Printf("[GO FONDY] URL: %v\n", url.String())
		log.Printf("[GO FONDY] Reservation data: %v\n", reservationData)
	}

	var key string
	if credit {
		key = merchantAccount.MerchantCreditKey
	} else {
		key = merchantAccount.MerchantKey
	}

	err := request.Sign(key, m.options.IsDebug)
	if err != nil {
		return nil, fmt.Errorf("cannot sign request: %v", err)
	}

	jsonValue, err := json.Marshal(models.NewFondyRequest(request))
	if err != nil {
		return nil, fmt.Errorf("cannot marshal request: %w", err)
	}

	if m.options.IsDebug {
		log.Printf("[GO FONDY] Request: %v\n", string(jsonValue))
	}

	req, err := http.NewRequest(methodPost, url.String(), bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %w", err)
	}

	req.Header.Set("User-Agent", "GOFONDY/"+consts.Version)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", requestID)
	req.Header.Set("X-API-Version", "1.0")

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot send request: %w", err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("cannot close response body: %v", err)
		}
	}()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response: %w", err)
	}

	if m.options.IsDebug {
		log.Printf("[GO FONDY] Response: %v\n", string(raw))
	}

	return &raw, nil
}
