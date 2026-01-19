package cli

func (o *rootOptions) mode() outputMode {
	if o == nil {
		return outputHuman
	}
	if o.noOutput {
		return outputPlain
	}
	if o.json {
		return outputJSON
	}
	if o.plain {
		return outputPlain
	}
	mode, err := parseOutputMode(o.output)
	if err != nil {
		return outputHuman
	}
	return mode
}

func (o *rootOptions) isJSON() bool  { return o.mode() == outputJSON }
func (o *rootOptions) isPlain() bool { return o.mode() == outputPlain }
