package core

type Vote map[string]VoteCount

type VoteCount struct {
	N             int                    `json:"n"`
	VoteDetailMap map[string]interface{} `json:"vote_detail"`
}

func NewVote() Vote {
	return make(map[string]VoteCount)
}

func (v Vote) addVote(voteKey string, address string, detail interface{}) {
	if _, ok := v[voteKey]; !ok {
		detailMap := make(map[string]interface{})
		detailMap[address] = detail
		v[voteKey] = VoteCount{
			N:             1,
			VoteDetailMap: detailMap,
		}
	} else {
		count := v[voteKey]
		if _, ok := count.VoteDetailMap[address]; !ok {
			count.VoteDetailMap[address] = detail
			count.N++
		} 	
		
		v[voteKey] = count
	}
}

func (v Vote) getVoteNum(voteKey string) int {
	if _, ok := v[voteKey]; ok {
		return v[voteKey].N
	} else {
		return 0
	}
}
