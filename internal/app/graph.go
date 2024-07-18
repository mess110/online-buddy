package app

func NewFriendGraph() map[string][]string {
	result := map[string][]string{}

	result["kiki"] = []string{"felix", "branzi"}
	result["branzi"] = []string{"felix", "kiki"}
	result["felix"] = []string{"branzi", "kiki"}

	return result
}
