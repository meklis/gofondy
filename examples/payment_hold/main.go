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

package main

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/stremovskyy/gofondy"
	"github.com/stremovskyy/gofondy/consts"
	"github.com/stremovskyy/gofondy/examples"
	"github.com/stremovskyy/gofondy/models"
)

func main() {
	fondyGateway := gofondy.New(models.DefaultOptions())

	merchAccount := &models.MerchantAccount{
		MerchantID:     examples.MerchantId,
		MerchantKey:    examples.MerchantKey,
		MerchantString: "Test Merchant",
	}

	invoiceId := uuid.New()

	holdAmount := float64(3)

	invoiceRequest := &models.InvoiceRequest{
		InvoiceID: invoiceId,
		Merchant:  merchAccount,
		Amount:    holdAmount,
	}

	paymentByToken, err := fondyGateway.Hold(invoiceRequest)
	if err != nil {
		log.Fatal(err)
	}

	if *paymentByToken.ResponseStatus == consts.FondyResponseStatusSuccess {
		fmt.Printf("Order (%s) status: %s\n", paymentByToken.OrderID, *paymentByToken.OrderStatus)
	} else {
		fmt.Printf("Error: %s\n", *paymentByToken.ErrorMessage)
	}
}
