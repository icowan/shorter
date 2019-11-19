/**
 * @Time : 19/11/2019 10:13 AM
 * @Author : solacowa@gmail.com
 * @File : model
 * @Software: GoLand
 */

package service

import "time"

type Redirect struct {
	Code      string    `json:"code"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
}
