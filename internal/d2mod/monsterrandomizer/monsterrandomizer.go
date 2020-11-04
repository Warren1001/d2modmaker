package monsterrandomizer

import (
	"math/rand"
	"time"
	"fmt"

	"github.com/tlentz/d2modmaker/internal/d2mod/config"

	"github.com/tlentz/d2modmaker/internal/d2fs"
	"github.com/tlentz/d2modmaker/internal/d2fs/txts/levels"
	"internal/d2fs/txts/monStats"
)

func Run(cfg *config.Data, d2files d2fs.Files) {
	
	fmt.Println("called")
	
	opts := getRandomOptions(cfg)
	rand.Seed(opts.Seed)
	
	levelsF := d2files.Get(levels.FileName)
	monStatsF := d2files.Get(monStats.FileName)
	
	for rowIdx := range levelsF.Rows {
		
		fmt.Println("modifying rowIdx=", rowIdx)
		
		for i := 0; i < 10; i++ {
				
			colIdx := levels.Mon1 + i;
			oldVal := levelsF.Rows[rowIdx][colIdx]
			randMon := rand.Intn(105)
			newVal := monStatsF.Rows[randMon][monStats.Id]
			levelsF.Rows[rowIdx][colIdx] = newVal
			
			fmt.Println("Actual oldVal=", oldVal, "/ newVal=", newVal)
			
		}		
		
	}
}

func getRandomOptions(cfg *config.Data) config.RandomOptions {
	defaultCfg := config.RandomOptions{
		Seed:     time.Now().UnixNano(),
		MinProps: -1,
		MaxProps: -1,
	}
	if cfg.RandomOptions.Seed >= 0 {
		defaultCfg.Seed = cfg.RandomOptions.Seed
	}
	defaultCfg.IsBalanced = cfg.RandomOptions.IsBalanced
	if cfg.RandomOptions.MaxProps >= 0 {
		defaultCfg.MaxProps = cfg.RandomOptions.MaxProps
	}
	if cfg.RandomOptions.MinProps >= 0 {
		defaultCfg.MinProps = cfg.RandomOptions.MinProps
	}
	defaultCfg.PerfectProps = cfg.RandomOptions.PerfectProps
	defaultCfg.UseOSkills = cfg.RandomOptions.UseOSkills
	defaultCfg.RandomizeMonsters = cfg.RandomOptions.RandomizeMonsters

	cfg.RandomOptions.Seed = defaultCfg.Seed
	return defaultCfg
}
