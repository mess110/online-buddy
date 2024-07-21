package app

func NewFriendGraph() map[string][]string {
	result := map[string][]string{}

	result["kiki"] = []string{"felix", "branzi"}
	result["branzi"] = []string{"felix", "kiki"}
	result["felix"] = []string{"branzi", "kiki", "cata", "spiri"}
	result["cata"] = []string{"spiri", "felix"}
	result["spiri"] = []string{"cata", "felix"}

	result["5af4f9ea-c543-4a11-a384-78bcd681f8ff"] = []string{"f2bbc332-ea80-4ba2-a413-040c3e2ed1d8", "8f16b837-0ad0-47e8-9935-032f87047efb"}
	result["f2bbc332-ea80-4ba2-a413-040c3e2ed1d8"] = []string{"5af4f9ea-c543-4a11-a384-78bcd681f8ff", "8f16b837-0ad0-47e8-9935-032f87047efb"}
	result["8f16b837-0ad0-47e8-9935-032f87047efb"] = []string{"5af4f9ea-c543-4a11-a384-78bcd681f8ff"}

	return result
}
