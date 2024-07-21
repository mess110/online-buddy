package datatypes

/*
This struct holds the relationships between users, in a real app
this would be stored somewhere else. For the sake of simplicity,
I opted for a hardcoded map stored in memory
*/
type FriendGraph struct {
	friends map[string][]string
}

func (f *FriendGraph) GetAllFriends(key string) []string {
	return f.friends[key]
}

func (f *FriendGraph) GetAll() map[string][]string {
	return f.friends
}

func NewFriendGraph() *FriendGraph {
	result := map[string][]string{}

	result["kiki"] = []string{"felix", "branzi"}
	result["branzi"] = []string{"felix", "kiki"}
	result["felix"] = []string{"branzi", "kiki", "cata", "spiri"}
	result["cata"] = []string{"spiri", "felix"}
	result["spiri"] = []string{"cata", "felix"}
	result["horea"] = []string{"branzi"}

	return &FriendGraph{
		friends: result,
	}
}
