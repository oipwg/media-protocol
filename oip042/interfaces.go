package oip042

type IOip042 interface {
	IOip042()
}

type OipAction interface {
	Validate(context OipContext) (OipAction, error)
	Store(context OipContext) error
}
