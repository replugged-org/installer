package middle

// FinderLocation is a BrowserLocation with the output of CheckDiscordLocation.
// The CCUpdaterUI code put the necessary information into GameInstance but this tool doesn't do that.
// So we use this instead to get the same data.
type FinderLocation struct {
	BrowserLocation
	Instance *DiscordInstance
}

func DiscordFinderVFSList(vfsPath string) []FinderLocation {
	var vfsEntries []FinderLocation
	for _, fi := range BrowserVFSList(vfsPath) {
		if fi.Dir {
			finderLocation := FinderLocation{
				BrowserLocation: fi,
				Instance:        CheckDiscordLocation(fi.Location),
			}
			vfsEntries = append(vfsEntries, finderLocation)
		}
	}

	return vfsEntries
}
