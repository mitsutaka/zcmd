package sync

import "github.com/mitsutaka/zcmd"

func findTargetSyncs(cfgs []*zcmd.SyncInfo, args []string) []*zcmd.SyncInfo {
	if len(args) == 0 {
		// Sync all paths
		return cfgs
	}

	targetCfgs := make([]*zcmd.SyncInfo, 0)

	for _, cfg := range cfgs {
		for _, arg := range args {
			if cfg.Name == arg {
				targetCfgs = append(targetCfgs, cfg)
				break
			}
		}
	}
	return targetCfgs
}
