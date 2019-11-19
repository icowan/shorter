/**
 * @Time : 19/11/2019 10:14 AM
 * @Author : solacowa@gmail.com
 * @File : repository
 * @Software: GoLand
 */

package service

type Repository interface {
	Find(code string) (redirect *Redirect, err error)
	Store(redirect *Redirect) error
}
