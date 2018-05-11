package oip042

import (
	"errors"
)

var ErrNotImplemented = errors.New("not implemented")

func (p Publish) Validate(ctx OipContext) (OipAction, error) {
	if p.Artifact != nil {
		return p.Artifact.Validate(ctx)
	}
	return nil, ErrNotImplemented
}

func (r Register) Validate(ctx OipContext) (OipAction, error) {
	if r.Autominer != nil {
		return r.Autominer.Validate(ctx)
	}
	if r.AutominerPool != nil {
		return r.AutominerPool.Validate(ctx)
	}
	if r.Influencer != nil {
		return r.Influencer.Validate(ctx)
	}
	if r.Pub != nil {
		return r.Pub.Validate(ctx)
	}
	if r.Platform != nil {
		return r.Platform.Validate(ctx)
	}
	return nil, ErrNotImplemented
}

func (e Edit) Validate(ctx OipContext) (OipAction, error) {
	if e.Autominer != nil {
		return e.Autominer.Validate(ctx)
	}
	if e.AutominerPool != nil {
		return e.AutominerPool.Validate(ctx)
	}
	if e.Influencer != nil {
		return e.Influencer.Validate(ctx)
	}
	if e.Pub != nil {
		return e.Pub.Validate(ctx)
	}
	if e.Platform != nil {
		return e.Platform.Validate(ctx)
	}
	if e.Artifact != nil {
		return e.Artifact.Validate(ctx)
	}
	return nil, ErrNotImplemented
}

func (d Deactivate) Validate(ctx OipContext) (OipAction, error) {
	//if d.Autominer != nil {
	//    return d.Autominer.Validate(ctx)
	//}
	//if d.AutominerPool != nil {
	//    return d.AutominerPool.Validate(ctx)
	//}
	//if d.Influencer != nil {
	//    return d.Influencer.Validate(ctx)
	//}
	if d.Pub != nil {
		return d.Pub.Validate(ctx)
	}
	//if d.Platform != nil {
	//    return d.Platform.Validate(ctx)
	//}
	if d.Artifact != nil {
		return d.Artifact.Validate(ctx)
	}
	return nil, ErrNotImplemented
}

func (t Transfer) Validate(ctx OipContext) (OipAction, error) {
	//if t.Autominer != nil {
	//    return t.Autominer.Validate(ctx)
	//}
	//if t.AutominerPool != nil {
	//    return t.AutominerPool.Validate(ctx)
	//}
	//if t.Influencer != nil {
	//    return t.Influencer.Validate(ctx)
	//}
	//if t.Pub != nil {
	//    return t.Pub.Validate(ctx)
	//}
	//if t.Platform != nil {
	//    return t.Platform.Validate(ctx)
	//}
	if t.Artifact != nil {
		return t.Artifact.Validate(ctx)
	}
	return nil, ErrNotImplemented
}
