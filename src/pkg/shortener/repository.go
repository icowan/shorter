package shortener

type Repository interface {
	Find(code string) (r *Redirect, err error)
	Store(redirect *Redirect) error
}
