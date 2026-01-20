package cli

func (o *rootOptions) isAgent() bool {
	if o == nil {
		return false
	}
	return o.agent
}
