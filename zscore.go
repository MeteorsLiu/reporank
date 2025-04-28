package main

var (
	mean = map[string]float64{
		"stars":        1.132366,
		"forks":        0.451522,
		"watchers":     1.199872,
		"networks":     0.101942,
		"issues":       0.000105,
		"contributors": 1.192904,
		"commit_rate":  0.159262,
	}

	std = map[string]float64{
		"stars":        0.990616,
		"forks":        0.844213,
		"watchers":     0.571998,
		"networks":     0.684791,
		"issues":       0.002566,
		"contributors": 1.407939,
		"commit_rate":  0.460403,
	}

	weight = map[string]float64{
		"stars":        0.229479,
		"forks":        0.227530,
		"watchers":     0.104085,
		"networks":     0.102685,
		"issues":       0.198170,
		"contributors": 0.913740,
		"commit_rate":  0.000304,
	}
)

type Data struct {
	Key   string
	Value float64
}

func sumData(datas []Data) (sum float64) {
	for _, data := range datas {
		sum += zScore(data) * weight[data.Key]
	}
	return
}

func zScore(data Data) float64 {
	return (data.Value - mean[data.Key]) / std[data.Key]
}
